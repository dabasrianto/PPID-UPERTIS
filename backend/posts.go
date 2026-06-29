package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

// ─── Posts Table Migration ───────────────────────────────────────────────────

func migratePostsTable() {
	// Create migrations tracking table
	db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS _migrations (
			name VARCHAR(255) PRIMARY KEY,
			executed_at TIMESTAMP DEFAULT NOW()
		)
	`)

	// Create posts table (WordPress-like structure)
	_, err := db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS posts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			post_type VARCHAR(50) NOT NULL DEFAULT 'news',
			title VARCHAR(255) NOT NULL,
			slug VARCHAR(255) UNIQUE NOT NULL,
			excerpt TEXT,
			content TEXT,
			featured_image VARCHAR(500),
			category VARCHAR(100),
			tags TEXT[],
			author VARCHAR(255) DEFAULT 'Admin',
			status VARCHAR(50) DEFAULT 'draft',
			featured BOOLEAN DEFAULT false,
			published_at TIMESTAMP,
			seo_title VARCHAR(255),
			seo_description TEXT,
			views INTEGER DEFAULT 0,
			created_by VARCHAR(255),
			updated_by VARCHAR(255),
			deleted_at TIMESTAMP,
			import_batch_id VARCHAR(100),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Printf("Warning: Failed to create posts table: %v", err)
		return
	}

	// Add columns if missing (safe for existing tables)
	alterStatements := []string{
		`ALTER TABLE posts ADD COLUMN IF NOT EXISTS created_by VARCHAR(255)`,
		`ALTER TABLE posts ADD COLUMN IF NOT EXISTS updated_by VARCHAR(255)`,
		`ALTER TABLE posts ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP`,
		`ALTER TABLE posts ADD COLUMN IF NOT EXISTS import_batch_id VARCHAR(100)`,
		`ALTER TABLE posts ADD COLUMN IF NOT EXISTS scheduled_at TIMESTAMP`,
	}
	for _, stmt := range alterStatements {
		db.ExecContext(context.Background(), stmt)
	}

	// Indexes (idempotent — IF NOT EXISTS)
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug)`,
		`CREATE INDEX IF NOT EXISTS idx_posts_post_type ON posts(post_type)`,
		`CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status)`,
		`CREATE INDEX IF NOT EXISTS idx_posts_type_status ON posts(post_type, status)`,
		`CREATE INDEX IF NOT EXISTS idx_posts_published_at ON posts(published_at DESC NULLS LAST)`,
		`CREATE INDEX IF NOT EXISTS idx_posts_featured ON posts(featured)`,
		`CREATE INDEX IF NOT EXISTS idx_posts_category ON posts(category)`,
		`CREATE INDEX IF NOT EXISTS idx_posts_deleted_at ON posts(deleted_at)`,
	}
	for _, idx := range indexes {
		db.ExecContext(context.Background(), idx)
	}

	// Run data migration only once (tracked)
	if !migrationCompleted("migrate_news_blog_to_posts") {
		zlog.Info().Msg("Migration: starting news+blog → posts")
		migrateNewsToPostsTable()
		migrateBlogToPostsTable()
		markMigrationCompleted("migrate_news_blog_to_posts")
		zlog.Info().Msg("Migration: completed")
	} else {
		zlog.Info().Msg("Migration: already completed, skipping")
	}
}

// migrationCompleted checks if a named migration has already run
func migrationCompleted(name string) bool {
	var count int
	db.QueryRowContext(context.Background(),
		"SELECT COUNT(*) FROM _migrations WHERE name = $1", name).Scan(&count)
	return count > 0
}

// markMigrationCompleted records that a migration has been executed
func markMigrationCompleted(name string) {
	db.ExecContext(context.Background(),
		"INSERT INTO _migrations (name) VALUES ($1) ON CONFLICT DO NOTHING", name)
}

