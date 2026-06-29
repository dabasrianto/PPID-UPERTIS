package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type dohDNSResponse struct {
	Answer []struct {
		Data string `json:"data"`
		Type int    `json:"type"`
	} `json:"Answer"`
}

func dohResolveA(ctx context.Context, hostname string) (string, error) {
	h := strings.TrimSpace(hostname)
	if h == "" {
		return "", fmt.Errorf("empty hostname")
	}
	h = strings.TrimSuffix(h, ".")

	u := "https://cloudflare-dns.com/dns-query?name=" + url.QueryEscape(h) + "&type=A"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/dns-json")
	req.Header.Set("User-Agent", "curl/8")

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 200 {
			snippet = snippet[:200]
		}
		return "", fmt.Errorf("doh HTTP %d: %s", resp.StatusCode, snippet)
	}

	var r dohDNSResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return "", err
	}
	for _, a := range r.Answer {
		if a.Type != 1 {
			continue
		}
		ip := strings.TrimSpace(a.Data)
		if net.ParseIP(ip) != nil {
			return ip, nil
		}
	}
	return "", fmt.Errorf("no A record")
}

func dialWithDoH(ctx context.Context, dialer *net.Dialer, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return dialer.DialContext(ctx, "tcp4", addr)
	}
	if ip := net.ParseIP(host); ip != nil {
		return dialer.DialContext(ctx, "tcp4", addr)
	}

	ip, err := dohResolveA(ctx, host)
	if err != nil {
		return dialer.DialContext(ctx, "tcp4", addr)
	}
	return dialer.DialContext(ctx, "tcp4", net.JoinHostPort(ip, port))
}

// Hero Slides
func getHeroSlides(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(),
		"SELECT id, title, eyebrow, description, image_url, cta_primary_label, cta_primary_href, cta_secondary_label, cta_secondary_href, COALESCE(video_url, ''), active, sort_order, created_at, updated_at FROM hero_slides WHERE active = true ORDER BY sort_order")
	if err != nil {
		log.Println("getHeroSlides query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch hero slides"})
	}
	defer rows.Close()

	slides := []HeroSlide{}
	for rows.Next() {
		var s HeroSlide
		if err := rows.Scan(&s.ID, &s.Title, &s.Eyebrow, &s.Description, &s.ImageURL, &s.CTA1Label, &s.CTA1Href, &s.CTA2Label, &s.CTA2Href, &s.VideoURL, &s.Active, &s.SortOrder, &s.CreatedAt, &s.UpdatedAt); err != nil {
			log.Println("getHeroSlides scan error:", err)
			continue
		}
		slides = append(slides, s)
	}
	return c.JSON(slides)
}

func createHeroSlide(c *fiber.Ctx) error {
	var slide HeroSlide
	if err := c.BodyParser(&slide); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := db.QueryRowContext(context.Background(),
		"INSERT INTO hero_slides (title, eyebrow, description, image_url, cta_primary_label, cta_primary_href, cta_secondary_label, cta_secondary_href, video_url, active, sort_order) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id, created_at, updated_at",
		slide.Title, slide.Eyebrow, slide.Description, slide.ImageURL, slide.CTA1Label, slide.CTA1Href, slide.CTA2Label, slide.CTA2Href, slide.VideoURL, slide.Active, slide.SortOrder,
	).Scan(&slide.ID, &slide.CreatedAt, &slide.UpdatedAt)

	if err != nil {
		log.Println("createHeroSlide error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create hero slide"})
	}
	return c.Status(http.StatusCreated).JSON(slide)
}

func updateHeroSlide(c *fiber.Ctx) error {
	id := c.Params("id")
	var slide HeroSlide
	if err := c.BodyParser(&slide); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := db.ExecContext(context.Background(),
		"UPDATE hero_slides SET title=$1, eyebrow=$2, description=$3, image_url=$4, cta_primary_label=$5, cta_primary_href=$6, cta_secondary_label=$7, cta_secondary_href=$8, video_url=$9, active=$10, sort_order=$11, updated_at=NOW() WHERE id=$12",
		slide.Title, slide.Eyebrow, slide.Description, slide.ImageURL, slide.CTA1Label, slide.CTA1Href, slide.CTA2Label, slide.CTA2Href, slide.VideoURL, slide.Active, slide.SortOrder, id)

	if err != nil {
		log.Println("updateHeroSlide error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update hero slide"})
	}
	return c.JSON(fiber.Map{"message": "Updated successfully"})
}

func deleteHeroSlide(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.ExecContext(context.Background(), "DELETE FROM hero_slides WHERE id=$1", id)
	if err != nil {
		log.Println("deleteHeroSlide error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete hero slide"})
	}
	return c.JSON(fiber.Map{"message": "Deleted successfully"})
}

// Faculties
func getFaculties(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(),
		`SELECT id, code, name, slug, COALESCE(description, ''), COALESCE(vision, ''),
		COALESCE(mission, ''), COALESCE(about_content, ''), COALESCE(hero_title, ''),
		COALESCE(hero_eyebrow, ''), COALESCE(hero_description, ''), COALESCE(cover_image_url, ''),
		COALESCE(about_image_url, ''), COALESCE(about_images, '[]'::jsonb), COALESCE(facility_images, '[]'::jsonb), COALESCE(cover_images, '[]'::jsonb), COALESCE(accent, 'navy'), active, programs, sort_order, facilities, contact_info,
		COALESCE(tujuan, ''), COALESCE(struktur_organisasi_url, ''), kerjasama,
		created_at, updated_at FROM faculties WHERE active = true ORDER BY sort_order`)
	if err != nil {
		log.Println("getFaculties query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch faculties"})
	}
	defer rows.Close()

	faculties := []Faculty{}
	for rows.Next() {
		var f Faculty
		if err := rows.Scan(&f.ID, &f.Code, &f.Name, &f.Slug, &f.Description, &f.Vision, &f.Mission, &f.AboutContent, &f.HeroTitle, &f.HeroEyebrow, &f.HeroDescription, &f.CoverImageURL, &f.AboutImageURL, &f.AboutImages, &f.FacilityImages, &f.CoverImages, &f.Accent, &f.Active, &f.Programs, &f.SortOrder, &f.Facilities, &f.ContactInfo, &f.Tujuan, &f.StrukturOrganisasiURL, &f.Kerjasama, &f.CreatedAt, &f.UpdatedAt); err != nil {
			log.Println("getFaculties scan error:", err)
			continue
		}
		faculties = append(faculties, f)
	}
	return c.JSON(faculties)
}

func getFacultyBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	var f Faculty
	err := db.QueryRowContext(context.Background(),
		`SELECT id, code, name, slug, COALESCE(description, ''), COALESCE(vision, ''),
		COALESCE(mission, ''), COALESCE(about_content, ''), COALESCE(hero_title, ''),
		COALESCE(hero_eyebrow, ''), COALESCE(hero_description, ''), COALESCE(cover_image_url, ''),
		COALESCE(about_image_url, ''), COALESCE(about_images, '[]'::jsonb), COALESCE(facility_images, '[]'::jsonb), COALESCE(cover_images, '[]'::jsonb), COALESCE(accent, 'navy'), active, programs, sort_order, facilities, contact_info,
		COALESCE(tujuan, ''), COALESCE(struktur_organisasi_url, ''), kerjasama,
		created_at, updated_at FROM faculties WHERE slug=$1`,
		slug).Scan(&f.ID, &f.Code, &f.Name, &f.Slug, &f.Description, &f.Vision, &f.Mission, &f.AboutContent, &f.HeroTitle, &f.HeroEyebrow, &f.HeroDescription, &f.CoverImageURL, &f.AboutImageURL, &f.AboutImages, &f.FacilityImages, &f.CoverImages, &f.Accent, &f.Active, &f.Programs, &f.SortOrder, &f.Facilities, &f.ContactInfo, &f.Tujuan, &f.StrukturOrganisasiURL, &f.Kerjasama, &f.CreatedAt, &f.UpdatedAt)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Faculty not found"})
	}
	return c.JSON(f)
}

// News
func getNews(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(),
		`SELECT id, title, COALESCE(slug,''), COALESCE(excerpt, ''), COALESCE(content, ''), COALESCE(image_url, ''), category, date, featured, active, COALESCE(author_name,'Admin'), COALESCE(seo_title,''), COALESCE(seo_description,''), created_at, updated_at 
		FROM news WHERE COALESCE(active, true) = true ORDER BY date DESC`)
	if err != nil {
		log.Println("getNews query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch news"})
	}
	defer rows.Close()

	news := []News{}
	for rows.Next() {
		var n News
		if err := rows.Scan(&n.ID, &n.Title, &n.Slug, &n.Excerpt, &n.Content, &n.ImageURL, &n.Category, &n.Date, &n.Featured, &n.Active, &n.AuthorName, &n.SeoTitle, &n.SeoDescription, &n.CreatedAt, &n.UpdatedAt); err != nil {
			log.Println("getNews scan error:", err)
			continue
		}
		news = append(news, n)
	}
	return c.JSON(news)
}

func getNewsByID(c *fiber.Ctx) error {
	param := c.Params("id")
	var n News
	// Support both UUID and slug lookup
	query := "SELECT id, title, COALESCE(slug,''), COALESCE(excerpt, ''), COALESCE(content, ''), COALESCE(image_url, ''), category, date, featured, active, COALESCE(author_name,'Admin'), COALESCE(seo_title,''), COALESCE(seo_description,''), created_at, updated_at FROM news WHERE id::text=$1 OR slug=$1"
	err := db.QueryRowContext(context.Background(), query, param).Scan(
		&n.ID, &n.Title, &n.Slug, &n.Excerpt, &n.Content, &n.ImageURL, &n.Category, &n.Date, &n.Featured, &n.Active, &n.AuthorName, &n.SeoTitle, &n.SeoDescription, &n.CreatedAt, &n.UpdatedAt)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "News not found"})
	}
	return c.JSON(n)
}

// Testimonials
func getTestimonials(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(),
		"SELECT id, name, role, quote, image_url, year, COALESCE(video_url, ''), COALESCE(youtube_url, ''), COALESCE(tiktok_url, ''), COALESCE(instagram_url, ''), active, sort_order, created_at, updated_at FROM testimonials WHERE active = true ORDER BY sort_order")
	if err != nil {
		log.Println("getTestimonials query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch testimonials"})
	}
	defer rows.Close()

	testimonials := []Testimonial{}
	for rows.Next() {
		var t Testimonial
		if err := rows.Scan(&t.ID, &t.Name, &t.Role, &t.Quote, &t.ImageURL, &t.Year, &t.VideoURL, &t.YoutubeURL, &t.TiktokURL, &t.InstagramURL, &t.Active, &t.SortOrder, &t.CreatedAt, &t.UpdatedAt); err != nil {
			log.Println("getTestimonials scan error:", err)
			continue
		}
		testimonials = append(testimonials, t)
	}
	return c.JSON(testimonials)
}

// Blog Posts
func getBlogPosts(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(),
		`SELECT id, title, slug, COALESCE(excerpt,''), COALESCE(cover_image_url,''), 
		 COALESCE(author_name,''), COALESCE(category,''), status, 
		 COALESCE(published_at, created_at) as pub_date, created_at
		 FROM blog_posts 
		 WHERE status ILIKE 'publish%' 
		 ORDER BY COALESCE(published_at, created_at) DESC`)
	if err != nil {
		log.Println("getBlogPosts query error:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch blog posts"})
	}
	defer rows.Close()

	posts := []map[string]interface{}{}
	for rows.Next() {
		var id, title, slug, excerpt, cover, author, category, status, publishedAt, createdAt string
		if err := rows.Scan(&id, &title, &slug, &excerpt, &cover, &author, &category, &status, &publishedAt, &createdAt); err != nil {
			log.Println("getBlogPosts scan error:", err)
			continue
		}
		posts = append(posts, map[string]interface{}{
			"id":              id,
			"title":           title,
			"slug":            slug,
			"excerpt":         excerpt,
			"cover_image_url": cover,
			"author_name":     author,
			"category":        category,
			"status":          status,
			"published_at":    publishedAt,
			"created_at":      createdAt,
		})
	}
	return c.JSON(posts)
}

func getBlogPostBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	var (
		id, title, pSlug, excerpt, content, coverURL string
		authorName, category, status, seoTitle, seoDesc string
		publishedAt, createdAt, updatedAt string
		facultyID                          *string
		tags                               = []string{}
	)

	err := db.QueryRowContext(context.Background(),
		`SELECT id, title, slug, COALESCE(excerpt,''), COALESCE(content,''), 
		 COALESCE(cover_image_url,''), COALESCE(author_name,''), COALESCE(category,''),
		 COALESCE(tags, ARRAY[]::text[]), faculty_id::text,
		 status, COALESCE(published_at, created_at), 
		 COALESCE(seo_title,''), COALESCE(seo_description,''),
		 created_at, updated_at
		 FROM blog_posts WHERE LOWER(slug) = LOWER($1)`,
		slug,
	).Scan(&id, &title, &pSlug, &excerpt, &content, &coverURL,
		&authorName, &category, pq.Array(&tags), &facultyID,
		&status, &publishedAt, &seoTitle, &seoDesc,
		&createdAt, &updatedAt)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Blog post not found"})
	}

	fid := ""
	if facultyID != nil {
		fid = *facultyID
	}

	return c.JSON(fiber.Map{
		"id":              id,
		"title":           title,
		"slug":            pSlug,
		"excerpt":         excerpt,
		"content":         content,
		"cover_image_url": coverURL,
		"author_name":     authorName,
		"category":        category,
		"tags":            tags,
		"faculty_id":      fid,
		"status":          status,
		"published_at":    publishedAt,
		"seo_title":       seoTitle,
		"seo_description": seoDesc,
		"created_at":      createdAt,
		"updated_at":      updatedAt,
	})
}

// Pages
func getPageBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	var p Page
	err := db.QueryRowContext(context.Background(),
		`SELECT id, title, COALESCE(subtitle,''), slug, COALESCE(content,''), COALESCE(cover_image_url,''),
		published, sort_order, COALESCE(seo_title,''), COALESCE(seo_description,''),
		created_at, updated_at FROM pages WHERE slug=$1 AND published=true`,
		slug).Scan(&p.ID, &p.Title, &p.Subtitle, &p.Slug, &p.Content, &p.CoverImageURL, &p.Published, &p.SortOrder, &p.SeoTitle, &p.SeoDescription, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		log.Println("getPageBySlug error:", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Page not found"})
	}
	return c.JSON(p)
}

// Contact
func submitContact(c *fiber.Ctx) error {
	var msg ContactMessage
	if err := c.BodyParser(&msg); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := db.ExecContext(context.Background(),
		"INSERT INTO contact_messages (name, email, phone, subject, message) VALUES ($1,$2,$3,$4,$5)",
		msg.Name, msg.Email, msg.Phone, msg.Subject, msg.Message)

	if err != nil {
		log.Println("submitContact error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send message"})
	}

	// Send notification to admin WhatsApp asynchronously
	notifyAdminNewMessage(msg)

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Message sent successfully"})
}

// Helper: Fetch all replies for a list of message UUIDs
func fetchMessageReplies(messageIDs []string) (map[string][]ContactMessageReply, error) {
	repliesMap := make(map[string][]ContactMessageReply)
	if len(messageIDs) == 0 {
		return repliesMap, nil
	}

	query := `
		SELECT id, message_id, sender_type, sender_name, message, created_at
		FROM contact_message_replies
		WHERE message_id = ANY($1)
		ORDER BY created_at ASC
	`
	rows, err := db.QueryContext(context.Background(), query, pq.Array(messageIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r ContactMessageReply
		var messageID string
		var createdAt time.Time
		if err := rows.Scan(&r.ID, &messageID, &r.SenderType, &r.SenderName, &r.Message, &createdAt); err != nil {
			log.Println("fetchMessageReplies scan error:", err)
			continue
		}
		r.MessageID = messageID
		r.CreatedAt = createdAt.Format(time.RFC3339)
		repliesMap[messageID] = append(repliesMap[messageID], r)
	}

	return repliesMap, nil
}

// User: Get Personal Messages
func getUserMessages(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var email string
	err := db.QueryRowContext(context.Background(), "SELECT email FROM users WHERE id = $1", userID).Scan(&email)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	rows, err := db.QueryContext(context.Background(),
		"SELECT id, name, email, phone, subject, message, is_read, COALESCE(is_read_user, true), COALESCE(reply_message,''), COALESCE(replied_at::text,''), created_at FROM contact_messages WHERE email = $1 ORDER BY created_at DESC",
		email)
	if err != nil {
		log.Println("getUserMessages error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch messages"})
	}
	defer rows.Close()

	messages := []ContactMessage{}
	messageIDs := []string{}
	for rows.Next() {
		var m ContactMessage
		if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Phone, &m.Subject, &m.Message, &m.IsRead, &m.IsReadUser, &m.ReplyMessage, &m.RepliedAt, &m.CreatedAt); err != nil {
			log.Println("getUserMessages scan error:", err)
			continue
		}
		m.Replies = []ContactMessageReply{}
		messages = append(messages, m)
		messageIDs = append(messageIDs, m.ID)
	}

	repliesMap, err := fetchMessageReplies(messageIDs)
	if err == nil {
		for i, m := range messages {
			if reps, ok := repliesMap[m.ID]; ok {
				messages[i].Replies = reps
			}
		}
	}

	return c.JSON(messages)
}

