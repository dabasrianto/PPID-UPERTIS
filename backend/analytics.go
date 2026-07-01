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

func getVisitorAnalytics(c *fiber.Ctx) error {
	ctx := context.Background()
	now := time.Now()

	// 1. Get total active visitors (last 5 minutes)
	var activeCount int
	_ = db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM (
			SELECT ip_address FROM active_visitors WHERE last_active > NOW() - INTERVAL '5 minutes'
			UNION
			SELECT ip_address FROM active_sessions WHERE last_active > NOW() - INTERVAL '5 minutes'
		) t
	`).Scan(&activeCount)

	// 2. Get active visitor list details
	activeVisitors := []fiber.Map{}
	
	// Active anonymous visitors
	vRows, err := db.QueryContext(ctx, `
		SELECT ip_address, last_page, browser, os, last_active, COALESCE(country_code, 'ID')
		FROM active_visitors
		WHERE last_active > NOW() - INTERVAL '5 minutes'
		ORDER BY last_active DESC LIMIT 30
	`)
	if err == nil {
		defer vRows.Close()
		for vRows.Next() {
			var ip, page, browser, os, country string
			var lastActive time.Time
			if vRows.Scan(&ip, &page, &browser, &os, &lastActive, &country) == nil {
				activeVisitors = append(activeVisitors, fiber.Map{
					"ip_address": ip,
					"last_page": page,
					"browser": browser,
					"os": os,
					"last_active": lastActive,
					"country_code": country,
					"type": "visitor",
				})
			}
		}
	}

	// Active admin/user sessions
	sRows, err := db.QueryContext(ctx, `
		SELECT s.ip_address, s.last_page, s.browser, s.os, s.last_active, COALESCE(s.country_code, 'ID'), u.full_name
		FROM active_sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.last_active > NOW() - INTERVAL '5 minutes'
		ORDER BY s.last_active DESC LIMIT 20
	`)
	if err == nil {
		defer sRows.Close()
		for sRows.Next() {
			var ip, page, browser, os, country, name string
			var lastActive time.Time
			if sRows.Scan(&ip, &page, &browser, &os, &lastActive, &country, &name) == nil {
				activeVisitors = append(activeVisitors, fiber.Map{
					"ip_address": ip,
					"last_page": page,
					"browser": browser,
					"os": os,
					"last_active": lastActive,
					"country_code": country,
					"full_name": name,
					"type": "admin",
				})
			}
		}
	}

	// 3. Realtime Timeline: last 30 minutes
	type MinuteBucket struct {
		Time  string `json:"time"`
		Count int    `json:"count"`
	}
	buckets30m := make([]MinuteBucket, 30)
	bucketsMap := make(map[string]int)
	for i := 29; i >= 0; i-- {
		t := now.Add(-time.Duration(i) * time.Minute)
		tStr := t.Format("15:04")
		buckets30m[29-i] = MinuteBucket{Time: tStr, Count: 0}
		bucketsMap[tStr] = 29 - i
	}

	rowsMin, err := db.QueryContext(ctx, `
		SELECT TO_CHAR(visited_at, 'HH24:MI') as min_bucket, COUNT(*)
		FROM visitor_logs
		WHERE visited_at >= NOW() - INTERVAL '30 minutes'
		GROUP BY min_bucket
	`)
	if err == nil {
		defer rowsMin.Close()
		for rowsMin.Next() {
			var minBucket string
			var count int
			if rowsMin.Scan(&minBucket, &count) == nil {
				if idx, found := bucketsMap[minBucket]; found {
					buckets30m[idx].Count = count
				}
			}
		}
	}

	// 4. Historical Timeline: last 30 days
	type DailyTimeline struct {
		Date      string `json:"date"`
		Pageviews int    `json:"pageviews"`
		Visitors  int    `json:"visitors"`
	}
	buckets30d := make([]DailyTimeline, 30)
	daysMap := make(map[string]int)
	for i := 29; i >= 0; i-- {
		d := now.AddDate(0, 0, -i)
		dStr := d.Format("2006-01-02")
		label := d.Format("02 Jan")
		buckets30d[29-i] = DailyTimeline{Date: label, Pageviews: 0, Visitors: 0}
		daysMap[dStr] = 29 - i
	}

	rowsDay, err := db.QueryContext(ctx, `
		SELECT DATE(visited_at) as date_bucket, COUNT(*), COUNT(DISTINCT ip_address)
		FROM visitor_logs
		WHERE visited_at >= NOW() - INTERVAL '30 days'
		GROUP BY date_bucket
	`)
	if err == nil {
		defer rowsDay.Close()
		for rowsDay.Next() {
			var dateBucket time.Time
			var pageviews, visitors int
			if rowsDay.Scan(&dateBucket, &pageviews, &visitors) == nil {
				dStr := dateBucket.Format("2006-01-02")
				if idx, found := daysMap[dStr]; found {
					buckets30d[idx].Pageviews = pageviews
					buckets30d[idx].Visitors = visitors
				}
			}
		}
	}

	// 5. Top Pages (last 30 days)
	type TopPage struct {
		Page  string `json:"page"`
		Count int    `json:"count"`
	}
	topPages := []TopPage{}
	rowsPages, err := db.QueryContext(ctx, `
		SELECT page_url, COUNT(*) as count
		FROM visitor_logs
		WHERE visited_at >= NOW() - INTERVAL '30 days'
		GROUP BY page_url
		ORDER BY count DESC LIMIT 10
	`)
	if err == nil {
		defer rowsPages.Close()
		for rowsPages.Next() {
			var page string
			var count int
			if rowsPages.Scan(&page, &count) == nil {
				topPages = append(topPages, TopPage{Page: page, Count: count})
			}
		}
	}

	// 6. Top Countries (last 30 days)
	type TopCountry struct {
		Country string `json:"country"`
		Count   int    `json:"count"`
	}
	topCountries := []TopCountry{}
	rowsCountries, err := db.QueryContext(ctx, `
		SELECT country_code, COUNT(*) as count
		FROM visitor_logs
		WHERE visited_at >= NOW() - INTERVAL '30 days'
		GROUP BY country_code
		ORDER BY count DESC LIMIT 10
	`)
	if err == nil {
		defer rowsCountries.Close()
		for rowsCountries.Next() {
			var country string
			var count int
			if rowsCountries.Scan(&country, &count) == nil {
				topCountries = append(topCountries, TopCountry{Country: country, Count: count})
			}
		}
	}

	// 7. Top Browsers (last 30 days)
	type TopBrowser struct {
		Browser string `json:"browser"`
		Count   int    `json:"count"`
	}
	topBrowsers := []TopBrowser{}
	rowsBrowsers, err := db.QueryContext(ctx, `
		SELECT COALESCE(browser, 'Other') as browser, COUNT(*) as count
		FROM visitor_logs
		WHERE visited_at >= NOW() - INTERVAL '30 days'
		GROUP BY browser
		ORDER BY count DESC LIMIT 5
	`)
	if err == nil {
		defer rowsBrowsers.Close()
		for rowsBrowsers.Next() {
			var browser string
			var count int
			if rowsBrowsers.Scan(&browser, &count) == nil {
				topBrowsers = append(topBrowsers, TopBrowser{Browser: browser, Count: count})
			}
		}
	}

	// 8. Top OS (last 30 days)
	type TopOS struct {
		OS    string `json:"os"`
		Count int    `json:"count"`
	}
	topOS := []TopOS{}
	rowsOS, err := db.QueryContext(ctx, `
		SELECT COALESCE(os, 'Other') as os, COUNT(*) as count
		FROM visitor_logs
		WHERE visited_at >= NOW() - INTERVAL '30 days'
		GROUP BY os
		ORDER BY count DESC LIMIT 5
	`)
	if err == nil {
		defer rowsOS.Close()
		for rowsOS.Next() {
			var os string
			var count int
			if rowsOS.Scan(&os, &count) == nil {
				topOS = append(topOS, TopOS{OS: os, Count: count})
			}
		}
	}

	return c.JSON(fiber.Map{
		"active_count":    activeCount,
		"active_visitors": activeVisitors,
		"timeline_30m":    buckets30m,
		"timeline_30d":    buckets30d,
		"top_pages":       topPages,
		"top_countries":   topCountries,
		"top_browsers":    topBrowsers,
		"top_os":          topOS,
	})
}

