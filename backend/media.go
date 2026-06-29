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

// ─── Media Library Migration ─────────────────────────────────────────────────

func migrateMediaTable() {
	db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS media (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			filename VARCHAR(255) NOT NULL,
			original_name VARCHAR(255),
			url VARCHAR(500) NOT NULL UNIQUE,
			mime_type VARCHAR(100),
			size_bytes BIGINT DEFAULT 0,
			folder VARCHAR(100) DEFAULT 'media',
			uploaded_by VARCHAR(255),
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_media_folder ON media(folder)`)
	db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_media_created ON media(created_at DESC)`)
	// Add unique constraint on url if not exists
	db.ExecContext(context.Background(), `DO $$ BEGIN
		ALTER TABLE media ADD CONSTRAINT media_url_unique UNIQUE (url);
	EXCEPTION WHEN duplicate_table THEN NULL;
	END $$`)

	// Sync existing files ONLY ONCE (tracked via _migrations)
	if !migrationCompleted("sync_existing_uploads_v2") {
		syncExistingUploads()
		markMigrationCompleted("sync_existing_uploads_v2")
	}
}

// syncExistingUploads scans public/uploads/ and records any files not yet in media table
func syncExistingUploads() {
	baseDir := filepath.Join("..", "public", "uploads")
	var synced int

	filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Only image/pdf files
		ext := strings.ToLower(filepath.Ext(info.Name()))
		if !allowedExtensions[ext] {
			return nil
		}

		// Build relative URL
		relPath := strings.TrimPrefix(path, baseDir)
		relPath = strings.ReplaceAll(relPath, "\\", "/")
		url := "/uploads" + relPath

		// Determine folder from path
		parts := strings.Split(strings.TrimPrefix(relPath, "/"), "/")
		folder := "media"
		if len(parts) > 1 {
			folder = parts[0]
		}

		// Determine MIME type
		mimeType := "application/octet-stream"
		switch ext {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".png":
			mimeType = "image/png"
		case ".gif":
			mimeType = "image/gif"
		case ".webp":
			mimeType = "image/webp"
		case ".pdf":
			mimeType = "application/pdf"
		case ".ico":
			mimeType = "image/x-icon"
		case ".docx":
			mimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		case ".xlsx":
			mimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		case ".zip":
			mimeType = "application/zip"
		case ".rar":
			mimeType = "application/x-rar-compressed"
		case ".mp4":
			mimeType = "video/mp4"
		case ".webm":
			mimeType = "video/webm"
		case ".ogg":
			mimeType = "video/ogg"
		case ".mov":
			mimeType = "video/quicktime"
		case ".m4v":
			mimeType = "video/x-m4v"
		case ".avi":
			mimeType = "video/x-msvideo"
		case ".mkv":
			mimeType = "video/x-matroska"
		}

		// Insert if not exists (ON CONFLICT DO NOTHING)
		res, _ := db.ExecContext(context.Background(),
			`INSERT INTO media (filename, original_name, url, mime_type, size_bytes, folder)
			 VALUES ($1, $2, $3, $4, $5, $6)
			 ON CONFLICT (url) DO NOTHING`,
			info.Name(), info.Name(), url, mimeType, info.Size(), folder,
		)
		if res != nil {
			rows, _ := res.RowsAffected()
			if rows > 0 {
				synced++
			}
		}
		return nil
	})

	if synced > 0 {
		zlog.Info().Int("synced", synced).Msg("Media library: synced existing uploads")
	}
}

// ─── Multi-Upload Endpoint ───────────────────────────────────────────────────

