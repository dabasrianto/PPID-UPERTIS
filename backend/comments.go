package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ─── Comments Migration ──────────────────────────────────────────────────────

func migrateCommentsTable() {
	db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS comments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			entity_type VARCHAR(20) NOT NULL,
			entity_id UUID NOT NULL,
			parent_id UUID,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			status VARCHAR(20) DEFAULT 'pending',
			ip_address VARCHAR(50),
			user_agent VARCHAR(500),
			approved_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_comments_entity ON comments(entity_type, entity_id, status)`)
	db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_comments_user ON comments(user_id)`)
	db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_comments_status ON comments(status, created_at DESC)`)
}

// ─── Banned keywords (configurable later via settings) ───────────────────────

var bannedKeywords = []string{"judi", "togel", "slot", "porn", "viagra", "casino"}

func containsBannedKeyword(content string) bool {
	lower := strings.ToLower(content)
	for _, kw := range bannedKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// ─── Public API ──────────────────────────────────────────────────────────────

// getComments returns approved comments for a given entity
func getComments(c *fiber.Ctx) error {
	entityType := c.Query("entity_type", "")
	entityID := c.Query("entity_id", "")
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	if entityType == "" || entityID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "entity_type and entity_id are required"})
	}
	if entityType != "post" && entityType != "event" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "entity_type must be 'post' or 'event'"})
	}
	if limit > 50 {
		limit = 50
	}

	rows, err := db.QueryContext(context.Background(),
		`SELECT c.id, c.content, c.parent_id, c.created_at, u.id, u.full_name, u.role
		 FROM comments c
		 JOIN users u ON u.id = c.user_id
		 WHERE c.entity_type = $1 AND c.entity_id = $2 AND c.status = 'approved'
		 ORDER BY c.created_at ASC
		 LIMIT $3 OFFSET $4`,
		entityType, entityID, limit, offset,
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch comments"})
	}
	defer rows.Close()

	comments := []fiber.Map{}
	for rows.Next() {
		var id, content, createdAt, userID, fullName, role string
		var parentID *string
		if rows.Scan(&id, &content, &parentID, &createdAt, &userID, &fullName, &role) == nil {
			comment := fiber.Map{
				"id":         id,
				"content":    content,
				"parent_id":  parentID,
				"created_at": createdAt,
				"user": fiber.Map{
					"id":        userID,
					"full_name": fullName,
					"role":      role,
				},
			}
			comments = append(comments, comment)
		}
	}

	c.Set("Cache-Control", "public, max-age=60")
	return c.JSON(comments)
}

// submitComment creates a new comment (requires auth + verified user)
func submitComment(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Login diperlukan untuk berkomentar"})
	}

	var req struct {
		EntityType string `json:"entity_type"`
		EntityID   string `json:"entity_id"`
		Content    string `json:"content"`
		ParentID   string `json:"parent_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validate
	if req.EntityType != "post" && req.EntityType != "event" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "entity_type must be 'post' or 'event'"})
	}
	if req.EntityID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "entity_id is required"})
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Komentar tidak boleh kosong"})
	}
	if len(content) > 2000 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Komentar maksimal 2000 karakter"})
	}

	// Check user is active and has WhatsApp verified (skip for admin)
	var userStatus, phoneWA, userRole string
	err := db.QueryRowContext(context.Background(),
		"SELECT COALESCE(status,'active'), COALESCE(phone_whatsapp,''), role FROM users WHERE id = $1",
		userID,
	).Scan(&userStatus, &phoneWA, &userRole)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}
	if userStatus != "active" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Akun belum aktif"})
	}
	// Admin/faculty_admin can comment without WhatsApp verification
	if phoneWA == "" && userRole != "admin" && userRole != "faculty_admin" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Verifikasi WhatsApp diperlukan untuk berkomentar"})
	}

	// Spam detection: banned keywords
	if containsBannedKeyword(content) {
		// Silently save as spam
		db.ExecContext(context.Background(),
			`INSERT INTO comments (entity_type, entity_id, user_id, content, status, ip_address, user_agent)
			 VALUES ($1, $2, $3, $4, 'spam', $5, $6)`,
			req.EntityType, req.EntityID, userID, content, c.IP(), truncate(c.Get("User-Agent"), 200),
		)
		logAuthEvent(fmt.Sprintf("%v", userID), "comment_spam", c.IP(), c.Get("User-Agent"), "banned_keyword")
		// Return success to not reveal detection
		return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Komentar terkirim, menunggu persetujuan admin."})
	}

	// Duplicate detection (same user, same entity, same content within 5 min)
	var dupCount int
	db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM comments WHERE user_id = $1 AND entity_id = $2 AND content = $3 AND created_at > NOW() - INTERVAL '5 minutes'`,
		userID, req.EntityID, content,
	).Scan(&dupCount)
	if dupCount > 0 {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Komentar duplikat terdeteksi"})
	}

	// Insert as pending
	var parentID interface{}
	if req.ParentID != "" {
		parentID = req.ParentID
	}

	var commentID string
	err = db.QueryRowContext(context.Background(),
		`INSERT INTO comments (entity_type, entity_id, parent_id, user_id, content, status, ip_address, user_agent)
		 VALUES ($1, $2, $3, $4, $5, 'pending', $6, $7) RETURNING id`,
		req.EntityType, req.EntityID, parentID, userID, content, c.IP(), truncate(c.Get("User-Agent"), 200),
	).Scan(&commentID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan komentar"})
	}

	zlog.Info().
		Str("comment_id", commentID).
		Str("user_id", fmt.Sprintf("%v", userID)).
		Str("entity", req.EntityType+"/"+req.EntityID).
		Msg("Comment submitted (pending)")

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Komentar terkirim, menunggu persetujuan admin.",
		"id":      commentID,
	})
}

// ─── Admin Moderation ────────────────────────────────────────────────────────

func getAdminComments(c *fiber.Ctx) error {
	status := c.Query("status", "pending")
	entityType := c.Query("entity_type", "")
	limit := c.QueryInt("limit", 50)

	query := `SELECT c.id, c.entity_type, c.entity_id, c.content, c.status, c.ip_address, c.created_at,
		 u.id, u.full_name, u.email
		 FROM comments c
		 JOIN users u ON u.id = c.user_id
		 WHERE c.status = $1`
	args := []interface{}{status}
	argIdx := 2

	if entityType != "" {
		query += fmt.Sprintf(" AND c.entity_type = $%d", argIdx)
		args = append(args, entityType)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY c.created_at DESC LIMIT $%d", argIdx)
	args = append(args, limit)

	rows, err := db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch comments"})
	}
	defer rows.Close()

	comments := []fiber.Map{}
	for rows.Next() {
		var id, eType, eID, content, st, ip, createdAt, uID, uName, uEmail string
		if rows.Scan(&id, &eType, &eID, &content, &st, &ip, &createdAt, &uID, &uName, &uEmail) == nil {
			comments = append(comments, fiber.Map{
				"id": id, "entity_type": eType, "entity_id": eID,
				"content": content, "status": st, "ip_address": ip, "created_at": createdAt,
				"user": fiber.Map{"id": uID, "full_name": uName, "email": uEmail},
			})
		}
	}
	return c.JSON(comments)
}

func moderateComment(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Action string `json:"action"` // approve, reject, spam
		Reply  string `json:"reply"`  // optional admin reply
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var newStatus string
	switch req.Action {
	case "approve":
		newStatus = "approved"
	case "reject":
		newStatus = "rejected"
	case "spam":
		newStatus = "spam"
	case "reply":
		// Admin reply — insert new comment as approved
		if req.Reply == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Reply content is required"})
		}
		adminID, _ := c.Locals("user_id").(string)
		// Get parent comment's entity info
		var entityType, entityID string
		err := db.QueryRowContext(context.Background(),
			"SELECT entity_type, entity_id FROM comments WHERE id = $1", id,
		).Scan(&entityType, &entityID)
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Comment not found"})
		}
		_, err = db.ExecContext(context.Background(),
			`INSERT INTO comments (entity_type, entity_id, parent_id, user_id, content, status, ip_address, approved_at)
			 VALUES ($1, $2, $3, $4, $5, 'approved', $6, NOW())`,
			entityType, entityID, id, adminID, req.Reply, c.IP(),
		)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save reply"})
		}
		logAuthEvent(adminID, "comment_reply", c.IP(), c.Get("User-Agent"), "parent="+id)
		return c.JSON(fiber.Map{"message": "Balasan berhasil dikirim"})
	default:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Action must be approve, reject, spam, or reply"})
	}

	approvedAt := ""
	if newStatus == "approved" {
		approvedAt = ", approved_at = NOW()"
	}

	_, err := db.ExecContext(context.Background(),
		fmt.Sprintf("UPDATE comments SET status = $1%s WHERE id = $2", approvedAt),
		newStatus, id,
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update comment"})
	}

	adminID, _ := c.Locals("user_id").(string)
	logAuthEvent(adminID, "comment_"+req.Action, c.IP(), c.Get("User-Agent"), "comment_id="+id)

	return c.JSON(fiber.Map{"message": fmt.Sprintf("Comment %sd", req.Action)})
}

// ─── Comment count helper (for post/event detail) ────────────────────────────

func getCommentCount(c *fiber.Ctx) error {
	entityType := c.Query("entity_type", "")
	entityID := c.Query("entity_id", "")
	if entityType == "" || entityID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "entity_type and entity_id required"})
	}

	var count int
	db.QueryRowContext(context.Background(),
		"SELECT COUNT(*) FROM comments WHERE entity_type = $1 AND entity_id = $2 AND status = 'approved'",
		entityType, entityID,
	).Scan(&count)

	return c.JSON(fiber.Map{"count": count})
}
