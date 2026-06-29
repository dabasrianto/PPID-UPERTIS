package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// Token blacklist for logout (in-memory, cleared on restart)
// TECHNICAL DEBT: This blacklist is lost on backend restart.
// Revoked tokens become valid again until they naturally expire (24h).
// Acceptable for: single VPS, single instance, low-medium traffic.
// Future upgrade path: Redis-backed blacklist or refresh token rotation.
var tokenBlacklist = struct {
	sync.RWMutex
	tokens map[string]time.Time
}{tokens: make(map[string]time.Time)}

func initJWTSecret() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("FATAL: JWT_SECRET environment variable is required. Set it in your .env file.")
	}
	if len(secret) < 32 {
		log.Fatal("FATAL: JWT_SECRET must be at least 32 characters long for production security")
	}
	jwtSecret = []byte(secret)

	// Start background goroutine to clean expired tokens from blacklist
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			cleanExpiredBlacklistTokens()
		}
	}()
}

func cleanExpiredBlacklistTokens() {
	tokenBlacklist.Lock()
	defer tokenBlacklist.Unlock()
	now := time.Now()
	for token, expiry := range tokenBlacklist.tokens {
		if now.After(expiry) {
			delete(tokenBlacklist.tokens, token)
		}
	}
}

func blacklistToken(tokenString string, expiry time.Time) {
	tokenBlacklist.Lock()
	defer tokenBlacklist.Unlock()
	tokenBlacklist.tokens[tokenString] = expiry
}

func isTokenBlacklisted(tokenString string) bool {
	tokenBlacklist.RLock()
	defer tokenBlacklist.RUnlock()
	_, exists := tokenBlacklist.tokens[tokenString]
	return exists
}

func generateToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func authMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		log.Printf("authMiddleware: Missing authorization header")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing authorization header",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" || tokenString == authHeader {
		log.Printf("authMiddleware: Invalid token format")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token format",
		})
	}

	// Check if token has been blacklisted (logged out)
	// TODO: Replace in-memory blacklist with Redis when scaling beyond single instance.
	// Current tradeoff: tokens revoked via logout/password-change are lost on restart.
	// Impact: revoked tokens valid again until natural expiry (24h). Acceptable for current scale.
	if isTokenBlacklisted(tokenString) {
		log.Printf("authMiddleware: Token has been revoked")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token has been revoked",
		})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || token == nil || !token.Valid {
		log.Printf("authMiddleware: Invalid or expired token: err=%v, token_valid=%v", err, token != nil && token.Valid)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("authMiddleware: Invalid token claims")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	c.Locals("user_id", claims["user_id"])
	c.Locals("role", claims["role"])
	c.Locals("token_string", tokenString)

	return c.Next()
}

func login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	var user User
	var passwordHash string
	var status string
	err := db.QueryRowContext(context.Background(),
		"SELECT id, email, password_hash, full_name, role, COALESCE(status, 'active') FROM users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Email, &passwordHash, &user.FullName, &user.Role, &status)

	if err != nil || !checkPassword(passwordHash, req.Password) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if status == "pending" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Akun Anda sedang menunggu persetujuan admin.",
		})
	}

	if status == "rejected" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Akun Anda telah ditolak.",
		})
	}

	token, err := generateToken(user.ID, user.Role)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}

func logout(c *fiber.Ctx) error {
	// Extract token from Authorization header and blacklist it
	authHeader := c.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString != "" && tokenString != authHeader {
		// Parse to get expiry time
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if exp, ok := claims["exp"].(float64); ok {
					blacklistToken(tokenString, time.Unix(int64(exp), 0))
				}
			}
		}
	}
	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

