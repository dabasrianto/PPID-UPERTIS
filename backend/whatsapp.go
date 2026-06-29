package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	qrcode "github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
	_ "github.com/lib/pq"
)

var (
	waClient     *whatsmeow.Client
	waContainer  *sqlstore.Container
	waQRString   string
	waConnected  bool
	waMutex      sync.RWMutex
	waInitMutex  sync.Mutex
)

func initWhatsApp() {
	waInitMutex.Lock()
	defer waInitMutex.Unlock()

	if waClient != nil {
		return
	}

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// We use the same connection string as the main DB
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "kampuspro"),
		getEnv("DB_SSLMODE", "disable"),
	)
	
	var container *sqlstore.Container
	var err error
	if waContainer != nil {
		container = waContainer
	} else {
		container, err = sqlstore.New(context.Background(), "postgres", connStr, dbLog)
		if err != nil {
			log.Printf("Failed to connect to WhatsApp store: %v", err)
			return
		}
		waContainer = container
	}

	deviceRes, err := container.GetFirstDevice(context.Background())
	if err != nil {
		log.Printf("Failed to get device: %v", err)
		return
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	waClient = whatsmeow.NewClient(deviceRes, clientLog)
	waClient.AddEventHandler(handler)

	if waClient.Store.ID == nil {
		// No device linked — start QR login flow
		// Use GetQRChannel as reliable source of QR strings
		qrChan, err2 := waClient.GetQRChannel(context.Background())
		if err2 != nil {
			log.Printf("Failed to get QR channel: %v", err2)
		}
		if err = waClient.Connect(); err != nil {
			log.Printf("Failed to connect to WhatsApp: %v", err)
			waClient = nil
			return
		}
		if qrChan != nil {
			go func() {
				for evt := range qrChan {
					if evt.Event == "code" {
						waMutex.Lock()
						waQRString = evt.Code
						waMutex.Unlock()
						log.Printf("WhatsApp QR Code ready: len=%d", len(evt.Code))
					} else if evt.Event == "success" {
						waMutex.Lock()
						waConnected = true
						waQRString = ""
						waMutex.Unlock()
						log.Println("WhatsApp QR: login success!")
					} else {
						log.Printf("WhatsApp QR channel event: %s", evt.Event)
					}
				}
			}()
		}
		log.Println("WhatsApp: waiting for QR scan...")
	} else {
		// Already logged in
		if err = waClient.Connect(); err != nil {
			log.Printf("Failed to connect to WhatsApp: %v", err)
			waClient = nil
			return
		}
		
		// Wait a moment for connection states to sync
		time.Sleep(1 * time.Second)
		
		if !waClient.IsLoggedIn() {
			log.Println("WhatsApp: Session is invalid/unlinked, deleting device store")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_ = waClient.Store.Delete(ctx)
			cancel()
			waClient.Disconnect()
			waClient = nil
			// Start QR flow instead
			go initWhatsApp()
			return
		}

		waConnected = true
		log.Println("WhatsApp: restored existing session")
	}
}


func handler(evt interface{}) {
	switch e := evt.(type) {
	case *events.QR:
		waMutex.Lock()
		if len(e.Codes) > 0 {
			waQRString = e.Codes[0]
			log.Printf("WhatsApp QR Code updated via event")
		}
		waMutex.Unlock()
	case *events.PairSuccess:
		waMutex.Lock()
		waConnected = true
		waQRString = ""
		waMutex.Unlock()
		log.Printf("WhatsApp paired successfully! ID: %s", e.ID.String())
	case *events.PairError:
		log.Printf("WhatsApp pairing failed: %v", e.Error)
	case *events.Connected:
		waMutex.Lock()
		waConnected = true
		waQRString = ""
		waMutex.Unlock()
		log.Println("WhatsApp connected successfully")
	case *events.LoggedOut:
		waMutex.Lock()
		waConnected = false
		waQRString = ""
		waMutex.Unlock()
		log.Printf("WhatsApp logged out, reason: %v", e.Reason)
		if waClient != nil && waClient.Store != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_ = waClient.Store.Delete(ctx)
			cancel()
		}
		waClient = nil
		go initWhatsApp()
	case *events.ConnectFailure:
		log.Printf("WhatsApp connection failure: %v (code: %d)", e.Message, e.Reason)
	case *events.Disconnected:
		waMutex.Lock()
		waConnected = false
		waMutex.Unlock()
		log.Println("WhatsApp disconnected")
	default:
		// Ignore other events
	}
}

func getWhatsAppStatus(c *fiber.Ctx) error {
	waMutex.RLock()
	connected := waConnected
	qrStr := waQRString
	waMutex.RUnlock()

	// Ensure the connection flag reflects the real client state
	if waClient != nil && waClient.IsLoggedIn() {
		waMutex.Lock()
		waConnected = true
		connected = true
		waMutex.Unlock()
	}

	// Generate QR image on the server side to avoid frontend encoding issues
	var qrImage string
	if !connected && qrStr != "" {
		png, err := qrcode.Encode(qrStr, qrcode.Low, 400)
		if err != nil {
			log.Printf("Failed to generate QR image: %v", err)
		} else {
			qrImage = "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
		}
	}

	// Get connected phone number
	var phoneNumber string
	if waClient != nil && waClient.Store != nil && waClient.Store.ID != nil {
		phoneNumber = waClient.Store.ID.User
	}

	return c.JSON(fiber.Map{
		"connected":    connected,
		"phone_number": phoneNumber,
		"qr":           qrStr,
		"qr_image":     qrImage,
	})
}

