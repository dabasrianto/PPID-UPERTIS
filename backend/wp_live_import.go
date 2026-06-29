package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

// WPLiveImportRequest — import from live WordPress site via REST API
type WPLiveImportRequest struct {
	WPURL          string `json:"wp_url"`
	PostType       string `json:"post_type"` // target post_type in our system
	Page           int    `json:"page"`
	PerPage        int    `json:"per_page"`
	DownloadImages *bool  `json:"download_images"`
	DateFrom       string `json:"date_from"`     // filter: start date YYYY-MM-DD
	DateTo         string `json:"date_to"`       // filter: end date YYYY-MM-DD
	OverrideDate   string `json:"override_date"` // optional: custom publish date YYYY-MM-DD
	Category       string `json:"category"`      // optional: custom category name
	ImportAll      bool   `json:"import_all"`
}

// WPAPIPost — WordPress REST API post response
type WPAPIPost struct {
	ID      int    `json:"id"`
	Date    string `json:"date"`
	Slug    string `json:"slug"`
	Status  string `json:"status"`
	Title   struct {
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
		FeaturedMedia []struct {
			SourceURL string `json:"source_url"`
		} `json:"wp:featuredmedia"`
	} `json:"_embedded"`
}

func importFromWPLive(c *fiber.Ctx) error {
	var req WPLiveImportRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	wpURL := strings.TrimRight(req.WPURL, "/")
	if wpURL == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "WordPress URL wajib diisi"})
	}

	// Prepend https:// if no protocol prefix is present
	if !strings.HasPrefix(wpURL, "http://") && !strings.HasPrefix(wpURL, "https://") {
		wpURL = "https://" + wpURL
	}

	if !isSafeURL(wpURL) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "WordPress URL is not allowed (SSRF protection)"})
	}

	postType := req.PostType
	if postType == "" {
		postType = "news"
	}
	if !validPostTypes[postType] {
		postType = "news"
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 || perPage > 50 {
		perPage = 10
	}

	// Determine download_images
	downloadImages := req.DownloadImages == nil || *req.DownloadImages

	importAll := req.ImportAll
	totalImported := 0
	totalSkipped := 0
	var errors []string
	var batchID string
	var totalPages string
	var totalPosts string

	startPage := page
	endPage := page
	if importAll {
		startPage = 1
	}

	userID, _ := c.Locals("user_id").(string)
	batchID = fmt.Sprintf("wp-live-%d", time.Now().Unix())

	p := startPage
	for {
		// Fetch from WordPress REST API
		apiURL := fmt.Sprintf("%s/wp-json/wp/v2/posts?_embed=true&per_page=%d&page=%d&status=publish", wpURL, perPage, p)
		if req.DateFrom != "" {
			apiURL += "&after=" + req.DateFrom + "T00:00:00"
		}
		if req.DateTo != "" {
			apiURL += "&before=" + req.DateTo + "T23:59:59"
		}

		client := &http.Client{Timeout: 60 * time.Second}
		httpReq, _ := http.NewRequest("GET", apiURL, nil)
		httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
		httpReq.Header.Set("Accept", "application/json, text/plain, */*")
		httpReq.Header.Set("Accept-Language", "id-ID,id;q=0.9,en;q=0.8")
		httpReq.Header.Set("Referer", wpURL+"/")
		httpReq.Header.Set("Origin", wpURL)

		resp, err := client.Do(httpReq)
		if err != nil {
			zlog.Error().Err(err).Str("wp_url", wpURL).Int("page", p).Msg("Failed to connect to WordPress")
			if p == startPage {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Gagal terhubung ke WordPress: " + err.Error()})
			}
			errors = append(errors, fmt.Sprintf("Gagal mengambil halaman %d: %v", p, err))
			break
		}

		if resp.StatusCode != 200 {
			bodySnippet, _ := io.ReadAll(io.LimitReader(resp.Body, 1000))
			resp.Body.Close()
			zlog.Warn().
				Int("status_code", resp.StatusCode).
				Str("wp_url", wpURL).
				Int("page", p).
				Str("response_body", string(bodySnippet)).
				Msg("WordPress API returned non-200 status")

			if p == startPage {
				var errMsg string
				switch resp.StatusCode {
				case 403:
					errMsg = "Akses ditolak (403). WordPress target memblokir akses API."
				case 404:
					errMsg = "REST API tidak ditemukan (404). Pastikan URL WordPress benar."
				default:
					errMsg = fmt.Sprintf("WordPress API error: HTTP %d", resp.StatusCode)
				}
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error":   errMsg,
					"details": string(bodySnippet),
				})
			}
			break // Stop looping on non-200 code
		}

		// Read total pagination headers on first page
		if p == startPage {
			totalPages = resp.Header.Get("X-WP-TotalPages")
			totalPosts = resp.Header.Get("X-WP-Total")
			if importAll && totalPages != "" {
				var totalPagesVal int
				fmt.Sscanf(totalPages, "%d", &totalPagesVal)
				if totalPagesVal > 1 {
					endPage = totalPagesVal
					if endPage > 100 {
						endPage = 100 // Cap loop to prevent timeouts
					}
				}
			}
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var wpPosts []WPAPIPost
		if err := json.Unmarshal(body, &wpPosts); err != nil {
			if p == startPage {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Gagal parse response WordPress"})
			}
			errors = append(errors, fmt.Sprintf("Gagal parse halaman %d", p))
			break
		}

		if len(wpPosts) == 0 {
			break
		}

		for _, wp := range wpPosts {
			title := sanitizeTitle(wp.Title.Rendered)
			slug := wp.Slug
			if slug == "" {
				slug = generateSlug(title)
			}

			// Check duplicate
			var existingID string
			db.QueryRowContext(c.Context(), "SELECT id FROM posts WHERE slug=$1", slug).Scan(&existingID)
			if existingID != "" {
				totalSkipped++
				continue
			}

			// Author
			author := "Admin"
			if len(wp.Embedded.Author) > 0 {
				author = wp.Embedded.Author[0].Name
			}

			// Featured image — download to local
			var featuredImage string
			if len(wp.Embedded.FeaturedMedia) > 0 && wp.Embedded.FeaturedMedia[0].SourceURL != "" {
				firstImg := wp.Embedded.FeaturedMedia[0].SourceURL
				if downloadImages {
					featuredImage = downloadImageToLocal(firstImg)
				} else {
					featuredImage = firstImg
				}
			}

			// Content — sanitize + download embedded images
			content := sanitizeWPContent(wp.Content.Rendered)
			if downloadImages {
				content = rewriteContentImages(content)
			}

			// If no featured image from media, extract from content
			if featuredImage == "" {
				firstImg := extractFirstImage(content)
				if downloadImages && firstImg != "" {
					featuredImage = downloadImageToLocal(firstImg)
				} else {
					featuredImage = firstImg
				}
			}

			// Excerpt
			excerpt := sanitizeWPContent(wp.Excerpt.Rendered)
			if excerpt == "" {
				excerpt = autoExcerpt(content, 200)
			}

			// Status
			status := "published"

			// Category
			category := "Umum"
			if req.Category != "" {
				category = req.Category
			}

			// Tags
			var tags []string

			publishDate := wp.Date
			if req.OverrideDate != "" {
				if len(req.OverrideDate) == 10 {
					originalTime := "12:00:00"
					if len(wp.Date) >= 19 {
						originalTime = wp.Date[11:19]
					}
					publishDate = req.OverrideDate + "T" + originalTime
				} else {
					publishDate = req.OverrideDate
				}
			}

			// Insert
			_, err := db.ExecContext(c.Context(), `
				INSERT INTO posts (post_type, title, slug, excerpt, content, featured_image, category, tags, author, status, published_at, created_by, import_batch_id, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW())`,
				postType, title, slug, excerpt, content, featuredImage,
				category, pq.Array(tags), author, status, publishDate,
				userID, batchID,
			)
			if err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", title, err))
				continue
			}
			totalImported++
		}

		p++
		if p > endPage {
			break
		}
	}

	zlog.Info().
		Str("batch_id", batchID).
		Str("wp_url", wpURL).
		Int("imported", totalImported).
		Int("skipped", totalSkipped).
		Int("errors", len(errors)).
		Msg("WordPress live import completed")

	msg := fmt.Sprintf("Import selesai: %d berhasil, %d dilewati, %d gagal", totalImported, totalSkipped, len(errors))
	if importAll {
		msg = fmt.Sprintf("Import Semua Halaman selesai: %d berhasil, %d dilewati, %d gagal", totalImported, totalSkipped, len(errors))
	}

	return c.JSON(fiber.Map{
		"message":     msg,
		"imported":    totalImported,
		"skipped":     totalSkipped,
		"failed":      len(errors),
		"errors":      errors,
		"batch_id":    batchID,
		"total_pages": totalPages,
		"total_posts": totalPosts,
		"page":        page,
	})
}
