package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	ipCountryCache = make(map[string]string)
	ipCountryMu    sync.RWMutex
)

func resolveCountry(ip string, cfCountry string) string {
	if cfCountry != "" {
		return strings.ToUpper(cfCountry)
	}
	// Default local networks to ID
	if ip == "127.0.0.1" || ip == "::1" || strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "172.16.") {
		return "ID"
	}
	return fetchCountryFromIP(ip)
}

func fetchCountryFromIP(ip string) string {
	ipCountryMu.RLock()
	if cCode, exists := ipCountryCache[ip]; exists {
		ipCountryMu.RUnlock()
		return cCode
	}
	ipCountryMu.RUnlock()

	countryCode := "ID" // default fallback
	client := http.Client{Timeout: 1500 * time.Millisecond}
	resp, err := client.Get("http://ip-api.com/json/" + ip)
	if err == nil {
		defer resp.Body.Close()
		var res struct {
			CountryCode string `json:"countryCode"`
			Status      string `json:"status"`
		}
		if json.NewDecoder(resp.Body).Decode(&res) == nil && res.Status == "success" {
			countryCode = res.CountryCode
		}
	}

	ipCountryMu.Lock()
	ipCountryCache[ip] = countryCode
	ipCountryMu.Unlock()

	return countryCode
}

func migrateLiveUsersTable() {
	db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS active_sessions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE UNIQUE,
			ip_address VARCHAR(50),
			user_agent VARCHAR(500),
			last_page VARCHAR(500),
			last_active TIMESTAMP DEFAULT NOW(),
			browser VARCHAR(100),
			os VARCHAR(100)
		)
	`)
	db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_active_sessions_last ON active_sessions(last_active DESC)`)

	db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS active_visitors (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			ip_address VARCHAR(50) UNIQUE,
			user_agent VARCHAR(500),
			last_page VARCHAR(500),
			last_active TIMESTAMP DEFAULT NOW(),
			browser VARCHAR(100),
			os VARCHAR(100)
		)
	`)
	db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_active_visitors_last ON active_visitors(last_active DESC)`)

	// Run migration to add country_code columns if missing
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE active_sessions ADD COLUMN IF NOT EXISTS country_code VARCHAR(10) DEFAULT 'ID'`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE active_visitors ADD COLUMN IF NOT EXISTS country_code VARCHAR(10) DEFAULT 'ID'`)
}