// getWhatsAppStatusPublic is a lightweight public endpoint for frontend polling.
// Returns minimal JSON, always HTTP 200, no sensitive data.
func getWhatsAppStatusPublic(c *fiber.Ctx) error {
	waMutex.RLock()
	connected := waConnected
	qrStr := waQRString
	waMutex.RUnlock()

	if waClient != nil && waClient.IsLoggedIn() {
		connected = true
	}

	return c.JSON(fiber.Map{
		"connected":     connected,
		"authenticated": connected,
		"waiting_qr":    !connected && qrStr != "",
	})
}

// getWhatsAppQR serves the QR code as PNG image or JSON status.
// No cache. Polling-safe.
func getWhatsAppQR(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-store, no-cache, must-revalidate")
	c.Set("Pragma", "no-cache")

	waMutex.RLock()
	connected := waConnected
	qrStr := waQRString
	waMutex.RUnlock()

	if waClient != nil && waClient.IsLoggedIn() {
		connected = true
	}

	// Already authenticated
	if connected {
		return c.JSON(fiber.Map{"authenticated": true})
	}

	// QR not ready yet
	if qrStr == "" {
		return c.JSON(fiber.Map{"waiting": true})
	}

	// Generate and serve QR as PNG
	png, err := qrcode.Encode(qrStr, qrcode.Medium, 400)
	if err != nil {
		return c.JSON(fiber.Map{"waiting": true})
	}

	c.Set("Content-Type", "image/png")
	return c.Send(png)
}

// Reset the WhatsApp connection and generate a new QR code.
func resetWhatsApp(c *fiber.Ctx) error {
	waMutex.Lock()
	waConnected = false
	waQRString = ""
	waMutex.Unlock()

	// 1. Disconnect and logout the active client if exists
	if waClient != nil {
		log.Println("WhatsApp: logging out active client...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_ = waClient.Logout(ctx) // Try clean logout
		cancel()
		waClient.Disconnect()
		waClient = nil
	}

	// 2. Completely delete all device records in the database store
	if waContainer != nil {
		log.Println("WhatsApp: clearing all device records in DB store...")
		ctxList, cancelList := context.WithTimeout(context.Background(), 5*time.Second)
		devices, err := waContainer.GetAllDevices(ctxList)
		cancelList()
		if err == nil {
			for _, dev := range devices {
				ctxDel, cancelDel := context.WithTimeout(context.Background(), 5*time.Second)
				errDel := dev.Delete(ctxDel)
				cancelDel()
				if errDel != nil {
					log.Printf("Failed to delete device record from DB: %v", errDel)
				}
			}
		}
	}

	// 3. Re-initialize a completely fresh WhatsApp connection
	go initWhatsApp()

	// Wait 3 seconds for QR generation to begin
	time.Sleep(3 * time.Second)

	waMutex.RLock()
	qrStr := waQRString
	connected := waConnected
	waMutex.RUnlock()

	return c.JSON(fiber.Map{
		"connected": connected,
		"qr":        qrStr,
		"message":   "Koneksi WhatsApp berhasil di-reset. Silakan scan QR code baru.",
	})
}

func sendOTP(phone string, code string) error {
	if waClient == nil || !waConnected {
		return fmt.Errorf("whatsapp not connected")
	}

	recipient, ok := parseJID(phone)
	if !ok {
		return fmt.Errorf("invalid phone number")
	}

	message := fmt.Sprintf("Kode OTP Anda adalah: *%s*\n\nJangan berikan kode ini kepada siapapun. Kode berlaku selama 5 menit.", code)
	_, err := waClient.SendMessage(context.Background(), recipient, &waE2E.Message{
		Conversation: proto.String(message),
	})
	return err
}

func SendWhatsAppMessage(phone string, message string) error {
	if waClient == nil || !waConnected {
		return fmt.Errorf("whatsapp not connected")
	}

	recipient, ok := parseJID(phone)
	if !ok {
		return fmt.Errorf("invalid phone number")
	}

	_, err := waClient.SendMessage(context.Background(), recipient, &waE2E.Message{
		Conversation: proto.String(message),
	})
	return err
}

func parseJID(arg string) (types.JID, bool) {
	if arg == "" {
		return types.NewJID("", types.DefaultUserServer), false
	}
	// Strip all non-digit characters
	cleaned := ""
	for _, ch := range arg {
		if ch >= '0' && ch <= '9' {
			cleaned += string(ch)
		}
	}
	if cleaned == "" {
		return types.NewJID("", types.DefaultUserServer), false
	}
	// Convert local format (08xx) to international (628xx)
	if len(cleaned) > 1 && cleaned[0] == '0' {
		cleaned = "62" + cleaned[1:]
	}
	return types.NewJID(cleaned, types.DefaultUserServer), true
}

// normalizePhoneNumber cleans and converts phone to international format (62xxx)
func normalizePhoneNumber(phone string) string {
	cleaned := ""
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			cleaned += string(ch)
		}
	}
	if cleaned == "" {
		return phone
	}
	if len(cleaned) > 1 && cleaned[0] == '0' {
		cleaned = "62" + cleaned[1:]
	}
	return cleaned
}