func uploadMultipleFiles(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form data"})
	}

	files := form.File["files"]
	if len(files) == 0 {
		// Fallback: try single file field
		files = form.File["file"]
	}
	if len(files) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "No files uploaded"})
	}

	folder := c.FormValue("folder", "media")
	folder = cleanPathPart(folder)
	userID, _ := c.Locals("user_id").(string)

	dir := filepath.Join("..", "public", "uploads", folder)
	os.MkdirAll(dir, 0755)

	var uploaded []fiber.Map
	var errors []string

	for _, file := range files {
		// Validate size
		if file.Size > maxFileSize {
			errors = append(errors, fmt.Sprintf("%s: too large (max 30MB)", file.Filename))
			continue
		}
		if file.Size == 0 {
			errors = append(errors, fmt.Sprintf("%s: empty file", file.Filename))
			continue
		}

		// Validate extension
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext == "" || !allowedExtensions[ext] {
			errors = append(errors, fmt.Sprintf("%s: type not allowed", file.Filename))
			continue
		}

		// Validate MIME
		contentType := file.Header.Get("Content-Type")
		validMIMEs, exists := allowedMIMETypes[ext]
		if !exists {
			errors = append(errors, fmt.Sprintf("%s: unsupported type", file.Filename))
			continue
		}
		mimeValid := false
		for _, m := range validMIMEs {
			if contentType == m {
				mimeValid = true
				break
			}
		}
		if !mimeValid {
			errors = append(errors, fmt.Sprintf("%s: MIME mismatch", file.Filename))
			continue
		}

		// Generate filename and save
		name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		target := filepath.Join(dir, name)

		if err := c.SaveFile(file, target); err != nil {
			errors = append(errors, fmt.Sprintf("%s: save failed", file.Filename))
			continue
		}

		// Auto-compress image (sync to ensure correct filename/extension in DB & response)
		compressImage(target)

		// Check if file was converted (PNG → JPG) and update size
		finalName := name
		finalMIME := contentType
		if ext == ".png" {
			jpgName := strings.TrimSuffix(name, ext) + ".jpg"
			jpgPath := filepath.Join(dir, jpgName)
			if _, err := os.Stat(jpgPath); err == nil {
				finalName = jpgName
				finalMIME = "image/jpeg"
			}
		}

		// Update actual size after compression
		finalPath := filepath.Join(dir, finalName)
		var finalSize int64 = file.Size
		if info, err := os.Stat(finalPath); err == nil {
			finalSize = info.Size()
		}

		url := "/uploads/" + folder + "/" + finalName

		// Record in media library
		var mediaID string
		db.QueryRowContext(context.Background(),
			`INSERT INTO media (filename, original_name, url, mime_type, size_bytes, folder, uploaded_by)
			 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			finalName, file.Filename, url, finalMIME, finalSize, folder, userID,
		).Scan(&mediaID)

		uploaded = append(uploaded, fiber.Map{
			"id":            mediaID,
			"url":           url,
			"filename":      finalName,
			"original_name": file.Filename,
			"size":          finalSize,
			"mime_type":     finalMIME,
		})
	}

	return c.JSON(fiber.Map{
		"uploaded": uploaded,
		"count":   len(uploaded),
		"errors":  errors,
	})
}

// ─── Media Library API ───────────────────────────────────────────────────────

func getMediaLibrary(c *fiber.Ctx) error {
	folder := c.Query("folder", "")
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)
	if limit > 100 {
		limit = 100
	}

	query := `SELECT id, filename, original_name, url, mime_type, size_bytes, folder, created_at
		 FROM media`
	args := []interface{}{}
	argIdx := 1

	if folder != "" {
		query += fmt.Sprintf(" WHERE folder = $%d", argIdx)
		args = append(args, folder)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch media"})
	}
	defer rows.Close()

	media := []fiber.Map{}
	for rows.Next() {
		var id, filename, originalName, url, mimeType, folder, createdAt string
		var size int64
		if rows.Scan(&id, &filename, &originalName, &url, &mimeType, &size, &folder, &createdAt) == nil {
			media = append(media, fiber.Map{
				"id": id, "filename": filename, "original_name": originalName,
				"url": url, "mime_type": mimeType, "size": size,
				"folder": folder, "created_at": createdAt,
			})
		}
	}

	// Total count
	var total int
	countQuery := "SELECT COUNT(*) FROM media"
	var countArgs []interface{}
	if folder != "" {
		countQuery += " WHERE folder = $1"
		countArgs = append(countArgs, folder)
	}
	db.QueryRowContext(context.Background(), countQuery, countArgs...).Scan(&total)

	return c.JSON(fiber.Map{
		"items": media,
		"total": total,
	})
}

func deleteMedia(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get file path before deleting record
	var url string
	db.QueryRowContext(context.Background(), "SELECT url FROM media WHERE id = $1", id).Scan(&url)

	// Delete record
	_, err := db.ExecContext(context.Background(), "DELETE FROM media WHERE id = $1", id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete"})
	}

	// Delete file from disk
	if url != "" {
		filePath := filepath.Join("..", "public", url)
		os.Remove(filePath)
	}

	return c.JSON(fiber.Map{"message": "Media deleted"})
}


// bulkDeleteMedia deletes multiple media items by IDs or all of them
func bulkDeleteMedia(c *fiber.Ctx) error {
	var req struct {
		IDs    []string `json:"ids"`
		All    bool     `json:"all"`
		Folder string   `json:"folder"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var ids []string
	if req.All {
		query := "SELECT id FROM media"
		var args []interface{}
		if req.Folder != "" {
			query += " WHERE folder = $1"
			args = append(args, req.Folder)
		}
		rows, err := db.QueryContext(context.Background(), query, args...)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to query media items"})
		}
		defer rows.Close()
		for rows.Next() {
			var id string
			if rows.Scan(&id) == nil {
				ids = append(ids, id)
			}
		}
	} else {
		ids = req.IDs
	}

	if len(ids) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "No media selected to delete"})
	}

	deleted := 0
	for _, id := range ids {
		var url string
		db.QueryRowContext(context.Background(), "SELECT url FROM media WHERE id = $1", id).Scan(&url)
		if url != "" {
			filePath := filepath.Join("..", "public", url)
			os.Remove(filePath)
		}
		res, _ := db.ExecContext(context.Background(), "DELETE FROM media WHERE id = $1", id)
		if res != nil {
			n, _ := res.RowsAffected()
			deleted += int(n)
		}
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("%d media dihapus", deleted),
		"deleted": deleted,
	})
}
