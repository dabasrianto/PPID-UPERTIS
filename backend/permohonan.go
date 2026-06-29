package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// PermohonanRequest represents the public form payload for requesting information
type PermohonanRequest struct {
	ApplicantType  string `json:"applicant_type"`
	Name           string `json:"name"`
	IdentityNumber string `json:"identity_number"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Address        string `json:"address"`
	Details        string `json:"details"`
	Purpose        string `json:"purpose"`
	ObtainMethod   string `json:"obtain_method"`
	DeliveryMethod string `json:"delivery_method"`
	AttachmentURL  string `json:"attachment_url"`
}

// Generate a random ticket suffix
func generateTicketNumber() string {
	now := time.Now()
	dateStr := now.Format("20060102")
	// Generate random number 1000 - 9999
	nBig, err := rand.Int(rand.Reader, big.NewInt(9000))
	num := 1000
	if err == nil {
		num = int(nBig.Int64()) + 1000
	}
	return fmt.Sprintf("PPID-REQ-%s-%d", dateStr, num)
}

func submitPermohonan(c *fiber.Ctx) error {
	var req PermohonanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	// Simple validation
	if req.ApplicantType == "" || req.Name == "" || req.IdentityNumber == "" || req.Email == "" || req.Phone == "" || req.Address == "" || req.Details == "" || req.Purpose == "" || req.ObtainMethod == "" || req.DeliveryMethod == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Semua field bertanda bintang (*) wajib diisi"})
	}

	ticket := generateTicketNumber()

	var id string
	err := db.QueryRowContext(context.Background(),
		`INSERT INTO permohonan_informasi 
		 (ticket_number, applicant_type, name, identity_number, email, phone, address, details, purpose, obtain_method, delivery_method, attachment_url) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id::text`,
		ticket, req.ApplicantType, req.Name, req.IdentityNumber, req.Email, req.Phone, req.Address, req.Details, req.Purpose, req.ObtainMethod, req.DeliveryMethod, req.AttachmentURL,
	).Scan(&id)

	if err != nil {
		zlog.Error().Err(err).Msg("Failed to submit permohonan informasi")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan permohonan informasi"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"id":            id,
		"ticket_number": ticket,
		"message":       "Permohonan informasi berhasil diajukan",
	})
}

func checkPermohonanStatus(c *fiber.Ctx) error {
	ticket := c.Params("ticket")
	if ticket == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Nomor tiket wajib diisi"})
	}

	var name, details, purpose, status, adminResponse, createdAt string
	err := db.QueryRowContext(context.Background(),
		`SELECT name, details, purpose, status, COALESCE(admin_response, ''), created_at::text 
		 FROM permohonan_informasi WHERE ticket_number = $1`,
		ticket,
	).Scan(&name, &details, &purpose, &status, &adminResponse, &createdAt)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Nomor tiket tidak ditemukan"})
	}

	return c.JSON(fiber.Map{
		"ticket_number":  ticket,
		"name":           name,
		"details":        details,
		"purpose":        purpose,
		"status":         status,
		"admin_response": adminResponse,
		"created_at":     createdAt,
	})
}

func uploadPermohonanAttachment(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "tidak ada file yang diunggah"})
	}

	// Limit to 5MB for public uploads
	if file.Size > 5*1024*1024 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ukuran file maksimal 5MB"})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "file harus memiliki ekstensi"})
	}

	// Whitelist KTP / image and pdf only for safety on public forms
	allowedPubExt := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".pdf": true,
	}

	if !allowedPubExt[ext] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "jenis file tidak diizinkan. Hanya diperbolehkan: jpg, jpeg, png, pdf",
		})
	}

	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dir := filepath.Join("..", "public", "uploads", "permohonan")
	if err := os.MkdirAll(dir, 0755); err != nil {
		zlog.Error().Err(err).Str("dir", dir).Msg("Permohonan upload mkdir failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "gagal membuat direktori upload"})
	}

	target := filepath.Join(dir, name)
	if err := c.SaveFile(file, target); err != nil {
		zlog.Error().Err(err).Str("target", target).Msg("Permohonan upload save failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "gagal menyimpan file"})
	}

	// Auto-compress image (if image)
	if ext != ".pdf" {
		compressImage(target)
		// Check if png was converted to jpg
		if ext == ".png" {
			jpgName := strings.TrimSuffix(name, ext) + ".jpg"
			jpgPath := filepath.Join(dir, jpgName)
			if _, err := os.Stat(jpgPath); err == nil {
				name = jpgName
			}
		}
	}

	url := "/uploads/permohonan/" + name
	return c.JSON(fiber.Map{"url": url})
}