func generateOTP(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		result[i] = digits[num.Int64()]
	}
	return string(result)
}

// API Handlers for OTP

func requestOTP(c *fiber.Ctx) error {
	var req struct {
		Phone string `json:"phone"`
	}
	if err := c.BodyParser(&req); err != nil || req.Phone == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Phone number is required"})
	}

	code := generateOTP(6)
	expiresAt := time.Now().Add(5 * time.Minute)

	// Save to DB
	_, err := db.ExecContext(context.Background(),
		"INSERT INTO otps (phone, code, expires_at) VALUES ($1, $2, $3)",
		req.Phone, code, expiresAt)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate OTP"})
	}

	// Send via WhatsApp
	err = sendOTP(req.Phone, code)
	if err != nil {
		log.Printf("Failed to send WhatsApp OTP to %s: %v", req.Phone, err)
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Gagal mengirim OTP via WhatsApp. Pastikan koneksi WhatsApp aktif.",
		})
	}

	return c.JSON(fiber.Map{"message": "OTP sent successfully"})
}

func verifyOTP(c *fiber.Ctx) error {
	var req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	if err := c.BodyParser(&req); err != nil || req.Phone == "" || req.Code == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Phone and Code are required"})
	}

	var id string
	err := db.QueryRowContext(context.Background(),
		"SELECT id FROM otps WHERE phone = $1 AND code = $2 AND expires_at > NOW() AND is_used = FALSE ORDER BY created_at DESC LIMIT 1",
		req.Phone, req.Code).Scan(&id)

	if err == sql.ErrNoRows {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or expired OTP"})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Mark as used
	_, _ = db.ExecContext(context.Background(), "UPDATE otps SET is_used = TRUE WHERE id = $1", id)

	return c.JSON(fiber.Map{"message": "OTP verified successfully"})
}

func getOTPLogs(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(),
		"SELECT id, phone, code, expires_at, is_used, created_at FROM otps ORDER BY created_at DESC LIMIT 100")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch OTP logs"})
	}
	defer rows.Close()

	var logs []fiber.Map
	for rows.Next() {
		var id, phone, code, expiresAt, createdAt string
		var isUsed bool
		if err := rows.Scan(&id, &phone, &code, &expiresAt, &isUsed, &createdAt); err != nil {
			continue
		}
		logs = append(logs, fiber.Map{
			"id":         id,
			"phone":      phone,
			"code":       code,
			"expires_at": expiresAt,
			"is_used":    isUsed,
			"created_at": createdAt,
		})
	}

	return c.JSON(logs)
}

// notifyAdminNewMessage sends a WhatsApp notification to the admin when a message is received
func notifyAdminNewMessage(msg ContactMessage) {
	go func() {
		// 1. Get admin's WhatsApp number from settings
		var settingJSON []byte
		err := db.QueryRowContext(context.Background(), "SELECT value FROM site_settings WHERE key = 'contact'").Scan(&settingJSON)
		if err != nil {
			log.Printf("[WA Notify] Failed to fetch contact settings: %v", err)
			return
		}

		var contactData map[string]interface{}
		if err := json.Unmarshal(settingJSON, &contactData); err != nil {
			log.Printf("[WA Notify] Failed to parse contact settings JSON: %v", err)
			return
		}

		adminWA, _ := contactData["whatsapp"].(string)
		if adminWA == "" {
			log.Printf("[WA Notify] Admin WhatsApp number is not configured in settings")
			return
		}

		// 2. Format the message content
		subject := msg.Subject
		if subject == "" {
			subject = "-"
		}
		phone := msg.Phone
		if phone == "" {
			phone = "-"
		}

		message := fmt.Sprintf(
			"🔔 *Notifikasi Pesan Baru UPERTIS*\n\n"+
				"👤 *Nama:* %s\n"+
				"📧 *Email:* %s\n"+
				"📞 *No. HP/WA:* %s\n"+
				"📝 *Subjek:* %s\n\n"+
				"💬 *Pesan:*\n%s",
			msg.Name,
			msg.Email,
			phone,
			subject,
			msg.Message,
		)

		// 3. Send the WhatsApp message
		err = SendWhatsAppMessage(adminWA, message)
		if err != nil {
			log.Printf("[WA Notify] Failed to send notification to admin %s: %v", adminWA, err)
		} else {
			log.Printf("[WA Notify] Successfully sent message notification to admin %s", adminWA)
		}
	}()
}