func migrateNewsToPostsTable() {
	rows, err := db.QueryContext(context.Background(),
		`SELECT id, title, COALESCE(slug, ''), COALESCE(excerpt,''), COALESCE(content,''), 
		 COALESCE(image_url,''), COALESCE(category,'Umum'), date, featured, 
		 COALESCE(author_name,'Admin'), COALESCE(seo_title,''), COALESCE(seo_description,''),
		 created_at, updated_at
		 FROM news WHERE COALESCE(active, true) = true`)
	if err != nil {
		log.Printf("Migration: failed to read news: %v", err)
		return
	}
	defer rows.Close()

	migrated := 0
	for rows.Next() {
		var id, title, slug, excerpt, content, imageURL, category, author, seoTitle, seoDesc, createdAt, updatedAt string
		var date string
		var featured bool

		if err := rows.Scan(&id, &title, &slug, &excerpt, &content, &imageURL, &category, &date, &featured, &author, &seoTitle, &seoDesc, &createdAt, &updatedAt); err != nil {
			log.Printf("Migration: scan error: %v", err)
			continue
		}

		if slug == "" {
			slug = generateSlug(title)
		}
		slug = ensureUniquePostSlug(slug)

		_, err := db.ExecContext(context.Background(),
			`INSERT INTO posts (id, post_type, title, slug, excerpt, content, featured_image, category, author, status, featured, published_at, seo_title, seo_description, created_at, updated_at)
			 VALUES ($1, 'news', $2, $3, $4, $5, $6, $7, $8, 'published', $9, $10, $11, $12, $13, $14)
			 ON CONFLICT (slug) DO NOTHING`,
			id, title, slug, excerpt, content, imageURL, category, author, featured, date, seoTitle, seoDesc, createdAt, updatedAt)
		if err != nil {
			log.Printf("Migration: insert news error [%s]: %v", title, err)
			continue
		}
		migrated++
	}
	log.Printf("Migration: %d news items migrated to posts table", migrated)
}

func migrateBlogToPostsTable() {
	rows, err := db.QueryContext(context.Background(),
		`SELECT id, title, slug, COALESCE(excerpt,''), COALESCE(content,''), 
		 COALESCE(cover_image_url,''), COALESCE(category,'Umum'), COALESCE(tags,'{}'),
		 COALESCE(author_name,'Admin'), status,
		 COALESCE(seo_title,''), COALESCE(seo_description,''), COALESCE(views,0),
		 COALESCE(published_at, created_at), created_at, updated_at
		 FROM blog_posts`)
	if err != nil {
		log.Printf("Migration: failed to read blog_posts: %v", err)
		return
	}
	defer rows.Close()

	migrated := 0
	for rows.Next() {
		var id, title, slug, excerpt, content, imageURL, category, author, status, seoTitle, seoDesc, publishedAt, createdAt, updatedAt string
		var tags []string
		var views int

		if err := rows.Scan(&id, &title, &slug, &excerpt, &content, &imageURL, &category, pq.Array(&tags), &author, &status, &seoTitle, &seoDesc, &views, &publishedAt, &createdAt, &updatedAt); err != nil {
			log.Printf("Migration: blog scan error: %v", err)
			continue
		}

		slug = ensureUniquePostSlug(slug)

		_, err := db.ExecContext(context.Background(),
			`INSERT INTO posts (id, post_type, title, slug, excerpt, content, featured_image, category, tags, author, status, published_at, seo_title, seo_description, views, created_at, updated_at)
			 VALUES ($1, 'article', $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
			 ON CONFLICT (slug) DO NOTHING`,
			id, title, slug, excerpt, content, imageURL, category, pq.Array(tags), author, status, publishedAt, seoTitle, seoDesc, views, createdAt, updatedAt)
		if err != nil {
			log.Printf("Migration: insert blog error [%s]: %v", title, err)
			continue
		}
		migrated++
	}
	log.Printf("Migration: %d blog posts migrated to posts table", migrated)
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == ' ' {
			return r
		}
		return -1
	}, slug)
	slug = strings.ReplaceAll(slug, " ", "-")
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "post"
	}
	if len(slug) > 80 {
		slug = slug[:80]
	}
	return slug
}

func ensureUniquePostSlug(slug string) string {
	var count int
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM posts WHERE slug=$1", slug).Scan(&count)
	if count == 0 {
		return slug
	}
	return fmt.Sprintf("%s-%d", slug, time.Now().Unix())
}

// sanitizeTitle removes HTML tags (like <br>, <strong>, etc.) and normalizes whitespace in the title.
func sanitizeTitle(title string) string {
	if title == "" {
		return ""
	}

	// 1. Unescape HTML entities first (so &lt;br&gt; becomes <br> before regex)
	title = html.UnescapeString(title)

	// 2. Replace br tags and newlines with space
	brRe := regexp.MustCompile(`(?i)</?br\s*/?>`)
	title = brRe.ReplaceAllString(title, " ")

	// 3. Remove all other HTML tags
	htmlTagRe := regexp.MustCompile(`<[^>]*>`)
	title = htmlTagRe.ReplaceAllString(title, "")

	// 4. Normalize multiple spaces
	multiSpaceRe := regexp.MustCompile(`\s+`)
	title = multiSpaceRe.ReplaceAllString(title, " ")

	// 5. Trim
	return strings.TrimSpace(title)
}

var validPostTypes = map[string]bool{
	"news": true, "article": true, "announcement": true,
}

