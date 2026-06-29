package main

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// KnowledgeFile represents one uploaded knowledge document
type KnowledgeFile struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Content string `json:"content"`
}

// handleKnowledgeUpload accepts a file upload, extracts text, and returns it
func handleKnowledgeUpload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
	}

	// Limit to 5MB
	if file.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File too large (max 5MB)"})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{
		".txt": true, ".md": true, ".csv": true,
		".pdf": true, ".docx": true, ".xlsx": true,
	}
	if !allowed[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unsupported file type. Allowed: TXT, MD, CSV, PDF, DOCX, XLSX",
		})
	}

	// Save file to uploads/chatbot-knowledge/
	uploadDir := "./uploads/chatbot-knowledge"
	os.MkdirAll(uploadDir, 0755)

	filename := fmt.Sprintf("%d_%s", makeTimestamp(), sanitizeFilename(file.Filename))
	savePath := filepath.Join(uploadDir, filename)
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	// Extract text content
	content, err := extractTextFromFile(savePath, ext)
	if err != nil {
		log.Printf("Knowledge extraction error for %s: %v", file.Filename, err)
		// Still return with empty content
		content = fmt.Sprintf("[Gagal mengekstrak teks dari %s]", file.Filename)
	}

	// Trim very long content (max ~30k chars to keep prompt manageable)
	if len(content) > 30000 {
		content = content[:30000] + "\n...[konten dipotong karena terlalu panjang]"
	}

	url := "/uploads/chatbot-knowledge/" + filename

	return c.JSON(fiber.Map{
		"name":    file.Filename,
		"url":     url,
		"content": content,
		"size":    file.Size,
	})
}

func sanitizeFilename(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	return re.ReplaceAllString(name, "_")
}

func makeTimestamp() int64 {
	return time.Now().UnixMilli()
}

// extractTextFromFile reads a file and extracts plain text based on type
func extractTextFromFile(path, ext string) (string, error) {
	switch ext {
	case ".txt", ".md":
		return readPlainText(path)
	case ".csv":
		return readCSV(path)
	case ".docx":
		return readDOCX(path)
	case ".xlsx":
		return readXLSX(path)
	case ".pdf":
		return readPDFBasic(path)
	default:
		return "", fmt.Errorf("unsupported format: %s", ext)
	}
}

func readPlainText(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func readCSV(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	var sb strings.Builder
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		sb.WriteString(strings.Join(record, " | "))
		sb.WriteString("\n")
	}
	return sb.String(), nil
}

// readDOCX extracts text from a .docx file (ZIP containing XML)
func readDOCX(path string) (string, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()
			data, err := io.ReadAll(rc)
			if err != nil {
				return "", err
			}
			return stripXMLTags(string(data)), nil
		}
	}
	return "", fmt.Errorf("document.xml not found in docx")
}

// readXLSX extracts text from a .xlsx file (ZIP containing shared strings XML)
func readXLSX(path string) (string, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var sb strings.Builder

	// Read shared strings first
	sharedStrings := []string{}
	for _, f := range r.File {
		if f.Name == "xl/sharedStrings.xml" {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			data, _ := io.ReadAll(rc)
			rc.Close()
			// Extract <t>...</t> values
			re := regexp.MustCompile(`<t[^>]*>([^<]+)</t>`)
			matches := re.FindAllStringSubmatch(string(data), -1)
			for _, m := range matches {
				sharedStrings = append(sharedStrings, m[1])
			}
		}
	}

	// Read sheet data
	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "xl/worksheets/sheet") && strings.HasSuffix(f.Name, ".xml") {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			data, _ := io.ReadAll(rc)
			rc.Close()

			text := stripXMLTags(string(data))
			if text != "" {
				sb.WriteString(text)
				sb.WriteString("\n")
			}
		}
	}

	// Add shared strings content
	if len(sharedStrings) > 0 {
		sb.WriteString("\n")
		sb.WriteString(strings.Join(sharedStrings, " | "))
	}

	return sb.String(), nil
}

// readPDFBasic does a basic text extraction from PDF by reading raw streams
func readPDFBasic(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	content := string(data)
	var sb strings.Builder

	// Extract text between BT...ET (text blocks in PDF)
	re := regexp.MustCompile(`BT\s*(.*?)\s*ET`)
	blocks := re.FindAllStringSubmatch(content, -1)

	for _, block := range blocks {
		// Extract text from Tj and TJ operators
		tjRe := regexp.MustCompile(`\(([^)]*)\)\s*Tj`)
		matches := tjRe.FindAllStringSubmatch(block[1], -1)
		for _, m := range matches {
			text := m[1]
			// Unescape basic PDF escapes
			text = strings.ReplaceAll(text, "\\(", "(")
			text = strings.ReplaceAll(text, "\\)", ")")
			text = strings.ReplaceAll(text, "\\\\", "\\")
			sb.WriteString(text)
		}

		// Also try TJ arrays
		tjArrayRe := regexp.MustCompile(`\[([^\]]*)\]\s*TJ`)
		arrayMatches := tjArrayRe.FindAllStringSubmatch(block[1], -1)
		for _, m := range arrayMatches {
			innerRe := regexp.MustCompile(`\(([^)]*)\)`)
			innerMatches := innerRe.FindAllStringSubmatch(m[1], -1)
			for _, im := range innerMatches {
				sb.WriteString(im[1])
			}
		}
		sb.WriteString(" ")
	}

	result := sb.String()
	if strings.TrimSpace(result) == "" {
		return "[PDF ini menggunakan encoding yang tidak bisa diekstrak secara otomatis. Silakan konversi ke format TXT/DOCX terlebih dahulu.]", nil
	}
	return result, nil
}

// stripXMLTags removes all XML tags and returns plain text
func stripXMLTags(s string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	text := re.ReplaceAllString(s, " ")
	// Clean up multiple spaces
	spaceRe := regexp.MustCompile(`\s+`)
	text = spaceRe.ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}

// init ensures upload directory exists
func init() {
	os.MkdirAll("./uploads/chatbot-knowledge", 0755)
}
