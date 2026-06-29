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

type RekapItem struct {
	No               int    `json:"no"`
	Nidn             string `json:"nidn"`
	NamaDosen        string `json:"nama_dosen"`
	ProgramStudi     string `json:"program_studi"`
	Serdos           string `json:"serdos"`
	GanjilPenelitian string `json:"ganjil_penelitian"`
	GanjilPkm        string `json:"ganjil_pkm"`
	GanjilLuaran     string `json:"ganjil_luaran"`
	PenulisKedua     string `json:"penulis_kedua"`
	GenapPenelitian  string `json:"genap_penelitian"`
	GenapPkm         string `json:"genap_pkm"`
	GenapLuaran      string `json:"genap_luaran"`
}

type DbLecturer struct {
	ID        string
	Name      string
	PddiktiID string
	Active    bool
}

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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func cleanName(name string) string {
	name = strings.ToLower(name)
	
	// Remove prefixes
	prefixes := []string{"dr.", "prof.", "apt.", "ns.", "drs.", "dra.", "hj.", "h.", "dr", "prof", "apt", "ns", "drs", "dra"}
	for _, p := range prefixes {
		if strings.HasPrefix(name, p+" ") {
			name = strings.TrimPrefix(name, p+" ")
		}
	}

	// Remove suffixes starting with comma
	if idx := strings.Index(name, ","); idx != -1 {
		name = name[:idx]
	}

	// Keep only letters and numbers
	reg := regexp.MustCompile("[^a-z0-9]")
	name = reg.ReplaceAllString(name, "")
	
	return strings.TrimSpace(name)
}