// sanitizeWPContent strips WordPress Gutenberg block comments, PageLayer markup,
// Microsoft Word artifacts, normalizes Gutenberg blocks, and leaves clean HTML.
func sanitizeWPContent(content string) string {
	if content == "" {
		return ""
	}

	// 1. Remove all <!-- wp:... --> and <!-- /wp:... --> block comments
	wpBlockRe := regexp.MustCompile(`<!--\s*/?wp:[^>]*-->`)
	content = wpBlockRe.ReplaceAllString(content, "")

	// 2. Remove remaining HTML comments
	htmlCommentRe := regexp.MustCompile(`<!--[\s\S]*?-->`)
	content = htmlCommentRe.ReplaceAllString(content, "")

	// 3. Remove Microsoft Word tags: <o:p>, </o:p>, <w:...>, etc.
	msTagRe := regexp.MustCompile(`</?[ovwx]:[^>]*>`)
	content = msTagRe.ReplaceAllString(content, "")

	// 4. Remove class="Mso..." attributes (Word classes)
	msoClassRe := regexp.MustCompile(`\s*class="[^"]*Mso[^"]*"`)
	content = msoClassRe.ReplaceAllString(content, "")

	// 5. Remove inline style attributes (Word/PageLayer junk)
	styleRe := regexp.MustCompile(`\s*style="[^"]*"`)
	content = styleRe.ReplaceAllString(content, "")

	// 6. Normalize Gutenberg classes
	wpClassRe := regexp.MustCompile(`\s*class="wp-[^"]*"`)
	content = wpClassRe.ReplaceAllString(content, "")

	// 7. Remove useless empty spans from Word paste
	emptySpanRe := regexp.MustCompile(`<span[^>]*>\s*</span>`)
	content = emptySpanRe.ReplaceAllString(content, "")
	bareSpanRe := regexp.MustCompile(`<span>([^<]*)</span>`)
	content = bareSpanRe.ReplaceAllString(content, "$1")

	// 8. Remove empty <p> tags, <p>&nbsp;</p>, <p><br></p>
	emptyParaRe := regexp.MustCompile(`<p[^>]*>\s*(&nbsp;|\s|<br\s*/?>)*\s*</p>`)
	content = emptyParaRe.ReplaceAllString(content, "")

	// 9. Collapse 3+ consecutive <br> into paragraph break
	multiBrRe := regexp.MustCompile(`(<br\s*/?>[\s]*){3,}`)
	content = multiBrRe.ReplaceAllString(content, "</p><p>")

	// 10. Replace &nbsp; with space
	content = strings.ReplaceAll(content, "&nbsp;", " ")

	// 11. Clean excessive whitespace
	multiNewlineRe := regexp.MustCompile(`\n{3,}`)
	content = multiNewlineRe.ReplaceAllString(content, "\n\n")
	multiSpaceRe := regexp.MustCompile(`  +`)
	content = multiSpaceRe.ReplaceAllString(content, " ")

	// 12. Trim
	content = strings.TrimSpace(content)

	return content
}

