package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Map of common flags as inline SVGs (18x12 pixels)
var flagSVGs = map[string]string{
	"ID": `<rect width="18" height="6" fill="#FF0000"/><rect y="6" width="18" height="6" fill="#FFFFFF"/>`,
	"SG": `<rect width="18" height="6" fill="#E52B50"/><rect y="6" width="18" height="6" fill="#FFFFFF"/><circle cx="4" cy="3" r="1.8" fill="#FFFFFF"/><circle cx="5" cy="3.5" r="1.8" fill="#E52B50"/><circle cx="7" cy="2" r="0.4" fill="#FFFFFF"/><circle cx="8" cy="3" r="0.4" fill="#FFFFFF"/><circle cx="7.5" cy="4.2" r="0.4" fill="#FFFFFF"/><circle cx="6.5" cy="4.2" r="0.4" fill="#FFFFFF"/><circle cx="6" cy="3" r="0.4" fill="#FFFFFF"/>`,
	"MY": `<rect width="18" height="12" fill="#FF0000"/><rect y="1.7" width="18" height="1.7" fill="#FFFFFF"/><rect y="5.1" width="18" height="1.7" fill="#FFFFFF"/><rect y="8.5" width="18" height="1.7" fill="#FFFFFF"/><rect width="9" height="7" fill="#002060"/><circle cx="4.5" cy="3.5" r="2.2" fill="#FFD700"/><circle cx="5.5" cy="3.5" r="2.2" fill="#002060"/><polygon points="5.5,2.5 5.8,3.2 6.5,3.2 6.0,3.7 6.2,4.4 5.5,4.0 4.8,4.4 5.0,3.7 4.5,3.2 5.2,3.2" fill="#FFD700"/>`,
	"US": `<rect width="18" height="12" fill="#B22234"/><rect y="1.7" width="18" height="1.7" fill="#FFFFFF"/><rect y="5.1" width="18" height="1.7" fill="#FFFFFF"/><rect y="8.5" width="18" height="1.7" fill="#FFFFFF"/><rect width="9" height="7" fill="#3C3B6E"/><circle cx="2.5" cy="2.5" r="0.5" fill="#FFFFFF"/><circle cx="6.5" cy="2.5" r="0.5" fill="#FFFFFF"/><circle cx="4.5" cy="4.5" r="0.5" fill="#FFFFFF"/>`,
	"GB": `<rect width="18" height="12" fill="#012169"/><path d="M0,0 L18,12 M18,0 L0,12" stroke="#FFFFFF" stroke-width="2"/><path d="M0,0 L18,12 M18,0 L0,12" stroke="#C8102E" stroke-width="1.2"/><path d="M9,0 L9,12 M0,6 L18,6" stroke="#FFFFFF" stroke-width="3"/><path d="M9,0 L9,12 M0,6 L18,6" stroke="#C8102E" stroke-width="1.8"/>`,
	"AU": `<rect width="18" height="12" fill="#012169"/><rect width="9" height="6" fill="#012169"/><path d="M0,0 L9,6 M9,0 L0,6" stroke="#FFFFFF" stroke-width="1"/><path d="M0,0 L9,6 M9,0 L0,6" stroke="#C8102E" stroke-width="0.6"/><path d="M4.5,0 L4.5,6 M0,3 L9,3" stroke="#FFFFFF" stroke-width="1.5"/><path d="M4.5,0 L4.5,6 M0,3 L9,3" stroke="#C8102E" stroke-width="0.9"/><circle cx="13.5" cy="3" r="0.4" fill="#FFFFFF"/><circle cx="15.5" cy="5" r="0.4" fill="#FFFFFF"/><circle cx="13.5" cy="7" r="0.4" fill="#FFFFFF"/><circle cx="11.5" cy="5" r="0.4" fill="#FFFFFF"/><circle cx="4.5" cy="9" r="0.7" fill="#FFFFFF"/>`,
	"JP": `<rect width="18" height="12" fill="#FFFFFF" stroke="#CCCCCC" stroke-width="0.5"/><circle cx="9" cy="6" r="3.5" fill="#BC002D"/>`,
	"KR": `<rect width="18" height="12" fill="#FFFFFF" stroke="#CCCCCC" stroke-width="0.5"/><circle cx="9" cy="6" r="3" fill="#CD2E3A"/><path d="M 6,6 A 1.5,1.5 0 0,0 9,6 A 1.5,1.5 0 0,1 12,6 A 3,3 0 0,1 6,6 Z" fill="#0047A0"/><rect x="4" y="3" width="1.5" height="1.5" fill="#000000"/><rect x="12.5" y="3" width="1.5" height="1.5" fill="#000000"/><rect x="4" y="7.5" width="1.5" height="1.5" fill="#000000"/><rect x="12.5" y="7.5" width="1.5" height="1.5" fill="#000000"/>`,
	"DE": `<rect width="18" height="4" fill="#000000"/><rect y="4" width="18" height="4" fill="#DD0000"/><rect y="8" width="18" height="4" fill="#FFCE00"/>`,
	"TW": `<rect width="18" height="12" fill="#FE3030"/><rect width="9" height="6" fill="#000095"/><circle cx="4.5" cy="3" r="1.5" fill="#FFFFFF"/>`,
	"CN": `<rect width="18" height="12" fill="#EE1C25"/><polygon points="3,1.5 3.3,2.5 4.3,2.5 3.5,3 3.8,4 3,3.4 2.2,4 2.5,3 1.7,2.5 2.7,2.5" fill="#FFDE00"/><polygon points="5.5,1 5.6,1.4 6,1.4 5.7,1.6 5.8,2 5.5,1.8 5.2,2 5.3,1.6 5,1.4 5.4,1.4" fill="#FFDE00"/><polygon points="6.5,2 6.6,2.4 7,2.4 6.7,2.6 6.8,3 6.5,2.8 6.2,3 6.3,2.6 6,2.4 6.4,2.4" fill="#FFDE00"/>`,
	"GLOBE": `<rect width="18" height="12" fill="#374151" rx="2" stroke="#4B5563" stroke-width="0.5"/><circle cx="9" cy="6" r="4" fill="none" stroke="#FFFFFF" stroke-width="0.8"/><line x1="9" y1="2" x2="9" y2="10" stroke="#FFFFFF" stroke-width="0.5"/><line x1="5" y1="6" x2="13" y2="6" stroke="#FFFFFF" stroke-width="0.5"/>`,
}