func toTitleCase(s string) string {
	words := strings.Fields(strings.ToLower(s))
	for i, w := range words {
		if len(w) > 0 {
			// Capitalize first letter of word
			words[i] = strings.ToUpper(string(w[0])) + w[1:]
		}
	}
	return strings.Join(words, " ")
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

func main() {
	// Load environment
	_ = godotenv.Load(".env")
	_ = godotenv.Load("backend/.env")
	_ = godotenv.Load("../.env")

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "")
	sslmode := getEnv("DB_SSLMODE", "disable")

	if dbname == "" {
		log.Fatal("DB_NAME required")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database ping failed: ", err)
	}

	// Query faculties
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

	// Load rekap JSON
	rekapJSONPath := "rekap_penelitian.json"
	rekapBytes, err := os.ReadFile(rekapJSONPath)
	if err != nil {
		rekapJSONPath = "backend/rekap_penelitian.json"
		rekapBytes, err = os.ReadFile(rekapJSONPath)
		if err != nil {
			log.Fatal("Failed to read rekap_penelitian.json: ", err)
		}
	}

	var rekapItems []RekapItem
	if err := json.Unmarshal(rekapBytes, &rekapItems); err != nil {
		log.Fatal("Failed to parse JSON: ", err)
	}

	// Fetch existing lecturers
	rows, err := db.QueryContext(context.Background(), "SELECT id, name, COALESCE(pddikti_id, ''), active FROM faculty_lecturers")
	if err != nil {
		log.Fatal("Query lecturers failed: ", err)
	}
	defer rows.Close()

	var dbLecturers []DbLecturer
	dbMapByName := make(map[string]DbLecturer)
	dbMapByNidn := make(map[string]DbLecturer)

	for rows.Next() {
		var l DbLecturer
		if err := rows.Scan(&l.ID, &l.Name, &l.PddiktiID, &l.Active); err != nil {
			log.Fatal("Scan failed: ", err)
		}
		dbLecturers = append(dbLecturers, l)
		
		cleanN := cleanName(l.Name)
		if cleanN != "" {
			dbMapByName[cleanN] = l
		}
		
		cleanNidn := strings.TrimSpace(l.PddiktiID)
		if cleanNidn != "" {
			dbMapByNidn[cleanNidn] = l
		}
	}

	fmt.Println("=== STARTING DOSEN SYNCHRONIZATION ===")
	fmt.Printf("Total Dosen in Master List: %d\n", len(rekapItems))
	fmt.Printf("Total Dosen in Database: %d\n\n", len(dbLecturers))

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal("Failed to begin transaction: ", err)
	}

	matchedIds := make(map[string]bool)
	insertedCount := 0
	updatedCount := 0
	deactivatedCount := 0

	for _, item := range rekapItems {
		var matchedDb DbLecturer
		found := false

		// Try NIDN match
		cleanRekapNidn := strings.TrimSpace(item.Nidn)
		if cleanRekapNidn != "" && cleanRekapNidn != "-" {
			if l, exists := dbMapByNidn[cleanRekapNidn]; exists {
				matchedDb = l
				found = true
			}
		}

		// Try Name match
		if !found {
			cleanRekapName := cleanName(item.NamaDosen)
			if l, exists := dbMapByName[cleanRekapName]; exists {
				matchedDb = l
				found = true
			}
		}

		if found {
			matchedIds[matchedDb.ID] = true
			
			// Update matching DB entry to match rekap settings
			// We update expertise to exact program studi name, ensure active=true
			var query string
			var args []interface{}
			
			if matchedDb.PddiktiID == "" && cleanRekapNidn != "" && cleanRekapNidn != "-" {
				query = "UPDATE faculty_lecturers SET active = true, expertise = $1, pddikti_id = $2, source = 'rekap_sync' WHERE id = $3"
				args = []interface{}{item.ProgramStudi, cleanRekapNidn, matchedDb.ID}
			} else {
				query = "UPDATE faculty_lecturers SET active = true, expertise = $1, source = 'rekap_sync' WHERE id = $2"
				args = []interface{}{item.ProgramStudi, matchedDb.ID}
			}

			_, err = tx.ExecContext(ctx, query, args...)
			if err != nil {
				tx.Rollback()
				log.Fatalf("Failed to update lecturer %s: %v", item.NamaDosen, err)
			}
			updatedCount++
		} else {
			// Insert new lecturer
			code, ok := prodiToFacultyCode[item.ProgramStudi]
			if !ok {
				tx.Rollback()
				log.Fatalf("Failed to resolve faculty mapping for major: %s", item.ProgramStudi)
			}
			facultyID := facultyMap[code]
			if facultyID == "" && code == "FEB" {
				facultyID = facultyMap["FEBIS"]
			}
			if facultyID == "" && code == "FEBIS" {
				facultyID = facultyMap["FEB"]
			}
			if facultyID == "" {
				tx.Rollback()
				log.Fatalf("Faculty code %s not found in DB", code)
			}

			formattedName := toTitleCase(item.NamaDosen)
			slugBase := makeSlug(formattedName)
			slug := slugBase
			
			// Simple slug duplicates check in memory
			slugCount := 1
			slugExists := true
			for slugExists {
				var count int
				err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM faculty_lecturers WHERE slug = $1", slug).Scan(&count)
				if err != nil {
					tx.Rollback()
					log.Fatal(err)
				}
				if count == 0 {
					slugExists = false
				} else {
					slugCount++
					slug = fmt.Sprintf("%s-%d", slugBase, slugCount)
				}
			}

			var pddiktiVal sql.NullString
			if cleanRekapNidn != "" && cleanRekapNidn != "-" {
				pddiktiVal = sql.NullString{String: cleanRekapNidn, Valid: true}
			}

			_, err = tx.ExecContext(ctx, `
				INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
				VALUES ($1, $2, $3, 'Dosen Pengajar', 'S2', $4, 'rekap_sync', $5, true, 0)`,
				facultyID, formattedName, slug, item.ProgramStudi, pddiktiVal)
			if err != nil {
				tx.Rollback()
				log.Fatalf("Failed to insert lecturer %s: %v", formattedName, err)
			}
			insertedCount++
			fmt.Printf("+ Inserted missing lecturer: %s | NIDN: %s | Prodi: %s\n", formattedName, item.Nidn, item.ProgramStudi)
		}
	}

	// Deactivate database lecturers not in rekap list
	for _, dbL := range dbLecturers {
		if !matchedIds[dbL.ID] {
			_, err = tx.ExecContext(ctx, "UPDATE faculty_lecturers SET active = false WHERE id = $1", dbL.ID)
			if err != nil {
				tx.Rollback()
				log.Fatalf("Failed to deactivate lecturer %s: %v", dbL.Name, err)
			}
			deactivatedCount++
			fmt.Printf("- Deactivated extra lecturer: %s | NIDN/PDDIKTI: %s\n", dbL.Name, dbL.PddiktiID)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal("Failed to commit transaction: ", err)
	}

	fmt.Println("\n=== SYNCHRONIZATION COMPLETED ===")
	fmt.Printf("✓ Match Updated:  %d\n", updatedCount)
	fmt.Printf("✓ Newly Inserted: %d\n", insertedCount)
	fmt.Printf("✓ Deactivated:    %d\n", deactivatedCount)
	fmt.Printf("Total Synchronized Active Lecturers: %d\n", updatedCount+insertedCount)
}
