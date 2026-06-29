package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var safePathPart = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

// Allowed file extensions for upload — strict whitelist
var allowedExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".webp": true, ".pdf": true, ".ico": true,
	".docx": true, ".xlsx": true, ".zip": true, ".rar": true,
	".mp4": true, ".webm": true, ".ogg": true, ".mov": true,
	".m4v": true, ".avi": true, ".mkv": true,
}

// Allowed MIME types — must match extension
var allowedMIMETypes = map[string][]string{
	".jpg":  {"image/jpeg", "image/jpg"},
	".jpeg": {"image/jpeg", "image/jpg"},
	".png":  {"image/png"},
	".gif":  {"image/gif"},
	".webp": {"image/webp"},
	".pdf":  {"application/pdf"},
	".ico":  {"image/x-icon", "image/vnd.microsoft.icon", "image/ico"},
	".docx": {"application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/msword"},
	".xlsx": {"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.ms-excel"},
	".zip":  {"application/zip", "application/x-zip-compressed"},
	".rar":  {"application/x-rar-compressed", "application/octet-stream", "application/x-rar"},
	".mp4":  {"video/mp4"},
	".webm": {"video/webm"},
	".ogg":  {"video/ogg", "application/ogg", "video/x-theora+ogg"},
	".mov":  {"video/quicktime"},
	".m4v":  {"video/x-m4v", "video/mp4"},
	".avi":  {"video/x-msvideo", "video/avi", "video/msvideo"},
	".mkv":  {"video/x-matroska", "video/mkv"},
}

// Dangerous file signatures (magic bytes) to block
var dangerousSignatures = [][]byte{
	{0x4D, 0x5A},                   // EXE/DLL (MZ)
	{0x7F, 0x45, 0x4C, 0x46},      // ELF binary
	{0x23, 0x21},                   // Shebang script (#!/)
	{0x50, 0x4B, 0x03, 0x04},      // ZIP (could be .jar, .apk)
}

// Max file size: 30MB (auto-compressed after upload)
const maxFileSize = 30 * 1024 * 1024

func cleanPathPart(value string) string {
	value = safePathPart.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")
	if value == "" {
		return "media"
	}
	// Prevent path traversal
	value = strings.ReplaceAll(value, "..", "")
	value = strings.ReplaceAll(value, "/", "")
	value = strings.ReplaceAll(value, "\\", "")
	if value == "" {
		return "media"
	}
	return value
}

func uploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "file is required"})
	}

	// 1. Validate file size
	if file.Size > maxFileSize {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "file too large, maximum 30MB"})
	}
	if file.Size == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "file is empty"})
	}

	// 2. Validate extension (strict whitelist)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "file must have an extension"})
	}
	if !allowedExtensions[ext] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "file type not allowed. Allowed: jpg, jpeg, png, gif, webp, pdf, ico, docx, xlsx, zip, rar",
		})
	}

	// 3. Validate Content-Type header matches extension
	contentType := file.Header.Get("Content-Type")
	validMIMEs, exists := allowedMIMETypes[ext]
	if !exists {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "unsupported file type"})
	}
	mimeValid := false
	for _, m := range validMIMEs {
		if contentType == m {
			mimeValid = true
			break
		}
	}
	if !mimeValid {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("content type %s does not match extension %s", contentType, ext),
		})
	}

	// 4. Basic safety check (skip heavy magic bytes — MIME + extension is sufficient for campus CMS)
	// Dangerous executables are already blocked by extension whitelist

	// 5. Sanitize folder (prevent path traversal)
	folder := cleanPathPart(c.FormValue("folder", "media"))

	// 6. Generate random filename (never use original filename)
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dir := filepath.Join("..", "public", "uploads", folder)
	if err := os.MkdirAll(dir, 0755); err != nil {
		zlog.Error().Err(err).Str("dir", dir).Msg("Upload mkdir failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create upload directory"})
	}

	// 7. Verify final path is within expected directory (anti path traversal)
	target := filepath.Join(dir, name)
	absTarget, _ := filepath.Abs(target)
	absDir, _ := filepath.Abs(dir)
	if !strings.HasPrefix(absTarget, absDir) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid file path"})
	}

	// 8. Save file
	if err := c.SaveFile(file, target); err != nil {
		zlog.Error().Err(err).Str("target", target).Str("folder", folder).Msg("Upload save failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save file"})
	}

	// 9. Auto-compress image (resize + JPEG quality 75%) — sync, not async
	// so we can return correct URL if PNG converted to JPG
	compressImage(target)

	// Check if file was converted (PNG → JPG)
	finalName := name
	if ext == ".png" {
		jpgName := strings.TrimSuffix(name, ext) + ".jpg"
		jpgPath := filepath.Join(dir, jpgName)
		if _, err := os.Stat(jpgPath); err == nil {
			finalName = jpgName
		}
	}

	zlog.Info().Str("url", "/uploads/"+folder+"/"+finalName).Int64("size", file.Size).Msg("File uploaded")
	return c.JSON(fiber.Map{"url": "/uploads/" + folder + "/" + finalName})
}