type CountryCount struct {
	CountryCode string
	Count       int
}

func getVisitorCounterSVG(c *fiber.Ctx) error {
	// 1. Determine Country Code from Cloudflare header or default
	country := c.Get("CF-IPCountry", "")
	if country == "" {
		country = "ID"
	}
	country = strings.ToUpper(strings.TrimSpace(country))
	if len(country) != 2 {
		country = "ID"
	}

	// 2. Increment count if cookie v_ctr_<country> does not exist (10-minute unique window)
	cookieName := "v_ctr_" + country
	if c.Cookies(cookieName) == "" {
		c.Cookie(&fiber.Cookie{
			Name:     cookieName,
			Value:    "1",
			Expires:  time.Now().Add(10 * time.Minute),
			HTTPOnly: true,
			SameSite: "Lax",
		})

		_, err := db.ExecContext(context.Background(), `
			INSERT INTO visitor_country_counts (country_code, count, updated_at)
			VALUES ($1, 1, NOW())
			ON CONFLICT (country_code) 
			DO UPDATE SET count = visitor_country_counts.count + 1, updated_at = NOW()
		`, country)
		if err != nil {
			log.Println("Warning: Failed to increment visitor count:", err)
		}
	}

	// 3. Query top 9 countries by count
	rows, err := db.QueryContext(context.Background(), `
		SELECT country_code, count 
		FROM visitor_country_counts 
		ORDER BY count DESC, country_code ASC 
		LIMIT 9
	`)
	
	var counts []CountryCount
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cc CountryCount
			if rows.Scan(&cc.CountryCode, &cc.Count) == nil {
				counts = append(counts, cc)
			}
		}
	}

	// Make sure we have at least 1 count to render
	if len(counts) == 0 {
		counts = append(counts, CountryCount{CountryCode: "ID", Count: 1})
	}

	// 4. Generate SVG XML
	width := 240
	height := 75
	
	// Start SVG template
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, width, height, width, height))
	// Background with a nice navy glassmorphic style
	sb.WriteString(`<rect width="100%" height="100%" rx="8" fill="#111827" fill-opacity="0.85" stroke="#1B3A6B" stroke-width="1"/>`)
	
	// Layout top 9 flags in a grid of 3 columns x 3 rows
	xOffsets := []int{12, 88, 164}
	yOffsets := []int{12, 32, 52}
	
	for idx, cc := range counts {
		col := idx % 3
		row := idx / 3
		
		x := xOffsets[col]
		y := yOffsets[row]
		
		// Render Flag SVG
		flagSVG := getFlagSVG(cc.CountryCode)
		
		// Wrap flag inside a group to position it
		sb.WriteString(fmt.Sprintf(`<g transform="translate(%d, %d)">`, x, y))
		sb.WriteString(flagSVG)
		sb.WriteString(`</g>`)
		
		// Render Country Code
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="#9CA3AF" font-family="system-ui, -apple-system, sans-serif" font-size="10" font-weight="bold" dominant-baseline="middle">%s</text>`, x+22, y+6, cc.CountryCode))
		
		// Format count: e.g. 1.2k if >= 1000
		countStr := fmt.Sprintf("%d", cc.Count)
		if cc.Count >= 1000 {
			countStr = fmt.Sprintf("%.1fk", float64(cc.Count)/1000.0)
			countStr = strings.Replace(countStr, ".0k", "k", 1)
		}
		
		// Render Count Value
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="#FFFFFF" font-family="system-ui, -apple-system, sans-serif" font-size="10" font-weight="bold" dominant-baseline="middle">%s</text>`, x+38, y+6, countStr))
	}
	
	sb.WriteString(`</svg>`)

	c.Set("Content-Type", "image/svg+xml")
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	return c.SendString(sb.String())
}