// Admin: Get Messages
func getMessages(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(),
		"SELECT id, name, email, phone, subject, message, is_read, COALESCE(is_read_user, true), COALESCE(reply_message,''), COALESCE(replied_at::text,''), created_at FROM contact_messages ORDER BY created_at DESC")
	if err != nil {
		log.Println("getMessages query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch messages"})
	}
	defer rows.Close()

	messages := []ContactMessage{}
	messageIDs := []string{}
	for rows.Next() {
		var m ContactMessage
		if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Phone, &m.Subject, &m.Message, &m.IsRead, &m.IsReadUser, &m.ReplyMessage, &m.RepliedAt, &m.CreatedAt); err != nil {
			log.Println("getMessages scan error:", err)
			continue
		}
		m.Replies = []ContactMessageReply{}
		messages = append(messages, m)
		messageIDs = append(messageIDs, m.ID)
	}

	repliesMap, err := fetchMessageReplies(messageIDs)
	if err == nil {
		for i, m := range messages {
			if reps, ok := repliesMap[m.ID]; ok {
				messages[i].Replies = reps
			}
		}
	}

	return c.JSON(messages)
}

// Admin: Reply to Message (direct WhatsApp & database storage)
func replyMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		ReplyMessage string `json:"reply_message"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	req.ReplyMessage = strings.TrimSpace(req.ReplyMessage)
	if req.ReplyMessage == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Pesan balasan tidak boleh kosong"})
	}

	// 1. Fetch message detail to get recipient's phone/whatsapp and details
	var name, phone, email, originalMsg string
	err := db.QueryRowContext(context.Background(),
		"SELECT name, COALESCE(phone,''), email, message FROM contact_messages WHERE id = $1", id).
		Scan(&name, &phone, &email, &originalMsg)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Pesan tidak ditemukan"})
	}

	// 2. Save reply to DB and mark as read
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Println("replyMessage tx error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memulai transaksi"})
	}
	defer tx.Rollback()

	// Update parent for legacy/compatibility support and mark unread for user
	_, err = tx.ExecContext(context.Background(),
		"UPDATE contact_messages SET reply_message = $1, replied_at = NOW(), is_read = true, is_read_user = false WHERE id = $2",
		req.ReplyMessage, id)
	if err != nil {
		log.Println("replyMessage update parent error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui pesan utama"})
	}

	// Insert into contact_message_replies
	_, err = tx.ExecContext(context.Background(),
		"INSERT INTO contact_message_replies (message_id, sender_type, sender_name, message) VALUES ($1, $2, $3, $4)",
		id, "admin", "Admin", req.ReplyMessage)
	if err != nil {
		log.Println("replyMessage insert reply error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan balasan ke thread"})
	}

	if err := tx.Commit(); err != nil {
		log.Println("replyMessage commit error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan transaksi balasan"})
	}

	// 3. Send direct reply to user's WhatsApp if phone number is available
	if phone != "" && phone != "-" {
		go func() {
			replyText := fmt.Sprintf(
				"✉️ *Balasan dari Admin UPERTIS*\n\n"+
					"Halo %s,\n"+
					"Terkait pesan Anda: \"%s\"\n\n"+
					"💬 *Jawaban:* \n%s\n\n"+
					"Terima kasih.",
				name, originalMsg, req.ReplyMessage,
			)
			errSend := SendWhatsAppMessage(phone, replyText)
			if errSend != nil {
				log.Printf("[WA Reply] Failed to send WhatsApp reply to %s: %v", phone, errSend)
			} else {
				log.Printf("[WA Reply] Successfully sent WhatsApp reply to %s", phone)
			}
		}()
	}

	return c.JSON(fiber.Map{"message": "Balasan berhasil dikirim"})
}

// User: Reply to Message Thread (sends WhatsApp notification to admin & database storage)
func replyMessageUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		ReplyMessage string `json:"reply_message"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	req.ReplyMessage = strings.TrimSpace(req.ReplyMessage)
	if req.ReplyMessage == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Pesan balasan tidak boleh kosong"})
	}

	// 1. Fetch user's email and details to ensure they own the message and get their name
	var userEmail, userFullName string
	err := db.QueryRowContext(context.Background(),
		"SELECT email, full_name FROM users WHERE id = $1", userID).
		Scan(&userEmail, &userFullName)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	// 2. Fetch parent message detail to verify owner matches userEmail and check WhatsApp notification details
	var parentEmail, parentSubject, parentPhone string
	err = db.QueryRowContext(context.Background(),
		"SELECT email, COALESCE(subject,''), COALESCE(phone,'') FROM contact_messages WHERE id = $1", id).
		Scan(&parentEmail, &parentSubject, &parentPhone)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Pesan tidak ditemukan"})
	}

	if strings.ToLower(parentEmail) != strings.ToLower(userEmail) {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Anda tidak memiliki akses ke pesan ini"})
	}

	// 3. Save reply to DB and mark parent message as UNREAD (is_read = false) so admin gets notified
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Println("replyMessageUser tx error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memulai transaksi"})
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(context.Background(),
		"UPDATE contact_messages SET is_read = false, is_read_user = true WHERE id = $1", id)
	if err != nil {
		log.Println("replyMessageUser update status error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui status pesan"})
	}

	_, err = tx.ExecContext(context.Background(),
		"INSERT INTO contact_message_replies (message_id, sender_type, sender_name, message) VALUES ($1, $2, $3, $4)",
		id, "user", userFullName, req.ReplyMessage)
	if err != nil {
		log.Println("replyMessageUser insert reply error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan balasan ke thread"})
	}

	if err := tx.Commit(); err != nil {
		log.Println("replyMessageUser commit error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan transaksi"})
	}

	// 4. Send WhatsApp notification to Admin configured in settings
	go func() {
		// Get admin's WhatsApp number from settings
		var settingJSON []byte
		err := db.QueryRowContext(context.Background(), "SELECT value FROM site_settings WHERE key = 'contact'").Scan(&settingJSON)
		if err != nil {
			log.Printf("[WA Reply Notify] Failed to fetch contact settings: %v", err)
			return
		}

		var contactData map[string]interface{}
		if err := json.Unmarshal(settingJSON, &contactData); err != nil {
			log.Printf("[WA Reply Notify] Failed to parse contact settings JSON: %v", err)
			return
		}

		adminWA, _ := contactData["whatsapp"].(string)
		if adminWA == "" {
			log.Printf("[WA Reply Notify] Admin WhatsApp number is not configured in settings")
			return
		}

		subject := parentSubject
		if subject == "" {
			subject = "-"
		}

		messageText := fmt.Sprintf(
			"💬 *Balasan Pesan Baru UPERTIS*\n\n"+
				"👤 *Nama:* %s\n"+
				"📧 *Email:* %s\n"+
				"📝 *Subjek:* %s\n\n"+
				"💬 *Balasan User:* \n%s",
			userFullName,
			userEmail,
			subject,
			req.ReplyMessage,
		)

		errSend := SendWhatsAppMessage(adminWA, messageText)
		if errSend != nil {
			log.Printf("[WA Reply Notify] Failed to send notification to admin %s: %v", adminWA, errSend)
		} else {
			log.Printf("[WA Reply Notify] Successfully sent notification to admin %s", adminWA)
		}
	}()

	return c.JSON(fiber.Map{"message": "Balasan berhasil dikirim"})
}

// User: Mark all messages for current user as read
func markUserMessagesRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var email string
	err := db.QueryRowContext(context.Background(), "SELECT email FROM users WHERE id = $1", userID).Scan(&email)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	_, err = db.ExecContext(context.Background(),
		"UPDATE contact_messages SET is_read_user = true WHERE email = $1 AND is_read_user = false",
		email)
	if err != nil {
		log.Println("markUserMessagesRead error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update messages status"})
	}

	return c.JSON(fiber.Map{"message": "All messages marked as read"})
}

func markMessageRead(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		IsRead *bool `json:"is_read"`
	}
	if err := c.BodyParser(&req); err != nil {
		log.Println("markMessageRead parse error:", err)
	}
	isRead := true
	if req.IsRead != nil {
		isRead = *req.IsRead
	}
	_, err := db.ExecContext(context.Background(), "UPDATE contact_messages SET is_read=$1 WHERE id=$2", isRead, id)
	if err != nil {
		log.Println("markMessageRead error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update message"})
	}
	return c.JSON(fiber.Map{"message": "Message marked as read"})
}

func getPages(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(),
		`SELECT id, title, COALESCE(subtitle,''), slug, COALESCE(content,''), COALESCE(cover_image_url,''),
		published, sort_order, COALESCE(seo_title,''), COALESCE(seo_description,''),
		created_at, updated_at FROM pages ORDER BY sort_order, title`)
	if err != nil {
		log.Println("getPages query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch pages"})
	}
	defer rows.Close()

	pages := []Page{}
	for rows.Next() {
		var p Page
		if err := rows.Scan(&p.ID, &p.Title, &p.Subtitle, &p.Slug, &p.Content, &p.CoverImageURL, &p.Published, &p.SortOrder, &p.SeoTitle, &p.SeoDescription, &p.CreatedAt, &p.UpdatedAt); err != nil {
			log.Println("getPages scan error:", err)
			continue
		}
		pages = append(pages, p)
	}
	return c.JSON(pages)
}

func getPrograms(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(),
		`SELECT id, faculty_id, name, COALESCE(slug, ''), level, description, accreditation, duration_years, active, sort_order,
		        COALESCE(visi, ''), COALESCE(misi, ''), COALESCE(tujuan, ''), COALESCE(gelar_akademik, ''),
		        COALESCE(kompetensi_lulusan, ''), COALESCE(fasilitas_laboratorium, ''), COALESCE(fasilitas_laboratorium_image, ''),
		        COALESCE(fasilitas_laboratorium_images, '[]'::jsonb), created_at, updated_at
		 FROM faculty_programs ORDER BY sort_order, name`)
	if err != nil {
		log.Println("getPrograms query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch programs"})
	}
	defer rows.Close()

	programs := []Program{}
	for rows.Next() {
		var p Program
		if err := rows.Scan(
			&p.ID, &p.FacultyID, &p.Name, &p.Slug, &p.Level, &p.Description, &p.Accreditation, &p.DurationYears, &p.Active, &p.SortOrder,
			&p.Visi, &p.Misi, &p.Tujuan, &p.GelarAkademik, &p.KompetensiLulusan, &p.FasilitasLaboratorium, &p.FasilitasLaboratoriumImage,
			&p.FasilitasLaboratoriumImages, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			log.Println("getPrograms scan error:", err)
			continue
		}
		programs = append(programs, p)
	}
	return c.JSON(programs)
}

func getProgramByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var progID, name, level, description, accreditation, facultyID, createdAt, slug string
	var visi, misi, tujuan, gelarAkademik, kompetensiLulusan, fasilitasLaboratorium, fasilitasLaboratoriumImage, coverImageURL, cardBgColor string
	var durationYears int
	var syllabusBytes, careerPathsBytes, strukturOrganisasiBytes, fasilitasLaboratoriumImagesBytes []byte

	err := db.QueryRowContext(context.Background(),
		`SELECT fp.id, fp.name, COALESCE(fp.level,''), COALESCE(fp.description,''), COALESCE(fp.accreditation,''),
		 fp.duration_years, fp.faculty_id, fp.created_at, COALESCE(fp.slug,''),
		 COALESCE(fp.syllabus, '[]'::jsonb), COALESCE(fp.career_paths, '[]'::jsonb),
		 COALESCE(fp.visi, ''), COALESCE(fp.misi, ''), COALESCE(fp.tujuan, ''),
		 COALESCE(fp.gelar_akademik, ''), COALESCE(fp.kompetensi_lulusan, ''), COALESCE(fp.fasilitas_laboratorium, ''),
		 COALESCE(fp.struktur_organisasi, '[]'::jsonb), COALESCE(fp.fasilitas_laboratorium_image, ''),
		 COALESCE(fp.fasilitas_laboratorium_images, '[]'::jsonb), COALESCE(fp.cover_image_url, ''), COALESCE(fp.card_bg_color, '')
		 FROM faculty_programs fp WHERE (fp.id::text = $1 OR fp.slug = $2) AND fp.active = true`, id, id,
	).Scan(&progID, &name, &level, &description, &accreditation, &durationYears, &facultyID, &createdAt, &slug, &syllabusBytes, &careerPathsBytes, &visi, &misi, &tujuan, &gelarAkademik, &kompetensiLulusan, &fasilitasLaboratorium, &strukturOrganisasiBytes, &fasilitasLaboratoriumImage, &fasilitasLaboratoriumImagesBytes, &coverImageURL, &cardBgColor)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Program studi tidak ditemukan"})
	}

	// Get faculty info
	var facultyName, facultySlug string
	db.QueryRowContext(context.Background(),
		"SELECT name, slug FROM faculties WHERE id = $1", facultyID,
	).Scan(&facultyName, &facultySlug)

	// Get lecturers for this program
	// Build the exact expertise string used in Excel import: e.g. "Gizi (D3)", "Farmasi (S1)"
	// The program name in DB is "S1 Farmasi", "D3 Gizi", etc. We need to strip the level prefix.
	keyword := name
	words := strings.Fields(strings.ToLower(name))
	for _, w := range words {
		if w != "s1" && w != "d3" && w != "d4" && w != "d-iii" && w != "d-iv" && w != "profesi" && w != "dan" && w != "program" && w != "studi" && w != "analis" && w != "kesehatan" && len(w) > 3 {
			keyword = w
			break
		}
	}

	// Build exact match: strip level prefix from program name to create "BaseName (Level)" format
	// e.g. "S1 Farmasi" + level "S1" -> "Farmasi (S1)"
	// e.g. "D3 Gizi" + level "D3" -> "Gizi (D3)"
	// e.g. "Profesi Ners" + level "Profesi" -> "Pendidikan Profesi Ners (Profesi)"
	baseName := name
	if level != "" {
		prefixes := []string{level + " ", strings.ToUpper(level) + " ", strings.ToLower(level) + " "}
		for _, pfx := range prefixes {
			if strings.HasPrefix(name, pfx) || strings.HasPrefix(strings.ToLower(name), strings.ToLower(pfx)) {
				baseName = strings.TrimSpace(name[len(pfx):])
				break
			}
		}
	}
	exactExpertise := baseName + " (" + level + ")"

	lecturers := []fiber.Map{}
	rows, _ := db.QueryContext(context.Background(),
		`SELECT id, COALESCE(slug,''), name, COALESCE(gelar,''), position, expertise, COALESCE(photo_url,''), COALESCE(pddikti_id,'')
		 FROM faculty_lecturers 
		 WHERE faculty_id = $1 
		   AND active = true 
		   AND (
		     LOWER(expertise) = LOWER($2)
		     OR LOWER(expertise) = LOWER($3)
		     OR LOWER(expertise) LIKE '%' || LOWER($4) || ' (' || LOWER($5) || ')%'
		   )
		 ORDER BY sort_order, name LIMIT 30`, facultyID, exactExpertise, name, baseName, level)
	// Fallback: if no exact match found, try keyword match
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var lid, lslug, lname, lgelar, lposition, lexpertise, lphoto, lpddikti string
			if rows.Scan(&lid, &lslug, &lname, &lgelar, &lposition, &lexpertise, &lphoto, &lpddikti) == nil {
				lecturers = append(lecturers, fiber.Map{
					"id": lid, "slug": lslug, "name": lname, "gelar": lgelar, "position": lposition, 
					"expertise": lexpertise, "photo_url": lphoto, "pddikti_id": lpddikti,
				})
			}
		}
	}
	// Fallback with broader keyword search if no results found
	if len(lecturers) == 0 {
		rows2, _ := db.QueryContext(context.Background(),
			`SELECT id, COALESCE(slug,''), name, COALESCE(gelar,''), position, expertise, COALESCE(photo_url,''), COALESCE(pddikti_id,'') 
			 FROM faculty_lecturers 
			 WHERE faculty_id = $1 
			   AND active = true 
			   AND (
			     LOWER(expertise) LIKE '%' || LOWER($2) || '%' 
			     OR LOWER($3) LIKE '%' || LOWER(expertise) || '%'
			   )
			 ORDER BY sort_order, name LIMIT 30`, facultyID, keyword, name)
		if rows2 != nil {
			defer rows2.Close()
			for rows2.Next() {
				var lid, lslug, lname, lgelar, lposition, lexpertise, lphoto, lpddikti string
				if rows2.Scan(&lid, &lslug, &lname, &lgelar, &lposition, &lexpertise, &lphoto, &lpddikti) == nil {
					lecturers = append(lecturers, fiber.Map{
						"id": lid, "slug": lslug, "name": lname, "gelar": lgelar, "position": lposition, 
						"expertise": lexpertise, "photo_url": lphoto, "pddikti_id": lpddikti,
					})
				}
			}
		}
	}

	return c.JSON(fiber.Map{
		"id": progID, "name": name, "level": level, "description": description,
		"accreditation": accreditation, "duration_years": durationYears,
		"faculty_id": facultyID, "faculty_name": facultyName, "faculty_slug": facultySlug,
		"lecturers": lecturers, "created_at": createdAt, "slug": slug,
		"syllabus": json.RawMessage(syllabusBytes), "career_paths": json.RawMessage(careerPathsBytes),
		"struktur_organisasi": json.RawMessage(strukturOrganisasiBytes),
		"visi": visi, "misi": misi, "tujuan": tujuan,
		"gelar_akademik": gelarAkademik, "kompetensi_lulusan": kompetensiLulusan,
		"fasilitas_laboratorium": fasilitasLaboratorium,
		"fasilitas_laboratorium_image": fasilitasLaboratoriumImage,
		"fasilitas_laboratorium_images": json.RawMessage(fasilitasLaboratoriumImagesBytes),
		"cover_image_url": coverImageURL,
		"card_bg_color": cardBgColor,
	})
}