// keepalive — authenticated users ping this every 30s
func keepalive(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		Page string `json:"page"`
	}
	c.BodyParser(&req)

	ua := c.Get("User-Agent")
	ip := c.IP()
	browser, os := parseUA(ua)
	page := req.Page
	if page == "" {
		page = "/"
	}

	cfCountry := c.Get("CF-IPCountry", "")
	country := resolveCountry(ip, cfCountry)

	var prevPage string
	_ = db.QueryRowContext(context.Background(), "SELECT last_page FROM active_sessions WHERE user_id = $1", userID).Scan(&prevPage)

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO active_sessions (user_id, ip_address, user_agent, last_page, last_active, browser, os, country_code)
		 VALUES ($1, $2, $3, $4, NOW(), $5, $6, $7)
		 ON CONFLICT (user_id) DO UPDATE SET ip_address=$2, user_agent=$3, last_page=$4, last_active=NOW(), browser=$5, os=$6, country_code=$7`,
		userID, ip, truncateStr(ua, 500), page, browser, os, country,
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed"})
	}

	if prevPage != page {
		_, _ = db.ExecContext(context.Background(),
			`INSERT INTO visitor_logs (ip_address, country_code, page_url, user_agent, browser, os, visited_at)
			 VALUES ($1, $2, $3, $4, $5, $6, NOW())`,
			ip, country, truncateStr(page, 500), truncateStr(ua, 500), browser, os,
		)
	}

	return c.JSON(fiber.Map{"ok": true})
}

// getLiveUsers — admin only, returns users active in last 5 minutes + anonymous visitors
func getLiveUsers(c *fiber.Ctx) error {
	// Authenticated users
	rows, err := db.QueryContext(context.Background(),
		`SELECT s.user_id, u.full_name, u.email, u.role, s.ip_address, s.last_page, s.browser, s.os, s.last_active, COALESCE(s.country_code, 'ID')
		 FROM active_sessions s
		 JOIN users u ON u.id = s.user_id
		 WHERE s.last_active > NOW() - INTERVAL '5 minutes'
		 ORDER BY s.last_active DESC`)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch"})
	}
	defer rows.Close()

	users := []fiber.Map{}
	for rows.Next() {
		var uid, name, email, role, ip, page, browser, os, country string
		var lastActive time.Time
		if rows.Scan(&uid, &name, &email, &role, &ip, &page, &browser, &os, &lastActive, &country) == nil {
			users = append(users, fiber.Map{
				"user_id": uid, "full_name": name, "email": email, "role": role,
				"ip_address": ip, "last_page": page, "browser": browser, "os": os,
				"last_active": lastActive, "country_code": country,
			})
		}
	}

	// Anonymous visitors
	vRows, _ := db.QueryContext(context.Background(),
		`SELECT ip_address, last_page, browser, os, last_active, COALESCE(country_code, 'ID')
		 FROM active_visitors
		 WHERE last_active > NOW() - INTERVAL '5 minutes'
		 AND ip_address NOT IN (SELECT ip_address FROM active_sessions WHERE last_active > NOW() - INTERVAL '5 minutes')
		 ORDER BY last_active DESC`)
	visitors := []fiber.Map{}
	if vRows != nil {
		defer vRows.Close()
		for vRows.Next() {
			var ip, page, browser, os, country string
			var lastActive time.Time
			if vRows.Scan(&ip, &page, &browser, &os, &lastActive, &country) == nil {
				visitors = append(visitors, fiber.Map{
					"ip_address": ip, "last_page": page, "browser": browser, "os": os,
					"last_active": lastActive, "country_code": country,
				})
			}
		}
	}

	// Country visitor counts
	vcRows, err := db.QueryContext(context.Background(),
		`SELECT country_code, count FROM visitor_country_counts ORDER BY count DESC, country_code ASC`)
	visitorCountries := []fiber.Map{}
	if err == nil {
		defer vcRows.Close()
		for vcRows.Next() {
			var code string
			var count int
			if vcRows.Scan(&code, &count) == nil {
				visitorCountries = append(visitorCountries, fiber.Map{
					"country_code": code,
					"count":        count,
				})
			}
		}
	}

	return c.JSON(fiber.Map{
		"count":             len(users),
		"users":             users,
		"visitor_count":     len(visitors),
		"visitors":          visitors,
		"visitor_countries": visitorCountries,
	})
}

// cleanup old sessions (run periodically)
func startSessionCleanup() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			db.ExecContext(context.Background(), "DELETE FROM active_sessions WHERE last_active < NOW() - INTERVAL '10 minutes'")
			db.ExecContext(context.Background(), "DELETE FROM active_visitors WHERE last_active < NOW() - INTERVAL '10 minutes'")
		}
	}()
}

// visitorPing — public, tracks anonymous visitors by IP
func visitorPing(c *fiber.Ctx) error {
	var req struct {
		Page string `json:"page"`
	}
	c.BodyParser(&req)

	ua := c.Get("User-Agent")
	ip := c.IP()
	browser, os := parseUA(ua)
	page := req.Page
	if page == "" {
		page = "/"
	}

	cfCountry := c.Get("CF-IPCountry", "")
	country := resolveCountry(ip, cfCountry)

	_, _ = db.ExecContext(context.Background(),
		`INSERT INTO active_visitors (ip_address, user_agent, last_page, last_active, browser, os, country_code)
		 VALUES ($1, $2, $3, NOW(), $4, $5, $6)
		 ON CONFLICT (ip_address) DO UPDATE SET user_agent=$2, last_page=$3, last_active=NOW(), browser=$4, os=$5, country_code=$6`,
		ip, truncateStr(ua, 500), page, browser, os, country,
	)

	_, _ = db.ExecContext(context.Background(),
		`INSERT INTO visitor_logs (ip_address, country_code, page_url, user_agent, browser, os, visited_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NOW())`,
		ip, country, truncateStr(page, 500), truncateStr(ua, 500), browser, os,
	)

	return c.JSON(fiber.Map{"ok": true})
}

// Simple UA parser
func parseUA(ua string) (browser, os string) {
	ua = strings.ToLower(ua)
	// Browser
	switch {
	case strings.Contains(ua, "edg"):
		browser = "Edge"
	case strings.Contains(ua, "chrome") && !strings.Contains(ua, "edg"):
		browser = "Chrome"
	case strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome"):
		browser = "Safari"
	case strings.Contains(ua, "firefox"):
		browser = "Firefox"
	case strings.Contains(ua, "opera") || strings.Contains(ua, "opr"):
		browser = "Opera"
	default:
		browser = "Other"
	}
	// OS
	switch {
	case strings.Contains(ua, "windows"):
		os = "Windows"
	case strings.Contains(ua, "mac os") || strings.Contains(ua, "macintosh"):
		os = "macOS"
	case strings.Contains(ua, "linux") && !strings.Contains(ua, "android"):
		os = "Linux"
	case strings.Contains(ua, "android"):
		os = "Android"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad"):
		os = "iOS"
	default:
		os = "Other"
	}
	return
}

func truncateStr(s string, max int) string {
	if len(s) > max {
		return s[:max]
	}
	return s
}

