package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ─── PMB Admin Dashboard Stats ───────────────────────────────────────────────

func getPMBDashboard(c *fiber.Ctx) error {
	var total, draft, waitingPayment, paymentVerified, documentSubmitted, documentVerified, passed, failed, enrolled int64

	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM pmb_candidates").Scan(&total)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM pmb_candidates WHERE status='DRAFT'").Scan(&draft)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM pmb_candidates WHERE status='WAITING_PAYMENT'").Scan(&waitingPayment)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM pmb_candidates WHERE status='PAYMENT_VERIFIED'").Scan(&paymentVerified)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM pmb_candidates WHERE status='DOCUMENT_SUBMITTED'").Scan(&documentSubmitted)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM pmb_candidates WHERE status='DOCUMENT_VERIFIED'").Scan(&documentVerified)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM pmb_candidates WHERE status='PASSED'").Scan(&passed)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM pmb_candidates WHERE status='FAILED'").Scan(&failed)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM pmb_candidates WHERE status='ENROLLED'").Scan(&enrolled)

	var pendingPayments int64
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM pmb_payments WHERE status='PENDING'").Scan(&pendingPayments)

	var activeBatch string
	db.QueryRowContext(context.Background(), "SELECT COALESCE(name,'') FROM pmb_batches WHERE is_active=true LIMIT 1").Scan(&activeBatch)

	return c.JSON(fiber.Map{
		"total":              total,
		"draft":              draft,
		"waiting_payment":    waitingPayment,
		"payment_verified":   paymentVerified,
		"document_submitted": documentSubmitted,
		"document_verified":  documentVerified,
		"passed":             passed,
		"failed":             failed,
		"enrolled":           enrolled,
		"pending_payments":   pendingPayments,
		"active_batch":       activeBatch,
	})
}

// ─── Payment Verification ────────────────────────────────────────────────────

func getPendingPayments(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(), `
		SELECT p.id, p.candidate_id, p.payment_type, p.amount, p.proof_image_url, p.status, p.created_at,
		       c.full_name, c.registration_number, c.phone_whatsapp
		FROM pmb_payments p
		JOIN pmb_candidates c ON c.id = p.candidate_id
		ORDER BY p.created_at DESC
		LIMIT 50
	`)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch payments"})
	}
	defer rows.Close()

	payments := []fiber.Map{}
	for rows.Next() {
		var id, candidateID, paymentType, proofURL, status, createdAt, fullName string
		var amount float64
		var regNumber, phone *string
		if rows.Scan(&id, &candidateID, &paymentType, &amount, &proofURL, &status, &createdAt, &fullName, &regNumber, &phone) == nil {
			payments = append(payments, fiber.Map{
				"id": id, "candidate_id": candidateID, "payment_type": paymentType,
				"amount": amount, "proof_image_url": proofURL, "status": status,
				"created_at": createdAt, "full_name": fullName,
				"registration_number": regNumber, "phone_whatsapp": phone,
			})
		}
	}
	return c.JSON(payments)
}

func verifyPayment(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Action string `json:"action"` // approve, reject
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Action != "approve" && req.Action != "reject" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Action must be approve or reject"})
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	defer tx.Rollback()

	// Update payment status
	newStatus := "VERIFIED"
	if req.Action == "reject" {
		newStatus = "REJECTED"
	}
	_, err = tx.ExecContext(context.Background(),
		"UPDATE pmb_payments SET status = $1, verified_at = NOW() WHERE id = $2", newStatus, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update payment"})
	}

	// If approved, update candidate status
	if req.Action == "approve" {
		var candidateID string
		tx.QueryRowContext(context.Background(), "SELECT candidate_id FROM pmb_payments WHERE id = $1", id).Scan(&candidateID)
		if candidateID != "" {
			tx.ExecContext(context.Background(),
				"UPDATE pmb_candidates SET status = 'PAYMENT_VERIFIED', updated_at = NOW() WHERE id = $1", candidateID)

			// Send WhatsApp notification
			var phone string
			tx.QueryRowContext(context.Background(), "SELECT COALESCE(phone_whatsapp,'') FROM pmb_candidates WHERE id = $1", candidateID).Scan(&phone)
			if phone != "" {
				go SendWhatsAppMessage(phone, "Selamat! Pembayaran pendaftaran Anda telah diverifikasi. Silakan lanjutkan proses pendaftaran di portal PMB UPERTIS.")
			}
		}
	} else {
		// Rejected — notify candidate
		var candidateID string
		tx.QueryRowContext(context.Background(), "SELECT candidate_id FROM pmb_payments WHERE id = $1", id).Scan(&candidateID)
		if candidateID != "" {
			tx.ExecContext(context.Background(),
				"UPDATE pmb_candidates SET status = 'WAITING_PAYMENT', updated_at = NOW() WHERE id = $1", candidateID)

			var phone string
			tx.QueryRowContext(context.Background(), "SELECT COALESCE(phone_whatsapp,'') FROM pmb_candidates WHERE id = $1", candidateID).Scan(&phone)
			if phone != "" {
				go SendWhatsAppMessage(phone, "Mohon maaf, bukti pembayaran Anda ditolak. Silakan upload ulang bukti pembayaran yang valid di portal PMB UPERTIS.")
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	adminID, _ := c.Locals("user_id").(string)
	logAuthEvent(adminID, "payment_"+req.Action, c.IP(), "", "payment_id="+id)

	return c.JSON(fiber.Map{"message": fmt.Sprintf("Payment %sd", req.Action)})
}

// ─── Candidate Status Update (admin) ─────────────────────────────────────────

func updateCandidateStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Status  string `json:"status"`
		Message string `json:"message"` // optional WA message
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	validStatuses := map[string]bool{
		"DRAFT": true, "WAITING_PAYMENT": true, "PAYMENT_VERIFIED": true,
		"DOCUMENT_SUBMITTED": true, "DOCUMENT_VERIFIED": true, "PASSED": true, "FAILED": true, "ENROLLED": true,
	}
	if !validStatuses[req.Status] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid status"})
	}

	_, err := db.ExecContext(context.Background(),
		"UPDATE pmb_candidates SET status = $1, updated_at = NOW() WHERE id = $2", req.Status, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update status"})
	}

	// Send WhatsApp notification if message provided
	if req.Message != "" {
		var phone string
		db.QueryRowContext(context.Background(), "SELECT COALESCE(phone_whatsapp,'') FROM pmb_candidates WHERE id = $1", id).Scan(&phone)
		if phone != "" {
			go SendWhatsAppMessage(phone, req.Message)
		}
	}

	return c.JSON(fiber.Map{"message": "Status updated to " + req.Status})
}

// ─── Document Upload (candidate) ─────────────────────────────────────────────

func uploadDocument(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		DocumentType string `json:"document_type"` // KTP, KK, IJAZAH, PAS_FOTO, RAPOR
		FileURL      string `json:"file_url"`
	}
	if err := c.BodyParser(&req); err != nil || req.DocumentType == "" || req.FileURL == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "document_type and file_url are required"})
	}

	validTypes := map[string]bool{"KTP": true, "KK": true, "IJAZAH": true, "PAS_FOTO": true, "RAPOR": true, "SERTIFIKAT": true}
	if !validTypes[req.DocumentType] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid document type"})
	}

	var candidateID string
	err := db.QueryRowContext(context.Background(), "SELECT id FROM pmb_candidates WHERE user_id = $1", userID).Scan(&candidateID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Candidate not found"})
	}

	// Upsert document (replace if same type exists)
	_, err = db.ExecContext(context.Background(), `
		INSERT INTO pmb_documents (candidate_id, document_type, file_url, status)
		VALUES ($1, $2, $3, 'PENDING')
		ON CONFLICT (candidate_id, document_type) DO UPDATE SET file_url = $3, status = 'PENDING', updated_at = NOW()
	`, candidateID, req.DocumentType, req.FileURL)

	if err != nil {
		// If unique constraint doesn't exist, just insert
		_, err = db.ExecContext(context.Background(), `
			INSERT INTO pmb_documents (candidate_id, document_type, file_url, status)
			VALUES ($1, $2, $3, 'PENDING')
		`, candidateID, req.DocumentType, req.FileURL)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save document"})
		}
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Document uploaded"})
}

