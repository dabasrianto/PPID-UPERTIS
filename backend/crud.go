package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

var tableColumns = map[string]map[string]bool{
	"popup_banners": {
		"title": true, "description": true, "image_url": true, "link_text": true,
		"link_url": true, "start_date": true, "end_date": true, "active": true, "sort_order": true,
		"show_image_only": true,
	},
	"hero_slides": {
		"title": true, "eyebrow": true, "description": true, "image_url": true,
		"cta_primary_label": true, "cta_primary_href": true, "cta_secondary_label": true,
		"cta_secondary_href": true, "video_url": true, "active": true, "sort_order": true,
	},
	"faculties": {
		"code": true, "name": true, "slug": true, "description": true, "vision": true,
		"mission": true, "about_content": true, "hero_title": true, "hero_eyebrow": true,
		"hero_description": true, "cover_image_url": true, "accent": true, "active": true,
		"programs": true, "sort_order": true, "facilities": true, "contact_info": true,
		"tujuan": true, "struktur_organisasi_url": true, "kerjasama": true, "about_image_url": true, "about_images": true, "facility_images": true, "cover_images": true,
	},
	"faculty_programs": {
		"faculty_id": true, "name": true, "slug": true, "level": true, "description": true,
		"accreditation": true, "duration_years": true, "active": true, "sort_order": true,
		"syllabus": true, "career_paths": true, "struktur_organisasi": true,
		"visi": true, "misi": true, "tujuan": true, "gelar_akademik": true, "kompetensi_lulusan": true, "fasilitas_laboratorium": true,
		"fasilitas_laboratorium_image": true, "fasilitas_laboratorium_images": true, "cover_image_url": true, "card_bg_color": true,
	},
	"faculty_lecturers": {
		"faculty_id": true, "name": true, "gelar": true, "position": true, "education": true,
		"expertise": true, "photo_url": true, "leadership_group": true, "active": true, "sort_order": true,
	},
	"news": {
		"title": true, "slug": true, "excerpt": true, "content": true, "image_url": true, "category": true,
		"date": true, "featured": true, "active": true, "author_name": true, "seo_title": true, "seo_description": true,
		"faculty_id": true,
	},
	"testimonials": {
		"name": true, "role": true, "quote": true, "image_url": true, "year": true,
		"video_url": true, "youtube_url": true, "tiktok_url": true, "instagram_url": true,
		"active": true, "sort_order": true,
	},
	"pmb_batches": {
		"name": true, "academic_year": true, "start_date": true, "end_date": true,
		"registration_fee": true, "is_active": true,
	},
	"pmb_candidates": {
		"user_id": true, "batch_id": true, "registration_number": true, "full_name": true,
		"nisn": true, "phone_whatsapp": true, "school_origin": true, 
		"first_choice_program_id": true, "second_choice_program_id": true, "status": true,
	},
	"blog_posts": {
		"title": true, "slug": true, "excerpt": true, "content": true, "cover_image_url": true,
		"author_name": true, "category": true, "tags": true, "faculty_id": true, "status": true,
		"published_at": true, "seo_title": true, "seo_description": true,
	},
	"posts": {
		"post_type": true, "title": true, "slug": true, "excerpt": true, "content": true,
		"featured_image": true, "category": true, "tags": true, "author": true, "status": true,
		"featured": true, "published_at": true, "seo_title": true, "seo_description": true,
		"scheduled_at": true, "faculty_id": true,
	},
	"pages": {
		"title": true, "subtitle": true, "slug": true, "content": true,
		"cover_image_url": true, "published": true, "sort_order": true,
		"seo_title": true, "seo_description": true,
	},
	"campus_events": {
		"title": true, "slug": true, "description": true, "image_url": true, "event_date": true,
		"start_time": true, "end_time": true, "location": true, "map_coordinates": true, "active": true, "sort_order": true,
		"category": true, "registration_url": true, "is_internal_registration": true, "capacity": true, "speakers": true, "wa_message_template": true,
		"faculty_id": true, "images": true,
	},
	"event_registrations": {
		"event_id": true, "full_name": true, "email": true, "whatsapp": true, "status": true,
	},
	"event_categories": {
		"name": true, "type": true, "active": true,
	},
	"downloads": {
		"title": true, "description": true, "file_url": true, "category": true, "active": true,
	},
	"faqs": {
		"question": true, "answer": true, "category": true, "sort_order": true, "active": true,
	},
	"gallery": {
		"title": true, "description": true, "media_url": true, "media_type": true, "active": true,
	},
	"permohonan_informasi": {
		"applicant_type": true, "name": true, "identity_number": true, "email": true, "phone": true,
		"address": true, "details": true, "purpose": true, "obtain_method": true, "delivery_method": true,
		"attachment_url": true, "status": true, "admin_response": true,
	},
}