func registerUser(c *fiber.Ctx) error {
	var req struct {
		Email         string `json:"email"`
		Password      string `json:"password"`
		FullName      string `json:"full_name"`
		PhoneWhatsapp string `json:"phone_whatsapp"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Email == "" || req.Password == "" || req.FullName == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Email, password, dan nama lengkap wajib diisi"})
	}

	if req.PhoneWhatsapp == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Nomor WhatsApp wajib diisi untuk verifikasi"})
	}

	if len(req.Password) < 6 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Password minimal 6 karakter"})
	}

	hashed, err := hashPassword(req.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	status := "pending_otp"
	if isOTPDisabled() {
		status = "active"
	}

	var newID string
	err = db.QueryRowContext(context.Background(),
		`INSERT INTO users (email, password_hash, full_name, phone_whatsapp, role, status) 
		 VALUES ($1, $2, $3, $4, 'user', $5) RETURNING id`,
		req.Email, hashed, req.FullName, normalizePhoneNumber(req.PhoneWhatsapp), status,
	).Scan(&newID)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Email sudah terdaftar"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mendaftar"})
	}

	if isOTPDisabled() {
		logAuthEvent(newID, "register_success_no_otp", c.IP(), c.Get("User-Agent"), req.PhoneWhatsapp)
		return c.Status(http.StatusCreated).JSON(fiber.Map{
			"message":      "Registrasi berhasil.",
			"user_id":      newID,
			"requires_otp": false,
		})
	}

	// Generate and send OTP via WhatsApp
	otp := generateOTP6()
	otpHash := hashOTP(otp)
	expiresAt := time.Now().Add(time.Duration(otpExpiryMinutes()) * time.Minute)

	db.ExecContext(context.Background(),
		`INSERT INTO auth_otps (user_id, otp_hash, purpose, expires_at, ip_address, user_agent)
		 VALUES ($1, $2, 'register', $3, $4, $5)`,
		newID, otpHash, expiresAt, c.IP(), truncate(c.Get("User-Agent"), 200),
	)

	message := fmt.Sprintf("Kode verifikasi akun UPERTIS Anda: *%s*\n\nBerlaku %d menit. Jangan bagikan kode ini.", otp, otpExpiryMinutes())
	if err := SendWhatsAppMessage(normalizePhoneNumber(req.PhoneWhatsapp), message); err != nil {
		zlog.Warn().Err(err).Str("phone", req.PhoneWhatsapp).Msg("Register OTP send failed")
		// Still return success — user can request resend
	}

	logAuthEvent(newID, "register_otp_sent", c.IP(), c.Get("User-Agent"), req.PhoneWhatsapp)

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":  "Kode OTP telah dikirim ke WhatsApp Anda.",
		"user_id":  newID,
		"requires_otp": true,
	})
}

// verifyRegisterOTP verifies OTP and activates the user account
func verifyRegisterOTP(c *fiber.Ctx) error {
	var req struct {
		UserID string `json:"user_id"`
		OTP    string `json:"otp"`
	}
	if err := c.BodyParser(&req); err != nil || req.UserID == "" || req.OTP == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "user_id dan otp wajib diisi"})
	}

	otpHash := hashOTP(req.OTP)

	// Atomic: find valid OTP and consume it
	var otpID string
	err := db.QueryRowContext(context.Background(),
		`UPDATE auth_otps SET consumed_at = NOW()
		 WHERE id = (
			SELECT id FROM auth_otps
			WHERE user_id = $1 AND purpose = 'register' AND otp_hash = $2
			AND consumed_at IS NULL AND expires_at > NOW() AND attempts < $3
			ORDER BY created_at DESC LIMIT 1
		 ) RETURNING id`,
		req.UserID, otpHash, otpMaxAttempts(),
	).Scan(&otpID)

	if err != nil {
		// Increment attempts
		db.ExecContext(context.Background(),
			`UPDATE auth_otps SET attempts = attempts + 1
			 WHERE user_id = $1 AND purpose = 'register' AND consumed_at IS NULL AND expires_at > NOW()`,
			req.UserID,
		)
		logAuthEvent(req.UserID, "register_otp_failed", c.IP(), c.Get("User-Agent"), "")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Kode OTP salah atau sudah kedaluwarsa"})
	}

	// Activate user
	db.ExecContext(context.Background(),
		`UPDATE users SET status = 'active' WHERE id = $1`, req.UserID)

	logAuthEvent(req.UserID, "register_verified", c.IP(), c.Get("User-Agent"), "")

	// Generate token so user is immediately logged in
	var email, fullName, role string
	db.QueryRowContext(context.Background(),
		`SELECT email, COALESCE(full_name,''), role FROM users WHERE id = $1`, req.UserID,
	).Scan(&email, &fullName, &role)

	token, _ := generateToken(req.UserID, role)

	return c.JSON(fiber.Map{
		"message": "Akun berhasil diverifikasi dan aktif.",
		"token":   token,
		"user": fiber.Map{
			"id": req.UserID, "email": email, "full_name": fullName, "role": role,
		},
	})
}

func getMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var user User
	err := db.QueryRowContext(context.Background(),
		"SELECT id, email, full_name, role, created_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Email, &user.FullName, &user.Role, &user.CreatedAt)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	rows, err := db.QueryContext(context.Background(), "SELECT faculty_id FROM faculty_admins WHERE user_id=$1", user.ID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var facultyID string
			if rows.Scan(&facultyID) == nil {
				user.FacultyIDs = append(user.FacultyIDs, facultyID)
			}
		}
	}

	return c.JSON(user)
}

func requireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}

		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}

func registerCandidate(c *fiber.Ctx) error {
	var req struct {
		FullName      string `json:"full_name"`
		Email         string `json:"email"`
		Password      string `json:"password"`
		PhoneWhatsapp string `json:"phone_whatsapp"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" || req.FullName == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Email, password, and full name are required",
		})
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	defer tx.Rollback()

	status := "pending_otp"
	if isOTPDisabled() {
		status = "active"
	}

	// 1. Create User
	hashed, hashErr := hashPassword(req.Password)
	if hashErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create account"})
	}
	var userID string
	err = tx.QueryRowContext(context.Background(),
		`INSERT INTO users (email, password_hash, full_name, role, status) 
		 VALUES ($1, $2, $3, 'candidate', $4) 
		 RETURNING id`,
		req.Email, hashed, req.FullName, status,
	).Scan(&userID)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
		}
		log.Println("Error creating user:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create account"})
	}

	// 2. Create Candidate Profile
	normalizedPhone := normalizePhoneNumber(req.PhoneWhatsapp)
	_, err = tx.ExecContext(context.Background(),
		`INSERT INTO pmb_candidates (user_id, full_name, phone_whatsapp) 
		 VALUES ($1, $2, $3)`,
		userID, req.FullName, normalizedPhone,
	)
	if err != nil {
		log.Println("Error creating candidate profile:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to setup candidate profile"})
	}

	if isOTPDisabled() {
		if err := tx.Commit(); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}
		logAuthEvent(userID, "register_success_no_otp", c.IP(), c.Get("User-Agent"), normalizedPhone)
		return c.Status(http.StatusCreated).JSON(fiber.Map{
			"message":      "Pendaftaran berhasil.",
			"user_id":      userID,
			"requires_otp": false,
		})
	}

	// 3. Generate and insert OTP
	otp := generateOTP6()
	otpHash := hashOTP(otp)
	expiresAt := time.Now().Add(time.Duration(otpExpiryMinutes()) * time.Minute)

	_, err = tx.ExecContext(context.Background(),
		`INSERT INTO auth_otps (user_id, otp_hash, purpose, expires_at, ip_address, user_agent)
		 VALUES ($1, $2, 'register', $3, $4, $5)`,
		userID, otpHash, expiresAt, c.IP(), truncate(c.Get("User-Agent"), 200),
	)
	if err != nil {
		log.Println("Error creating OTP:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to setup verification"})
	}

	if err := tx.Commit(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	// 4. Send OTP via WhatsApp
	message := fmt.Sprintf("Kode verifikasi akun PMB UPERTIS Anda: *%s*\n\nBerlaku %d menit. Jangan bagikan kode ini.", otp, otpExpiryMinutes())
	if err := SendWhatsAppMessage(normalizedPhone, message); err != nil {
		zlog.Warn().Err(err).Str("phone", normalizedPhone).Msg("Candidate Register OTP send failed")
	}

	logAuthEvent(userID, "register_otp_sent", c.IP(), c.Get("User-Agent"), normalizedPhone)

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":      "Kode OTP verifikasi telah dikirim ke WhatsApp Anda.",
		"user_id":      userID,
		"requires_otp": true,
	})
}