func getCandidateDocuments(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var candidateID string
	db.QueryRowContext(context.Background(), "SELECT id FROM pmb_candidates WHERE user_id = $1", userID).Scan(&candidateID)
	if candidateID == "" {
		return c.JSON([]fiber.Map{})
	}

	rows, err := db.QueryContext(context.Background(),
		"SELECT id, document_type, file_url, status, created_at FROM pmb_documents WHERE candidate_id = $1 ORDER BY created_at", candidateID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch documents"})
	}
	defer rows.Close()

	docs := []fiber.Map{}
	for rows.Next() {
		var id, docType, fileURL, status, createdAt string
		if rows.Scan(&id, &docType, &fileURL, &status, &createdAt) == nil {
			docs = append(docs, fiber.Map{"id": id, "document_type": docType, "file_url": fileURL, "status": status, "created_at": createdAt})
		}
	}
	return c.JSON(docs)
}

// ─── Batch Auto-Close (scheduled) ───────────────────────────────────────────

func startBatchAutoClose() {
	go func() {
		ticker := time.NewTicker(60 * time.Minute) // Check every hour
		defer ticker.Stop()
		for range ticker.C {
			res, err := db.ExecContext(context.Background(),
				`UPDATE pmb_batches SET is_active = false WHERE is_active = true AND end_date < CURRENT_DATE`)
			if err == nil {
				count, _ := res.RowsAffected()
				if count > 0 {
					zlog.Info().Int64("count", count).Msg("PMB batches auto-closed (past end_date)")
				}
			}
		}
	}()
}

// ─── PMB Export ──────────────────────────────────────────────────────────────

func exportCandidates(c *fiber.Ctx) error {
	status := c.Query("status", "")

	query := `SELECT c.registration_number, c.full_name, c.nisn, c.phone_whatsapp, c.school_origin, c.status, c.created_at,
		 COALESCE(p1.name,'') as first_choice, COALESCE(p2.name,'') as second_choice
		 FROM pmb_candidates c
		 LEFT JOIN faculty_programs p1 ON p1.id = c.first_choice_program_id
		 LEFT JOIN faculty_programs p2 ON p2.id = c.second_choice_program_id`

	var args []interface{}
	if status != "" {
		query += " WHERE c.status = $1"
		args = append(args, status)
	}
	query += " ORDER BY c.created_at DESC"

	rows, err := db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to export"})
	}
	defer rows.Close()

	candidates := []fiber.Map{}
	for rows.Next() {
		var regNum, fullName, nisn, phone, school, st, createdAt, first, second string
		if rows.Scan(&regNum, &fullName, &nisn, &phone, &school, &st, &createdAt, &first, &second) == nil {
			candidates = append(candidates, fiber.Map{
				"registration_number": regNum, "full_name": fullName, "nisn": nisn,
				"phone_whatsapp": phone, "school_origin": school, "status": st,
				"created_at": createdAt, "first_choice": first, "second_choice": second,
			})
		}
	}
	return c.JSON(candidates)
}