// Tables that have an updated_at column
var tablesWithUpdatedAt = map[string]bool{
	"popup_banners":    true,
	"hero_slides":      true,
	"faculties":        true,
	"faculty_programs": true,
	"faculty_lecturers": true,
	"news":             true,
	"testimonials":     true,
	"pmb_batches":      true,
	"pmb_candidates":   true,
	"blog_posts":       true,
	"posts":            true,
	"pages":            true,
	"campus_events":    true,
	"downloads":        true,
	"faqs":             true,
	"gallery":          true,
	"permohonan_informasi": true,
}

func normalizeValue(column string, value interface{}) (interface{}, error) {
	if value == nil || value == "" {
		if column == "faculty_id" || strings.HasSuffix(column, "_at") {
			return nil, nil
		}
		if column == "capacity" || column == "sort_order" {
			return 0, nil
		}
	}

	switch column {
	case "facilities", "contact_info", "value", "speakers", "syllabus", "career_paths", "struktur_organisasi", "kerjasama", "fasilitas_laboratorium_images", "about_images", "facility_images", "cover_images", "images":
		// If value is already raw JSON bytes, use as-is
		if raw, ok := value.(*json.RawMessage); ok && raw != nil {
			return string(*raw), nil
		}
		if raw, ok := value.(json.RawMessage); ok {
			return string(raw), nil
		}
		b, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		return string(b), nil
	case "tags":
		items := []string{}
		if arr, ok := value.([]interface{}); ok {
			for _, item := range arr {
				items = append(items, fmt.Sprint(item))
			}
		}
		return pq.Array(items), nil
	default:
		return value, nil
	}
}

func cleanPayload(table string, payload map[string]interface{}) (map[string]interface{}, error) {
	allowed, ok := tableColumns[table]
	if !ok {
		return nil, fmt.Errorf("table is not allowed")
	}

	cleaned := map[string]interface{}{}
	for key, value := range payload {
		if !allowed[key] {
			continue
		}
		normalized, err := normalizeValue(key, value)
		if err != nil {
			return nil, err
		}
		cleaned[key] = normalized
	}

	if table == "faculty_programs" {
		slugVal, hasSlug := cleaned["slug"]
		nameVal, hasName := cleaned["name"]
		levelVal, hasLevel := cleaned["level"]
		
		var rawSlug string
		if hasSlug && slugVal != nil {
			rawSlug, _ = slugVal.(string)
		}
		
		if rawSlug == "" && hasName && nameVal != nil {
			if nameStr, ok := nameVal.(string); ok && nameStr != "" {
				levelStr := ""
				if hasLevel && levelVal != nil {
					if l, ok := levelVal.(string); ok {
						levelStr = l
					}
				}
				combined := nameStr
				if levelStr != "" && !strings.HasPrefix(strings.ToLower(strings.TrimSpace(nameStr)), strings.ToLower(strings.TrimSpace(levelStr))) {
					combined = levelStr + "-" + nameStr
				}
				cleaned["slug"] = generateSlug(combined)
			}
		} else if rawSlug != "" {
			cleaned["slug"] = generateSlug(rawSlug)
		}
	}

	return cleaned, nil
}

func getRows(c *fiber.Ctx) error {
	table := c.Params("table")

	// Ensure table is whitelisted
	if _, ok := tableColumns[table]; !ok {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "table not allowed"})
	}

	// Tables that support faculty_id filtering
	facultyFilteredTables := map[string]bool{
		"posts": true, "campus_events": true, "news": true, "blog_posts": true,
		"faculty_programs": true, "faculty_lecturers": true,
	}

	query := fmt.Sprintf("SELECT COALESCE(json_agg(row_to_json(t)), '[]') FROM %s t", table)

	// Apply faculty filter for faculty_admin on supported tables
	role, _ := c.Locals("role").(string)
	args := []interface{}{}
	if role == "faculty_admin" && facultyFilteredTables[table] {
		userID, _ := c.Locals("user_id").(string)
		facultyIDs := getFacultyIDsForUser(userID)
		if len(facultyIDs) == 0 {
			// No faculties assigned — return empty
			c.Set("Content-Type", "application/json")
			return c.Send([]byte("[]"))
		}
		placeholders := ""
		for i, fid := range facultyIDs {
			if i > 0 {
				placeholders += ","
			}
			placeholders += fmt.Sprintf("$%d", i+1)
			args = append(args, fid)
		}
		query += fmt.Sprintf(" WHERE (t.faculty_id IN (%s) OR t.faculty_id IS NULL)", placeholders)
	}

	var result []byte
	var err error
	if len(args) > 0 {
		err = db.QueryRowContext(context.Background(), query, args...).Scan(&result)
	} else {
		err = db.QueryRowContext(context.Background(), query).Scan(&result)
	}
	if err != nil {
		log.Println("getRows error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch rows"})
	}

	c.Set("Content-Type", "application/json")
	return c.Send(result)
}