// extractFirstImage finds the first <img src="..."> in HTML content
func extractFirstImage(content string) string {
	imgRe := regexp.MustCompile(`<img[^>]+src="([^"]+)"`)
	matches := imgRe.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// downloadImageToLocal downloads a remote image and saves it to local uploads.
// Returns local URL path or original URL if download fails.
func downloadImageToLocal(remoteURL string) string {
	if remoteURL == "" {
		return ""
	}
	// Skip if already local
	if strings.HasPrefix(remoteURL, "/uploads/") {
		return remoteURL
	}
	// Only download http/https
	if !strings.HasPrefix(remoteURL, "http://") && !strings.HasPrefix(remoteURL, "https://") {
		return remoteURL
	}
	// SSRF protection
	if !isSafeURL(remoteURL) {
		zlog.Warn().Str("url", remoteURL).Msg("WP import: blocked unsafe URL (SSRF)")
		return remoteURL
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", remoteURL, nil)
	if err != nil {
		return remoteURL
	}
	// Browser-like headers to bypass Cloudflare/WAF
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	req.Header.Set("Accept-Language", "id-ID,id;q=0.9,en;q=0.8")
	req.Header.Set("Referer", remoteURL)

	resp, err := client.Do(req)
	if err != nil {
		zlog.Warn().Str("url", remoteURL).Err(err).Msg("WP import: failed to download image")
		return remoteURL
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		zlog.Warn().Str("url", remoteURL).Int("status", resp.StatusCode).Msg("WP import: image download non-200")
		return remoteURL
	}

	// Determine extension from URL or content-type
	ext := filepath.Ext(remoteURL)
	if ext == "" || strings.Contains(ext, "?") {
		ct := resp.Header.Get("Content-Type")
		switch {
		case strings.Contains(ct, "jpeg"), strings.Contains(ct, "jpg"):
			ext = ".jpg"
		case strings.Contains(ct, "png"):
			ext = ".png"
		case strings.Contains(ct, "webp"):
			ext = ".webp"
		case strings.Contains(ct, "gif"):
			ext = ".gif"
		default:
			ext = ".jpg"
		}
	}
	// Clean extension (remove query params)
	if idx := strings.Index(ext, "?"); idx > 0 {
		ext = ext[:idx]
	}

	// Save to uploads/posts/
	dir := filepath.Join("..", "public", "uploads", "posts")
	os.MkdirAll(dir, 0755)

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	target := filepath.Join(dir, filename)

	file, err := os.Create(target)
	if err != nil {
		return remoteURL
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		os.Remove(target)
		return remoteURL
	}

	localURL := "/uploads/posts/" + filename
	zlog.Info().Str("remote", remoteURL).Str("local", localURL).Msg("WP import: image downloaded")

	// Get file size
	var size int64
	if info, err := os.Stat(target); err == nil {
		size = info.Size()
	}

	// Determine MIME type
	mimeType := "image/jpeg"
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".webp":
		mimeType = "image/webp"
	case ".gif":
		mimeType = "image/gif"
	}

	// Insert into media library database so it appears in Admin -> Galeri Media
	_, _ = db.ExecContext(context.Background(),
		`INSERT INTO media (filename, original_name, url, mime_type, size_bytes, folder, uploaded_by)
		 VALUES ($1, $2, $3, $4, $5, 'posts', 'WP_Import')
		 ON CONFLICT (url) DO NOTHING`,
		filename, filename, localURL, mimeType, size,
	)

	return localURL
}

// rewriteContentImages downloads all images in HTML content to local storage
func rewriteContentImages(content string) string {
	if content == "" {
		return ""
	}
	imgRe := regexp.MustCompile(`<img([^>]+)src="([^"]+)"`)
	return imgRe.ReplaceAllStringFunc(content, func(match string) string {
		submatches := imgRe.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}
		remoteURL := submatches[2]
		if strings.HasPrefix(remoteURL, "/uploads/") || !strings.HasPrefix(remoteURL, "http") {
			return match // already local
		}
		localURL := downloadImageToLocal(remoteURL)
		return strings.Replace(match, remoteURL, localURL, 1)
	})
}

// autoExcerpt generates an excerpt from HTML content if none provided
func autoExcerpt(content string, maxLen int) string {
	tagRe := regexp.MustCompile(`<[^>]+>`)
	text := tagRe.ReplaceAllString(content, "")
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.TrimSpace(text)
	spaceRe := regexp.MustCompile(`\s+`)
	text = spaceRe.ReplaceAllString(text, " ")
	if len(text) <= maxLen {
		return text
	}
	cut := text[:maxLen]
	lastSpace := strings.LastIndex(cut, " ")
	if lastSpace > maxLen/2 {
		cut = cut[:lastSpace]
	}
	return cut + "..."
}

// ─── Posts API Handlers ──────────────────────────────────────────────────────

// startScheduledPublisher runs a background goroutine that publishes
// posts with scheduled_at <= NOW() every 60 seconds.
func startScheduledPublisher() {
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			res, err := db.ExecContext(context.Background(),
				`UPDATE posts SET status = 'published', published_at = scheduled_at, scheduled_at = NULL
				 WHERE status = 'scheduled' AND scheduled_at IS NOT NULL AND scheduled_at <= NOW() AND deleted_at IS NULL`)
			if err != nil {
				zlog.Error().Err(err).Msg("Scheduled publisher error")
				continue
			}
			count, _ := res.RowsAffected()
			if count > 0 {
				zlog.Info().Int64("count", count).Msg("Scheduled posts published")
			}
		}
	}()
	zlog.Info().Msg("Scheduled publisher started (interval: 60s)")
}

// Public: get active post categories
func getPostCategories(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(),
		`SELECT id, name, COALESCE(type, 'all') FROM event_categories WHERE active = true ORDER BY name`)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch categories"})
	}
	defer rows.Close()

	cats := []fiber.Map{}
	for rows.Next() {
		var id, name, catType string
		if rows.Scan(&id, &name, &catType) == nil {
			cats = append(cats, fiber.Map{"id": id, "name": name, "type": catType})
		}
	}
	return c.JSON(cats)
}

