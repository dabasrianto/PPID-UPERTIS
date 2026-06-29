package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ─── Config (from env with defaults) ─────────────────────────────────────────

func otpExpiryMinutes() int {
	v, _ := strconv.Atoi(os.Getenv("AUTH_OTP_EXPIRY_MINUTES"))
	if v <= 0 {
		return 5
	}
	return v
}

func otpMaxAttempts() int {
	v, _ := strconv.Atoi(os.Getenv("AUTH_OTP_MAX_ATTEMPTS"))
	if v <= 0 {
		return 5
	}
	return v
}

func trustedDeviceDays() int {
	v, _ := strconv.Atoi(os.Getenv("AUTH_TRUSTED_DEVICE_DAYS"))
	if v <= 0 {
		return 30
	}
	return v
}

func isOTPDisabled() bool {
	return os.Getenv("AUTH_DISABLE_OTP") == "true"
}

// ─── Auth OTP Tables Migration ───────────────────────────────────────────────

func migrateAuthOTPTables() {
	db.ExecContext(context.Background(), `ALTER TABLE users ADD COLUMN IF NOT EXISTS phone_whatsapp VARCHAR(50)`)

	db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS auth_otps (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			otp_hash VARCHAR(255) NOT NULL,
			purpose VARCHAR(50) DEFAULT 'login',
			expires_at TIMESTAMP NOT NULL,
			attempts INTEGER DEFAULT 0,
			consumed_at TIMESTAMP,
			ip_address VARCHAR(50),
			user_agent VARCHAR(500),
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_auth_otps_user ON auth_otps(user_id, purpose, consumed_at)`)

	db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS trusted_devices (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			device_token VARCHAR(255) UNIQUE NOT NULL,
			device_name VARCHAR(255),
			ua_hash VARCHAR(64),
			last_used_at TIMESTAMP DEFAULT NOW(),
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	db.ExecContext(context.Background(), `ALTER TABLE trusted_devices ADD COLUMN IF NOT EXISTS ua_hash VARCHAR(64)`)
	db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_trusted_devices_token ON trusted_devices(device_token)`)
	db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_trusted_devices_user ON trusted_devices(user_id)`)

	// Auth audit log
	db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS auth_audit_log (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id VARCHAR(255),
			event VARCHAR(50) NOT NULL,
			ip_address VARCHAR(50),
			user_agent VARCHAR(500),
			details TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_auth_audit_user ON auth_audit_log(user_id, created_at DESC)`)
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func generateOTP6() string {
	const digits = "0123456789"
	result := make([]byte, 6)
	for i := 0; i < 6; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		result[i] = digits[num.Int64()]
	}
	return string(result)
}

func hashOTP(otp string) string {
	h := sha256.Sum256([]byte(otp))
	return hex.EncodeToString(h[:])
}

func hashUA(ua string) string {
	h := sha256.Sum256([]byte(ua))
	return hex.EncodeToString(h[:16]) // short hash
}

func logAuthEvent(userID, event, ip, ua, details string) {
	db.ExecContext(context.Background(),
		`INSERT INTO auth_audit_log (user_id, event, ip_address, user_agent, details) VALUES ($1, $2, $3, $4, $5)`,
		userID, event, ip, truncate(ua, 200), details,
	)
	zlog.Info().
		Str("user_id", userID).
		Str("event", event).
		Str("ip", ip).
		Str("details", details).
		Msg("auth_audit")
}

// ─── Login with 2FA ──────────────────────────────────────────────────────────

