package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Lecturer struct {
	Reg        string `json:"reg"`
	Nama       string `json:"nama"`
	NIP        string `json:"nip"`
	Gelar      string `json:"gelar"`
	Jabatan    string `json:"jabatan"`
	Pendidikan string `json:"pendidikan"`
	Prodi      string `json:"prodi"`
}

// Faculty mapping
// Faculty mapping to code
var prodiToFacultyCode = map[string]string{
	"Farmasi (S1)":                            "FF",
	"Pendidikan Profesi Apoteker (Profesi)":    "FF",
	"Bisnis Digital (S1)":                      "FEB",
	"Ilmu Komunikasi (S1)":                     "FEB",
	"Gizi (D3)":                                "FK",
	"Gizi (S1)":                                "FK",
	"Kebidanan (D3)":                           "FK",
	"Keperawatan (D3)":                         "FK",
	"Keperawatan (S1)":                         "FK",
	"Pendidikan Profesi Ners (Profesi)":        "FK",
	"Teknologi Laboratorium Medis (D3)":        "FK",
	"Teknologi Laboratorium Medis (D4)":        "FK",
}

func makeSlug(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	s = re.ReplaceAllString(s, "")
	re2 := regexp.MustCompile(`[\s]+`)
	s = re2.ReplaceAllString(s, "-")
	re3 := regexp.MustCompile(`-+`)
	s = re3.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	// Load env
	_ = godotenv.Load(".env")
	_ = godotenv.Load("backend/.env")
	_ = godotenv.Load("../.env")

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "")
	if dbname == "" {
		log.Fatal("DB_NAME environment variable is required")
	}
	sslmode := getEnv("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	fmt.Println("✓ Connected to database")

	// Query actual faculty IDs from DB dynamically
	facultyMap := make(map[string]string)
	dbRows, err := db.Query("SELECT id, code FROM faculties")
	if err != nil {
		log.Fatal("Failed to query faculties: ", err)
	}
	defer dbRows.Close()
	for dbRows.Next() {
		var id, code string
		if err := dbRows.Scan(&id, &code); err != nil {
			log.Fatal(err)
		}
		facultyMap[code] = id
	}

	// Read JSON data
	jsonData, err := os.ReadFile("dosen_data.json")
	if err != nil {
		log.Fatal("Failed to read dosen_data.json: ", err)
	}

	var lecturers []Lecturer
	if err := json.Unmarshal(jsonData, &lecturers); err != nil {
		log.Fatal("Failed to parse JSON: ", err)
	}
	fmt.Printf("✓ Loaded %d lecturers from JSON\n", len(lecturers))

	ctx := context.Background()

	inserted := 0
	updated := 0
	skipped := 0
	slugCounts := make(map[string]int)

	for _, l := range lecturers {
		code, ok := prodiToFacultyCode[l.Prodi]
		if !ok {
			fmt.Printf("⚠ No faculty mapping for: %s (dosen: %s)\n", l.Prodi, l.Nama)
			skipped++
			continue
		}
		facultyID := facultyMap[code]
		if facultyID == "" && code == "FEB" {
			facultyID = facultyMap["FEBIS"]
		}
		if facultyID == "" && code == "FEBIS" {
			facultyID = facultyMap["FEB"]
		}
		if facultyID == "" {
			fmt.Printf("⚠ Faculty code %s (or alternative) not found in database (dosen: %s)\n", code, l.Nama)
			skipped++
			continue
		}

		// Build education string
		eduParts := []string{}
		if l.Pendidikan != "" {
			eduParts = append(eduParts, l.Pendidikan)
		}
		if l.Gelar != "" {
			eduParts = append(eduParts, l.Gelar)
		}
		education := strings.Join(eduParts, ", ")

		// Format name with degree (gelar)
		fullName := l.Nama
		if l.Gelar != "" {
			fullName = fmt.Sprintf("%s, %s", l.Nama, l.Gelar)
		}

		// Build slug
		slugBase := makeSlug(l.Nama)
		slug := slugBase
		if count, exists := slugCounts[slugBase]; exists {
			slugCounts[slugBase] = count + 1
			slug = fmt.Sprintf("%s-%d", slugBase, count+1)
		} else {
			slugCounts[slugBase] = 1
		}

		// Check if lecturer already exists by pddikti_id or name+faculty
		var existingID string
		if l.Reg != "" {
			err = db.QueryRowContext(ctx,
				"SELECT id FROM faculty_lecturers WHERE pddikti_id = $1", l.Reg).Scan(&existingID)
		}
		if existingID == "" {
			// Try by name + faculty (either exact name or name with degrees appended)
			err = db.QueryRowContext(ctx,
				"SELECT id FROM faculty_lecturers WHERE (UPPER(name) = UPPER($1) OR UPPER(name) LIKE UPPER($1) || ',%') AND faculty_id = $2",
				l.Nama, facultyID).Scan(&existingID)
		}

		if existingID != "" {
			// Update existing
			_, err = db.ExecContext(ctx, `
				UPDATE faculty_lecturers SET
					name = $1,
					position = COALESCE(NULLIF($2, ''), position),
					education = COALESCE(NULLIF($3, ''), education),
					expertise = COALESCE(NULLIF($4, ''), expertise),
					pddikti_id = COALESCE(NULLIF($5, ''), pddikti_id),
					source = 'excel_import'
				WHERE id = $6`,
				fullName, l.Jabatan, education, l.Prodi, l.Reg, existingID)
			if err != nil {
				fmt.Printf("⚠ Failed to update %s: %v\n", l.Nama, err)
				skipped++
				continue
			}
			updated++
		} else {
			// Insert new
			var pddiktiID *string
			if l.Reg != "" {
				pddiktiID = &l.Reg
			}

			_, err = db.ExecContext(ctx, `
				INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
				VALUES ($1, $2, $3, $4, $5, $6, 'excel_import', $7, true, 0)`,
				facultyID, fullName, slug, l.Jabatan, education, l.Prodi, pddiktiID)
			if err != nil {
				// If slug conflict, try with different suffix
				slug = slug + "-2"
				_, err = db.ExecContext(ctx, `
					INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
					VALUES ($1, $2, $3, $4, $5, $6, 'excel_import', $7, true, 0)`,
					facultyID, fullName, slug, l.Jabatan, education, l.Prodi, pddiktiID)
				if err != nil {
					fmt.Printf("⚠ Failed to insert %s: %v\n", l.Nama, err)
					skipped++
					continue
				}
			}
			inserted++
		}
	}

	fmt.Println("\n--- Import Summary ---")
	fmt.Printf("✓ Inserted: %d\n", inserted)
	fmt.Printf("✓ Updated:  %d\n", updated)
	fmt.Printf("⚠ Skipped:  %d\n", skipped)
	fmt.Printf("  Total:    %d\n", inserted+updated+skipped)

	// Verify
	var count int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM faculty_lecturers WHERE source = 'excel_import'").Scan(&count)
	fmt.Printf("\n✓ Total lecturers with source='excel_import': %d\n", count)

	var total int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM faculty_lecturers").Scan(&total)
	fmt.Printf("✓ Total lecturers in database: %d\n", total)
}
