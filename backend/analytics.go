package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DailyStat struct {
	Date      string `json:"date"`
	Label     string `json:"label"`
	Pesan     int    `json:"pesan"`
	Konten    int    `json:"konten"`
	Kunjungan int    `json:"kunjungan"`
}

func getAnalytics(c *fiber.Ctx) error {
	// Initialize 30 days of data
	days := 30
	today := time.Now()
	
	bucketsMap := make(map[string]*DailyStat)
	var buckets []DailyStat

	for i := days - 1; i >= 0; i-- {
		d := today.AddDate(0, 0, -i)
		key := d.Format("2006-01-02")
		label := d.Format("02 Jan")
		
		stat := DailyStat{
			Date:  key,
			Label: label,
		}
		buckets = append(buckets, stat)
		bucketsMap[key] = &buckets[len(buckets)-1]
	}

	// Fetch Messages
	msgRows, err := db.QueryContext(context.Background(), "SELECT DATE(created_at), COUNT(*) FROM contact_messages WHERE created_at >= NOW() - INTERVAL '30 days' GROUP BY DATE(created_at)")
	if err == nil {
		defer msgRows.Close()
		for msgRows.Next() {
			var d time.Time
			var count int
			if err := msgRows.Scan(&d, &count); err == nil {
				key := d.Format("2006-01-02")
				if stat, exists := bucketsMap[key]; exists {
					stat.Pesan += count
				}
			}
		}
	}

	// Fetch Content (Unified Posts)
	contentRows, err := db.QueryContext(context.Background(), "SELECT DATE(created_at), COUNT(*) FROM posts WHERE deleted_at IS NULL AND created_at >= NOW() - INTERVAL '30 days' GROUP BY DATE(created_at)")
	if err == nil {
		defer contentRows.Close()
		for contentRows.Next() {
			var d time.Time
			var count int
			if err := contentRows.Scan(&d, &count); err == nil {
				key := d.Format("2006-01-02")
				if stat, exists := bucketsMap[key]; exists {
					stat.Konten += count
				}
			}
		}
	}

	// Fetch Real Views (from blog_post_views)
	viewRows, err := db.QueryContext(context.Background(), "SELECT DATE(viewed_at), COUNT(*) FROM blog_post_views WHERE viewed_at >= NOW() - INTERVAL '30 days' GROUP BY DATE(viewed_at)")
	if err == nil {
		defer viewRows.Close()
		for viewRows.Next() {
			var d time.Time
			var count int
			if err := viewRows.Scan(&d, &count); err == nil {
				key := d.Format("2006-01-02")
				if stat, exists := bucketsMap[key]; exists {
					stat.Kunjungan += count
				}
			}
		}
	} else {
		log.Println("Analytics views error (maybe table missing viewed_at?):", err)
	}

	return c.JSON(buckets)
}