func getLecturers(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(),
		`SELECT id, COALESCE(faculty_id::text,''), COALESCE(slug,''), name, COALESCE(gelar,''),
		        COALESCE(position,''), COALESCE(education,''), COALESCE(expertise,''),
		        COALESCE(photo_url,''), COALESCE(leadership_group,''), active, sort_order,
		        COALESCE(source,''), COALESCE(pddikti_id,''), COALESCE(pddikti_uuid,''),
		        COALESCE(scholar_id,''), scholar_data,
		        created_at, updated_at
		 FROM faculty_lecturers ORDER BY sort_order, name`)
	if err != nil {
		log.Println("getLecturers query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch lecturers"})
	}
	defer rows.Close()

	lecturers := []Lecturer{}
	for rows.Next() {
		var l Lecturer
		var scholarRaw []byte
		if err := rows.Scan(
			&l.ID, &l.FacultyID, &l.Slug, &l.Name, &l.Gelar, &l.Position, &l.Education, &l.Expertise,
			&l.PhotoURL, &l.LeadershipGroup, &l.Active, &l.SortOrder, &l.Source, &l.PDDIKTIID, &l.PDDIKTIUUID,
			&l.ScholarID, &scholarRaw,
			&l.CreatedAt, &l.UpdatedAt,
		); err != nil {
			log.Println("getLecturers scan error:", err)
			continue
		}
		if len(scholarRaw) > 2 { // not null or {}
			raw := json.RawMessage(scholarRaw)
			l.ScholarData = &raw
		}
		lecturers = append(lecturers, l)
	}
	return c.JSON(lecturers)
}


func getLecturerByID(c *fiber.Ctx) error {
	identifier := strings.TrimSpace(c.Params("id"))
	if identifier == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid lecturer id"})
	}

	var (
		lid, facultyID, slug, name, gelar, position, education, expertise, photoURL, leadershipGroup string
		active                                                                        bool
		sortOrder                                                                     int
		createdAt, updatedAt                                                         string
		pddiktiID, pddiktiUUID                                                       string
		portfolioText                                                                 string
		facultyName, facultySlug                                                      string
		scholarIDVal                                                                  string
		scholarRaw                                                                    []byte
	)

	baseQuery := `
		SELECT l.id,
		       COALESCE(l.faculty_id::text,''),
		       COALESCE(l.slug,''),
		       l.name,
		       COALESCE(l.gelar,''),
		       COALESCE(l.position,''),
		       COALESCE(l.education,''),
		       COALESCE(l.expertise,''),
		       COALESCE(l.photo_url,''),
		       COALESCE(l.leadership_group,''),
		       l.active,
		       l.sort_order,
		       l.created_at::text,
		       l.updated_at::text,
		       COALESCE(l.pddikti_id,''),
		       COALESCE(l.pddikti_uuid,''),
		       COALESCE(l.portfolio::text,'{}'),
		       COALESCE(f.name,''),
		       COALESCE(f.slug,''),
		       COALESCE(l.scholar_id,''),
		       l.scholar_data
		FROM faculty_lecturers l
		LEFT JOIN faculties f ON f.id = l.faculty_id
		WHERE %s = $1
		LIMIT 1
	`

	tryByID := regexp.MustCompile(`^[0-9a-fA-F-]{36}$`).MatchString(identifier)
	var err error
	if tryByID {
		err = db.QueryRowContext(context.Background(), fmt.Sprintf(baseQuery, "l.id"), identifier).Scan(&lid, &facultyID, &slug, &name, &gelar, &position, &education, &expertise, &photoURL, &leadershipGroup, &active, &sortOrder, &createdAt, &updatedAt, &pddiktiID, &pddiktiUUID, &portfolioText, &facultyName, &facultySlug, &scholarIDVal, &scholarRaw)
		if err == sql.ErrNoRows {
			err = db.QueryRowContext(context.Background(), fmt.Sprintf(baseQuery, "l.slug"), identifier).Scan(&lid, &facultyID, &slug, &name, &gelar, &position, &education, &expertise, &photoURL, &leadershipGroup, &active, &sortOrder, &createdAt, &updatedAt, &pddiktiID, &pddiktiUUID, &portfolioText, &facultyName, &facultySlug, &scholarIDVal, &scholarRaw)
		}
	} else {
		err = db.QueryRowContext(context.Background(), fmt.Sprintf(baseQuery, "l.slug"), identifier).Scan(&lid, &facultyID, &slug, &name, &gelar, &position, &education, &expertise, &photoURL, &leadershipGroup, &active, &sortOrder, &createdAt, &updatedAt, &pddiktiID, &pddiktiUUID, &portfolioText, &facultyName, &facultySlug, &scholarIDVal, &scholarRaw)
		if err == sql.ErrNoRows {
			err = db.QueryRowContext(context.Background(), fmt.Sprintf(baseQuery, "l.id"), identifier).Scan(&lid, &facultyID, &slug, &name, &gelar, &position, &education, &expertise, &photoURL, &leadershipGroup, &active, &sortOrder, &createdAt, &updatedAt, &pddiktiID, &pddiktiUUID, &portfolioText, &facultyName, &facultySlug, &scholarIDVal, &scholarRaw)
		}
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Dosen tidak ditemukan"})
		}
		log.Println("getLecturerByID query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch lecturer"})
	}

	var portfolio interface{} = map[string]interface{}{}
	if strings.TrimSpace(portfolioText) != "" {
		var v interface{}
		if err := json.Unmarshal([]byte(portfolioText), &v); err == nil && v != nil {
			portfolio = v
		}
	}

	var scholarDataOut interface{}
	if len(scholarRaw) > 2 {
		var sd interface{}
		if jerr := json.Unmarshal(scholarRaw, &sd); jerr == nil {
			scholarDataOut = sd
		}
	}

	return c.JSON(fiber.Map{
		"id":               lid,
		"faculty_id":       facultyID,
		"faculty_name":     facultyName,
		"faculty_slug":     facultySlug,
		"slug":             slug,
		"name":             name,
		"gelar":            gelar,
		"position":         position,
		"education":        education,
		"expertise":        expertise,
		"photo_url":        photoURL,
		"leadership_group": leadershipGroup,
		"active":           active,
		"sort_order":       sortOrder,
		"created_at":       createdAt,
		"updated_at":       updatedAt,
		"pddikti_id":       pddiktiID,
		"pddikti_uuid":     pddiktiUUID,
		"portfolio":        portfolio,
		"scholar_id":       scholarIDVal,
		"scholar_data":     scholarDataOut,
	})
}

type seedPDDIKTIUPERTISLecturersRequest struct {
	PTURL         string `json:"pt_url"`
	Overwrite     bool   `json:"overwrite"`
	Offset        int    `json:"offset"`
	Limit         int    `json:"limit"`
	Semester      string `json:"semester"`
	Strategy      string `json:"strategy"`       // auto | ptid_homebase | ptname_search
	ProdiContains string `json:"prodi_contains"` // optional substring filter
	FullPortfolio bool   `json:"full_portfolio"`
}

type pddiktiProdi struct {
	ID   string
	Name string
}

type pddiktiLecturer struct {
	ID        string
	Name      string
	NIDN      string
	Position  string
	Education string
	ProdiName string
	TeachingCount    int
	PublicationsCount int
	PengabdianCount  int
	Portfolio map[string]interface{} `json:"-"`
}

type pddiktiDosenSearchItem struct {
	ID        string `json:"id"`
	Name      string `json:"nama"`
	NIDN      string `json:"nidn"`
	PTName    string `json:"nama_pt"`
	ProdiName string `json:"nama_prodi"`
}

type pddiktiDosenProfile struct {
	IDSDM               string `json:"id_sdm"`
	NamaDosen           string `json:"nama_dosen"`
	NamaPT              string `json:"nama_pt"`
	NamaProdi           string `json:"nama_prodi"`
	JabatanAkademik     string `json:"jabatan_akademik"`
	PendidikanTertinggi string `json:"pendidikan_tertinggi"`
}

func extractPDDIKTIPTID(ptURL string) (string, error) {
	raw := strings.TrimSpace(ptURL)
	if raw == "" {
		return "", fmt.Errorf("empty pt_url")
	}

	if strings.Contains(raw, "detail-pt/") {
		re := regexp.MustCompile(`detail-pt/([^/?#]+)`)
		m := re.FindStringSubmatch(raw)
		if len(m) == 2 {
			return m[1], nil
		}
	}

	if u, err := url.Parse(raw); err == nil {
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		for i := 0; i < len(parts)-1; i++ {
			if parts[i] == "detail-pt" && strings.TrimSpace(parts[i+1]) != "" {
				return parts[i+1], nil
			}
		}
	}

	if strings.Contains(raw, "==") && !strings.Contains(raw, "/") {
		return raw, nil
	}

	return "", fmt.Errorf("unsupported pt_url format")
}

func pddiktiCandidateSemesters(now time.Time) []string {
	years := []int{now.Year(), now.Year() - 1, now.Year() - 2, now.Year() - 3}
	var semesters []string
	for _, y := range years {
		semesters = append(semesters, fmt.Sprintf("%d2", y))
		semesters = append(semesters, fmt.Sprintf("%d1", y))
	}
	return semesters
}

func pddiktiGET(ctx context.Context, endpoint string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, err
	}

	xUserIP := getEnv("PDDIKTI_X_USER_IP", "103.47.132.29")
	timeoutSeconds := 25
	if v := strings.TrimSpace(getEnv("PDDIKTI_TIMEOUT_SECONDS", "")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 5 && n <= 120 {
			timeoutSeconds = n
		}
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("Origin", "https://pddikti.kemdiktisaintek.go.id")
	req.Header.Set("Referer", "https://pddikti.kemdiktisaintek.go.id/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("X-User-IP", xUserIP)
	req.Header.Set("sec-ch-ua", "\"Microsoft Edge\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"Windows\"")

	dialer := &net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}
	headerTimeoutSeconds := timeoutSeconds - 5
	if headerTimeoutSeconds < 1 {
		headerTimeoutSeconds = 1
	}
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
			return dialWithDoH(ctx, dialer, addr)
		},
		ForceAttemptHTTP2:     true,
		TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12},
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: time.Duration(headerTimeoutSeconds) * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second, Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body, resp.StatusCode, nil
}

func pddiktiFrontendGET(ctx context.Context, endpoint string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	dialer := &net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
			return dialWithDoH(ctx, dialer, addr)
		},
		ForceAttemptHTTP2:     true,
		TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12},
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{Timeout: 20 * time.Second, Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body, resp.StatusCode, nil
}

func pddiktiUnmarshal(body []byte) (interface{}, error) {
	var v interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func pddiktiAsSlice(v interface{}) ([]interface{}, bool) {
	switch t := v.(type) {
	case []interface{}:
		return t, true
	default:
		return nil, false
	}
}

func pddiktiAsMap(v interface{}) (map[string]interface{}, bool) {
	switch t := v.(type) {
	case map[string]interface{}:
		return t, true
	default:
		return nil, false
	}
}

func pddiktiGetString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			switch vv := v.(type) {
			case string:
				if strings.TrimSpace(vv) != "" {
					return vv
				}
			default:
				s := strings.TrimSpace(fmt.Sprint(vv))
				if s != "" && s != "<nil>" {
					return s
				}
			}
		}
	}
	return ""
}

func lecturerSlug(name string, id string) string {
	base := strings.ToLower(strings.TrimSpace(name))
	base = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(base, "-")
	base = strings.Trim(base, "-")
	if base == "" {
		base = "dosen"
	}

	suffix := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(id), "-", ""))
	if len(suffix) > 8 {
		suffix = suffix[:8]
	}
	if suffix == "" {
		return base
	}
	return base + "-" + suffix
}

func pddiktiSearchPT(c *fiber.Ctx) error {
	keyword := strings.TrimSpace(c.Query("q"))
	if keyword == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Parameter 'q' tidak boleh kosong."})
	}

	fallbackURL := "https://api-pddikti.ridwaanhall.com/search/pt?q=" + url.QueryEscape(keyword)
	{
		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		fbBody, fbStatus, fbErr := pddiktiFrontendGET(ctx, fallbackURL)
		cancel()
		if fbErr == nil && fbStatus == http.StatusOK {
			var fb struct {
				Data []map[string]interface{} `json:"data"`
			}
			if err := json.Unmarshal(fbBody, &fb); err == nil && len(fb.Data) > 0 {
				ptList := make([]map[string]interface{}, 0, len(fb.Data))
				for _, it := range fb.Data {
					id := pddiktiGetString(it, "id")
					nama := pddiktiGetString(it, "nama", "name")
					kode := pddiktiGetString(it, "kode")
					namaSingkat := pddiktiGetString(it, "nama_singkat", "namaSingkat", "singkatan")
					ptList = append(ptList, map[string]interface{}{
						"id":           id,
						"kode":         kode,
						"nama":         nama,
						"nama_singkat": namaSingkat,
					})
				}
				return c.JSON(fiber.Map{
					"status":       "success",
					"total_kampus": len(ptList),
					"data":         ptList,
					"fallback":     true,
				})
			}
		}
	}

	endpoint := "https://api-pddikti.kemdiktisaintek.go.id/pencarian/pt/" + url.PathEscape(keyword)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	body, status, err := pddiktiGET(ctx, endpoint)
	cancel()
	if err != nil {
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{
			"error":             "Gagal terhubung ke PDDIKTI: " + err.Error(),
			"upstream":          endpoint,
			"fallback_upstream": fallbackURL,
			"fallback_status":   0,
		})
	}
	if status != http.StatusOK {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 300 {
			snippet = snippet[:300]
		}
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{
			"error":              fmt.Sprintf("PDDIKTI API error: HTTP %d: %s", status, snippet),
			"upstream":           endpoint,
			"upstream_body":      snippet,
			"upstream_status":    status,
			"fallback_upstream":  fallbackURL,
			"fallback_status":    0,
		})
	}

	root, err := pddiktiUnmarshal(body)
	if err != nil {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 300 {
			snippet = snippet[:300]
		}
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Gagal parsing JSON.", "upstream_body": snippet})
	}

	ptList, ok := pddiktiAsSlice(root)
	if !ok {
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Format response PDDIKTI tidak dikenali."})
	}

	return c.JSON(fiber.Map{
		"status":       "success",
		"total_kampus": len(ptList),
		"data":         ptList,
	})
}

func fetchPDDIKTIPTName(ctx context.Context, ptID string) (string, error) {
	endpoint := fmt.Sprintf("https://api-pddikti.kemdiktisaintek.go.id/pt/detail/%s", url.PathEscape(ptID))
	body, status, err := pddiktiGET(ctx, endpoint)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 300 {
			snippet = snippet[:300]
		}
		return "", fmt.Errorf("pddikti pt/detail HTTP %d: %s", status, snippet)
	}

	root, err := pddiktiUnmarshal(body)
	if err != nil {
		return "", err
	}
	obj, ok := pddiktiAsMap(root)
	if !ok {
		return "", fmt.Errorf("unexpected pt/detail response format")
	}
	name := strings.TrimSpace(pddiktiGetString(obj, "nama_pt", "nama", "nm_pt"))
	if name == "" {
		return "", fmt.Errorf("pt name not found in response")
	}
	return name, nil
}

func fetchPDDIKTILecturersByPTName(ctx context.Context, ptName string) ([]pddiktiLecturer, map[string]int, []string, error) {
	endpoint := "https://api-pddikti.kemdiktisaintek.go.id/pencarian/dosen/" + url.PathEscape(ptName)
	body, status, err := pddiktiGET(ctx, endpoint)
	if err != nil {
		return nil, nil, nil, err
	}
	if status != http.StatusOK {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 300 {
			snippet = snippet[:300]
		}
		return nil, nil, nil, fmt.Errorf("pddikti pencarian/dosen HTTP %d: %s", status, snippet)
	}

	var items []pddiktiDosenSearchItem
	if err := json.Unmarshal(body, &items); err != nil {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 300 {
			snippet = snippet[:300]
		}
		return nil, nil, nil, fmt.Errorf("failed to parse pencarian/dosen JSON: %v (body=%s)", err, snippet)
	}

	lecturerMap := map[string]pddiktiLecturer{}
	perProdi := map[string]int{}
	fetchErrors := []string{}

	ptNameNorm := strings.ToLower(strings.TrimSpace(ptName))
	for _, it := range items {
		name := strings.TrimSpace(it.Name)
		if name == "" {
			continue
		}
		if strings.TrimSpace(it.PTName) != "" && strings.ToLower(strings.TrimSpace(it.PTName)) != ptNameNorm {
			continue
		}

		l := pddiktiLecturer{
			ID:        strings.TrimSpace(it.ID),
			Name:      name,
			NIDN:      strings.TrimSpace(it.NIDN),
			Position:  "Dosen",
			Education: "",
			ProdiName: strings.TrimSpace(it.ProdiName),
		}

		key := strings.ToLower(strings.TrimSpace(l.NIDN + "|" + l.Name))
		if key == "|" {
			continue
		}
		if _, ok := lecturerMap[key]; ok {
			continue
		}
		lecturerMap[key] = l
		if l.ProdiName != "" {
			perProdi[l.ProdiName]++
		}
	}

	lecturers := make([]pddiktiLecturer, 0, len(lecturerMap))
	for _, v := range lecturerMap {
		lecturers = append(lecturers, v)
	}
	sort.Slice(lecturers, func(i, j int) bool {
		ai := strings.ToLower(lecturers[i].ProdiName + " " + lecturers[i].Name)
		aj := strings.ToLower(lecturers[j].ProdiName + " " + lecturers[j].Name)
		return ai < aj
	})

	return lecturers, perProdi, fetchErrors, nil
}