// Public: get posts with pagination and filters
func getPosts(c *fiber.Ctx) error {
	postType := c.Query("type", "")
	category := c.Query("category", "")
	featured := c.Query("featured", "")
	search := c.Query("search", "")
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	siteSlug := getSiteID(c)
	if limit > 100 {
		limit = 100
	}

	query := `SELECT id, post_type, title, slug, COALESCE(excerpt,''), COALESCE(featured_image,''),
		 COALESCE(category,''), COALESCE(author,'Admin'), status, featured,
		 COALESCE(published_at, created_at), COALESCE(views,0), created_at
		 FROM posts WHERE status = 'published' AND deleted_at IS NULL`
	args := []interface{}{}
	argIdx := 1

	if postType != "" {
		query += fmt.Sprintf(" AND post_type = $%d", argIdx)
		args = append(args, postType)
		argIdx++
	}
	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}
	if featured == "true" {
		query += " AND featured = true"
	}
	if search != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR excerpt ILIKE $%d OR content ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	// Multi-site: filter by faculty
	query, args, argIdx = buildFacultyFilter(query, args, argIdx, siteSlug)

	query += fmt.Sprintf(" ORDER BY COALESCE(published_at, created_at) DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := db.QueryContext(context.Background(), query, args...)
	if err != nil {
		log.Println("getPosts query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch posts"})
	}
	defer rows.Close()

	posts := []fiber.Map{}
	for rows.Next() {
		var id, pt, title, slug, excerpt, image, cat, author, status, publishedAt, createdAt string
		var feat bool
		var views int
		if err := rows.Scan(&id, &pt, &title, &slug, &excerpt, &image, &cat, &author, &status, &feat, &publishedAt, &views, &createdAt); err != nil {
			continue
		}
		posts = append(posts, fiber.Map{
			"id": id, "post_type": pt, "title": title, "slug": slug,
			"excerpt": excerpt, "featured_image": image, "category": cat,
			"author": author, "status": status, "featured": feat,
			"published_at": publishedAt, "views": views, "created_at": createdAt,
		})
	}

	// No-cache for real-time posts listing
	c.Set("Cache-Control", "no-store, no-cache, must-revalidate")
	return c.JSON(posts)
}

// Public: get single post by slug
func getPostBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	var id, postType, title, pSlug, excerpt, content, image, category, author, status, seoTitle, seoDesc, publishedAt, createdAt, updatedAt string
	var tags []string
	var featured bool
	var views int

	err := db.QueryRowContext(context.Background(),
		`SELECT id, post_type, title, slug, COALESCE(excerpt,''), COALESCE(content,''),
		 COALESCE(featured_image,''), COALESCE(category,''), COALESCE(tags,'{}'),
		 COALESCE(author,'Admin'), status, featured,
		 COALESCE(seo_title,''), COALESCE(seo_description,''),
		 COALESCE(published_at, created_at), COALESCE(views,0), created_at, updated_at
		 FROM posts WHERE slug = $1 AND deleted_at IS NULL`,
		slug,
	).Scan(&id, &postType, &title, &pSlug, &excerpt, &content, &image, &category, pq.Array(&tags), &author, &status, &featured, &seoTitle, &seoDesc, &publishedAt, &views, &createdAt, &updatedAt)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}

	// Increment views async
	go func() {
		db.ExecContext(context.Background(), "UPDATE posts SET views = views + 1 WHERE slug = $1", slug)
	}()

	// No-cache for real-time post details
	c.Set("Cache-Control", "no-store, no-cache, must-revalidate")

	return c.JSON(fiber.Map{
		"id": id, "post_type": postType, "title": title, "slug": pSlug,
		"excerpt": excerpt, "content": content, "featured_image": image,
		"category": category, "tags": tags, "author": author,
		"status": status, "featured": featured,
		"seo_title": seoTitle, "seo_description": seoDesc,
		"published_at": publishedAt, "views": views,
		"created_at": createdAt, "updated_at": updatedAt,
	})
}

// ─── WordPress WXR XML Importer ──────────────────────────────────────────────

type WXRFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel WXRChannel `xml:"channel"`
}

type WXRChannel struct {
	Items []WXRItem `xml:"item"`
}