func checkFacultyWriteAccess(c *fiber.Ctx, table string, recordID string, payload map[string]interface{}) (bool, error) {
	role, _ := c.Locals("role").(string)
	if role == "admin" {
		// Super admin has full write/delete permissions on all tables
		return true, nil
	}

	if role == "ppid_admin" {
		// PPID Admin can write/delete from downloads, faqs, gallery, pages, posts, and permohonan_informasi tables
		if table == "downloads" || table == "faqs" || table == "gallery" || table == "pages" || table == "posts" || table == "permohonan_informasi" {
			return true, nil
		}
		return false, nil
	}

	if role != "faculty_admin" {
		// Only admin and faculty_admin can write/delete using generic CRUD endpoints
		return false, nil
	}

	// For faculty_admin, check table whitelist
	facultyFilteredTables := map[string]bool{
		"faculty_programs":  true,
		"faculty_lecturers": true,
		"posts":             true,
		"news":              true,
		"campus_events":     true,
		"blog_posts":        true,
		"faculties":         true, // only update allowed, checked below
	}

	if !facultyFilteredTables[table] {
		// Writing to other tables is not allowed for faculty_admin
		return false, nil
	}

	userID, _ := c.Locals("user_id").(string)
	facultyIDs := getFacultyIDsForUser(userID)
	if len(facultyIDs) == 0 {
		return false, nil
	}

	// Helper to check if a string is in a slice
	contains := func(slice []string, s string) bool {
		for _, x := range slice {
			if x == s {
				return true
			}
		}
		return false
	}

	// 1. Check if table is "faculties"
	if table == "faculties" {
		// For faculties, the recordID is the faculty_id itself
		if recordID != "" {
			return contains(facultyIDs, recordID), nil
		}
		// Creating/deleting a faculty is not allowed for faculty_admin
		return false, nil
	}

	// 2. If it's a create action (recordID is empty), check payload's faculty_id
	if recordID == "" && payload != nil {
		facultyIDVal, exists := payload["faculty_id"]
		if !exists || facultyIDVal == nil || facultyIDVal == "" {
			// If payload has no faculty_id or it's empty, and the user has exactly 1 faculty, we can auto-inject it
			if len(facultyIDs) == 1 {
				payload["faculty_id"] = facultyIDs[0]
				return true, nil
			}
			// Otherwise they must provide one of their assigned faculty IDs
			return false, nil
		}
		facultyIDStr, ok := facultyIDVal.(string)
		if !ok {
			return false, nil
		}
		return contains(facultyIDs, facultyIDStr), nil
	}

	// 3. For update/delete (recordID is not empty), fetch the existing record's faculty_id from DB
	if recordID != "" {
		var existingFacultyID sql.NullString
		query := fmt.Sprintf("SELECT faculty_id FROM %s WHERE id = $1", table)
		err := db.QueryRowContext(context.Background(), query, recordID).Scan(&existingFacultyID)
		if err != nil {
			if err == sql.ErrNoRows {
				return false, fmt.Errorf("record not found")
			}
			return false, err
		}

		// If the existing record is not assigned to any faculty (NULL), only super admins can modify it
		if !existingFacultyID.Valid {
			return false, nil
		}

		if !contains(facultyIDs, existingFacultyID.String) {
			return false, nil
		}

		// If it's an update, and they are trying to change the faculty_id, make sure the new faculty_id is also owned
		if payload != nil {
			if newFacultyIDVal, exists := payload["faculty_id"]; exists && newFacultyIDVal != nil && newFacultyIDVal != "" {
				newFacultyIDStr, ok := newFacultyIDVal.(string)
				if !ok {
					return false, nil
				}
				if !contains(facultyIDs, newFacultyIDStr) {
					return false, nil
				}
			}
		}

		return true, nil
	}

	return false, nil
}

func createRow(c *fiber.Ctx) error {
	table := c.Params("table")
	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Multi-tenant check
	allowed, err := checkFacultyWriteAccess(c, table, "", payload)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if !allowed {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Forbidden: Hak akses tidak memadai"})
	}

	cleaned, err := cleanPayload(table, payload)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if len(cleaned) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "no valid fields"})
	}

	columns := make([]string, 0, len(cleaned))
	for key := range cleaned {
		columns = append(columns, key)
	}
	sort.Strings(columns)

	placeholders := make([]string, len(columns))
	values := make([]interface{}, len(columns))
	for i, column := range columns {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		values[i] = cleaned[column]
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	var id string
	if err := db.QueryRowContext(context.Background(), query, values...).Scan(&id); err != nil {
		log.Println("createRow error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create record"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"id": id})
}