// getAdminProfile returns the current authenticated user's profile
func getAdminProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var user User
	err := db.QueryRowContext(context.Background(),
		"SELECT id, email, full_name, role, COALESCE(status,'active'), created_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Email, &user.FullName, &user.Role, &user.Status, &user.CreatedAt)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

// updateAdminProfile updates the current authenticated user's profile info
func updateAdminProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		FullName string `json:"full_name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	req.FullName = strings.TrimSpace(req.FullName)
	if req.FullName == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Nama lengkap tidak boleh kosong"})
	}

	_, err := db.ExecContext(context.Background(),
		"UPDATE users SET full_name = $1 WHERE id = $2",
		req.FullName, userID,
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui profil"})
	}

	return c.JSON(fiber.Map{
		"message": "Profil berhasil diperbarui",
		"full_name": req.FullName,
	})
}


// validatePasswordPolicy checks minimum password requirements:
// at least 1 uppercase, 1 lowercase, 1 number
func validatePasswordPolicy(password string) error {
	var hasUpper, hasLower, hasNumber bool
	for _, ch := range password {
		switch {
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= 'a' && ch <= 'z':
			hasLower = true
		case ch >= '0' && ch <= '9':
			hasNumber = true
		}
	}
	if !hasUpper {
		return fmt.Errorf("Password harus mengandung minimal 1 huruf besar")
	}
	if !hasLower {
		return fmt.Errorf("Password harus mengandung minimal 1 huruf kecil")
	}
	if !hasNumber {
		return fmt.Errorf("Password harus mengandung minimal 1 angka")
	}
	return nil
}

// changePassword allows authenticated user to change their own password
func changePassword(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	rid, _ := c.Locals(RequestIDKey).(string)

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if req.CurrentPassword == "" || req.NewPassword == "" || req.ConfirmPassword == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Semua field wajib diisi"})
	}

	// Validate new password strength
	if len(req.NewPassword) < 8 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Password baru minimal 8 karakter"})
	}

	// Password policy: at least 1 uppercase, 1 lowercase, 1 number
	if err := validatePasswordPolicy(req.NewPassword); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Reject same password as current
	if req.NewPassword == req.CurrentPassword {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Password baru tidak boleh sama dengan password saat ini"})
	}

	// Validate confirm match
	if req.NewPassword != req.ConfirmPassword {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Konfirmasi password tidak cocok"})
	}

	// Fetch current password hash and user info for name/email check
	var currentHash, email, fullName string
	err := db.QueryRowContext(context.Background(),
		"SELECT password_hash, email, COALESCE(full_name,'') FROM users WHERE id = $1", userID,
	).Scan(&currentHash, &email, &fullName)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Reject password containing email or name
	lowerPw := strings.ToLower(req.NewPassword)
	if email != "" && strings.Contains(lowerPw, strings.ToLower(strings.Split(email, "@")[0])) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Password tidak boleh mengandung email"})
	}
	if fullName != "" {
		nameParts := strings.Fields(strings.ToLower(fullName))
		for _, part := range nameParts {
			if len(part) >= 3 && strings.Contains(lowerPw, part) {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Password tidak boleh mengandung nama"})
			}
		}
	}

	// Verify current password
	if !checkPassword(currentHash, req.CurrentPassword) {
		zlog.Warn().
			Str("request_id", rid).
			Str("user_id", fmt.Sprintf("%v", userID)).
			Msg("Password change failed: wrong current password")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Password saat ini salah"})
	}

	// Hash new password
	newHash, err := hashPassword(req.NewPassword)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Update password in database
	_, err = db.ExecContext(context.Background(),
		"UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2",
		newHash, userID,
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengubah password"})
	}

	// Blacklist current token to force re-login
	tokenString, _ := c.Locals("token_string").(string)
	if tokenString != "" {
		blacklistToken(tokenString, time.Now().Add(24*time.Hour))
	}

	// Log password change event (no plaintext password logged)
	zlog.Info().
		Str("request_id", rid).
		Str("user_id", fmt.Sprintf("%v", userID)).
		Msg("Password changed successfully")

	return c.JSON(fiber.Map{
		"message": "Password berhasil diubah. Silakan login ulang.",
	})
}

type googleTokenInfo struct {
	Email         string      `json:"email"`
	EmailVerified interface{} `json:"email_verified"`
	Name          string      `json:"name"`
	Picture       string      `json:"picture"`
	Aud           string      `json:"aud"`
	Iss           string      `json:"iss"`
	Sub           string      `json:"sub"`
}

func loginWithGoogle(c *fiber.Ctx) error {
	var req struct {
		IDToken string `json:"id_token"`
		Role    string `json:"role"` // "user" or "candidate"
	}

	if err := c.BodyParser(&req); err != nil || req.IDToken == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "id_token is required",
		})
	}

	// Fetch token info from Google API
	tokenURL := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(req.IDToken)
	resp, err := http.Get(tokenURL)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to contact Google API",
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid Google token credentials",
		})
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read Google API response",
		})
	}

	var tokenInfo googleTokenInfo
	if err := json.Unmarshal(bodyBytes, &tokenInfo); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse Google API response",
		})
	}

	// Verify audience matches our GOOGLE_CLIENT_ID (fetch dynamically from site_settings with fallback to env)
	var expectedClientID string
	var settingJSON []byte
	err = db.QueryRowContext(context.Background(), "SELECT value FROM site_settings WHERE key = 'integrations'").Scan(&settingJSON)
	if err == nil {
		var config map[string]interface{}
		if err := json.Unmarshal(settingJSON, &config); err == nil {
			if cid, ok := config["google_client_id"].(string); ok && cid != "" {
				expectedClientID = cid
			}
		}
	}

	if expectedClientID == "" {
		expectedClientID = os.Getenv("GOOGLE_CLIENT_ID")
	}

	if expectedClientID == "" {
		zlog.Error().Msg("GOOGLE_CLIENT_ID is not configured in database or environment variables")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "OAuth server misconfiguration",
		})
	}

	// Google allows aud to be prefix/match. Let's do a strict verification
	if tokenInfo.Aud != expectedClientID {
		zlog.Warn().Str("token_aud", tokenInfo.Aud).Str("expected_aud", expectedClientID).Msg("Audience mismatch in Google Login")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token audience mismatch",
		})
	}

	// Verify issuer is Google
	if tokenInfo.Iss != "accounts.google.com" && tokenInfo.Iss != "https://accounts.google.com" {
		zlog.Warn().Str("token_iss", tokenInfo.Iss).Msg("Issuer mismatch in Google Login")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token issuer mismatch",
		})
	}

	// Verify email is verified
	verified := false
	if evStr, ok := tokenInfo.EmailVerified.(string); ok {
		verified = evStr == "true" || evStr == "1"
	} else if evBool, ok := tokenInfo.EmailVerified.(bool); ok {
		verified = evBool
	}

	if !verified {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Google account email is not verified",
		})
	}

	// Check if user exists in database
	var user User
	var status string
	err = db.QueryRowContext(context.Background(),
		"SELECT id, email, full_name, role, COALESCE(status, 'active') FROM users WHERE email = $1",
		tokenInfo.Email,
	).Scan(&user.ID, &user.Email, &user.FullName, &user.Role, &status)

	if err == nil {
		// User exists! Log them in
		if status == "pending" {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "Akun Anda sedang menunggu persetujuan admin.",
			})
		}
		if status == "rejected" {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "Akun Anda telah ditolak.",
			})
		}

		token, err := generateToken(user.ID, user.Role)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

		return c.JSON(fiber.Map{
			"token": token,
			"user":  user,
		})
	}

	// User does not exist, let's register a new user!
	roleToSet := "user"
	if req.Role == "candidate" {
		roleToSet = "candidate"
	}

	var newID string
	placeholderPassword, _ := hashPassword("google_oauth_placeholder_not_usable_password")

	ctx := context.Background()
	if roleToSet == "candidate" {
		// Create candidate requires a transaction to link to pmb_candidates
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to begin transaction"})
		}
		defer tx.Rollback()

		// 1. Insert user
		err = tx.QueryRowContext(ctx,
			`INSERT INTO users (email, password_hash, full_name, role, status)
			 VALUES ($1, $2, $3, $4, 'active') RETURNING id`,
			tokenInfo.Email, placeholderPassword, tokenInfo.Name, roleToSet,
		).Scan(&newID)
		if err != nil {
			zlog.Error().Err(err).Msg("Failed to create candidate user via Google")
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create account"})
		}

		// 2. Insert pmb_candidate record (empty phone_whatsapp)
		_, err = tx.ExecContext(ctx,
			`INSERT INTO pmb_candidates (user_id, full_name, phone_whatsapp)
			 VALUES ($1, $2, '')`,
			newID, tokenInfo.Name,
		)
		if err != nil {
			zlog.Error().Err(err).Msg("Failed to setup pmb_candidate profile via Google")
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to setup candidate profile"})
		}

		if err := tx.Commit(); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Transaction commit failed"})
		}
	} else {
		// General user
		err = db.QueryRowContext(ctx,
			`INSERT INTO users (email, password_hash, full_name, role, status)
			 VALUES ($1, $2, $3, $4, 'active') RETURNING id`,
			tokenInfo.Email, placeholderPassword, tokenInfo.Name, roleToSet,
		).Scan(&newID)
		if err != nil {
			zlog.Error().Err(err).Msg("Failed to create user via Google")
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create account"})
		}
	}

	user.ID = newID
	user.Email = tokenInfo.Email
	user.FullName = tokenInfo.Name
	user.Role = roleToSet

	token, err := generateToken(user.ID, user.Role)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	logAuthEvent(user.ID, "google_register_success", c.IP(), c.Get("User-Agent"), "")

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}



