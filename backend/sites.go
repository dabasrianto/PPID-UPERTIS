package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ─── Sites Table Migration ───────────────────────────────────────────────────

func migrateSitesTable() {
	db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS sites (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			slug VARCHAR(100) UNIQUE NOT NULL,
			domain VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			faculty_id UUID REFERENCES faculties(id) ON DELETE SET NULL,
			logo_url VARCHAR(500),
			accent_color VARCHAR(50) DEFAULT 'blue',
			menu JSONB DEFAULT '[]'::jsonb,
			settings JSONB DEFAULT '{}'::jsonb,
			active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	db.ExecContext(context.Background(), `ALTER TABLE sites ADD COLUMN IF NOT EXISTS favicon_url VARCHAR(500)`)
	db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_sites_domain ON sites(domain)`)
	db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_sites_slug ON sites(slug)`)
	zlog.Info().Msg("Sites table ensured")
}

// ─── Public: Get Site Config by Hostname ─────────────────────────────────────

func getSiteConfig(c *fiber.Ctx) error {
	host := c.Query("host", "")
	if host == "" {
		host = c.Get("X-Forwarded-Host", c.Hostname())
	}
	// Strip port if present
	if idx := strings.Index(host, ":"); idx > 0 {
		host = host[:idx]
	}

	var id, slug, domain, name, description, logoURL, faviconURL, accentColor, menuJSON, settingsJSON string
	var facultyID *string
	var active bool

	err := db.QueryRowContext(context.Background(),
		`SELECT id, slug, domain, name, COALESCE(description,''), COALESCE(faculty_id::text,''), 
		 COALESCE(logo_url,''), COALESCE(favicon_url,''), accent_color, COALESCE(menu::text,'[]'), COALESCE(settings::text,'{}'), active
		 FROM sites WHERE domain = $1 AND active = true`,
		host,
	).Scan(&id, &slug, &domain, &name, &description, &facultyID, &logoURL, &faviconURL, &accentColor, &menuJSON, &settingsJSON, &active)

	if err != nil {
		// Fallback: search by slug (subdomain part, e.g. 'ppid' from 'ppid.dabas.web.id')
		parts := strings.Split(host, ".")
		if len(parts) > 1 {
			subdomain := parts[0]
			err = db.QueryRowContext(context.Background(),
				`SELECT id, slug, domain, name, COALESCE(description,''), COALESCE(faculty_id::text,''), 
				 COALESCE(logo_url,''), COALESCE(favicon_url,''), accent_color, COALESCE(menu::text,'[]'), COALESCE(settings::text,'{}'), active
				 FROM sites WHERE slug = $1 AND active = true`,
				subdomain,
			).Scan(&id, &slug, &domain, &name, &description, &facultyID, &logoURL, &faviconURL, &accentColor, &menuJSON, &settingsJSON, &active)
		}
	}

	if err != nil {
		// Not a subdomain site — return main site config
		return c.JSON(fiber.Map{
			"is_main":      true,
			"slug":         "main",
			"name":         "",
			"accent_color": "blue",
		})
	}

	var menu interface{}
	json.Unmarshal([]byte(menuJSON), &menu)
	var settings interface{}
	json.Unmarshal([]byte(settingsJSON), &settings)

	return c.JSON(fiber.Map{
		"is_main":      false,
		"id":           id,
		"slug":         slug,
		"domain":       domain,
		"name":         name,
		"description":  description,
		"faculty_id":   facultyID,
		"logo_url":     logoURL,
		"favicon_url":  faviconURL,
		"accent_color": accentColor,
		"menu":         menu,
		"settings":     settings,
		"active":       active,
	})
}

// ─── Admin: CRUD Sites ───────────────────────────────────────────────────────

func getAdminSites(c *fiber.Ctx) error {
	rows, err := db.QueryContext(context.Background(),
		`SELECT id, slug, domain, name, COALESCE(description,''), COALESCE(faculty_id::text,''),
		 COALESCE(logo_url,''), COALESCE(favicon_url,''), accent_color, COALESCE(menu::text,'[]'), COALESCE(settings::text,'{}'), active, created_at
		 FROM sites ORDER BY created_at`)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch sites"})
	}
	defer rows.Close()

	sites := []fiber.Map{}
	for rows.Next() {
		var id, slug, domain, name, description, facultyID, logoURL, faviconURL, accentColor, menuStr, settingsStr, createdAt string
		var active bool
		if rows.Scan(&id, &slug, &domain, &name, &description, &facultyID, &logoURL, &faviconURL, &accentColor, &menuStr, &settingsStr, &active, &createdAt) == nil {
			var menu, settings interface{}
			json.Unmarshal([]byte(menuStr), &menu)
			json.Unmarshal([]byte(settingsStr), &settings)
			sites = append(sites, fiber.Map{
				"id": id, "slug": slug, "domain": domain, "name": name,
				"description": description, "faculty_id": facultyID, "logo_url": logoURL,
				"favicon_url": faviconURL,
				"accent_color": accentColor, "menu": menu, "settings": settings,
				"active": active, "created_at": createdAt,
			})
		}
	}
	return c.JSON(sites)
}