func fetchPDDIKTILecturersByPTID(ctx context.Context, ptID string, semesterOverride string) ([]pddiktiLecturer, string, map[string]int, []string, error) {
	semester := strings.TrimSpace(semesterOverride)
	semCandidates := pddiktiCandidateSemesters(time.Now())

	if semester == "" {
		var lastErr error
		for _, sem := range semCandidates {
			reqCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
			prodis, err := fetchPDDIKTIProdiList(reqCtx, ptID, sem)
			cancel()
			if err != nil {
				lastErr = err
				continue
			}
			if len(prodis) > 0 {
				semester = sem
				break
			}
		}
		if semester == "" {
			if lastErr != nil {
				return nil, "", nil, nil, fmt.Errorf("failed to resolve semester via pt/detail: %v", lastErr)
			}
			return nil, "", nil, nil, fmt.Errorf("failed to resolve semester via pt/detail")
		}
	}

	prodis, err := fetchPDDIKTIProdiList(ctx, ptID, semester)
	if err != nil {
		return nil, "", nil, nil, err
	}
	if len(prodis) == 0 {
		return nil, "", nil, nil, fmt.Errorf("no prodi found for pt/detail semester %s", semester)
	}

	lecturerMap := map[string]pddiktiLecturer{}
	perProdi := map[string]int{}
	fetchErrors := []string{}

	var mu sync.Mutex
	sem := make(chan struct{}, 6)
	var wg sync.WaitGroup

	for _, p := range prodis {
		prodi := p
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			trySemesters := []string{semester}
			if semesterOverride == "" {
				for _, s := range semCandidates {
					if s != semester {
						trySemesters = append(trySemesters, s)
					}
					if len(trySemesters) >= 4 {
						break
					}
				}
			}

			var lects []pddiktiLecturer
			var lastErr error
			for _, semTry := range trySemesters {
				reqCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
				lects2, err := fetchPDDIKTIHomebaseLecturers(reqCtx, prodi.ID, semTry)
				cancel()
				if err != nil {
					lastErr = err
					continue
				}
				lects = lects2
				break
			}

			if lastErr != nil && len(lects) == 0 {
				mu.Lock()
				fetchErrors = append(fetchErrors, fmt.Sprintf("homebase %s: %v", prodi.Name, lastErr))
				mu.Unlock()
				return
			}

			mu.Lock()
			for _, l := range lects {
				if strings.TrimSpace(l.ProdiName) == "" {
					l.ProdiName = strings.TrimSpace(prodi.Name)
				}
				key := strings.ToLower(strings.TrimSpace(l.NIDN + "|" + l.Name))
				if key == "|" {
					key = strings.ToLower(strings.TrimSpace(l.ID + "|" + l.Name))
				}
				if key == "|" {
					continue
				}
				if _, ok := lecturerMap[key]; ok {
					continue
				}
				lecturerMap[key] = l
				if strings.TrimSpace(l.ProdiName) != "" {
					perProdi[l.ProdiName]++
				}
			}
			mu.Unlock()
		}()
	}
	wg.Wait()

	lecturers := make([]pddiktiLecturer, 0, len(lecturerMap))
	for _, v := range lecturerMap {
		lecturers = append(lecturers, v)
	}
	sort.Slice(lecturers, func(i, j int) bool {
		ai := strings.ToLower(lecturers[i].ProdiName + " " + lecturers[i].Name)
		aj := strings.ToLower(lecturers[j].ProdiName + " " + lecturers[j].Name)
		return ai < aj
	})

	return lecturers, semester, perProdi, fetchErrors, nil
}

func fetchPDDIKTIJSONBody(ctx context.Context, endpoint string) ([]byte, error) {
	body, status, err := pddiktiGET(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 300 {
			snippet = snippet[:300]
		}
		return nil, fmt.Errorf("pddikti endpoint HTTP %d: %s", status, snippet)
	}
	return body, nil
}

func pddiktiJSONArrayLenFromValue(v interface{}) int {
	if v == nil {
		return 0
	}
	if list, ok := pddiktiAsSlice(v); ok {
		return len(list)
	}
	if obj, ok := pddiktiAsMap(v); ok {
		for _, key := range []string{"data", "result", "items"} {
			if vv, ok := obj[key]; ok {
				if list, ok := pddiktiAsSlice(vv); ok {
					return len(list)
				}
			}
		}
	}
	return 0
}

func fetchPDDIKTIDosenProfile(ctx context.Context, dosenID string) (pddiktiDosenProfile, error) {
	endpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/profile/" + url.PathEscape(dosenID)
	body, status, err := pddiktiGET(ctx, endpoint)
	if err != nil {
		return pddiktiDosenProfile{}, err
	}
	if status != http.StatusOK {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 300 {
			snippet = snippet[:300]
		}
		return pddiktiDosenProfile{}, fmt.Errorf("pddikti dosen/profile HTTP %d: %s", status, snippet)
	}

	var p pddiktiDosenProfile
	if err := json.Unmarshal(body, &p); err != nil {
		return pddiktiDosenProfile{}, err
	}
	return p, nil
}

func fetchPDDIKTIJSONArrayCount(ctx context.Context, endpoint string) (int, error) {
	body, status, err := pddiktiGET(ctx, endpoint)
	if err != nil {
		return 0, err
	}
	if status != http.StatusOK {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 300 {
			snippet = snippet[:300]
		}
		return 0, fmt.Errorf("pddikti endpoint HTTP %d: %s", status, snippet)
	}

	var v interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return 0, err
	}
	if v == nil {
		return 0, nil
	}
	if list, ok := pddiktiAsSlice(v); ok {
		return len(list), nil
	}
	return 0, nil
}

func enrichPDDIKTILecturer(ctx context.Context, l pddiktiLecturer) (pddiktiLecturer, []string) {
	errs := []string{}

	portfolio := map[string]interface{}{}

	profileEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/profile/" + url.PathEscape(l.ID)
	if body, err := fetchPDDIKTIJSONBody(ctx, profileEndpoint); err != nil {
		errs = append(errs, fmt.Sprintf("%s profile: %v", l.Name, err))
	} else {
		var p pddiktiDosenProfile
		if err := json.Unmarshal(body, &p); err != nil {
			errs = append(errs, fmt.Sprintf("%s profile parse: %v", l.Name, err))
		} else {
			if strings.TrimSpace(p.JabatanAkademik) != "" {
				l.Position = strings.TrimSpace(p.JabatanAkademik)
			}
			if strings.TrimSpace(p.PendidikanTertinggi) != "" {
				l.Education = strings.TrimSpace(p.PendidikanTertinggi)
			}
			if strings.TrimSpace(p.NamaProdi) != "" && strings.TrimSpace(l.ProdiName) == "" {
				l.ProdiName = strings.TrimSpace(p.NamaProdi)
			}
		}

		var raw interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			errs = append(errs, fmt.Sprintf("%s profile raw: %v", l.Name, err))
		} else {
			if m, ok := raw.(map[string]interface{}); ok {
				if l.NIDN != "" {
					m["nidn"] = l.NIDN
				}
			}
			portfolio["profile"] = raw
		}
	}

	teachEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/teaching-history/" + url.PathEscape(l.ID)
	if body, err := fetchPDDIKTIJSONBody(ctx, teachEndpoint); err != nil {
		errs = append(errs, fmt.Sprintf("%s teaching-history: %v", l.Name, err))
	} else {
		var raw interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			errs = append(errs, fmt.Sprintf("%s teaching-history parse: %v", l.Name, err))
		} else {
			portfolio["teaching_history"] = raw
			l.TeachingCount = pddiktiJSONArrayLenFromValue(raw)
		}
	}

	studyEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/study-history/" + url.PathEscape(l.ID)
	if body, err := fetchPDDIKTIJSONBody(ctx, studyEndpoint); err != nil {
		errs = append(errs, fmt.Sprintf("%s study-history: %v", l.Name, err))
	} else {
		var raw interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			errs = append(errs, fmt.Sprintf("%s study-history parse: %v", l.Name, err))
		} else {
			portfolio["study_history"] = raw
		}
	}

	penelitianEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/penelitian/" + url.PathEscape(l.ID)
	if body, err := fetchPDDIKTIJSONBody(ctx, penelitianEndpoint); err != nil {
		errs = append(errs, fmt.Sprintf("%s portofolio/penelitian: %v", l.Name, err))
	} else {
		var raw interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			errs = append(errs, fmt.Sprintf("%s portofolio/penelitian parse: %v", l.Name, err))
		} else {
			portfolio["penelitian"] = raw
			l.PublicationsCount = pddiktiJSONArrayLenFromValue(raw)
		}
	}

	pengabdianEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/pengabdian/" + url.PathEscape(l.ID)
	if body, err := fetchPDDIKTIJSONBody(ctx, pengabdianEndpoint); err != nil {
		errs = append(errs, fmt.Sprintf("%s portofolio/pengabdian: %v", l.Name, err))
	} else {
		var raw interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			errs = append(errs, fmt.Sprintf("%s portofolio/pengabdian parse: %v", l.Name, err))
		} else {
			portfolio["pengabdian"] = raw
			l.PengabdianCount = pddiktiJSONArrayLenFromValue(raw)
		}
	}

	karyaEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/karya/" + url.PathEscape(l.ID)
	if body, err := fetchPDDIKTIJSONBody(ctx, karyaEndpoint); err != nil {
		errs = append(errs, fmt.Sprintf("%s portofolio/karya: %v", l.Name, err))
	} else {
		var raw interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			errs = append(errs, fmt.Sprintf("%s portofolio/karya parse: %v", l.Name, err))
		} else {
			portfolio["karya"] = raw
		}
	}

	patenEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/paten/" + url.PathEscape(l.ID)
	if body, err := fetchPDDIKTIJSONBody(ctx, patenEndpoint); err != nil {
		errs = append(errs, fmt.Sprintf("%s portofolio/paten: %v", l.Name, err))
	} else {
		var raw interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			errs = append(errs, fmt.Sprintf("%s portofolio/paten parse: %v", l.Name, err))
		} else {
			portfolio["paten"] = raw
		}
	}

	if len(portfolio) > 0 {
		l.Portfolio = portfolio
	}

	parts := []string{}
	if strings.TrimSpace(l.Education) != "" {
		parts = append(parts, strings.TrimSpace(l.Education))
	}
	if l.TeachingCount > 0 {
		parts = append(parts, fmt.Sprintf("Mengajar:%d", l.TeachingCount))
	}
	if l.PublicationsCount > 0 {
		parts = append(parts, fmt.Sprintf("Publikasi:%d", l.PublicationsCount))
	}
	if l.PengabdianCount > 0 {
		parts = append(parts, fmt.Sprintf("Pengabdian:%d", l.PengabdianCount))
	}
	if len(parts) > 0 {
		l.Education = strings.Join(parts, " | ")
	}

	return l, errs
}

func enrichPDDIKTILecturerLight(ctx context.Context, l pddiktiLecturer) (pddiktiLecturer, []string) {
	errs := []string{}

	profileEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/profile/" + url.PathEscape(l.ID)
	if body, err := fetchPDDIKTIJSONBody(ctx, profileEndpoint); err != nil {
		errs = append(errs, fmt.Sprintf("%s profile: %v", l.Name, err))
		return l, errs
	} else {
		var p pddiktiDosenProfile
		if err := json.Unmarshal(body, &p); err != nil {
			errs = append(errs, fmt.Sprintf("%s profile parse: %v", l.Name, err))
		} else {
			if strings.TrimSpace(p.JabatanAkademik) != "" {
				l.Position = strings.TrimSpace(p.JabatanAkademik)
			}
			if strings.TrimSpace(p.PendidikanTertinggi) != "" {
				l.Education = strings.TrimSpace(p.PendidikanTertinggi)
			}
			if strings.TrimSpace(p.NamaProdi) != "" && strings.TrimSpace(l.ProdiName) == "" {
				l.ProdiName = strings.TrimSpace(p.NamaProdi)
			}
		}
	}

	return l, errs
}

func fetchPDDIKTIProdiList(ctx context.Context, ptID string, semester string) ([]pddiktiProdi, error) {
	ptName, err := fetchPDDIKTIPTName(ctx, ptID)
	if err != nil {
		return nil, err
	}

	endpoint := "https://api-pddikti.kemdiktisaintek.go.id/pencarian/prodi/" + url.PathEscape(ptName)
	body, status, err := pddiktiGET(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 300 {
			snippet = snippet[:300]
		}
		return nil, fmt.Errorf("pddikti pencarian/prodi HTTP %d: %s", status, snippet)
	}

	root, err := pddiktiUnmarshal(body)
	if err != nil {
		return nil, err
	}

	var items []interface{}
	if list, ok := pddiktiAsSlice(root); ok {
		items = list
	} else if obj, ok := pddiktiAsMap(root); ok {
		for _, key := range []string{"data", "result", "items"} {
			if v, ok := obj[key]; ok {
				if list, ok := pddiktiAsSlice(v); ok {
					items = list
					break
				}
			}
		}
	}

	ptNameLower := strings.ToLower(ptName)
	var prodis []pddiktiProdi
	for _, it := range items {
		m, ok := pddiktiAsMap(it)
		if !ok {
			continue
		}
		id := pddiktiGetString(m, "id", "id_sms", "id_sp", "id_prodi", "id_sms_prodi")
		name := pddiktiGetString(m, "nama", "nm_lemb", "nama_prodi", "nm_prodi", "nama_program_studi", "program_studi")
		jenjang := pddiktiGetString(m, "jenjang", "nama_jenjang_didik")
		ptField := pddiktiGetString(m, "pt", "nama_pt", "nm_pt")
		ptSingkat := pddiktiGetString(m, "pt_singkat", "nm_pt_singkat")

		if id == "" || name == "" {
			continue
		}

		ptFieldLower := strings.ToLower(ptField)
		ptSingkatLower := strings.ToLower(ptSingkat)
		match := false
		if ptFieldLower == ptNameLower || ptSingkatLower == ptNameLower || 
			strings.Contains(ptFieldLower, ptNameLower) || strings.Contains(ptNameLower, ptFieldLower) {
			match = true
		} else if strings.Contains(ptNameLower, "perintis") && (strings.Contains(ptFieldLower, "perintis") || strings.Contains(ptSingkatLower, "upertis")) {
			match = true
		}

		if !match {
			continue
		}

		fullName := name
		if jenjang != "" {
			fullName = fmt.Sprintf("%s (%s)", name, jenjang)
		}

		prodis = append(prodis, pddiktiProdi{ID: id, Name: fullName})
	}

	return prodis, nil
}

func fetchPDDIKTIHomebaseLecturers(ctx context.Context, prodiID string, semester string) ([]pddiktiLecturer, error) {
	endpoint := fmt.Sprintf("https://api-pddikti.kemdiktisaintek.go.id/dosen/homebase/%s?semester=%s", url.PathEscape(prodiID), url.QueryEscape(semester))
	body, status, err := pddiktiGET(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 300 {
			snippet = snippet[:300]
		}
		return nil, fmt.Errorf("pddikti dosen/homebase HTTP %d: %s", status, snippet)
	}

	root, err := pddiktiUnmarshal(body)
	if err != nil {
		return nil, err
	}

	var items []interface{}
	if list, ok := pddiktiAsSlice(root); ok {
		items = list
	} else if obj, ok := pddiktiAsMap(root); ok {
		for _, key := range []string{"dosen", "data", "result"} {
			if v, ok := obj[key]; ok {
				if list, ok := pddiktiAsSlice(v); ok {
					items = list
					break
				}
			}
		}
	}

	var lecturers []pddiktiLecturer
	for _, it := range items {
		m, ok := pddiktiAsMap(it)
		if !ok {
			continue
		}
		name := pddiktiGetString(m, "nm_sdm", "nama", "nm_dosen")
		if name == "" {
			continue
		}
		lecturers = append(lecturers, pddiktiLecturer{
			ID:        pddiktiGetString(m, "id_sdm", "id_dosen", "id"),
			Name:      name,
			NIDN:      pddiktiGetString(m, "nidn", "nidk"),
			Position:  pddiktiGetString(m, "jabatan_fungsional", "jabatan", "jabatan_akademik", "jabatan_fungsional_akademik"),
			Education: pddiktiGetString(m, "pendidikan_terakhir", "pendidikan", "jenjang"),
		})
	}

	return lecturers, nil
}

func applyPDDIKTILecturerFilters(lecturers []pddiktiLecturer, prodiContains string, offset int, limit int) ([]pddiktiLecturer, map[string]int) {
	filtered := lecturers
	if strings.TrimSpace(prodiContains) != "" {
		q := strings.ToLower(strings.TrimSpace(prodiContains))
		tmp := make([]pddiktiLecturer, 0, len(filtered))
		for _, l := range filtered {
			if strings.Contains(strings.ToLower(l.ProdiName), q) {
				tmp = append(tmp, l)
			}
		}
		filtered = tmp
	}

	if offset > 0 && len(filtered) > 0 {
		if offset >= len(filtered) {
			filtered = []pddiktiLecturer{}
		} else {
			filtered = filtered[offset:]
		}
	}

	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}

	perProdi := map[string]int{}
	for _, l := range filtered {
		if strings.TrimSpace(l.ProdiName) != "" {
			perProdi[l.ProdiName]++
		}
	}
	return filtered, perProdi
}

func pddiktiLecturerKey(l pddiktiLecturer) string {
	name := strings.ToLower(strings.TrimSpace(l.Name))
	nidn := strings.ToLower(strings.TrimSpace(l.NIDN))
	if nidn != "" && name != "" {
		return nidn + "|" + name
	}
	id := strings.ToLower(strings.TrimSpace(l.ID))
	if id != "" && name != "" {
		return id + "|" + name
	}
	return name
}