var (
	flagCache   = make(map[string]string)
	flagCacheMu sync.RWMutex
)

func getFlagSVG(countryCode string) string {
	countryCode = strings.ToUpper(strings.TrimSpace(countryCode))
	
	// 1. Check hardcoded SVGs first
	if svg, found := flagSVGs[countryCode]; found {
		return svg
	}
	
	// 2. Check dynamic cache
	flagCacheMu.RLock()
	cached, found := flagCache[countryCode]
	flagCacheMu.RUnlock()
	if found {
		return cached
	}
	
	// 3. Fetch from flagcdn.com on the fly with a 1-second timeout
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	
	url := fmt.Sprintf("https://flagcdn.com/%s.svg", strings.ToLower(countryCode))
	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		// Fallback to GLOBE (don't cache failure so we can retry later)
		return flagSVGs["GLOBE"]
	}
	defer resp.Body.Close()
	
	svgBytes, err := io.ReadAll(resp.Body)
	if err != nil || len(svgBytes) == 0 {
		return flagSVGs["GLOBE"]
	}
	
	// Convert to Base64 data URI image tag
	b64 := base64.StdEncoding.EncodeToString(svgBytes)
	imgTag := fmt.Sprintf(`<image width="18" height="12" href="data:image/svg+xml;base64,%s"/>`, b64)
	
	// Save to cache
	flagCacheMu.Lock()
	flagCache[countryCode] = imgTag
	flagCacheMu.Unlock()
	
	return imgTag
}
