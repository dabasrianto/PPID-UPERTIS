package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type WPSyncRequest struct {
	WPUrl      string `json:"wp_url"`
	CategoryID string `json:"category_id"`
	DateFrom   string `json:"date_from"`
	DateTo     string `json:"date_to"`
	Page       int    `json:"page"`
}

type WPPost struct {
	ID    int `json:"id"`
	Date  string `json:"date"`
	Slug  string `json:"slug"`
	Title struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
	Content struct {
		Rendered string `json:"rendered"`
	} `json:"content"`
	Excerpt struct {
		Rendered string `json:"rendered"`
	} `json:"excerpt"`
	Categories []int `json:"categories"`
	Tags       []int `json:"tags"`
	Embedded   struct {
		Author []struct {
			Name string `json:"name"`
		} `json:"author"`
		WpFeaturedmedia []struct {
			SourceUrl string `json:"source_url"`
		} `json:"wp:featuredmedia"`
	} `json:"_embedded"`
}

// Ensure slug is unique
func generateUniqueSlug(base string) string {
	var count int
	err := db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM blog_posts WHERE slug=$1", base).Scan(&count)
	if err == nil && count == 0 {
		return base
	}
	return fmt.Sprintf("%s-%d", base, time.Now().Unix())
}

// isSafeURL checks if a URL is safe from SSRF attacks (only allows http/https, blocks local/private IPs and loopbacks)
func isSafeURL(targetURL string) bool {
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		return false
	}

	lowerURL := strings.ToLower(targetURL)
	blockedPrefixes := []string{
		"http://localhost", "https://localhost",
		"http://127.", "https://127.",
		"http://10.", "https://10.",
		"http://172.16.", "https://172.16.",
		"http://192.168.", "https://192.168.",
		"http://0.0.0.0", "https://0.0.0.0",
		"http://[::1]", "https://[::1]",
	}
	for _, prefix := range blockedPrefixes {
		if strings.HasPrefix(lowerURL, prefix) {
			return false
		}
	}
	return true
}

// Download image and return local URL
func downloadWPImage(sourceUrl string) string {
	if sourceUrl == "" {
		return ""
	}

	if !isSafeURL(sourceUrl) {
		log.Println("downloadWPImage: blocked unsafe URL:", sourceUrl)
		return ""
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(sourceUrl)
	if err != nil {
		log.Println("Error downloading WP image:", err)
		return sourceUrl // fallback to external URL
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return sourceUrl
	}

	ext := filepath.Ext(sourceUrl)
	if ext == "" || strings.Contains(ext, "?") {
		ext = ".jpg"
	} else {
		// Clean extension
		parts := strings.Split(ext, "?")
		ext = parts[0]
	}

	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dir := filepath.Join("..", "public", "uploads", "blog")
	os.MkdirAll(dir, 0755)

	target := filepath.Join(dir, fileName)
	file, err := os.Create(target)
	if err != nil {
		log.Println("Error creating local image file:", err)
		return sourceUrl
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return sourceUrl
	}

	return "/uploads/blog/" + fileName
}

func syncWordPressPosts(c *fiber.Ctx) error {
	var req WPSyncRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	wpUrl := strings.TrimRight(req.WPUrl, "/")
	if wpUrl == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "WordPress URL is required"})
	}

	// Prepend https:// if no protocol prefix is present
	if !strings.HasPrefix(wpUrl, "http://") && !strings.HasPrefix(wpUrl, "https://") {
		wpUrl = "https://" + wpUrl
	}

	if !isSafeURL(wpUrl) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "WordPress URL is not allowed (SSRF protection)"})
	}

	page := req.Page
	if page < 1 {
		page = 1
	}

	apiEndpoint := fmt.Sprintf("%s/wp-json/wp/v2/posts?_embed=true&per_page=10&page=%d", wpUrl, page)
	if req.CategoryID != "" {
		apiEndpoint += "&categories=" + req.CategoryID
	}
	if req.DateFrom != "" {
		apiEndpoint += "&after=" + req.DateFrom + "T00:00:00"
	}
	if req.DateTo != "" {
		apiEndpoint += "&before=" + req.DateTo + "T23:59:59"
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	reqHttp, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat request ke WordPress"})
	}
	
	// Add User-Agent to avoid being blocked by security plugins
	reqHttp.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(reqHttp)
	if err != nil {
		log.Println("Error fetching WP posts:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal terhubung ke WordPress: " + err.Error()})
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("WP API returned status:", resp.StatusCode)
		if resp.StatusCode == 403 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Akses ditolak (403). WordPress Anda mungkin memblokir akses REST API."})
		}
		if resp.StatusCode == 404 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "REST API tidak ditemukan (404). Pastikan URL WordPress sudah benar."})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("WordPress API mengembalikan error status: %d", resp.StatusCode)})
	}

	totalPagesStr := resp.Header.Get("X-WP-TotalPages")
	totalPages, _ := strconv.Atoi(totalPagesStr)

	var wpPosts []WPPost
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &wpPosts); err != nil {
		bodyLen := len(bodyBytes)
		snippetLen := 500
		if bodyLen < snippetLen {
			snippetLen = bodyLen
		}
		log.Printf("Error decoding WP posts: %v. Body snippet: %s", err, string(bodyBytes[:snippetLen]))
		
		errorLen := 200
		if bodyLen < errorLen {
			errorLen = bodyLen
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to parse WordPress data",
			"details": string(bodyBytes[:errorLen]),
		})
	}

	importedCount := 0

	for _, wp := range wpPosts {
		// Check if post already exists by slug (WP slug)
		var existingId string
		err := db.QueryRowContext(context.Background(), "SELECT id FROM blog_posts WHERE slug=$1", wp.Slug).Scan(&existingId)
		if err == nil && existingId != "" {
			// Skip existing to avoid duplicates
			continue
		}

		title := sanitizeTitle(wp.Title.Rendered)
		content := wp.Content.Rendered
		excerpt := wp.Excerpt.Rendered
		slug := generateUniqueSlug(wp.Slug)
		
		var authorName string
		if len(wp.Embedded.Author) > 0 {
			authorName = wp.Embedded.Author[0].Name
		} else {
			authorName = "Admin"
		}

		var coverImageURL string
		if len(wp.Embedded.WpFeaturedmedia) > 0 {
			sourceUrl := wp.Embedded.WpFeaturedmedia[0].SourceUrl
			coverImageURL = downloadWPImage(sourceUrl)
		}

		// Insert into db
		_, err = db.ExecContext(context.Background(),
			`INSERT INTO blog_posts (title, slug, excerpt, content, cover_image_url, author_name, category, status, published_at, seo_title, seo_description)
			VALUES ($1, $2, $3, $4, $5, $6, $7, 'published', $8, $9, $10)`,
			title, slug, excerpt, content, coverImageURL, authorName, "Berita Kampus", wp.Date, title, excerpt)
		
		if err != nil {
			log.Printf("Error inserting WP post [%s]: %v", wp.Slug, err)
		} else {
			importedCount++
		}
	}

	totalInPage := len(wpPosts)
	skippedCount := totalInPage - importedCount

	return c.JSON(fiber.Map{
		"imported":    importedCount,
		"skipped":     skippedCount,
		"total_pages": totalPages,
		"page":        page,
		"message":     fmt.Sprintf("Berhasil memproses %d artikel (Baru: %d, Lewati: %d)", totalInPage, importedCount, skippedCount),
	})
}