func mergePDDIKTILecturers(primary []pddiktiLecturer, secondary []pddiktiLecturer) []pddiktiLecturer {
	m := map[string]pddiktiLecturer{}
	for _, l := range primary {
		k := pddiktiLecturerKey(l)
		if strings.TrimSpace(k) == "" {
			continue
		}
		m[k] = l
	}
	for _, l := range secondary {
		k := pddiktiLecturerKey(l)
		if strings.TrimSpace(k) == "" {
			continue
		}
		if existing, ok := m[k]; ok {
			if strings.TrimSpace(existing.ProdiName) == "" && strings.TrimSpace(l.ProdiName) != "" {
				existing.ProdiName = l.ProdiName
				m[k] = existing
			}
			continue
		}
		m[k] = l
	}
	out := make([]pddiktiLecturer, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		ai := strings.ToLower(out[i].ProdiName + " " + out[i].Name)
		aj := strings.ToLower(out[j].ProdiName + " " + out[j].Name)
		return ai < aj
	})
	return out
}

func lecturerLeadershipGroupFromPosition(position string) string {
	p := strings.ToLower(strings.TrimSpace(position))
	if p == "" {
		return ""
	}
	if strings.Contains(p, "rektor") || strings.Contains(p, "wakil rektor") {
		return "university"
	}
	if strings.Contains(p, "dekan") || strings.Contains(p, "wakil dekan") {
		return "faculty"
	}
	if strings.Contains(p, "kaprodi") || strings.Contains(p, "ketua program studi") || strings.Contains(p, "sekprodi") || strings.Contains(p, "sekretaris program studi") {
		return "program"
	}
	return ""
}

func autoDetectLeadershipFromLecturerName(name string, currentPosition string) (string, string) {
	n := strings.ToLower(name)
	
	// Rektor / Plt. Rektor
	if strings.Contains(n, "yaslina") || strings.Contains(n, "yendrizal jafri") {
		return "Rektor", "university"
	}
	
	// Dekan Fakultas Farmasi
	if strings.Contains(n, "eka fitrianda") {
		return "Dekan Fakultas Farmasi", "faculty"
	}
	
	// Kaprodi S1 Farmasi
	if strings.Contains(n, "revi yenti") {
		return "Ketua Program Studi S1 Farmasi", "program"
	}
	
	// If the currentPosition contains structural words
	p := strings.ToLower(currentPosition)
	if strings.Contains(p, "rektor") {
		return currentPosition, "university"
	}
	if strings.Contains(p, "dekan") {
		return currentPosition, "faculty"
	}
	if strings.Contains(p, "kaprodi") || strings.Contains(p, "ketua program studi") || strings.Contains(p, "sekprodi") || strings.Contains(p, "sekretaris program studi") {
		return currentPosition, "program"
	}
	
	return currentPosition, ""
}

func previewPDDIKTIUPERTISLecturers(c *fiber.Ctx) error {
	var req seedPDDIKTIUPERTISLecturersRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if strings.TrimSpace(req.PTURL) == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "pt_url is required"})
	}

	ctx := context.Background()
	ptID, err := extractPDDIKTIPTID(req.PTURL)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid pt_url format"})
	}

	ptName, err := fetchPDDIKTIPTName(ctx, ptID)
	if err != nil {
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Failed to fetch PDDIKTI PT detail: " + err.Error()})
	}

	strategyRequest := strings.TrimSpace(strings.ToLower(req.Strategy))
	if strategyRequest == "" {
		strategyRequest = "auto"
	}

	strategy := ""
	semester := strings.TrimSpace(req.Semester)
	var lecturers []pddiktiLecturer
	var perProdi map[string]int
	var fetchErrors []string
	homebaseFound := 0
	searchFound := 0

	if strategyRequest == "auto" {
		strategy = "hybrid"
		lectHB, semHB, _, errsHB, errHB := fetchPDDIKTILecturersByPTID(ctx, ptID, semester)
		if errHB != nil {
			errsHB = append(errsHB, "ptid_homebase: "+errHB.Error())
		}
		homebaseFound = len(lectHB)
		if homebaseFound > 0 {
			semester = semHB
		} else {
			semester = ""
		}

		lectSearch, _, errsSearch, errSearch := fetchPDDIKTILecturersByPTName(ctx, ptName)
		if errSearch != nil {
			errsSearch = append(errsSearch, "ptname_search: "+errSearch.Error())
		}
		searchFound = len(lectSearch)

		lecturers = mergePDDIKTILecturers(lectHB, lectSearch)
		fetchErrors = append(fetchErrors, errsHB...)
		fetchErrors = append(fetchErrors, errsSearch...)
		if len(lecturers) == 0 {
			if errHB != nil {
				return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Failed to fetch PDDIKTI lecturers: " + errHB.Error()})
			}
			if errSearch != nil {
				return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Failed to fetch PDDIKTI lecturers: " + errSearch.Error()})
			}
		}
	} else if strategyRequest == "ptname_search" {
		strategy = "ptname_search"
		semester = ""
		lecturers2, perProdi2, fetchErrors2, err2 := fetchPDDIKTILecturersByPTName(ctx, ptName)
		if err2 != nil {
			return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Failed to fetch PDDIKTI lecturers: " + err2.Error()})
		}
		lecturers, perProdi, fetchErrors = lecturers2, perProdi2, fetchErrors2
		searchFound = len(lecturers)
	} else {
		strategy = "ptid_homebase"
		lecturers2, semester2, perProdi2, fetchErrors2, err2 := fetchPDDIKTILecturersByPTID(ctx, ptID, semester)
		if err2 != nil || len(lecturers2) == 0 {
			if strategyRequest == "ptid_homebase" {
				if err2 != nil {
					return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Failed to fetch PDDIKTI lecturers (homebase): " + err2.Error()})
				}
				return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "No lecturers found via homebase strategy"})
			}
			strategy = "ptname_search"
			semester = ""
			lecturers3, perProdi3, fetchErrors3, err3 := fetchPDDIKTILecturersByPTName(ctx, ptName)
			if err3 != nil {
				if err2 != nil {
					return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Failed to fetch PDDIKTI lecturers: " + err3.Error() + " (homebase err: " + err2.Error() + ")"})
				}
				return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Failed to fetch PDDIKTI lecturers: " + err3.Error()})
			}
			lecturers, perProdi, fetchErrors = lecturers3, perProdi3, fetchErrors3
			searchFound = len(lecturers)
		} else {
			lecturers, perProdi, fetchErrors = lecturers2, perProdi2, fetchErrors2
			semester = semester2
			homebaseFound = len(lecturers)
		}
	}

	totalFound := len(lecturers)
	lecturers, perProdi = applyPDDIKTILecturerFilters(lecturers, req.ProdiContains, req.Offset, req.Limit)
	prodiSet := map[string]struct{}{}
	for k := range perProdi {
		prodiSet[k] = struct{}{}
	}

	return c.JSON(fiber.Map{
		"pt_url":          req.PTURL,
		"pt_id":           ptID,
		"pt_name":         ptName,
		"semester":        semester,
		"strategy":        strategy,
		"strategy_request": strategyRequest,
		"offset":          req.Offset,
		"limit":           req.Limit,
		"prodi_contains":  req.ProdiContains,
		"lecturers_found_homebase": homebaseFound,
		"lecturers_found_search":   searchFound,
		"prodi_count":     len(prodiSet),
		"lecturers_found_total": totalFound,
		"lecturers_found":       len(lecturers),
		"per_prodi":       perProdi,
		"fetch_errors":    fetchErrors,
		"lecturers":       lecturers,
	})
}

func fetchUPERTISCoreLeaders(ctx context.Context) []pddiktiLecturer {
	names := []string{"Yaslina", "Yendrizal Jafri", "Eka Fitrianda", "Revi Yenti"}
	var leaders []pddiktiLecturer
	
	client := &http.Client{Timeout: 15 * time.Second}
	
	for _, name := range names {
		endpoint := "https://api-pddikti.kemdiktisaintek.go.id/pencarian/dosen/" + url.PathEscape(name)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			continue
		}
		
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9,id;q=0.8")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Origin", "https://pddikti.kemdiktisaintek.go.id")
		req.Header.Set("Referer", "https://pddikti.kemdiktisaintek.go.id/")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-site")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
		req.Header.Set("X-User-IP", "103.47.132.29")
		
		resp, err := client.Do(req)
		if err != nil {
			zlog.Warn().Err(err).Str("name", name).Msg("fetchUPERTISCoreLeaders: Failed to fetch core leader search page")
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		
		var items []pddiktiDosenSearchItem
		if err := json.Unmarshal(body, &items); err != nil {
			zlog.Warn().Err(err).Str("name", name).Str("body", string(body)).Msg("fetchUPERTISCoreLeaders: Failed to unmarshal search result")
			continue
		}
		
		for _, it := range items {
			pt := strings.ToLower(it.PTName)
			if strings.Contains(pt, "perintis") {
				l := pddiktiLecturer{
					ID:        strings.TrimSpace(it.ID),
					Name:      strings.TrimSpace(it.Name),
					NIDN:      strings.TrimSpace(it.NIDN),
					Position:  "Dosen",
					Education: "",
					ProdiName: strings.TrimSpace(it.ProdiName),
				}
				leaders = append(leaders, l)
			}
		}
	}
	return leaders
}

func seedPDDIKTIUPERTISLecturers(c *fiber.Ctx) error {
	var req seedPDDIKTIUPERTISLecturersRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if strings.TrimSpace(req.PTURL) == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "pt_url is required"})
	}

	ctx := context.Background()

	ptID, err := extractPDDIKTIPTID(req.PTURL)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid pt_url format"})
	}

	farmasiFacultyID := ""
	if err := db.QueryRowContext(ctx, "SELECT id FROM faculties WHERE active = true AND LOWER(name) LIKE '%farmasi%' ORDER BY sort_order LIMIT 1").Scan(&farmasiFacultyID); err != nil {
		_ = db.QueryRowContext(ctx, "SELECT id FROM faculties WHERE active = true ORDER BY sort_order LIMIT 1").Scan(&farmasiFacultyID)
	}
	if farmasiFacultyID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "No active faculty found to attach lecturers"})
	}

	bisnisFacultyID := ""
	if err := db.QueryRowContext(ctx, "SELECT id FROM faculties WHERE active = true AND (LOWER(name) LIKE '%bisnis%' OR LOWER(name) LIKE '%ekonomi%' OR LOWER(name) LIKE '%manajemen%') ORDER BY sort_order LIMIT 1").Scan(&bisnisFacultyID); err != nil {
		bisnisFacultyID = farmasiFacultyID
	}
	if bisnisFacultyID == "" {
		bisnisFacultyID = farmasiFacultyID
	}

	kesehatanFacultyID := ""
	if err := db.QueryRowContext(ctx, "SELECT id FROM faculties WHERE active = true AND (LOWER(name) LIKE '%kesehatan%' OR LOWER(name) LIKE '%keperawatan%' OR LOWER(name) LIKE '%fikes%') ORDER BY sort_order LIMIT 1").Scan(&kesehatanFacultyID); err != nil {
		kesehatanFacultyID = farmasiFacultyID
	}
	if kesehatanFacultyID == "" {
		kesehatanFacultyID = farmasiFacultyID
	}

	ptName, err := fetchPDDIKTIPTName(ctx, ptID)
	if err != nil {
		zlog.Error().Err(err).Str("pt_id", ptID).Msg("seedPDDIKTIUPERTISLecturers: Failed to fetch PDDIKTI PT detail")
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Failed to fetch PDDIKTI PT detail: " + err.Error()})
	}

	strategyRequest := strings.TrimSpace(strings.ToLower(req.Strategy))
	if strategyRequest == "" {
		strategyRequest = "auto"
	}

	strategy := ""
	semester := strings.TrimSpace(req.Semester)
	var lecturers []pddiktiLecturer
	var perProdi map[string]int
	var fetchErrors []string
	homebaseFound := 0
	searchFound := 0

	if strategyRequest == "auto" {
		strategy = "hybrid"
		lectHB, semHB, _, errsHB, errHB := fetchPDDIKTILecturersByPTID(ctx, ptID, semester)
		if errHB != nil {
			errsHB = append(errsHB, "ptid_homebase: "+errHB.Error())
		}
		homebaseFound = len(lectHB)
		if homebaseFound > 0 {
			semester = semHB
		} else {
			semester = ""
		}

		lectSearch, _, errsSearch, errSearch := fetchPDDIKTILecturersByPTName(ctx, ptName)
		if errSearch != nil {
			errsSearch = append(errsSearch, "ptname_search: "+errSearch.Error())
		}
		searchFound = len(lectSearch)

		lecturers = mergePDDIKTILecturers(lectHB, lectSearch)
		fetchErrors = append(fetchErrors, errsHB...)
		fetchErrors = append(fetchErrors, errsSearch...)
	} else if strategyRequest == "ptname_search" {
		strategy = "ptname_search"
		semester = ""
		lecturers2, perProdi2, fetchErrors2, err2 := fetchPDDIKTILecturersByPTName(ctx, ptName)
		if err2 != nil {
			return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Failed to fetch PDDIKTI lecturers: " + err2.Error()})
		}
		lecturers, perProdi, fetchErrors = lecturers2, perProdi2, fetchErrors2
		searchFound = len(lecturers)
	} else {
		strategy = "ptid_homebase"
		lecturers2, semester2, perProdi2, fetchErrors2, err2 := fetchPDDIKTILecturersByPTID(ctx, ptID, semester)
		if err2 != nil || len(lecturers2) == 0 {
			if strategyRequest == "ptid_homebase" {
				if err2 != nil {
					return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Failed to fetch PDDIKTI lecturers (homebase): " + err2.Error()})
				}
				return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "No lecturers found via homebase strategy"})
			}
			strategy = "ptname_search"
			semester = ""
			lecturers3, perProdi3, fetchErrors3, err3 := fetchPDDIKTILecturersByPTName(ctx, ptName)
			if err3 != nil {
				if err2 != nil {
					return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Failed to fetch PDDIKTI lecturers: " + err3.Error() + " (homebase err: " + err2.Error() + ")"})
				}
				return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Failed to fetch PDDIKTI lecturers: " + err3.Error()})
			}
			lecturers, perProdi, fetchErrors = lecturers3, perProdi3, fetchErrors3
			searchFound = len(lecturers)
		} else {
			lecturers, perProdi, fetchErrors = lecturers2, perProdi2, fetchErrors2
			semester = semester2
			homebaseFound = len(lecturers)
		}
	}

	totalFound := len(lecturers)
	lecturers, perProdi = applyPDDIKTILecturerFilters(lecturers, req.ProdiContains, req.Offset, req.Limit)
	
	// Inject UPERTIS Core Leaders so they are ALWAYS imported
	coreLeaders := fetchUPERTISCoreLeaders(ctx)
	for _, cl := range coreLeaders {
		exists := false
		for _, l := range lecturers {
			if l.ID == cl.ID || (cl.NIDN != "" && l.NIDN == cl.NIDN) {
				exists = true
				break
			}
		}
		if !exists {
			lecturers = append(lecturers, cl)
			// Add count to perProdi mapping if it has a prodi name
			if cl.ProdiName != "" {
				perProdi[cl.ProdiName]++
			}
		}
	}

	if len(lecturers) == 0 {
		zlog.Warn().Str("pt_name", ptName).Interface("fetch_errors", fetchErrors).Msg("seedPDDIKTIUPERTISLecturers: No lecturers found after filters/strategy execution")
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "No lecturers found from PDDIKTI after filter", "pt_name": ptName})
	}

	enriched := make([]pddiktiLecturer, len(lecturers))
	copy(enriched, lecturers)
	var mu sync.Mutex
	sem := make(chan struct{}, 8)
	var wg sync.WaitGroup
	for i := range enriched {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			reqCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
			defer cancel()

			var updated pddiktiLecturer
			var errs []string
			if req.FullPortfolio {
				updated, errs = enrichPDDIKTILecturer(reqCtx, enriched[idx])
			} else {
				updated, errs = enrichPDDIKTILecturerLight(reqCtx, enriched[idx])
			}
			enriched[idx] = updated
			if len(errs) > 0 {
				mu.Lock()
				fetchErrors = append(fetchErrors, errs...)
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()
	lecturers = enriched

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
	}
	defer tx.Rollback()

	if req.Overwrite {
		if _, err := tx.ExecContext(ctx, "DELETE FROM faculty_lecturers WHERE source = 'pddikti' OR leadership_group = 'pddikti'"); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to clear existing PDDIKTI lecturers"})
		}
	}

	insertSQL := `
		INSERT INTO faculty_lecturers (faculty_id, slug, name, position, education, expertise, photo_url, leadership_group, source, active, sort_order, pddikti_id, pddikti_uuid, portfolio)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,true,$10,NULLIF($11, ''),NULLIF($12, ''),$13::jsonb)
		ON CONFLICT (pddikti_uuid) DO UPDATE SET
			faculty_id = EXCLUDED.faculty_id,
			name = EXCLUDED.name,
			position = CASE
				WHEN LOWER(faculty_lecturers.position) LIKE '%rektor%'
				  OR LOWER(faculty_lecturers.position) LIKE '%dekan%'
				  OR LOWER(faculty_lecturers.position) LIKE '%kaprodi%'
				  OR LOWER(faculty_lecturers.position) LIKE '%prodi%'
				  OR LOWER(faculty_lecturers.position) LIKE '%yayasan%'
				  OR LOWER(faculty_lecturers.position) LIKE '%founder%'
				  THEN faculty_lecturers.position
				ELSE EXCLUDED.position
			END,
			education = EXCLUDED.education,
			expertise = EXCLUDED.expertise,
			photo_url = CASE
				WHEN COALESCE(faculty_lecturers.photo_url, '') != '' THEN faculty_lecturers.photo_url
				ELSE EXCLUDED.photo_url
			END,
			source = EXCLUDED.source,
			active = true,
			sort_order = EXCLUDED.sort_order,
			updated_at = NOW(),
			pddikti_id = EXCLUDED.pddikti_id,
			portfolio = EXCLUDED.portfolio,
			leadership_group = CASE
				WHEN COALESCE(faculty_lecturers.leadership_group,'') = '' THEN EXCLUDED.leadership_group
				ELSE faculty_lecturers.leadership_group
			END
		RETURNING id
	`
	updateSlugSQL := "UPDATE faculty_lecturers SET slug = $1 WHERE id = $2 AND COALESCE(slug,'') = ''"

	insertedFarmasi := 0
	insertedBisnis := 0
	insertedKesehatan := 0
	sortOrder := 1

	for _, l := range lecturers {
		prodiLower := strings.ToLower(strings.TrimSpace(l.ProdiName))
		facultyID := farmasiFacultyID
		if strings.Contains(prodiLower, "bisnis") || strings.Contains(prodiLower, "ekonomi") || strings.Contains(prodiLower, "manajemen") || strings.Contains(prodiLower, "digital") {
			facultyID = bisnisFacultyID
		} else if strings.Contains(prodiLower, "farmasi") {
			facultyID = farmasiFacultyID
		} else if strings.Contains(prodiLower, "keperawatan") || strings.Contains(prodiLower, "kebidanan") || strings.Contains(prodiLower, "gizi") || strings.Contains(prodiLower, "kesehatan") || strings.Contains(prodiLower, "analis") || strings.Contains(prodiLower, "medis") || strings.Contains(prodiLower, "fisioterapi") || strings.Contains(prodiLower, "ners") || strings.Contains(prodiLower, "klinik") {
			facultyID = kesehatanFacultyID
		}

		// Auto-sync the program (prodi) into faculty_programs if it doesn't exist
		if l.ProdiName != "" {
			prodiNameTrimmed := strings.TrimSpace(l.ProdiName)
			var programExists bool
			err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM faculty_programs WHERE LOWER(name) = LOWER($1))", prodiNameTrimmed).Scan(&programExists)
			if err == nil && !programExists {
				// Detect level
				level := "S1"
				pLower := strings.ToLower(prodiNameTrimmed)
				if strings.Contains(pLower, "d3") || strings.Contains(pLower, "diii") {
					level = "D3"
				} else if strings.Contains(pLower, "d4") || strings.Contains(pLower, "div") {
					level = "D4"
				} else if strings.Contains(pLower, "profesi") || strings.Contains(pLower, "ners") || strings.Contains(pLower, "apoteker") {
					level = "Profesi"
				}

				// Detect duration
				duration := 4
				if level == "D3" {
					duration = 3
				} else if level == "Profesi" {
					duration = 1
				}

				// Slugify
				pSlug := strings.ToLower(prodiNameTrimmed)
				pSlug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(pSlug, "-")
				pSlug = strings.Trim(pSlug, "-")
				if pSlug == "" {
					pSlug = "prodi-" + strings.ToLower(level)
				}

				// Insert program
				_, _ = tx.ExecContext(ctx, `
					INSERT INTO faculty_programs (faculty_id, name, slug, level, description, accreditation, duration_years, active, sort_order)
					VALUES ($1, $2, $3, $4, $5, 'Baik', $6, true, 99)
					ON CONFLICT (slug) DO NOTHING
				`, facultyID, prodiNameTrimmed, pSlug, level, "Program Studi "+prodiNameTrimmed+" Universitas Perintis Indonesia.", duration)
			}
		}

		position := strings.TrimSpace(l.Position)
		if position == "" {
			position = "Dosen"
		}
		var leadershipGroup string
		position, leadershipGroup = autoDetectLeadershipFromLecturerName(l.Name, position)

		education := strings.TrimSpace(l.Education)
		expertise := strings.TrimSpace(l.ProdiName)
		if expertise == "" {
			expertise = "PDDIKTI"
		}

		name := strings.TrimSpace(l.Name)
		if name == "" {
			continue
		}

	portfolioJSON := "{}"
	if req.FullPortfolio && l.Portfolio != nil {
		if b, err := json.Marshal(l.Portfolio); err == nil {
			portfolioJSON = string(b)
		} else {
			fetchErrors = append(fetchErrors, fmt.Sprintf("%s portfolio marshal: %v", name, err))
		}
	}

		pddiktiID := strings.TrimSpace(l.NIDN)
		pddiktiUUID := strings.TrimSpace(l.ID)
		if pddiktiUUID == "" {
			continue
		}
		var insertedID string
		if err := tx.QueryRowContext(ctx, insertSQL, facultyID, nil, name, position, education, expertise, "", leadershipGroup, "pddikti", sortOrder, pddiktiID, pddiktiUUID, portfolioJSON).Scan(&insertedID); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert PDDIKTI lecturers"})
		}
		slug := lecturerSlug(name, insertedID)
		if _, err := tx.ExecContext(ctx, updateSlugSQL, slug, insertedID); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to set lecturer slug"})
		}
		if facultyID == bisnisFacultyID {
			insertedBisnis++
		} else if facultyID == kesehatanFacultyID {
			insertedKesehatan++
		} else {
			insertedFarmasi++
		}
		sortOrder++
	}

	if err := tx.Commit(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	prodiSet := map[string]struct{}{}
	for k := range perProdi {
		prodiSet[k] = struct{}{}
	}

	return c.JSON(fiber.Map{
		"pt_url":                 req.PTURL,
		"pt_id":                  ptID,
		"pt_name":                ptName,
		"semester":               semester,
		"strategy":               strategy,
		"strategy_request":       strategyRequest,
		"offset":                 req.Offset,
		"limit":                  req.Limit,
		"prodi_contains":         req.ProdiContains,
		"lecturers_found_homebase": homebaseFound,
		"lecturers_found_search":   searchFound,
		"lecturers_found_total":  totalFound,
		"prodi_count":            len(prodiSet),
		"lecturers_found":        len(lecturers),
		"per_prodi":              perProdi,
		"fetch_errors":           fetchErrors,
		"faculty_farmasi_id":     farmasiFacultyID,
		"faculty_bisnis_id":      bisnisFacultyID,
		"faculty_kesehatan_id":   kesehatanFacultyID,
		"inserted_farmasi":       insertedFarmasi,
		"inserted_bisnis_digital": insertedBisnis,
		"inserted_kesehatan":     insertedKesehatan,
		"total_inserted":         insertedFarmasi + insertedBisnis + insertedKesehatan,
	})
}