type WXRItem struct {
	Title      string        `xml:"title"`
	Link       string        `xml:"link"`
	Creator    string        `xml:"http://purl.org/dc/elements/1.1/ creator"`
	Content    string        `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	Excerpt    string        `xml:"http://wordpress.org/export/1.2/excerpt/ encoded"`
	PostDate   string        `xml:"http://wordpress.org/export/1.2/ post_date"`
	PostType   string        `xml:"http://wordpress.org/export/1.2/ post_type"`
	Status     string        `xml:"http://wordpress.org/export/1.2/ status"`
	PostName   string        `xml:"http://wordpress.org/export/1.2/ post_name"`
	Categories []WXRCategory `xml:"category"`
}

type WXRCategory struct {
	Domain string `xml:"domain,attr"`
	Value  string `xml:",chardata"`
}

type ImportResult struct {
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors,omitempty"`
	BatchID  string   `json:"batch_id"`
}

func importWPXMLToPosts(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "File tidak ditemukan"})
	}

	// Validate
	if file.Size > 50*1024*1024 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "File terlalu besar (max 50MB)"})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuka file"})
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca file"})
	}

	var wxr WXRFeed
	if err := xml.Unmarshal(data, &wxr); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Format XML tidak valid"})
	}

	// Import mode
	mode := c.FormValue("mode", "import") // "preview" or "import"
	postTypeMapping := c.FormValue("post_type", "article")
	if !validPostTypes[postTypeMapping] {
		postTypeMapping = "article"
	}
	downloadImages := c.FormValue("download_images", "true") != "false"
	targetCategory := c.FormValue("category", "")

	batchID := fmt.Sprintf("wp-import-%d", time.Now().Unix())
	result := ImportResult{BatchID: batchID}

	// Get user for audit
	userID, _ := c.Locals("user_id").(string)

	var previewItems []fiber.Map

	for _, item := range wxr.Channel.Items {
		if item.PostType != "post" {
			continue
		}
		if item.Status != "publish" && item.Status != "draft" {
			continue
		}

		item.Title = sanitizeTitle(item.Title)

		slug := item.PostName
		if slug == "" {
			slug = generateSlug(item.Title)
		}

		// Check duplicate
		var existingID string
		db.QueryRowContext(context.Background(), "SELECT id FROM posts WHERE slug=$1", slug).Scan(&existingID)
		if existingID != "" {
			result.Skipped++
			continue
		}

		// Get category
		categoryName := "Umum"
		var tags []string
		for _, cat := range item.Categories {
			if cat.Domain == "category" {
				categoryName = cat.Value
			} else if cat.Domain == "post_tag" {
				tags = append(tags, cat.Value)
			}
		}
		if targetCategory != "" {
			categoryName = targetCategory
		}

		status := "published"
		if item.Status == "draft" {
			status = "draft"
		}

		if mode == "preview" {
			previewItems = append(previewItems, fiber.Map{
				"title":    item.Title,
				"slug":     slug,
				"category": categoryName,
				"tags":     tags,
				"author":   item.Creator,
				"status":   status,
				"date":     item.PostDate,
			})
			result.Imported++
			continue
		}

		// Insert
		_, err = db.ExecContext(context.Background(), `
			INSERT INTO posts (post_type, title, slug, excerpt, content, featured_image, category, tags, author, status, published_at, created_by, import_batch_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW())`,
			postTypeMapping, item.Title, slug,
			func() string {
				e := sanitizeWPContent(item.Excerpt)
				if e == "" {
					return autoExcerpt(sanitizeWPContent(item.Content), 200)
				}
				return e
			}(),
			func() string {
				c := sanitizeWPContent(item.Content)
				if downloadImages {
					return rewriteContentImages(c)
				}
				return c
			}(),
			func() string {
				firstImg := extractFirstImage(item.Content)
				if downloadImages && firstImg != "" {
					return downloadImageToLocal(firstImg)
				}
				return firstImg
			}(),
			categoryName, pq.Array(tags), item.Creator, status, item.PostDate,
			userID, batchID,
		)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", item.Title, err))
			continue
		}
		result.Imported++
	}

	if mode == "preview" {
		return c.JSON(fiber.Map{
			"mode":    "preview",
			"items":   previewItems,
			"total":   len(previewItems),
			"message": fmt.Sprintf("Preview: %d item siap diimport", len(previewItems)),
		})
	}

	// Log import event
	zlog.Info().
		Str("batch_id", batchID).
		Int("imported", result.Imported).
		Int("skipped", result.Skipped).
		Int("failed", result.Failed).
		Str("user_id", userID).
		Msg("WordPress XML import completed")

	return c.JSON(fiber.Map{
		"message":  fmt.Sprintf("Import selesai: %d berhasil, %d dilewati, %d gagal", result.Imported, result.Skipped, result.Failed),
		"imported": result.Imported,
		"skipped":  result.Skipped,
		"failed":   result.Failed,
		"batch_id": batchID,
		"errors":   result.Errors,
	})
}