func loginWithOTP(c *fiber.Ctx) error {
	var req struct {
		Email       string `json:"email" form:"email"`
		Password    string `json:"password" form:"password"`
		DeviceToken string `json:"device_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
	}

	ip := c.IP()
	ua := c.Get("User-Agent")

	// Validate credentials
	var userID, email, passwordHash, fullName, role, status, phoneWA string
	err := db.QueryRowContext(context.Background(),
		"SELECT id, email, password_hash, full_name, role, COALESCE(status,'active'), COALESCE(phone_whatsapp,'') FROM users WHERE email = $1",
		req.Email,
	).Scan(&userID, &email, &passwordHash, &fullName, &role, &status, &phoneWA)

	credentialsValid := false
	var upProfile *UpertisProfile

	if err == nil && checkPassword(passwordHash, req.Password) {
		credentialsValid = true
	} else {
		// Try authenticating via live UPERTIS website
		profile, dbErr := authenticateUpertisLive(req.Email, req.Password)
		if dbErr == nil && profile != nil {
			credentialsValid = true
			upProfile = profile

			// Map roles: if target has admin role, map to admin, otherwise map accordingly
			assignedRole := upProfile.Role

			if userID != "" {
				// User exists locally, update their password hash and display name/role
				_, _ = db.ExecContext(context.Background(),
					"UPDATE users SET password_hash = $1, full_name = $2, role = $3, status = $4, phone_whatsapp = $5, updated_at = NOW() WHERE id = $6",
					upProfile.PasswordHash, upProfile.FullName, assignedRole, upProfile.Status, upProfile.PhoneWhatsapp, userID,
				)
				passwordHash = upProfile.PasswordHash
				fullName = upProfile.FullName
				role = assignedRole
				status = upProfile.Status
				phoneWA = upProfile.PhoneWhatsapp
			} else {
				// Create new local user
				_, dbErr := db.ExecContext(context.Background(),
					"INSERT INTO users (id, email, password_hash, full_name, role, status, phone_whatsapp) VALUES ($1, $2, $3, $4, $5, $6, $7)",
					upProfile.ID, upProfile.Email, upProfile.PasswordHash, upProfile.FullName, assignedRole, upProfile.Status, upProfile.PhoneWhatsapp,
				)
				if dbErr == nil {
					userID = upProfile.ID
					email = upProfile.Email
					fullName = upProfile.FullName
					role = assignedRole
					status = upProfile.Status
					phoneWA = upProfile.PhoneWhatsapp
				} else {
					logAuthEvent("", "login_failed", ip, ua, "jit_provision_error="+dbErr.Error())
					return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyinkronkan akun lokal"})
				}
			}
		}
	}

	if !credentialsValid {
		logAuthEvent("", "login_failed", ip, ua, "email="+req.Email)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Kredensial salah (Lokal & UPERTIS)"})
	}

	if status == "pending" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Akun Anda sedang menunggu persetujuan admin."})
	}
	if status == "rejected" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Akun Anda telah ditolak."})
	}

	// OTP disabled globally
	if isOTPDisabled() {
		logAuthEvent(userID, "login_success", ip, ua, "otp_disabled")
		return issueLoginToken(c, userID, role, fullName, email)
	}

	// Non-admin roles skip OTP
	if role != "admin" && role != "faculty_admin" {
		logAuthEvent(userID, "login_success", ip, ua, "non_admin_role")
		return issueLoginToken(c, userID, role, fullName, email)
	}

	// Check trusted device (bound to user + token + UA hash)
	if req.DeviceToken != "" {
		var deviceID string
		currentUAHash := hashUA(ua)
		err := db.QueryRowContext(context.Background(),
			`SELECT id FROM trusted_devices 
			 WHERE device_token = $1 AND user_id = $2 AND expires_at > NOW()
			 AND (ua_hash IS NULL OR ua_hash = $3)`,
			req.DeviceToken, userID, currentUAHash,
		).Scan(&deviceID)
		if err == nil {
			db.ExecContext(context.Background(), "UPDATE trusted_devices SET last_used_at = NOW() WHERE id = $1", deviceID)
			logAuthEvent(userID, "login_trusted_device", ip, ua, "device="+deviceID)
			return issueLoginToken(c, userID, role, fullName, email)
		}
	}

	// No WhatsApp number — skip OTP
	if phoneWA == "" {
		logAuthEvent(userID, "login_success", ip, ua, "no_whatsapp_configured")
		return issueLoginToken(c, userID, role, fullName, email)
	}

	// Cooldown check (60s)
	var lastOTPCreated time.Time
	err = db.QueryRowContext(context.Background(),
		"SELECT created_at FROM auth_otps WHERE user_id = $1 AND purpose = 'login' AND consumed_at IS NULL ORDER BY created_at DESC LIMIT 1",
		userID,
	).Scan(&lastOTPCreated)
	if err == nil && time.Since(lastOTPCreated) < 60*time.Second {
		return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
			"error":       "Tunggu 60 detik sebelum meminta OTP baru",
			"retry_after": 60 - int(time.Since(lastOTPCreated).Seconds()),
		})
	}

	// Generate OTP
	otp := generateOTP6()
	otpHash := hashOTP(otp)
	expiresAt := time.Now().Add(time.Duration(otpExpiryMinutes()) * time.Minute)

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO auth_otps (user_id, otp_hash, purpose, expires_at, ip_address, user_agent)
		 VALUES ($1, $2, 'login', $3, $4, $5)`,
		userID, otpHash, expiresAt, ip, truncate(ua, 200),
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Send via WhatsApp (async-safe with timeout)
	message := fmt.Sprintf("Kode OTP login UPERTIS CMS Anda: *%s*\n\nBerlaku %d menit. Jangan bagikan kode ini kepada siapa pun.", otp, otpExpiryMinutes())
	waErr := SendWhatsAppMessage(phoneWA, message)
	if waErr != nil {
		logAuthEvent(userID, "otp_send_failed", ip, ua, waErr.Error())
		// Fallback: allow login if WA gateway is down (production safety)
		return issueLoginToken(c, userID, role, fullName, email)
	}

	logAuthEvent(userID, "otp_sent", ip, ua, "")

	return c.JSON(fiber.Map{
		"status":  "otp_required",
		"user_id": userID,
		"message": "Kode OTP telah dikirim ke WhatsApp Anda",
	})
}

// verifyLoginOTP verifies OTP with atomic consume (replay protection)
func verifyLoginOTP(c *fiber.Ctx) error {
	var req struct {
		UserID      string `json:"user_id"`
		OTP         string `json:"otp"`
		TrustDevice bool   `json:"trust_device"`
		DeviceName  string `json:"device_name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.UserID == "" || req.OTP == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "User ID and OTP are required"})
	}

	ip := c.IP()
	ua := c.Get("User-Agent")
	maxAttempts := otpMaxAttempts()

	// Atomic: find + increment attempts in one query (prevents race condition)
	var otpID, storedHash string
	var expiresAt time.Time
	var attempts int
	err := db.QueryRowContext(context.Background(),
		`UPDATE auth_otps SET attempts = attempts + 1
		 WHERE id = (
		   SELECT id FROM auth_otps
		   WHERE user_id = $1 AND purpose = 'login' AND consumed_at IS NULL
		   ORDER BY created_at DESC LIMIT 1
		   FOR UPDATE
		 )
		 RETURNING id, otp_hash, expires_at, attempts`,
		req.UserID,
	).Scan(&otpID, &storedHash, &expiresAt, &attempts)

	if err != nil {
		logAuthEvent(req.UserID, "otp_verify_failed", ip, ua, "no_pending_otp")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Kode OTP tidak ditemukan atau sudah kedaluwarsa"})
	}

	// Check expiry
	if time.Now().After(expiresAt) {
		db.ExecContext(context.Background(), "UPDATE auth_otps SET consumed_at = NOW() WHERE id = $1", otpID)
		logAuthEvent(req.UserID, "otp_expired", ip, ua, "")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Kode OTP sudah kedaluwarsa"})
	}

	// Check max attempts (attempts already incremented)
	if attempts > maxAttempts {
		db.ExecContext(context.Background(), "UPDATE auth_otps SET consumed_at = NOW() WHERE id = $1", otpID)
		logAuthEvent(req.UserID, "otp_max_attempts", ip, ua, fmt.Sprintf("attempts=%d", attempts))
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Terlalu banyak percobaan. Silakan minta OTP baru."})
	}

	// Verify hash
	if hashOTP(req.OTP) != storedHash {
		logAuthEvent(req.UserID, "otp_verify_failed", ip, ua, fmt.Sprintf("attempt=%d", attempts))
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":              "Kode OTP salah",
			"attempts_remaining": maxAttempts - attempts,
		})
	}

	// Atomic consume (prevents replay — only consumes if not already consumed)
	res, _ := db.ExecContext(context.Background(),
		"UPDATE auth_otps SET consumed_at = NOW() WHERE id = $1 AND consumed_at IS NULL", otpID)
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		// Already consumed by parallel request — replay attack blocked
		logAuthEvent(req.UserID, "otp_replay_blocked", ip, ua, "")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Kode OTP sudah digunakan"})
	}

	// Get user info
	var role, fullName, email string
	db.QueryRowContext(context.Background(),
		"SELECT role, full_name, email FROM users WHERE id = $1", req.UserID,
	).Scan(&role, &fullName, &email)

	// Create trusted device
	var deviceToken string
	if req.TrustDevice {
		deviceToken = uuid.New().String()
		deviceName := req.DeviceName
		if deviceName == "" {
			deviceName = truncate(ua, 100)
		}
		db.ExecContext(context.Background(),
			`INSERT INTO trusted_devices (user_id, device_token, device_name, ua_hash, expires_at)
			 VALUES ($1, $2, $3, $4, $5)`,
			req.UserID, deviceToken, deviceName, hashUA(ua),
			time.Now().Add(time.Duration(trustedDeviceDays())*24*time.Hour),
		)
		logAuthEvent(req.UserID, "trusted_device_added", ip, ua, "device="+deviceToken[:8])
	}

	logAuthEvent(req.UserID, "login_success", ip, ua, "otp_verified")

	// Issue token
	token, err := generateToken(req.UserID, role)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	response := fiber.Map{
		"status": "success",
		"token":  token,
		"user": fiber.Map{
			"id":        req.UserID,
			"email":     email,
			"full_name": fullName,
			"role":      role,
		},
	}
	if deviceToken != "" {
		response["device_token"] = deviceToken
	}

	return c.JSON(response)
}

// resendLoginOTP resends OTP with cooldown
func resendLoginOTP(c *fiber.Ctx) error {
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&req); err != nil || req.UserID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
	}

	// Cooldown
	var lastCreated time.Time
	err := db.QueryRowContext(context.Background(),
		"SELECT created_at FROM auth_otps WHERE user_id = $1 AND purpose = 'login' ORDER BY created_at DESC LIMIT 1",
		req.UserID,
	).Scan(&lastCreated)
	if err == nil && time.Since(lastCreated) < 60*time.Second {
		remaining := 60 - int(time.Since(lastCreated).Seconds())
		return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
			"error":       fmt.Sprintf("Tunggu %d detik", remaining),
			"retry_after": remaining,
		})
	}

	// Get phone
	var phoneWA string
	err = db.QueryRowContext(context.Background(),
		"SELECT COALESCE(phone_whatsapp,'') FROM users WHERE id = $1", req.UserID,
	).Scan(&phoneWA)
	if err != nil || phoneWA == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "WhatsApp number not configured"})
	}

	// Invalidate old
	db.ExecContext(context.Background(),
		"UPDATE auth_otps SET consumed_at = NOW() WHERE user_id = $1 AND purpose = 'login' AND consumed_at IS NULL",
		req.UserID,
	)

	// New OTP
	otp := generateOTP6()
	expiresAt := time.Now().Add(time.Duration(otpExpiryMinutes()) * time.Minute)

	db.ExecContext(context.Background(),
		`INSERT INTO auth_otps (user_id, otp_hash, purpose, expires_at, ip_address, user_agent)
		 VALUES ($1, $2, 'login', $3, $4, $5)`,
		req.UserID, hashOTP(otp), expiresAt, c.IP(), truncate(c.Get("User-Agent"), 200),
	)

	message := fmt.Sprintf("Kode OTP login UPERTIS CMS Anda: *%s*\n\nBerlaku %d menit. Jangan bagikan kode ini kepada siapa pun.", otp, otpExpiryMinutes())
	if err := SendWhatsAppMessage(phoneWA, message); err != nil {
		logAuthEvent(req.UserID, "otp_resend_failed", c.IP(), c.Get("User-Agent"), err.Error())
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{"error": "Gagal mengirim OTP"})
	}

	logAuthEvent(req.UserID, "otp_resent", c.IP(), c.Get("User-Agent"), "")
	return c.JSON(fiber.Map{"message": "OTP baru telah dikirim"})
}

// revokeAllTrustedDevices removes all trusted devices for a user
func revokeAllTrustedDevices(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	res, err := db.ExecContext(context.Background(),
		"DELETE FROM trusted_devices WHERE user_id = $1", userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to revoke devices"})
	}

	count, _ := res.RowsAffected()
	logAuthEvent(fmt.Sprintf("%v", userID), "devices_revoked", c.IP(), c.Get("User-Agent"), fmt.Sprintf("count=%d", count))

	return c.JSON(fiber.Map{"message": fmt.Sprintf("%d perangkat terpercaya dihapus", count)})
}

// getAuthAuditLog returns recent auth events for admin
func getAuthAuditLog(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(),
		`SELECT id, COALESCE(user_id,''), event, COALESCE(ip_address,''), COALESCE(details,''), created_at
		 FROM auth_audit_log ORDER BY created_at DESC LIMIT 100`)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch audit log"})
	}
	defer rows.Close()

	var logs []fiber.Map
	for rows.Next() {
		var id, userID, event, ip, details, createdAt string
		if rows.Scan(&id, &userID, &event, &ip, &details, &createdAt) == nil {
			logs = append(logs, fiber.Map{
				"id": id, "user_id": userID, "event": event,
				"ip": ip, "details": details, "created_at": createdAt,
			})
		}
	}
	return c.JSON(logs)
}

// getTrustedDevices returns user's trusted devices
func getTrustedDevices(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	rows, err := db.QueryContext(context.Background(),
		`SELECT id, device_name, last_used_at, expires_at, created_at
		 FROM trusted_devices WHERE user_id = $1 AND expires_at > NOW()
		 ORDER BY last_used_at DESC`, userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch devices"})
	}
	defer rows.Close()

	var devices []fiber.Map
	for rows.Next() {
		var id, name, lastUsed, expires, created string
		if rows.Scan(&id, &name, &lastUsed, &expires, &created) == nil {
			devices = append(devices, fiber.Map{
				"id": id, "device_name": name, "last_used_at": lastUsed,
				"expires_at": expires, "created_at": created,
			})
		}
	}
	return c.JSON(devices)
}

// issueLoginToken generates JWT and returns login success
func issueLoginToken(c *fiber.Ctx, userID, role, fullName, email string) error {
	token, err := generateToken(userID, role)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"token":  token,
		"user": fiber.Map{
			"id":        userID,
			"email":     email,
			"full_name": fullName,
			"role":      role,
		},
	})
}

// resendRegisterOTP resends registration OTP with cooldown
func resendRegisterOTP(c *fiber.Ctx) error {
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&req); err != nil || req.UserID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
	}

	// Cooldown (60s)
	var lastCreated time.Time
	err := db.QueryRowContext(context.Background(),
		"SELECT created_at FROM auth_otps WHERE user_id = $1 AND purpose = 'register' ORDER BY created_at DESC LIMIT 1",
		req.UserID,
	).Scan(&lastCreated)
	if err == nil && time.Since(lastCreated) < 60*time.Second {
		remaining := 60 - int(time.Since(lastCreated).Seconds())
		return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
			"error":       fmt.Sprintf("Tunggu %d detik", remaining),
			"retry_after": remaining,
		})
	}

	// Get phone and role (to customize message: candidate vs user)
	var phoneWA, role string
	err = db.QueryRowContext(context.Background(),
		"SELECT COALESCE(phone_whatsapp,''), role FROM users WHERE id = $1", req.UserID,
	).Scan(&phoneWA, &role)
	if err != nil || phoneWA == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "WhatsApp number not configured"})
	}

	// Invalidate old
	db.ExecContext(context.Background(),
		"UPDATE auth_otps SET consumed_at = NOW() WHERE user_id = $1 AND purpose = 'register' AND consumed_at IS NULL",
		req.UserID,
	)

	// New OTP
	otp := generateOTP6()
	expiresAt := time.Now().Add(time.Duration(otpExpiryMinutes()) * time.Minute)

	db.ExecContext(context.Background(),
		`INSERT INTO auth_otps (user_id, otp_hash, purpose, expires_at, ip_address, user_agent)
		 VALUES ($1, $2, 'register', $3, $4, $5)`,
		req.UserID, hashOTP(otp), expiresAt, c.IP(), truncate(c.Get("User-Agent"), 200),
	)

	var message string
	if role == "candidate" {
		message = fmt.Sprintf("Kode verifikasi akun PMB UPERTIS Anda: *%s*\n\nBerlaku %d menit. Jangan bagikan kode ini.", otp, otpExpiryMinutes())
	} else {
		message = fmt.Sprintf("Kode verifikasi akun UPERTIS Anda: *%s*\n\nBerlaku %d menit. Jangan bagikan kode ini.", otp, otpExpiryMinutes())
	}

	if err := SendWhatsAppMessage(phoneWA, message); err != nil {
		logAuthEvent(req.UserID, "otp_resend_failed", c.IP(), c.Get("User-Agent"), err.Error())
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{"error": "Gagal mengirim OTP"})
	}

	logAuthEvent(req.UserID, "otp_resent", c.IP(), c.Get("User-Agent"), "")
	return c.JSON(fiber.Map{"message": "OTP baru telah dikirim"})
}

// ─── UPERTIS Database Auth Helpers ──────────────────────────────────────────

// ─── UPERTIS Live Website Auth Helpers ──────────────────────────────────────

type UpertisProfile struct {
	ID            string
	Email         string
	PasswordHash  string
	FullName      string
	Role          string
	Status        string
	PhoneWhatsapp string
}

type LiveAuthResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
	User   struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		FullName string `json:"full_name"`
		Role     string `json:"role"`
	} `json:"user"`
}

func authenticateUpertisLive(email, password string) (*UpertisProfile, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	loginPayload, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.Post("https://upertis.ac.id/api/v1/auth/login", "application/json", bytes.NewBuffer(loginPayload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid credentials on live server")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var authResp LiveAuthResponse
	if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
		return nil, err
	}

	if authResp.User.ID == "" {
		return nil, fmt.Errorf("empty user profile from live server")
	}

	hashedPassword, _ := hashPassword(password)

	return &UpertisProfile{
		ID:            authResp.User.ID,
		Email:         authResp.User.Email,
		PasswordHash:  hashedPassword,
		FullName:      authResp.User.FullName,
		Role:          authResp.User.Role,
		Status:        "active",
		PhoneWhatsapp: "",
	}, nil
}