// Settings
func getSettings(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(), "SELECT key, value, updated_at FROM site_settings WHERE key NOT IN ('chatbot', 'whatsapp', 'smtp', 'credentials')")
	if err != nil {
		log.Println("getSettings query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch settings"})
	}
	defer rows.Close()

	settings := []SiteSetting{}
	for rows.Next() {
		var s SiteSetting
		if err := rows.Scan(&s.Key, &s.Value, &s.UpdatedAt); err != nil {
			log.Println("getSettings scan error:", err)
			continue
		}
		settings = append(settings, s)
	}
	return c.JSON(settings)
}

func getFileSizeAndType(fileURL string) (string, string) {
	ext := strings.ToLower(filepath.Ext(fileURL))
	fileType := "UNKNOWN"
	if ext != "" {
		fileType = strings.ToUpper(strings.TrimPrefix(ext, "."))
	}

	fileSize := "Unknown Size"
	
	// Normalize path by extracting it if it's a full URL containing "/uploads/"
	path := fileURL
	if idx := strings.Index(fileURL, "/uploads/"); idx != -1 {
		path = fileURL[idx:]
	}

	if strings.HasPrefix(path, "/uploads/") {
		var info os.FileInfo
		var err error = fmt.Errorf("file not found")
		
		// Try multiple possible paths to locate the file
		localPaths := []string{
			filepath.Join("..", "public", path), // running from backend/, files in public/uploads/
			filepath.Join(".", "public", path),  // running from root/, files in public/uploads/
			filepath.Join("..", path),           // running from backend/, files in uploads/
			filepath.Join(".", path),            // running from root/, files in uploads/
		}
		
		for _, lp := range localPaths {
			if stat, statErr := os.Stat(lp); statErr == nil {
				info = stat
				err = nil
				break
			}
		}

		if err == nil && info != nil {
			sz := info.Size()
			if sz < 1024 {
				fileSize = fmt.Sprintf("%d B", sz)
			} else if sz < 1024*1024 {
				fileSize = fmt.Sprintf("%.1f KB", float64(sz)/1024)
			} else {
				fileSize = fmt.Sprintf("%.1f MB", float64(sz)/(1024*1024))
			}
		}
	}
	return fileType, fileSize
}

func getDownloads(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(), `
		SELECT id, title, description, file_url, category, active, COALESCE(downloads_count, 0), created_at, updated_at
		FROM downloads
		WHERE active = true
		ORDER BY created_at DESC
	`)
	if err != nil {
		log.Println("getDownloads query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch downloads"})
	}
	defer rows.Close()

	downloads := []fiber.Map{}
	for rows.Next() {
		var id, title, description, fileURL, category string
		var active bool
		var downloadsCount int
		var createdAt, updatedAt time.Time
		err := rows.Scan(&id, &title, &description, &fileURL, &category, &active, &downloadsCount, &createdAt, &updatedAt)
		if err != nil {
			log.Println("getDownloads scan error:", err)
			continue
		}

		fileType, fileSize := getFileSizeAndType(fileURL)

		downloads = append(downloads, fiber.Map{
			"id":              id,
			"title":           title,
			"description":     description,
			"file_url":        fileURL,
			"category":        category,
			"active":          active,
			"downloads_count": downloadsCount,
			"file_type":       fileType,
			"file_size":       fileSize,
			"created_at":      createdAt.Format(time.RFC3339),
			"updated_at":      updatedAt.Format(time.RFC3339),
		})
	}

	c.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	return c.JSON(downloads)
}

func incrementDownloadCount(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}

	result, err := db.ExecContext(context.Background(), `
		UPDATE downloads
		SET downloads_count = COALESCE(downloads_count, 0) + 1, updated_at = NOW()
		WHERE id::text = $1
	`, id)
	if err != nil {
		log.Println("Error incrementing download count:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to increment download count"})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
	} else if rowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Download not found"})
	}

	return c.JSON(fiber.Map{"success": true})
}

func getToolUsages(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(), "SELECT slug, count FROM tool_usages")
	if err != nil {
		log.Println("Error fetching tool usages:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tool usages"})
	}
	defer rows.Close()

	usages := make(map[string]int)
	for rows.Next() {
		var slug string
		var count int
		if err := rows.Scan(&slug, &count); err == nil {
			usages[slug] = count
		}
	}
	c.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	return c.JSON(usages)
}

func incrementToolUsage(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Slug is required"})
	}

	_, err := db.ExecContext(context.Background(), `
		INSERT INTO tool_usages (slug, count, updated_at)
		VALUES ($1, 1, NOW())
		ON CONFLICT (slug)
		DO UPDATE SET count = tool_usages.count + 1, updated_at = NOW()
	`, slug)
	if err != nil {
		log.Println("Error incrementing tool usage count:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to increment tool usage count"})
	}

	return c.JSON(fiber.Map{"success": true})
}

// getChatbotPublicConfig returns only safe chatbot fields (no API key)
func getChatbotPublicConfig(c *fiber.Ctx) error {
	var settingJSON []byte
	err := db.QueryRowContext(context.Background(), "SELECT value FROM site_settings WHERE key = 'chatbot'").Scan(&settingJSON)
	if err != nil {
		return c.JSON(fiber.Map{"enabled": false})
	}

	var full map[string]interface{}
	if err := json.Unmarshal(settingJSON, &full); err != nil {
		return c.JSON(fiber.Map{"enabled": false})
	}

	// Only expose safe public fields — never expose api_key
	public := fiber.Map{
		"enabled":         full["enabled"],
		"greeting":        full["greeting"],
		"bot_avatar":      full["bot_avatar"],
		"bot_avatar_size": full["bot_avatar_size"],
		"bot_name":        full["bot_name"],
	}

	c.Set("Cache-Control", "public, max-age=300")
	return c.JSON(public)
}

func getAdminSettings(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(), "SELECT key, value, updated_at FROM site_settings")
	if err != nil {
		log.Println("getAdminSettings query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch settings"})
	}
	defer rows.Close()

	settings := []SiteSetting{}
	for rows.Next() {
		var s SiteSetting
		if err := rows.Scan(&s.Key, &s.Value, &s.UpdatedAt); err != nil {
			log.Println("getAdminSettings scan error:", err)
			continue
		}
		settings = append(settings, s)
	}
	return c.JSON(settings)
}

func getDashboardStats(c *fiber.Ctx) error {
	var totalPages, totalPosts, totalDownloads, totalPermohonan, pendingPermohonan int64

	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM pages").Scan(&totalPages)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM posts WHERE deleted_at IS NULL").Scan(&totalPosts)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM downloads").Scan(&totalDownloads)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM permohonan_informasi").Scan(&totalPermohonan)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM permohonan_informasi WHERE LOWER(status) = 'pending'").Scan(&pendingPermohonan)

	return c.JSON(fiber.Map{
		"total_pages":        totalPages,
		"total_posts":        totalPosts,
		"total_downloads":    totalDownloads,
		"total_permohonan":    totalPermohonan,
		"pending_permohonan":  pendingPermohonan,
	})
}

// getNotificationBadges returns unread/pending counts for admin badge notifications
func getNotificationBadges(c *fiber.Ctx) error {
	var pendingComments, unreadMessages, pendingRegistrations, pendingPayments, pendingUsers int64

	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM comments WHERE status='pending'").Scan(&pendingComments)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM contact_messages WHERE is_read=false").Scan(&unreadMessages)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM pmb_candidates WHERE status='DRAFT' OR status='WAITING_PAYMENT'").Scan(&pendingRegistrations)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM pmb_payments WHERE status='PENDING'").Scan(&pendingPayments)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM users WHERE status='pending'").Scan(&pendingUsers)

	total := pendingComments + unreadMessages + pendingRegistrations + pendingPayments + pendingUsers

	c.Set("Cache-Control", "no-cache")
	return c.JSON(fiber.Map{
		"total":                 total,
		"pending_comments":      pendingComments,
		"unread_messages":       unreadMessages,
		"pending_registrations": pendingRegistrations,
		"pending_payments":      pendingPayments,
		"pending_users":         pendingUsers,
	})
}

// Campus Events (public)
func getEvents(c *fiber.Ctx) error {
	siteSlug := getSiteID(c)

	query := `SELECT id, title, slug, COALESCE(description,''), COALESCE(image_url,''), event_date,
		 COALESCE(start_time::text,'00:00'), COALESCE(end_time::text,'23:59'),
		 COALESCE(location,''), COALESCE(map_coordinates,''), active, sort_order,
		 COALESCE(category,''), COALESCE(registration_url,''), is_internal_registration, capacity, COALESCE(speakers::text,'[]'),
		 COALESCE(images::text,'[]'), created_at, updated_at
		 FROM campus_events WHERE active = true`
	args := []interface{}{}
	argIdx := 1

	// Multi-site: filter by faculty
	query, args, argIdx = buildFacultyFilter(query, args, argIdx, siteSlug)
	_ = argIdx

	query += " ORDER BY event_date ASC, sort_order"

	rows, err := db.QueryContext(context.Background(), query, args...)
	if err != nil {
		log.Println("getEvents query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch events"})
	}
	defer rows.Close()

	events := []map[string]interface{}{}
	for rows.Next() {
		var id, title, slug, desc, imageURL, eventDate, startTime, endTime, location, mapCoords, category, registrationURL, speakers, imagesStr, createdAt, updatedAt string
		var active, isInternal bool
		var sortOrder, capacity int
		if err := rows.Scan(&id, &title, &slug, &desc, &imageURL, &eventDate, &startTime, &endTime, &location, &mapCoords, &active, &sortOrder, &category, &registrationURL, &isInternal, &capacity, &speakers, &imagesStr, &createdAt, &updatedAt); err != nil {
			log.Println("getEvents scan error:", err)
			continue
		}
		
		var speakersJSON, imagesJSON interface{}
		json.Unmarshal([]byte(speakers), &speakersJSON)
		json.Unmarshal([]byte(imagesStr), &imagesJSON)

		events = append(events, map[string]interface{}{
			"id": id, "title": title, "slug": slug, "description": desc, "image_url": imageURL,
			"event_date": eventDate, "start_time": startTime, "end_time": endTime,
			"location": location, "map_coordinates": mapCoords, "active": active, "sort_order": sortOrder,
			"category": category, "registration_url": registrationURL, "is_internal_registration": isInternal,
			"capacity": capacity, "speakers": speakersJSON, "images": imagesJSON,
			"created_at": createdAt, "updated_at": updatedAt,
		})
	}
	return c.JSON(events)
}

