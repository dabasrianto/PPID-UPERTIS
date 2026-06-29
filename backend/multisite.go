package main

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// ─── Multi-Site Helpers ──────────────────────────────────────────────────────

// getSiteID extracts site identifier from query param or X-Site-ID header
// Returns empty string for main site (no filtering)
func getSiteID(c *fiber.Ctx) string {
	site := c.Query("site", "")
	if site == "" {
		site = c.Get("X-Site-ID", "")
	}
	if site == "main" || site == "" {
		return ""
	}
	return site
}

// getFacultyIDFromSite resolves a site slug to a faculty UUID
// Returns empty string if not found or if site is main
func getFacultyIDFromSite(siteSlug string) string {
	if siteSlug == "" || siteSlug == "main" {
		return ""
	}
	var facultyID string
	err := db.QueryRowContext(context.Background(),
		"SELECT id FROM faculties WHERE slug = $1", siteSlug).Scan(&facultyID)
	if err != nil {
		return ""
	}
	return facultyID
}

// buildFacultyFilter appends a WHERE clause for faculty_id filtering
// Returns the updated query, args, and argIdx
// Logic: show content that belongs to this faculty OR is shared (faculty_id IS NULL)
func buildFacultyFilter(query string, args []interface{}, argIdx int, siteSlug string) (string, []interface{}, int) {
	if siteSlug == "" || siteSlug == "ppid" {
		return query, args, argIdx
	}
	facultyID := getFacultyIDFromSite(siteSlug)
	if facultyID == "" {
		return query, args, argIdx
	}
	query += fmt.Sprintf(" AND (faculty_id = $%d OR faculty_id IS NULL)", argIdx)
	args = append(args, facultyID)
	argIdx++
	return query, args, argIdx
}

// getFacultyIDsForUser returns the faculty IDs assigned to a faculty_admin user
func getFacultyIDsForUser(userID string) []string {
	rows, err := db.QueryContext(context.Background(),
		"SELECT faculty_id FROM faculty_admins WHERE user_id = $1", userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if rows.Scan(&id) == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// buildAdminFacultyFilter restricts data access for faculty_admin users
// Super admin (role=admin) sees everything, faculty_admin sees only their faculty data
func buildAdminFacultyFilter(c *fiber.Ctx, query string, args []interface{}, argIdx int) (string, []interface{}, int) {
	role, _ := c.Locals("role").(string)
	if role == "admin" {
		// Super admin sees all
		return query, args, argIdx
	}

	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return query, args, argIdx
	}

	facultyIDs := getFacultyIDsForUser(userID)
	if len(facultyIDs) == 0 {
		// Faculty admin with no assigned faculties — show nothing
		query += " AND 1=0"
		return query, args, argIdx
	}

	// Show data belonging to their faculties + shared (NULL)
	placeholders := ""
	for i, fid := range facultyIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += fmt.Sprintf("$%d", argIdx)
		args = append(args, fid)
		argIdx++
	}
	query += fmt.Sprintf(" AND (faculty_id IN (%s) OR faculty_id IS NULL)", placeholders)
	return query, args, argIdx
}

// ensureMultisiteColumns adds faculty_id columns if missing (safe migration)
func ensureMultisiteColumns() {
	migrations := []string{
		`ALTER TABLE posts ADD COLUMN IF NOT EXISTS faculty_id UUID REFERENCES faculties(id) ON DELETE SET NULL`,
		`ALTER TABLE campus_events ADD COLUMN IF NOT EXISTS faculty_id UUID REFERENCES faculties(id) ON DELETE SET NULL`,
		`ALTER TABLE news ADD COLUMN IF NOT EXISTS faculty_id UUID REFERENCES faculties(id) ON DELETE SET NULL`,
		`CREATE INDEX IF NOT EXISTS idx_posts_faculty_id ON posts(faculty_id)`,
		`CREATE INDEX IF NOT EXISTS idx_events_faculty_id ON campus_events(faculty_id)`,
		`CREATE INDEX IF NOT EXISTS idx_news_faculty_id ON news(faculty_id)`,
		`ALTER TABLE hero_slides ADD COLUMN IF NOT EXISTS video_url TEXT`,
	}
	for _, m := range migrations {
		if _, err := db.ExecContext(context.Background(), m); err != nil {
			// Get safe display of migration string
			displayStr := m
			if len(displayStr) > 50 {
				displayStr = displayStr[:50]
			}
			zlog.Warn().Err(err).Str("migration", displayStr).Msg("Multisite migration warning (may already exist)")
		}
	}
	zlog.Info().Msg("Multisite columns ensured")
}