// Rollback import batch
func rollbackImportBatch(c *fiber.Ctx) error {
	batchID := c.Params("batchId")
	if batchID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "batch_id is required"})
	}

	res, err := db.ExecContext(context.Background(),
		"DELETE FROM posts WHERE import_batch_id = $1", batchID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to rollback"})
	}

	count, _ := res.RowsAffected()
	zlog.Info().Str("batch_id", batchID).Int64("deleted", count).Msg("Import batch rolled back")

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Rollback berhasil: %d posts dihapus", count),
		"deleted": count,
	})
}

// Soft delete post
func softDeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.ExecContext(context.Background(),
		"UPDATE posts SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete"})
	}
	return c.JSON(fiber.Map{"message": "Post moved to trash"})
}

// Restore soft-deleted post
func restorePost(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.ExecContext(context.Background(),
		"UPDATE posts SET deleted_at = NULL WHERE id = $1", id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to restore"})
	}
	return c.JSON(fiber.Map{"message": "Post restored"})
}


// ─── Migration Integrity Checker ─────────────────────────────────────────────

func checkMigrationIntegrity(c *fiber.Ctx) error {
	var postsCount, newsCount, blogCount, duplicateSlugs, emptyContent, orphanImages int

	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM posts WHERE deleted_at IS NULL").Scan(&postsCount)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM news WHERE COALESCE(active, true) = true").Scan(&newsCount)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM blog_posts").Scan(&blogCount)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM (SELECT slug FROM posts GROUP BY slug HAVING COUNT(*) > 1) x").Scan(&duplicateSlugs)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM posts WHERE (content IS NULL OR content = '') AND deleted_at IS NULL").Scan(&emptyContent)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM posts WHERE featured_image LIKE '/uploads/%' AND deleted_at IS NULL").Scan(&orphanImages)

	healthy := duplicateSlugs == 0 && postsCount >= newsCount

	return c.JSON(fiber.Map{
		"healthy":          healthy,
		"posts_count":      postsCount,
		"legacy_news":      newsCount,
		"legacy_blog":      blogCount,
		"expected_minimum": newsCount + blogCount,
		"duplicate_slugs":  duplicateSlugs,
		"empty_content":    emptyContent,
		"local_images":     orphanImages,
		"migration_status": func() string {
			if postsCount >= newsCount+blogCount {
				return "complete"
			}
			return "partial"
		}(),
	})
}