func getEventBySlug(c *fiber.Ctx) error {
	slugParam := c.Params("id") // Param is still named :id in router mapping, but it's actually slug
	var evID, title, slug, desc, imageURL, eventDate, startTime, endTime, location, mapCoords, category, registrationURL, speakers, waTemplate, imagesStr, createdAt, updatedAt string
	var active, isInternal bool
	var sortOrder, capacity int
	err := db.QueryRowContext(context.Background(),
		`SELECT id, title, slug, COALESCE(description,''), COALESCE(image_url,''), event_date,
		 COALESCE(start_time::text,'00:00'), COALESCE(end_time::text,'23:59'),
		 COALESCE(location,''), COALESCE(map_coordinates,''), active, sort_order,
		 COALESCE(category,''), COALESCE(registration_url,''), is_internal_registration, capacity, COALESCE(speakers::text,'[]'), COALESCE(wa_message_template,''),
		 COALESCE(images::text,'[]'), created_at, updated_at
		 FROM campus_events WHERE slug=$1`, slugParam,
	).Scan(&evID, &title, &slug, &desc, &imageURL, &eventDate, &startTime, &endTime, &location, &mapCoords, &active, &sortOrder, &category, &registrationURL, &isInternal, &capacity, &speakers, &waTemplate, &imagesStr, &createdAt, &updatedAt)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
	}

	var speakersJSON, imagesJSON interface{}
	json.Unmarshal([]byte(speakers), &speakersJSON)
	json.Unmarshal([]byte(imagesStr), &imagesJSON)

	return c.JSON(fiber.Map{
		"id": evID, "title": title, "slug": slug, "description": desc, "image_url": imageURL,
		"event_date": eventDate, "start_time": startTime, "end_time": endTime,
		"location": location, "map_coordinates": mapCoords, "active": active, "sort_order": sortOrder,
		"category": category, "registration_url": registrationURL, "is_internal_registration": isInternal,
		"capacity": capacity, "speakers": speakersJSON, "wa_message_template": waTemplate,
		"images": imagesJSON, "created_at": createdAt, "updated_at": updatedAt,
	})
}

type EventRegistrationPayload struct {
	EventID  string `json:"event_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Whatsapp string `json:"whatsapp"`
}

func registerForEvent(c *fiber.Ctx) error {
	var payload EventRegistrationPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Fetch event details
	var evID, title, eventDate, waTemplate string
	var capacity int
	var isInternal bool
	err := db.QueryRowContext(context.Background(),
		`SELECT id, title, event_date, capacity, is_internal_registration, COALESCE(wa_message_template, '') FROM campus_events WHERE id::text=$1 OR slug=$1`, c.Params("id"),
	).Scan(&evID, &title, &eventDate, &capacity, &isInternal, &waTemplate)

	if err != nil {
		log.Println("Error fetching event for registration:", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
	}

	if !isInternal {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Event does not accept internal registrations"})
	}

	// Check capacity
	if capacity > 0 {
		var count int
		db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM event_registrations WHERE event_id=$1`, evID).Scan(&count)
		if count >= capacity {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Mohon maaf, kuota peserta telah penuh."})
		}
	}

	// Insert registration
	_, err = db.ExecContext(context.Background(),
		`INSERT INTO event_registrations (event_id, full_name, email, whatsapp) VALUES ($1, $2, $3, $4)`,
		evID, payload.FullName, payload.Email, payload.Whatsapp,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Anda sudah terdaftar untuk event ini dengan email atau nomor WA yang sama."})
		}
		log.Println("Registration error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan pendaftaran"})
	}

	// Send WA message asynchronously
	go func() {
		msg := waTemplate
		if msg == "" {
			msg = fmt.Sprintf("Halo *{nama}*, pendaftaran Anda untuk event *{judul}* pada tanggal *{tanggal}* telah berhasil. Sampai jumpa di acara!")
		}
		// Replace tags
		msg = strings.ReplaceAll(msg, "{nama}", payload.FullName)
		msg = strings.ReplaceAll(msg, "{judul}", title)
		msg = strings.ReplaceAll(msg, "{tanggal}", eventDate[:10])
		
		SendWhatsAppMessage(payload.Whatsapp, msg)
	}()

	return c.JSON(fiber.Map{"success": true, "message": "Pendaftaran berhasil"})
}

// Placeholder handlers for remaining routes
func createNews(c *fiber.Ctx) error        { return c.JSON(fiber.Map{"message": "Not implemented"}) }
func updateNews(c *fiber.Ctx) error        { return c.JSON(fiber.Map{"message": "Not implemented"}) }
func deleteNews(c *fiber.Ctx) error        { return c.JSON(fiber.Map{"message": "Not implemented"}) }
func createTestimonial(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "Not implemented"}) }
func updateTestimonial(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "Not implemented"}) }
func deleteTestimonial(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "Not implemented"}) }
func createBlogPost(c *fiber.Ctx) error    { return c.JSON(fiber.Map{"message": "Not implemented"}) }
func updateBlogPost(c *fiber.Ctx) error    { return c.JSON(fiber.Map{"message": "Not implemented"}) }
func deleteBlogPost(c *fiber.Ctx) error    { return c.JSON(fiber.Map{"message": "Not implemented"}) }
func createPage(c *fiber.Ctx) error        { return c.JSON(fiber.Map{"message": "Not implemented"}) }
func updatePage(c *fiber.Ctx) error        { return c.JSON(fiber.Map{"message": "Not implemented"}) }
func deletePage(c *fiber.Ctx) error        { return c.JSON(fiber.Map{"message": "Not implemented"}) }
func deleteMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.ExecContext(context.Background(), "DELETE FROM contact_messages WHERE id=$1", id)
	if err != nil {
		log.Println("deleteMessage error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete message"})
	}
	return c.JSON(fiber.Map{"message": "Deleted successfully"})
}

func updateSettings(c *fiber.Ctx) error {
	var req SiteSetting
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	value, err := normalizeValue("value", req.Value)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	_, err = db.ExecContext(context.Background(),
		"INSERT INTO site_settings (key, value, updated_at) VALUES ($1, $2, NOW()) ON CONFLICT (key) DO UPDATE SET value=$2, updated_at=NOW()",
		req.Key, value)
	if err != nil {
		log.Println("updateSettings error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update settings"})
	}

	if req.Key == "campus_info" {
		go func(val interface{}) {
			var data []byte
			switch v := val.(type) {
			case string:
				data = []byte(v)
			case []byte:
				data = v
			default:
				var err error
				data, err = json.Marshal(v)
				if err != nil {
					log.Println("Failed to marshal campus_info value:", err)
					return
				}
			}
			var campusInfo struct {
				Name      string `json:"name"`
				ShortName string `json:"short_name"`
				Tagline   string `json:"tagline"`
			}
			if err := json.Unmarshal(data, &campusInfo); err != nil {
				log.Println("Failed to unmarshal campus_info struct:", err)
				return
			}
			if campusInfo.Name != "" || campusInfo.ShortName != "" || campusInfo.Tagline != "" {
				updateIndexHTMLMetadata(campusInfo.Name, campusInfo.ShortName, campusInfo.Tagline)
			}
		}(value)
	}

	return c.JSON(fiber.Map{"message": "Settings updated"})
}

// updateIndexHTMLMetadata updates the static index.html metadata on changes or startup
func updateIndexHTMLMetadata(name, shortName, tagline string) {
	if name == "" {
		name = "Universitas Perintis Indonesia"
	}
	if shortName == "" {
		shortName = "UPERTIS"
	}
	if tagline == "" {
		tagline = "Unggul, Professional, Etika, Rasional, Tangguh, Integritas, dan Solidaritas"
	}

	paths := []string{
		filepath.Join("..", "index.html"),
		filepath.Join("..", "dist", "index.html"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err != nil {
			continue
		}

		content, err := os.ReadFile(p)
		if err != nil {
			log.Printf("Failed to read %s: %v\n", p, err)
			continue
		}

		html := string(content)
		modified := false

		// Replace <title>...</title>
		reTitle := regexp.MustCompile(`(?i)<title>[^<]*</title>`)
		newTitle := fmt.Sprintf("<title>%s</title>", name)
		if reTitle.MatchString(html) {
			html = reTitle.ReplaceAllString(html, newTitle)
			modified = true
		}

		// Replace <meta name="description" content="..." />
		reDesc := regexp.MustCompile(`(?i)<meta\s+name="description"\s+content="[^"]*"\s*/?>`)
		newDesc := fmt.Sprintf(`<meta name="description" content="%s — %s" />`, shortName, tagline)
		if reDesc.MatchString(html) {
			html = reDesc.ReplaceAllString(html, newDesc)
			modified = true
		}

		// Replace <meta property="og:site_name" content="..." />
		reSiteName := regexp.MustCompile(`(?i)<meta\s+property="og:site_name"\s+content="[^"]*"\s*/?>`)
		newSiteName := fmt.Sprintf(`<meta property="og:site_name" content="%s" />`, name)
		if reSiteName.MatchString(html) {
			html = reSiteName.ReplaceAllString(html, newSiteName)
			modified = true
		}

		if modified {
			err = os.WriteFile(p, []byte(html), 0644)
			if err != nil {
				log.Printf("Failed to write to %s: %v\n", p, err)
			} else {
				log.Printf("Successfully updated metadata in %s\n", p)
			}
		}
	}
}

// syncIndexHTMLMetadata reads settings and runs metadata updates on startup
func syncIndexHTMLMetadata() {
	var value string
	err := db.QueryRowContext(context.Background(), "SELECT value FROM site_settings WHERE key = 'campus_info'").Scan(&value)
	if err != nil {
		log.Println("syncIndexHTMLMetadata: campus_info settings not found in DB or error:", err)
		return
	}

	var campusInfo struct {
		Name      string `json:"name"`
		ShortName string `json:"short_name"`
		Tagline   string `json:"tagline"`
	}
	if err := json.Unmarshal([]byte(value), &campusInfo); err != nil {
		log.Println("syncIndexHTMLMetadata: failed to unmarshal campus_info:", err)
		return
	}

	if campusInfo.Name != "" || campusInfo.ShortName != "" || campusInfo.Tagline != "" {
		updateIndexHTMLMetadata(campusInfo.Name, campusInfo.ShortName, campusInfo.Tagline)
	}
}

func getUsers(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(),
		`SELECT u.id, u.email, COALESCE(u.full_name, ''), u.role, COALESCE(u.status, 'active'), u.created_at,
			COALESCE(array_agg(fa.faculty_id) FILTER (WHERE fa.faculty_id IS NOT NULL), '{}') AS faculty_ids
		FROM users u
		LEFT JOIN faculty_admins fa ON fa.user_id = u.id
		GROUP BY u.id
		ORDER BY u.created_at DESC`)
	if err != nil {
		log.Println("getUsers query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.FullName, &u.Role, &u.Status, &u.CreatedAt, pq.Array(&u.FacultyIDs)); err != nil {
			log.Println("getUsers scan error:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process user data"})
		}
		users = append(users, u)
	}
	return c.JSON(users)
}

func createUser(c *fiber.Ctx) error {
	var req struct {
		Email      string   `json:"email"`
		Password   string   `json:"password"`
		FullName   string   `json:"full_name"`
		Role       string   `json:"role"`
		FacultyIDs []string `json:"faculty_ids"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if req.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}
	if req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Password is required"})
	}
	if len(req.Password) < 8 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Password must be at least 8 characters"})
	}
	if req.Role == "" {
		req.Role = "user"
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		log.Println("createUser hash error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process request"})
	}

	var id string
	err = db.QueryRowContext(context.Background(),
		`INSERT INTO users (email, password_hash, full_name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		req.Email, hashedPassword, req.FullName, req.Role).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Email sudah terdaftar"})
		}
		log.Println("createUser error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	if req.Role == "faculty_admin" {
		for _, facultyID := range req.FacultyIDs {
			if _, err := db.ExecContext(context.Background(), "INSERT INTO faculty_admins (user_id, faculty_id) VALUES ($1, $2)", id, facultyID); err != nil {
				log.Println("createUser faculty_admins insert error:", err)
			}
		}
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"id": id})
}

func updateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		FacultyIDs []string `json:"faculty_ids"`
		Status     string   `json:"status"`
		FullName   string   `json:"full_name"`
		Email      string   `json:"email"`
		Password   string   `json:"password"`
		Role       string   `json:"role"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	
	if req.Status != "" {
		if _, err := db.ExecContext(context.Background(), "UPDATE users SET status=$1 WHERE id=$2", req.Status, id); err != nil {
			log.Println("updateUser status error:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user status"})
		}
	}

	if req.Role != "" {
		if _, err := db.ExecContext(context.Background(), "UPDATE users SET role=$1 WHERE id=$2", req.Role, id); err != nil {
			log.Println("updateUser role error:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user role"})
		}
	}

	if req.FullName != "" || req.Email != "" || req.Password != "" {
		// Build dynamic query
		query := "UPDATE users SET updated_at=NOW()"
		var args []interface{}
		argId := 1

		if req.FullName != "" {
			query += fmt.Sprintf(", full_name=$%d", argId)
			args = append(args, req.FullName)
			argId++
		}
		if req.Email != "" {
			query += fmt.Sprintf(", email=$%d", argId)
			args = append(args, req.Email)
			argId++
		}
		if req.Password != "" {
			hashedPw, err := hashPassword(req.Password)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
			}
			query += fmt.Sprintf(", password_hash=$%d", argId)
			args = append(args, hashedPw)
			argId++
		}

		query += fmt.Sprintf(" WHERE id=$%d", argId)
		args = append(args, id)

		if _, err := db.ExecContext(context.Background(), query, args...); err != nil {
			log.Println("updateUser basic details error:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user details: " + err.Error()})
		}
	}
	
	if req.FacultyIDs != nil {
		if _, err := db.ExecContext(context.Background(), "DELETE FROM faculty_admins WHERE user_id=$1", id); err != nil {
			log.Println("updateUser faculty cleanup error:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
		}
		for _, facultyID := range req.FacultyIDs {
			if _, err := db.ExecContext(context.Background(), "INSERT INTO faculty_admins (user_id, faculty_id) VALUES ($1, $2)", id, facultyID); err != nil {
				log.Println("updateUser faculty insert error:", err)
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user faculty"})
			}
		}
	}
	return c.JSON(fiber.Map{"message": "User updated"})
}

func deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.ExecContext(context.Background(), "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		log.Println("deleteUser error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}
	return c.JSON(fiber.Map{"message": "User deleted"})
}

func createFaculty(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "Not implemented"}) }
func updateFaculty(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "Not implemented"}) }
func deleteFaculty(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "Not implemented"}) }

func getPopups(c *fiber.Ctx) error {
	rows, err := db.QueryContext(c.Context(), `
		SELECT id, title, description, image_url, link_text, link_url, start_date, end_date, active, sort_order, show_image_only
		FROM popup_banners
		WHERE active = true
		  AND (start_date IS NULL OR start_date <= NOW())
		  AND (end_date IS NULL OR end_date >= NOW())
		ORDER BY sort_order ASC, created_at DESC
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var popups []map[string]interface{}
	for rows.Next() {
		var id, title, linkText, linkUrl string
		var description, imageUrl sql.NullString
		var startDate, endDate sql.NullTime
		var active, showImageOnly bool
		var sortOrder int

		if err := rows.Scan(&id, &title, &description, &imageUrl, &linkText, &linkUrl, &startDate, &endDate, &active, &sortOrder, &showImageOnly); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		
		popup := map[string]interface{}{
			"id":              id,
			"title":           title,
			"description":     description.String,
			"image_url":       imageUrl.String,
			"link_text":       linkText,
			"link_url":        linkUrl,
			"active":          active,
			"sort_order":      sortOrder,
			"show_image_only": showImageOnly,
		}
		
		if startDate.Valid {
			popup["start_date"] = startDate.Time.Format(time.RFC3339)
		}
		if endDate.Valid {
			popup["end_date"] = endDate.Time.Format(time.RFC3339)
		}

		popups = append(popups, popup)
	}

	return c.JSON(popups)
}