func createSite(c *fiber.Ctx) error {
	var req struct {
		Slug        string      `json:"slug"`
		Domain      string      `json:"domain"`
		Name        string      `json:"name"`
		Description string      `json:"description"`
		FacultyID   string      `json:"faculty_id"`
		LogoURL     string      `json:"logo_url"`
		FaviconURL  string      `json:"favicon_url"`
		AccentColor string      `json:"accent_color"`
		Menu        interface{} `json:"menu"`
		Settings    interface{} `json:"settings"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if req.Slug == "" || req.Domain == "" || req.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "slug, domain, dan name wajib diisi"})
	}
	if req.AccentColor == "" {
		req.AccentColor = "blue"
	}

	menuJSON, _ := json.Marshal(req.Menu)
	settingsJSON, _ := json.Marshal(req.Settings)

	var facultyParam interface{} = nil
	if req.FacultyID != "" {
		facultyParam = req.FacultyID
	}

	var id string
	err := db.QueryRowContext(context.Background(),
		`INSERT INTO sites (slug, domain, name, description, faculty_id, logo_url, favicon_url, accent_color, menu, settings)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		req.Slug, req.Domain, req.Name, req.Description, facultyParam, req.LogoURL, req.FaviconURL, req.AccentColor,
		string(menuJSON), string(settingsJSON),
	).Scan(&id)

	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Slug atau domain sudah digunakan"})
		}
		zlog.Error().Err(err).Msg("Create site failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat site"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"id": id, "message": "Site created"})
}

func updateSite(c *fiber.Ctx) error {
	id := c.Params("id")
	
	// Resolve slug to UUID if needed
	var siteID string
	if len(id) == 36 && strings.Contains(id, "-") {
		siteID = id
	} else {
		_ = db.QueryRowContext(context.Background(), "SELECT id FROM sites WHERE slug = $1", id).Scan(&siteID)
	}
	if siteID == "" {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Site tidak ditemukan"})
	}

	var req struct {
		Slug        string      `json:"slug"`
		Domain      string      `json:"domain"`
		Name        string      `json:"name"`
		Description string      `json:"description"`
		FacultyID   string      `json:"faculty_id"`
		LogoURL     string      `json:"logo_url"`
		FaviconURL  string      `json:"favicon_url"`
		AccentColor string      `json:"accent_color"`
		Menu        interface{} `json:"menu"`
		Settings    interface{} `json:"settings"`
		Active      *bool       `json:"active"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	role, _ := c.Locals("role").(string)
	if role == "ppid_admin" {
		var targetSlug string
		err := db.QueryRowContext(context.Background(), "SELECT slug FROM sites WHERE id = $1", siteID).Scan(&targetSlug)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to verify site ID"})
		}
		if targetSlug != "ppid" {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Forbidden: Admin PPID hanya boleh mengedit site PPID"})
		}
		// Enforce original slug and domain for PPID site to avoid breaking multi-site routes
		var existingSlug, existingDomain string
		_ = db.QueryRowContext(context.Background(), "SELECT slug, domain FROM sites WHERE id = $1", siteID).Scan(&existingSlug, &existingDomain)
		req.Slug = existingSlug
		req.Domain = existingDomain
		activeBool := true
		req.Active = &activeBool
	}

	menuJSON, _ := json.Marshal(req.Menu)
	settingsJSON, _ := json.Marshal(req.Settings)

	var facultyParam interface{} = nil
	if req.FacultyID != "" {
		facultyParam = req.FacultyID
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	// If domain is empty (e.g. sent by frontend using site-config which lacks domain),
	// preserve the existing domain to avoid unique constraint violation.
	if req.Domain == "" {
		_ = db.QueryRowContext(context.Background(), "SELECT domain FROM sites WHERE id = $1", siteID).Scan(&req.Domain)
	}
	// If slug is empty, preserve existing slug
	if req.Slug == "" {
		_ = db.QueryRowContext(context.Background(), "SELECT slug FROM sites WHERE id = $1", siteID).Scan(&req.Slug)
	}

	_, err := db.ExecContext(context.Background(),
		`UPDATE sites SET slug=$1, domain=$2, name=$3, description=$4, faculty_id=$5, 
		 logo_url=$6, favicon_url=$7, accent_color=$8, menu=$9, settings=$10, active=$11, updated_at=NOW()
		 WHERE id=$12`,
		req.Slug, req.Domain, req.Name, req.Description, facultyParam,
		req.LogoURL, req.FaviconURL, req.AccentColor, string(menuJSON), string(settingsJSON), active, siteID,
	)
	if err != nil {
		zlog.Error().Err(err).Str("site_id", siteID).Str("slug", req.Slug).Str("name", req.Name).Msg("updateSite SQL error")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal update site"})
	}
	return c.JSON(fiber.Map{"message": "Site updated"})
}

func deleteSite(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.ExecContext(context.Background(), "DELETE FROM sites WHERE id = $1", id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal hapus site"})
	}
	return c.JSON(fiber.Map{"message": "Site deleted"})
}
