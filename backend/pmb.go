package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// saveCandidateForm updates the PMB Candidate details
func saveCandidateForm(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		NISN                   string `json:"nisn"`
		SchoolOrigin           string `json:"school_origin"`
		FirstChoiceProgramID   string `json:"first_choice_program_id"`
		SecondChoiceProgramID  string `json:"second_choice_program_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validate required fields
	if req.NISN == "" || req.SchoolOrigin == "" || req.FirstChoiceProgramID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "NISN, School Origin, and First Choice Program are required"})
	}

	var p1, p2 interface{}
	p1 = req.FirstChoiceProgramID
	p2 = nil
	if req.SecondChoiceProgramID != "" {
		p2 = req.SecondChoiceProgramID
	}

	// Update candidate profile
	query := `
		UPDATE pmb_candidates 
		SET nisn = $1, school_origin = $2, first_choice_program_id = $3, second_choice_program_id = $4, status = 'WAITING_PAYMENT', updated_at = NOW()
		WHERE user_id = $5
		RETURNING id, status
	`
	
	var candidateID, status string
	err := db.QueryRowContext(context.Background(), query, req.NISN, req.SchoolOrigin, p1, p2, userID).Scan(&candidateID, &status)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update candidate profile"})
	}

	return c.JSON(fiber.Map{
		"message": "Candidate form saved successfully",
		"status": status,
	})
}

// getCandidateProfile returns the candidate's current PMB details
func getCandidateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var profile struct {
		ID                     string  `json:"id"`
		RegistrationNumber     *string `json:"registration_number"`
		FullName               string  `json:"full_name"`
		NISN                   *string `json:"nisn"`
		PhoneWhatsapp          *string `json:"phone_whatsapp"`
		SchoolOrigin           *string `json:"school_origin"`
		FirstChoiceProgramID   *string `json:"first_choice_program_id"`
		SecondChoiceProgramID  *string `json:"second_choice_program_id"`
		Status                 string  `json:"status"`
	}

	query := `
		SELECT id, registration_number, full_name, nisn, phone_whatsapp, school_origin, 
		       first_choice_program_id, second_choice_program_id, status 
		FROM pmb_candidates 
		WHERE user_id = $1
	`
	
	err := db.QueryRowContext(context.Background(), query, userID).Scan(
		&profile.ID, &profile.RegistrationNumber, &profile.FullName, &profile.NISN, 
		&profile.PhoneWhatsapp, &profile.SchoolOrigin, &profile.FirstChoiceProgramID, 
		&profile.SecondChoiceProgramID, &profile.Status,
	)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Candidate profile not found"})
	}

	return c.JSON(profile)
}

// uploadPaymentProof handles uploading the payment receipt for the PMB registration fee
func uploadPaymentProof(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		ImageURL string `json:"image_url"`
	}

	if err := c.BodyParser(&req); err != nil || req.ImageURL == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Upload bukti pembayaran terlebih dahulu"})
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	defer tx.Rollback()

	// 1. Get candidate ID
	var candidateID string
	err = tx.QueryRowContext(context.Background(), "SELECT id FROM pmb_candidates WHERE user_id = $1", userID).Scan(&candidateID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Candidate not found"})
	}

	// 2. Check if already has pending payment (prevent duplicate)
	var existingPayment string
	tx.QueryRowContext(context.Background(),
		"SELECT id FROM pmb_payments WHERE candidate_id = $1 AND status = 'PENDING'", candidateID,
	).Scan(&existingPayment)
	if existingPayment != "" {
		// Update existing payment instead of creating new
		_, err = tx.ExecContext(context.Background(),
			`UPDATE pmb_payments SET proof_image_url = $1, updated_at = NOW() WHERE id = $2`,
			req.ImageURL, existingPayment,
		)
	} else {
		// Insert new payment
		_, err = tx.ExecContext(context.Background(),
			`INSERT INTO pmb_payments (candidate_id, payment_type, amount, proof_image_url, status) 
			 VALUES ($1, 'REGISTRATION', 250000, $2, 'PENDING')`,
			candidateID, req.ImageURL,
		)
	}
	if err != nil {
		zlog.Error().Err(err).Str("candidate_id", candidateID).Str("image_url", req.ImageURL).Msg("PMB payment insert failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save payment record"})
	}

	// 3. Update Candidate Status to WAITING_PAYMENT (admin will verify)
	_, err = tx.ExecContext(context.Background(),
		`UPDATE pmb_candidates SET status = 'WAITING_PAYMENT', updated_at = NOW() WHERE id = $1 AND status != 'PAYMENT_VERIFIED'`,
		candidateID,
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update candidate status"})
	}

	if err := tx.Commit(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.JSON(fiber.Map{
		"message": "Bukti pembayaran berhasil diupload. Menunggu verifikasi admin.",
		"status":  "PENDING",
	})
}

// candidateUploadFile handles file upload for PMB candidates (not admin-only)
func candidateUploadFile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "file is required"})
	}

	// Validate file size (max 5MB for documents)
	if file.Size > 5*1024*1024 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "File terlalu besar, maksimal 5MB"})
	}
	if file.Size == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "File kosong"})
	}

	// Validate extension (images + PDF only)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".pdf": true, ".webp": true,
	}
	if !allowedExts[ext] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Format file tidak didukung. Gunakan JPG, PNG, PDF, atau WEBP."})
	}

	// Validate MIME type
	contentType := file.Header.Get("Content-Type")
	validMIMEs := map[string][]string{
		".jpg":  {"image/jpeg", "image/jpg"},
		".jpeg": {"image/jpeg", "image/jpg"},
		".png":  {"image/png"},
		".pdf":  {"application/pdf"},
		".webp": {"image/webp"},
	}
	mimes, ok := validMIMEs[ext]
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Tipe file tidak didukung"})
	}
	mimeValid := false
	for _, m := range mimes {
		if contentType == m {
			mimeValid = true
			break
		}
	}
	if !mimeValid {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Content type %s tidak sesuai dengan ekstensi %s", contentType, ext)})
	}

	// Generate random filename
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dir := filepath.Join("..", "public", "uploads", "pmb")
	if err := os.MkdirAll(dir, 0755); err != nil {
		zlog.Error().Err(err).Str("dir", dir).Msg("PMB upload mkdir failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat direktori upload"})
	}

	target := filepath.Join(dir, name)

	// Anti path traversal
	absTarget, _ := filepath.Abs(target)
	absDir, _ := filepath.Abs(dir)
	if !strings.HasPrefix(absTarget, absDir) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid file path"})
	}

	if err := c.SaveFile(file, target); err != nil {
		zlog.Error().Err(err).Str("target", target).Msg("PMB upload save failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan file"})
	}

	// Auto-compress if image
	compressImage(target)

	// Check if PNG was converted to JPG
	finalName := name
	if ext == ".png" {
		jpgName := strings.TrimSuffix(name, ext) + ".jpg"
		jpgPath := filepath.Join(dir, jpgName)
		if _, err := os.Stat(jpgPath); err == nil {
			finalName = jpgName
		}
	}

	url := "/uploads/pmb/" + finalName
	zlog.Info().Str("url", url).Int64("size", file.Size).Str("user_id", fmt.Sprintf("%v", userID)).Msg("PMB candidate file uploaded")

	return c.JSON(fiber.Map{"url": url})
}

// submitAllDocuments marks all documents as submitted and updates candidate status
func submitAllDocuments(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var candidateID string
	err := db.QueryRowContext(context.Background(), "SELECT id FROM pmb_candidates WHERE user_id = $1", userID).Scan(&candidateID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Candidate not found"})
	}

	// Check candidate has at least the required documents
	var docCount int
	db.QueryRowContext(context.Background(),
		"SELECT COUNT(*) FROM pmb_documents WHERE candidate_id = $1", candidateID).Scan(&docCount)
	if docCount < 3 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Minimal 3 dokumen harus diunggah (KK, Ijazah, Pas Foto)"})
	}

	// Update candidate status to DOCUMENT_SUBMITTED
	_, err = db.ExecContext(context.Background(),
		`UPDATE pmb_candidates SET status = 'DOCUMENT_SUBMITTED', updated_at = NOW() WHERE id = $1 AND status = 'PAYMENT_VERIFIED'`,
		candidateID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui status"})
	}

	return c.JSON(fiber.Map{"message": "Berkas berhasil dikirim. Menunggu verifikasi panitia.", "status": "DOCUMENT_SUBMITTED"})
}