func updateRow(c *fiber.Ctx) error {
	table := c.Params("table")
	id := c.Params("id")
	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Multi-tenant check
	allowed, err := checkFacultyWriteAccess(c, table, id, payload)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Record not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if !allowed {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Forbidden: Hak akses tidak memadai"})
	}

	cleaned, err := cleanPayload(table, payload)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if len(cleaned) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "no valid fields"})
	}

	columns := make([]string, 0, len(cleaned))
	for key := range cleaned {
		columns = append(columns, key)
	}
	sort.Strings(columns)

	assignments := make([]string, len(columns))
	values := make([]interface{}, len(columns)+1)
	for i, column := range columns {
		assignments[i] = fmt.Sprintf("%s=$%d", column, i+1)
		values[i] = cleaned[column]
	}
	values[len(values)-1] = id

	// Only add updated_at for tables that have the column
	updatedAtClause := ""
	if tablesWithUpdatedAt[table] {
		updatedAtClause = ", updated_at=NOW()"
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s%s WHERE id=$%d",
		table,
		strings.Join(assignments, ", "),
		updatedAtClause,
		len(values),
	)

	if _, err := db.ExecContext(context.Background(), query, values...); err != nil {
		log.Println("updateRow error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update record"})
	}
	return c.JSON(fiber.Map{"message": "Updated successfully"})
}

func deleteRow(c *fiber.Ctx) error {
	table := c.Params("table")
	id := c.Params("id")
	if _, ok := tableColumns[table]; !ok {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "table is not allowed"})
	}

	// Multi-tenant check
	allowed, err := checkFacultyWriteAccess(c, table, id, nil)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Record not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if !allowed {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Forbidden: Hak akses tidak memadai"})
	}

	if _, err := db.ExecContext(context.Background(), fmt.Sprintf("DELETE FROM %s WHERE id=$1", table), id); err != nil {
		log.Println("deleteRow error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete record"})
	}
	return c.JSON(fiber.Map{"message": "Deleted successfully"})
}

func bulkAction(c *fiber.Ctx) error {
	table := c.Params("table")
	if _, ok := tableColumns[table]; !ok {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "table is not allowed"})
	}

	var req struct {
		Action string   `json:"action"` // "delete", "publish", "unpublish"
		IDs    []string `json:"ids"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if len(req.IDs) == 0 {
		return c.JSON(fiber.Map{"message": "No IDs provided"})
	}

	// Multi-tenant check for each ID
	for _, id := range req.IDs {
		allowed, err := checkFacultyWriteAccess(c, table, id, nil)
		if err != nil {
			if err.Error() == "record not found" {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Record %s not found", id)})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if !allowed {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": fmt.Sprintf("Forbidden: Hak akses tidak memadai untuk ID %s", id)})
		}
	}

	var query string
	switch req.Action {
	case "delete":
		query = fmt.Sprintf("DELETE FROM %s WHERE id = ANY($1)", table)
		_, err := db.ExecContext(context.Background(), query, pq.Array(req.IDs))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	case "publish":
		col := "status"
		var val interface{} = "published"
		if table == "pages" {
			col = "published"
			val = true
		}
		query = fmt.Sprintf("UPDATE %s SET %s = $1 WHERE id = ANY($2)", table, col)
		_, err := db.ExecContext(context.Background(), query, val, pq.Array(req.IDs))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	case "unpublish":
		col := "status"
		var val interface{} = "draft"
		if table == "pages" {
			col = "published"
			val = false
		}
		query = fmt.Sprintf("UPDATE %s SET %s = $1 WHERE id = ANY($2)", table, col)
		_, err := db.ExecContext(context.Background(), query, val, pq.Array(req.IDs))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	default:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid action"})
	}

	return c.JSON(fiber.Map{"message": fmt.Sprintf("Bulk %s completed for %d items", req.Action, len(req.IDs))})
}

func scanRows(rows *sql.Rows) ([]fiber.Map, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := []fiber.Map{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		row := fiber.Map{}
		for i, column := range columns {
			switch v := values[i].(type) {
			case []byte:
				var decoded interface{}
				if json.Valid(v) && json.Unmarshal(v, &decoded) == nil {
					row[column] = decoded
				} else {
					row[column] = string(v)
				}
			default:
				row[column] = v
			}
		}
		result = append(result, row)
	}
	return result, rows.Err()
}
