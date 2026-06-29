package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"context"
)

type WXR struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title      string `xml:"title"`
	Link       string `xml:"link"`
	Creator    string `xml:"dc:creator"`
	Content    string `xml:"content:encoded"`
	Excerpt    string `xml:"excerpt:encoded"`
	PostDate   string `xml:"wp:post_date"`
	PostType   string `xml:"wp:post_type"`
	Status     string `xml:"wp:status"`
	PostName   string `xml:"wp:post_name"`
	Categories []WPCategory `xml:"category"`
}

type WPCategory struct {
	Domain string `xml:"domain,attr"`
	Value  string `xml:",chardata"`
}

func importWordPressXML(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "File tidak ditemukan"})
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

	var wxr WXR
	if err := xml.Unmarshal(data, &wxr); err != nil {
		log.Println("XML Unmarshal error:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Format XML tidak valid atau bukan format WordPress"})
	}

	imported := 0
	skipped := 0

	for _, item := range wxr.Channel.Items {
		// Hanya ambil post (bukan page atau attachment)
		if item.PostType != "post" {
			continue
		}

		// Check status (bisa disesuaikan jika ingin hanya yang publish)
		if item.Status != "publish" && item.Status != "draft" {
			continue
		}

		title := sanitizeTitle(item.Title)
		slug := item.PostName
		if slug == "" {
			// fallback to slugify title
			slug = strings.ToLower(strings.ReplaceAll(title, " ", "-"))
		}

		// Check duplicate
		var existingId string
		err := db.QueryRowContext(context.Background(), "SELECT id FROM blog_posts WHERE slug=$1", slug).Scan(&existingId)
		if err == nil && existingId != "" {
			skipped++
			continue
		}

		// Get category
		categoryName := "Umum"
		for _, cat := range item.Categories {
			if cat.Domain == "category" {
				categoryName = cat.Value
				break
			}
		}

		status := "published"

		// Insert
		_, err = db.ExecContext(context.Background(), `
			INSERT INTO blog_posts (
				title, slug, excerpt, content, cover_image_url, 
				author_name, category, status, published_at, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		`,
			title,
			slug,
			item.Excerpt,
			item.Content,
			"", // XML WordPress biasanya simpan attachment terpisah, kita biarkan kosong dulu
			item.Creator,
			categoryName,
			status,
			item.PostDate,
		)

		if err != nil {
			log.Println("Error inserting XML post:", err)
			continue
		}
		imported++
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Berhasil mengimpor %d artikel dari XML (%d dilewati).", imported, skipped),
		"imported": imported,
		"skipped": skipped,
	})
}