type LiveAPIPost struct {
	ID            string    `json:"id"`
	PostType      string    `json:"post_type"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Excerpt       string    `json:"excerpt"`
	Content       string    `json:"content"`
	FeaturedImage string    `json:"featured_image"`
	Category      string    `json:"category"`
	Tags          []string  `json:"tags"`
	Author        string    `json:"author"`
	Status        string    `json:"status"`
	Featured      bool      `json:"featured"`
	PublishedAt   string    `json:"published_at"`
	Views         int       `json:"views"`
	CreatedAt     string    `json:"created_at"`
}

// syncPostsFromUpertis copies all posts from the live UPERTIS website (https://upertis.ac.id)
// to the local PPID database (kampuspro_ppid).
func syncPostsFromUpertis() (int, int, error) {
	zlog.Info().Msg("PostSync: Starting background synchronization from live UPERTIS website (https://upertis.ac.id)...")

	client := &http.Client{Timeout: 30 * time.Second}
	limit := 50
	offset := 0
	addedCount := 0
	skippedCount := 0

	for {
		url := fmt.Sprintf("https://upertis.ac.id/api/v1/posts?limit=%d&offset=%d", limit, offset)
		resp, err := client.Get(url)
		if err != nil {
			zlog.Error().Err(err).Msg("PostSync: Failed to fetch post list from live UPERTIS")
			return addedCount, skippedCount, err
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			zlog.Error().Err(err).Msg("PostSync: Failed to read post list body")
			return addedCount, skippedCount, err
		}

		var postsList []LiveAPIPost
		if err := json.Unmarshal(bodyBytes, &postsList); err != nil {
			zlog.Error().Err(err).Msg("PostSync: Failed to unmarshal post list JSON")
			return addedCount, skippedCount, err
		}

		if len(postsList) == 0 {
			break // No more posts
		}

		for _, postSummary := range postsList {
			// Check if post already exists locally
			var exists int
			_ = db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM posts WHERE slug = $1", postSummary.Slug).Scan(&exists)
			if exists > 0 {
				// Already exists, skip fetching details
				skippedCount++
				continue
			}

			// Safe delay between detail requests
			time.Sleep(300 * time.Millisecond)

			// Fetch the full details (with content field) for each post
			detailURL := fmt.Sprintf("https://upertis.ac.id/api/v1/posts/%s", postSummary.Slug)
			detailResp, err := client.Get(detailURL)
			if err != nil {
				zlog.Warn().Err(err).Str("slug", postSummary.Slug).Msg("PostSync: Failed to fetch post detail")
				continue
			}

			detailBytes, err := io.ReadAll(detailResp.Body)
			detailResp.Body.Close()
			if err != nil {
				zlog.Warn().Err(err).Str("slug", postSummary.Slug).Msg("PostSync: Failed to read post detail body")
				continue
			}

			var fullPost LiveAPIPost
			if err := json.Unmarshal(detailBytes, &fullPost); err != nil {
				zlog.Warn().Err(err).Str("slug", postSummary.Slug).Msg("PostSync: Failed to unmarshal post detail JSON")
				continue
			}

			// Parse time string to sql.NullTime
			var publishedAt sql.NullTime
			if fullPost.PublishedAt != "" {
				t, err := time.Parse(time.RFC3339, fullPost.PublishedAt)
				if err == nil {
					publishedAt.Time = t
					publishedAt.Valid = true
				}
			}

			featuredImage := fullPost.FeaturedImage
			if featuredImage != "" && strings.HasPrefix(featuredImage, "/") {
				featuredImage = "https://upertis.ac.id" + featuredImage
			}

			content := fullPost.Content
			content = strings.ReplaceAll(content, `src="/uploads/`, `src="https://upertis.ac.id/uploads/`)
			content = strings.ReplaceAll(content, `src='/uploads/`, `src='https://upertis.ac.id/uploads/`)

			_, err = db.ExecContext(context.Background(), `
				INSERT INTO posts (
					id, post_type, title, slug, excerpt, content,
					featured_image, category, tags, author,
					status, featured, published_at, views,
					created_at, updated_at
				) VALUES (
					$1, $2, $3, $4, $5, $6,
					$7, $8, $9, $10,
					$11, $12, $13, $14,
					$15, NOW()
				) ON CONFLICT (slug) DO UPDATE SET
					post_type = EXCLUDED.post_type,
					title = EXCLUDED.title,
					excerpt = EXCLUDED.excerpt,
					content = EXCLUDED.content,
					featured_image = EXCLUDED.featured_image,
					category = EXCLUDED.category,
					tags = EXCLUDED.tags,
					author = EXCLUDED.author,
					status = EXCLUDED.status,
					featured = EXCLUDED.featured,
					published_at = EXCLUDED.published_at,
					views = EXCLUDED.views,
					updated_at = NOW()
			`, fullPost.ID, fullPost.PostType, fullPost.Title, fullPost.Slug, fullPost.Excerpt, content,
				featuredImage, fullPost.Category, pq.Array(fullPost.Tags), fullPost.Author,
				fullPost.Status, fullPost.Featured, publishedAt, fullPost.Views,
				fullPost.CreatedAt,
			)
			if err == nil {
				addedCount++
			} else {
				zlog.Warn().Err(err).Str("slug", fullPost.Slug).Msg("PostSync: Failed to save post locally")
			}
		}

		offset += len(postsList)
	}

	// Clean up any remaining relative image paths in database
	_, _ = db.ExecContext(context.Background(), `
		UPDATE posts 
		SET featured_image = 'https://upertis.ac.id' || featured_image 
		WHERE featured_image LIKE '/uploads/%' AND NOT featured_image LIKE 'http%'
	`)
	_, _ = db.ExecContext(context.Background(), `
		UPDATE posts 
		SET content = REPLACE(content, 'src="/uploads/', 'src="https://upertis.ac.id/uploads/')
		WHERE content LIKE '%src="/uploads/%'
	`)

	zlog.Info().Int("added", addedCount).Int("skipped", skippedCount).Msg("PostSync: Completed synchronization from UPERTIS successfully!")
	return addedCount, skippedCount, nil
}

// StartBackgroundPostSync starts a background scheduler to sync posts every 5 minutes.
func StartBackgroundPostSync() {
	// Sync once immediately on start
	go func() {
		// Wait a few seconds for DB setup to complete
		time.Sleep(5 * time.Second)
		_, _, _ = syncPostsFromUpertis()
		
		// Run every 5 minutes
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			_, _, _ = syncPostsFromUpertis()
		}
	}()
}

// handleManualLiveSync triggers post synchronization from UPERTIS live site manually.
func handleManualLiveSync(c *fiber.Ctx) error {
	added, skipped, err := syncPostsFromUpertis()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyinkronkan berita: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Sinkronisasi selesai dengan sukses!",
		"added":   added,
		"skipped": skipped,
	})
}