func uploadLogo(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
	}

	logoType := c.FormValue("type", "logo") // "favicon" or "logo"

	var actualSavePath string
	var actualDistPath string

	// Dynamically determine paths based on whether process CWD is project root or backend folder
	if _, err := os.Stat("public"); err == nil {
		if logoType == "favicon" {
			actualSavePath = filepath.Join("public", "favicon.ico")
			actualDistPath = filepath.Join("dist", "favicon.ico")
		} else {
			actualSavePath = filepath.Join("public", "logo.png")
			actualDistPath = filepath.Join("dist", "logo.png")
		}
	} else {
		if logoType == "favicon" {
			actualSavePath = filepath.Join("..", "public", "favicon.ico")
			actualDistPath = filepath.Join("..", "dist", "favicon.ico")
		} else {
			actualSavePath = filepath.Join("..", "public", "logo.png")
			actualDistPath = filepath.Join("..", "dist", "logo.png")
		}
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(actualSavePath), 0755); err != nil {
		log.Println("uploadLogo MkdirAll public error:", err)
	}

	if err := c.SaveFile(file, actualSavePath); err != nil {
		log.Println("uploadLogo save error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	// Also copy to dist if the dist folder exists, so it reflects immediately in Nginx
	distDir := filepath.Dir(actualDistPath)
	if _, err := os.Stat(distDir); err == nil {
		input, err := os.ReadFile(actualSavePath)
		if err == nil {
			err = os.WriteFile(actualDistPath, input, 0644)
			if err != nil {
				log.Println("uploadLogo copy to dist error:", err)
			} else {
				log.Println("uploadLogo successfully copied to:", actualDistPath)
			}
		} else {
			log.Println("uploadLogo read savePath error:", err)
		}
	}

	publicPath := "/favicon.ico"
	if logoType != "favicon" {
		publicPath = "/logo.png"
	}

	return c.JSON(fiber.Map{"message": "File uploaded successfully", "path": publicPath})
}

// getFAQs returns all active FAQs ordered by sort_order
func getFAQs(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(), `
		SELECT id, question, answer, category, sort_order, active, created_at, updated_at
		FROM faqs
		WHERE active = true
		ORDER BY sort_order ASC, created_at ASC
	`)
	if err != nil {
		log.Println("getFAQs query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch FAQs"})
	}
	defer rows.Close()

	type FAQ struct {
		ID        string    `json:"id"`
		Question  string    `json:"question"`
		Answer    string    `json:"answer"`
		Category  string    `json:"category"`
		SortOrder int       `json:"sort_order"`
		Active    bool      `json:"active"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	faqs := []FAQ{}
	for rows.Next() {
		var f FAQ
		if err := rows.Scan(&f.ID, &f.Question, &f.Answer, &f.Category, &f.SortOrder, &f.Active, &f.CreatedAt, &f.UpdatedAt); err != nil {
			log.Println("getFAQs scan error:", err)
			continue
		}
		faqs = append(faqs, f)
	}
	return c.JSON(faqs)
}

// handlePageRequest serves the index.html with dynamically injected SEO meta tags
func handlePageRequest(c *fiber.Ctx) error {
	path := c.Path()

	// Skip API and uploads routes
	if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/uploads") {
		return c.Next()
	}

	// Read index.html
	indexPath := filepath.Join("..", "dist", "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		// Fallback to source index.html
		indexPath = filepath.Join("..", "index.html")
	}

	content, err := os.ReadFile(indexPath)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("index.html not found")
	}

	html := string(content)

	// Default values
	var title = "Universitas Perintis Indonesia"
	var description = "UPERTIS — Unggul, Professional, Etika, Rasional, Tangguh, Integritas, dan Solidaritas"
	var image = "https://upertis.ac.id/logo.png"

	// Fetch default campus info from db
	var value string
	err = db.QueryRowContext(context.Background(), "SELECT value FROM site_settings WHERE key = 'campus_info'").Scan(&value)
	if err == nil {
		var campusInfo struct {
			Name      string `json:"name"`
			ShortName string `json:"short_name"`
			Tagline   string `json:"tagline"`
		}
		if json.Unmarshal([]byte(value), &campusInfo) == nil {
			if campusInfo.Name != "" {
				title = campusInfo.Name
			}
			if campusInfo.ShortName != "" && campusInfo.Tagline != "" {
				description = fmt.Sprintf("%s — %s", campusInfo.ShortName, campusInfo.Tagline)
			}
		}
	}

	// If route is a post detail page
	// Format: /posts/:slug
	if strings.HasPrefix(path, "/posts/") {
		slug := strings.TrimPrefix(path, "/posts/")
		// Fetch post from database
		var postTitle, postExcerpt, postImage sql.NullString
		err := db.QueryRowContext(context.Background(),
			`SELECT title, excerpt, featured_image FROM posts WHERE slug = $1 AND status = 'published' AND deleted_at IS NULL`,
			slug,
		).Scan(&postTitle, &postExcerpt, &postImage)

		if err == nil {
			if postTitle.Valid && postTitle.String != "" {
				title = fmt.Sprintf("%s — UPERTIS", postTitle.String)
			}
			if postExcerpt.Valid && postExcerpt.String != "" {
				description = postExcerpt.String
			}
			if postImage.Valid && postImage.String != "" {
				image = postImage.String
				// If relative image path, prefix with domain
				if strings.HasPrefix(image, "/") {
					domain := getEnv("SITE_DOMAIN", "https://upertis.ac.id")
					image = domain + image
				}
			}
		}
	} else if strings.HasPrefix(path, "/berita/") {
		// Format: /berita/:id
		id := strings.TrimPrefix(path, "/berita/")
		var newsTitle, newsExcerpt, newsImage sql.NullString
		err := db.QueryRowContext(context.Background(),
			`SELECT title, excerpt, image_url FROM news WHERE id::text = $1 AND deleted_at IS NULL`,
			id,
		).Scan(&newsTitle, &newsExcerpt, &newsImage)

		if err == nil {
			if newsTitle.Valid && newsTitle.String != "" {
				title = fmt.Sprintf("%s — UPERTIS", newsTitle.String)
			}
			if newsExcerpt.Valid && newsExcerpt.String != "" {
				description = newsExcerpt.String
			}
			if newsImage.Valid && newsImage.String != "" {
				image = newsImage.String
				if strings.HasPrefix(image, "/") {
					domain := getEnv("SITE_DOMAIN", "https://upertis.ac.id")
					image = domain + image
				}
			}
		}
	} else if strings.HasPrefix(path, "/events/") {
		// Format: /events/:slug
		slug := strings.TrimPrefix(path, "/events/")
		var eventTitle, eventDesc, eventImage sql.NullString
		err := db.QueryRowContext(context.Background(),
			`SELECT title, description, image_url FROM campus_events WHERE slug = $1 AND deleted_at IS NULL`,
			slug,
		).Scan(&eventTitle, &eventDesc, &eventImage)

		if err == nil {
			if eventTitle.Valid && eventTitle.String != "" {
				title = fmt.Sprintf("%s — UPERTIS", eventTitle.String)
			}
			if eventDesc.Valid && eventDesc.String != "" {
				description = eventDesc.String
				description = cleanHTMLTags(description)
				if len(description) > 200 {
					description = description[:200] + "..."
				}
			}
			if eventImage.Valid && eventImage.String != "" {
				image = eventImage.String
				if strings.HasPrefix(image, "/") {
					domain := getEnv("SITE_DOMAIN", "https://upertis.ac.id")
					image = domain + image
				}
			}
		}
	}

	// Dynamically replace the meta tags in HTML
	// Replace <title>...</title>
	reTitle := regexp.MustCompile(`(?i)<title>[^<]*</title>`)
	html = reTitle.ReplaceAllString(html, fmt.Sprintf("<title>%s</title>", title))

	// Replace description
	reDesc := regexp.MustCompile(`(?i)<meta\s+name="description"\s+content="[^"]*"\s*/?>`)
	html = reDesc.ReplaceAllString(html, fmt.Sprintf(`<meta name="description" content="%s" />`, description))

	// Replace Open Graph meta tags
	reOgTitle := regexp.MustCompile(`(?i)<meta\s+property="og:title"\s+content="[^"]*"\s*/?>`)
	if reOgTitle.MatchString(html) {
		html = reOgTitle.ReplaceAllString(html, fmt.Sprintf(`<meta property="og:title" content="%s" />`, title))
	} else {
		html = strings.Replace(html, "</head>", fmt.Sprintf(`<meta property="og:title" content="%s" /></head>`, title), 1)
	}

	reOgDesc := regexp.MustCompile(`(?i)<meta\s+property="og:description"\s+content="[^"]*"\s*/?>`)
	if reOgDesc.MatchString(html) {
		html = reOgDesc.ReplaceAllString(html, fmt.Sprintf(`<meta property="og:description" content="%s" />`, description))
	} else {
		html = strings.Replace(html, "</head>", fmt.Sprintf(`<meta property="og:description" content="%s" /></head>`, description), 1)
	}

	reOgImage := regexp.MustCompile(`(?i)<meta\s+property="og:image"\s+content="[^"]*"\s*/?>`)
	if reOgImage.MatchString(html) {
		html = reOgImage.ReplaceAllString(html, fmt.Sprintf(`<meta property="og:image" content="%s" />`, image))
	} else {
		html = strings.Replace(html, "</head>", fmt.Sprintf(`<meta property="og:image" content="%s" /></head>`, image), 1)
	}

	c.Set("Content-Type", "text/html")
	return c.SendString(html)
}

func cleanHTMLTags(src string) string {
	re := regexp.MustCompile("<[^>]*>")
	return re.ReplaceAllString(src, "")
}

func pddiktiJSONArrayLen(body []byte) int {
	var list []interface{}
	if err := json.Unmarshal(body, &list); err == nil {
		return len(list)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(body, &obj); err == nil {
		for _, key := range []string{"data", "result", "items"} {
			if vv, ok := obj[key]; ok {
				if subList, ok := vv.([]interface{}); ok {
					return len(subList)
				}
			}
		}
	}
	return 0
}

type searchPDDIKTILecturerItem struct {
	ID        string `json:"id"`
	Nama      string `json:"nama"`
	NIDN      string `json:"nidn"`
	PTName    string `json:"nama_pt"`
	ProdiName string `json:"nama_prodi"`
}

func searchPDDIKTILecturer(c *fiber.Ctx) error {
	keyword := strings.TrimSpace(c.Query("q"))
	if keyword == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Parameter 'q' tidak boleh kosong."})
	}

	endpoint := "https://api-pddikti.kemdiktisaintek.go.id/pencarian/dosen/" + url.PathEscape(keyword)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	body, status, err := pddiktiGET(ctx, endpoint)
	if err != nil {
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Gagal terhubung ke PDDIKTI: " + err.Error()})
	}
	if status != http.StatusOK {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 300 {
			snippet = snippet[:300]
		}
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": fmt.Sprintf("PDDIKTI API error: HTTP %d: %s", status, snippet)})
	}

	var items []searchPDDIKTILecturerItem
	if err := json.Unmarshal(body, &items); err != nil {
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Gagal mengurai response PDDIKTI"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   items,
	})
}

type importSingleLecturerRequest struct {
	UUID  string `json:"uuid"`
	NIDN  string `json:"nidn"`
	Name  string `json:"name"`
	Prodi string `json:"prodi"`
}

func importSinglePDDIKTILecturer(c *fiber.Ctx) error {
	var req importSingleLecturerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if strings.TrimSpace(req.UUID) == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "uuid is required"})
	}

	ctx := context.Background()

	// Fetch profile
	profileEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/profile/" + url.PathEscape(req.UUID)
	pBody, pStatus, pErr := pddiktiGET(ctx, profileEndpoint)
	if pErr != nil || pStatus != http.StatusOK {
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": fmt.Sprintf("Gagal mengambil profil dosen dari PDDIKTI (HTTP %d)", pStatus)})
	}

	var rawProfile interface{}
	if err := json.Unmarshal(pBody, &rawProfile); err != nil {
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": "Gagal mengurai profil dosen dari PDDIKTI"})
	}

	portfolio := map[string]interface{}{}
	var highestEdu, academicPosition string

	if m, ok := rawProfile.(map[string]interface{}); ok {
		m["nidn"] = req.NIDN
		portfolio["profile"] = m
		if ja, ok := m["jabatan_akademik"].(string); ok && ja != "" {
			academicPosition = ja
		}
		if pt, ok := m["pendidikan_tertinggi"].(string); ok && pt != "" {
			highestEdu = pt
		}
	}

	// Fetch other details
	var teachingCount, pubCount, pengCount int

	// 1. Teaching
	teachEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/teaching-history/" + url.PathEscape(req.UUID)
	if tBody, tStatus, tErr := pddiktiGET(ctx, teachEndpoint); tErr == nil && tStatus == 200 {
		var raw interface{}
		json.Unmarshal(tBody, &raw)
		portfolio["teaching_history"] = raw
		teachingCount = pddiktiJSONArrayLen(tBody)
	}

	// 2. Study
	studyEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/study-history/" + url.PathEscape(req.UUID)
	if sBody, sStatus, sErr := pddiktiGET(ctx, studyEndpoint); sErr == nil && sStatus == 200 {
		var raw interface{}
		json.Unmarshal(sBody, &raw)
		portfolio["study_history"] = raw
	}

	// 3. Penelitian
	penelitianEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/penelitian/" + url.PathEscape(req.UUID)
	if rBody, rStatus, rErr := pddiktiGET(ctx, penelitianEndpoint); rErr == nil && rStatus == 200 {
		var raw interface{}
		json.Unmarshal(rBody, &raw)
		portfolio["penelitian"] = raw
		pubCount = pddiktiJSONArrayLen(rBody)
	}

	// 4. Pengabdian
	pengabdianEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/pengabdian/" + url.PathEscape(req.UUID)
	if aBody, aStatus, aErr := pddiktiGET(ctx, pengabdianEndpoint); aErr == nil && aStatus == 200 {
		var raw interface{}
		json.Unmarshal(aBody, &raw)
		portfolio["pengabdian"] = raw
		pengCount = pddiktiJSONArrayLen(aBody)
	}

	// 5. Karya
	karyaEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/karya/" + url.PathEscape(req.UUID)
	if kBody, kStatus, kErr := pddiktiGET(ctx, karyaEndpoint); kErr == nil && kStatus == 200 {
		var raw interface{}
		json.Unmarshal(kBody, &raw)
		portfolio["karya"] = raw
	}

	// 6. Paten
	patenEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/paten/" + url.PathEscape(req.UUID)
	if ptBody, ptStatus, ptErr := pddiktiGET(ctx, patenEndpoint); ptErr == nil && ptStatus == 200 {
		var raw interface{}
		json.Unmarshal(ptBody, &raw)
		portfolio["paten"] = raw
	}

	portfolioJSON, err := json.Marshal(portfolio)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat JSON portfolio"})
	}

	// Resolve Faculty ID
	prodiLower := strings.ToLower(strings.TrimSpace(req.Prodi))
	var facultyID string
	var query string

	if strings.Contains(prodiLower, "bisnis") || strings.Contains(prodiLower, "ekonomi") || strings.Contains(prodiLower, "manajemen") || strings.Contains(prodiLower, "digital") {
		query = "SELECT id FROM faculties WHERE active = true AND (LOWER(name) LIKE '%bisnis%' OR LOWER(name) LIKE '%ekonomi%' OR LOWER(name) LIKE '%manajemen%') ORDER BY sort_order LIMIT 1"
	} else if strings.Contains(prodiLower, "farmasi") {
		query = "SELECT id FROM faculties WHERE active = true AND LOWER(name) LIKE '%farmasi%' ORDER BY sort_order LIMIT 1"
	} else {
		query = "SELECT id FROM faculties WHERE active = true AND (LOWER(name) LIKE '%kesehatan%' OR LOWER(name) LIKE '%keperawatan%' OR LOWER(name) LIKE '%fikes%') ORDER BY sort_order LIMIT 1"
	}

	err = db.QueryRowContext(ctx, query).Scan(&facultyID)
	if err != nil {
		_ = db.QueryRowContext(ctx, "SELECT id FROM faculties WHERE active = true ORDER BY sort_order LIMIT 1").Scan(&facultyID)
	}

	if facultyID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Tidak ada fakultas aktif untuk menampung dosen"})
	}

	// Format Education text: e.g. "S2 | Mengajar:15 | Publikasi:2"
	parts := []string{}
	if highestEdu != "" {
		parts = append(parts, highestEdu)
	}
	if teachingCount > 0 {
		parts = append(parts, fmt.Sprintf("Mengajar:%d", teachingCount))
	}
	if pubCount > 0 {
		parts = append(parts, fmt.Sprintf("Publikasi:%d", pubCount))
	}
	if pengCount > 0 {
		parts = append(parts, fmt.Sprintf("Pengabdian:%d", pengCount))
	}
	formattedEdu := strings.Join(parts, " | ")

	slug := lecturerSlug(req.Name, "")
	position := academicPosition
	if position == "" {
		position = "Dosen"
	}

	insertSQL := `
		INSERT INTO faculty_lecturers (faculty_id, slug, name, position, education, expertise, photo_url, leadership_group, source, active, sort_order, pddikti_id, pddikti_uuid, portfolio)
		VALUES ($1, $2, $3, $4, $5, $6, '', '', 'pddikti', true, 0, NULLIF($7, ''), NULLIF($8, ''), $9::jsonb)
		ON CONFLICT (pddikti_uuid) DO UPDATE SET
			faculty_id = EXCLUDED.faculty_id,
			name = EXCLUDED.name,
			position = CASE
				WHEN LOWER(faculty_lecturers.position) LIKE '%rektor%'
				  OR LOWER(faculty_lecturers.position) LIKE '%dekan%'
				  OR LOWER(faculty_lecturers.position) LIKE '%kaprodi%'
				  OR LOWER(faculty_lecturers.position) LIKE '%prodi%'
				  OR LOWER(faculty_lecturers.position) LIKE '%yayasan%'
				  OR LOWER(faculty_lecturers.position) LIKE '%founder%'
				  THEN faculty_lecturers.position
				ELSE EXCLUDED.position
			END,
			education = EXCLUDED.education,
			expertise = EXCLUDED.expertise,
			photo_url = CASE
				WHEN COALESCE(faculty_lecturers.photo_url, '') != '' THEN faculty_lecturers.photo_url
				ELSE EXCLUDED.photo_url
			END,
			source = EXCLUDED.source,
			active = true,
			updated_at = NOW(),
			pddikti_id = EXCLUDED.pddikti_id,
			portfolio = EXCLUDED.portfolio,
			leadership_group = CASE
				WHEN COALESCE(faculty_lecturers.leadership_group, '') = '' THEN EXCLUDED.leadership_group
				ELSE faculty_lecturers.leadership_group
			END
		RETURNING id
	`

	var insertedID string
	err = db.QueryRowContext(ctx, insertSQL, facultyID, slug, req.Name, position, formattedEdu, req.Prodi, req.NIDN, req.UUID, string(portfolioJSON)).Scan(&insertedID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan dosen ke database: " + err.Error()})
	}

	// Ensure unique slug
	uniqueSlug := lecturerSlug(req.Name, insertedID)
	_, _ = db.ExecContext(ctx, "UPDATE faculty_lecturers SET slug = $1 WHERE id = $2", uniqueSlug, insertedID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Dosen '%s' berhasil di-import/di-update.", req.Name),
		"id":      insertedID,
		"slug":    uniqueSlug,
	})
}

func getGallery(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(), `
		SELECT id, title, description, media_url, media_type, active, created_at, updated_at
		FROM gallery
		WHERE active = true
		ORDER BY created_at DESC
	`)
	if err != nil {
		log.Println("getGallery query error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch gallery items"})
	}
	defer rows.Close()

	items := []fiber.Map{}
	for rows.Next() {
		var id, title, description, mediaURL, mediaType string
		var active bool
		var createdAt, updatedAt time.Time
		err := rows.Scan(&id, &title, &description, &mediaURL, &mediaType, &active, &createdAt, &updatedAt)
		if err != nil {
			log.Println("getGallery scan error:", err)
			continue
		}

		items = append(items, fiber.Map{
			"id":          id,
			"title":       title,
			"description": description,
			"media_url":   mediaURL,
			"media_type":  mediaType,
			"active":      active,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		})
	}

	return c.JSON(items)
}


