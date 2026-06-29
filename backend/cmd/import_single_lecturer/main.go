package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type pddiktiDosenSearchItem struct {
	ID        string `json:"id"`
	Name      string `json:"nama"`
	NIDN      string `json:"nidn"`
	PTName    string `json:"nama_pt"`
	ProdiName string `json:"nama_prodi"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func pddiktiGET(ctx context.Context, endpoint string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Origin", "https://pddikti.kemdiktisaintek.go.id")
	req.Header.Set("Referer", "https://pddikti.kemdiktisaintek.go.id/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("X-User-IP", "103.47.132.29")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body, resp.StatusCode, nil
}

func pddiktiJSONArrayLen(body []byte) int {
	var list []interface{}
	if err := json.Unmarshal(body, &list); err == nil {
		return len(list)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(body, &obj); err == nil {
		for _, key := range []string{"data", "result", "items"} {
			if vv, ok := obj[key]; ok {
				if subList, ok := vv.([]interface{}); ok {
					return len(subList)
				}
			}
		}
	}
	return 0
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
	queryFlag := flag.String("query", "", "Search query (NIDN or Name of the lecturer)")
	flag.Parse()

	if *queryFlag == "" {
		fmt.Println("Usage:")
		fmt.Println("  go run ./cmd/import_single_lecturer/main.go -query <NIDN-atau-Nama-Dosen>")
		return
	}

	// Load env files
	_ = godotenv.Load(".env")
	_ = godotenv.Load("backend/.env")
	_ = godotenv.Load("../.env")

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "")
	if dbname == "" {
		log.Fatal("DB_NAME is required in environment")
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

	ctx := context.Background()

	// 1. Search PDDIKTI for the lecturer
	fmt.Printf("Searching PDDIKTI for query: '%s'...\n", *queryFlag)
	searchURL := "https://api-pddikti.kemdiktisaintek.go.id/pencarian/dosen/" + url.PathEscape(*queryFlag)
	body, status, err := pddiktiGET(ctx, searchURL)
	if err != nil || status != 200 {
		log.Fatalf("✗ Failed to search PDDIKTI: status=%d, err=%v\n", status, err)
	}

	var searchItems []pddiktiDosenSearchItem
	if err := json.Unmarshal(body, &searchItems); err != nil {
		log.Fatalf("✗ Failed to unmarshal search items: %v\n", err)
	}

	var matchedLecturers []pddiktiDosenSearchItem
	for _, it := range searchItems {
		ptLower := strings.ToLower(it.PTName)
		if strings.Contains(ptLower, "perintis") || strings.Contains(ptLower, "upertis") {
			matchedLecturers = append(matchedLecturers, it)
		}
	}

	if len(matchedLecturers) == 0 {
		fmt.Println("No lecturers found matching 'Universitas Perintis Indonesia' (UPERTIS).")
		fmt.Println("All matches from other campuses:")
		for i, it := range searchItems {
			fmt.Printf("[%d] Name: %s | NIDN: %s | PT: %s | Prodi: %s\n", i+1, it.Name, it.NIDN, it.PTName, it.ProdiName)
		}
		
		if len(searchItems) > 0 {
			fmt.Println("\nIf you want to force import one of the above anyway, please double check their PT homebase on PDDIKTI.")
		}
		return
	}

	// Use the first matched lecturer
	target := matchedLecturers[0]
	fmt.Printf("\n✓ Found match: %s (NIDN: %s, PT: %s, Prodi: %s)\n", target.Name, target.NIDN, target.PTName, target.ProdiName)

	uuid := target.ID

	// 2. Fetch full portfolio from PDDIKTI endpoints
	portfolio := map[string]interface{}{}
	var teachingCount, pubCount, pengCount int
	var highestEdu, academicPosition string

	fmt.Println("Fetching profile...")
	profileEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/profile/" + url.PathEscape(uuid)
	if pBody, pStatus, pErr := pddiktiGET(ctx, profileEndpoint); pErr == nil && pStatus == 200 {
		var raw interface{}
		json.Unmarshal(pBody, &raw)
		
		if m, ok := raw.(map[string]interface{}); ok {
			m["nidn"] = target.NIDN
			portfolio["profile"] = m
			if ja, ok := m["jabatan_akademik"].(string); ok && ja != "" {
				academicPosition = ja
			}
			if pt, ok := m["pendidikan_tertinggi"].(string); ok && pt != "" {
				highestEdu = pt
			}
		}
	} else {
		log.Fatalf("✗ Failed to fetch lecturer profile: status=%d, err=%v\n", pStatus, pErr)
	}

	fmt.Println("Fetching teaching history...")
	teachEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/teaching-history/" + url.PathEscape(uuid)
	if tBody, tStatus, tErr := pddiktiGET(ctx, teachEndpoint); tErr == nil && tStatus == 200 {
		var raw interface{}
		json.Unmarshal(tBody, &raw)
		portfolio["teaching_history"] = raw
		teachingCount = pddiktiJSONArrayLen(tBody)
	}

	fmt.Println("Fetching study history...")
	studyEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/study-history/" + url.PathEscape(uuid)
	if sBody, sStatus, sErr := pddiktiGET(ctx, studyEndpoint); sErr == nil && sStatus == 200 {
		var raw interface{}
		json.Unmarshal(sBody, &raw)
		portfolio["study_history"] = raw
	}

	fmt.Println("Fetching penelitian/publications...")
	penelitianEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/penelitian/" + url.PathEscape(uuid)
	if rBody, rStatus, rErr := pddiktiGET(ctx, penelitianEndpoint); rErr == nil && rStatus == 200 {
		var raw interface{}
		json.Unmarshal(rBody, &raw)
		portfolio["penelitian"] = raw
		pubCount = pddiktiJSONArrayLen(rBody)
	}

	fmt.Println("Fetching pengabdian...")
	pengabdianEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/pengabdian/" + url.PathEscape(uuid)
	if aBody, aStatus, aErr := pddiktiGET(ctx, pengabdianEndpoint); aErr == nil && aStatus == 200 {
		var raw interface{}
		json.Unmarshal(aBody, &raw)
		portfolio["pengabdian"] = raw
		pengCount = pddiktiJSONArrayLen(aBody)
	}

	fmt.Println("Fetching karya...")
	karyaEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/karya/" + url.PathEscape(uuid)
	if kBody, kStatus, kErr := pddiktiGET(ctx, karyaEndpoint); kErr == nil && kStatus == 200 {
		var raw interface{}
		json.Unmarshal(kBody, &raw)
		portfolio["karya"] = raw
	}

	fmt.Println("Fetching paten...")
	patenEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/paten/" + url.PathEscape(uuid)
	if ptBody, ptStatus, ptErr := pddiktiGET(ctx, patenEndpoint); ptErr == nil && ptStatus == 200 {
		var raw interface{}
		json.Unmarshal(ptBody, &raw)
		portfolio["paten"] = raw
	}

	portfolioJSON, err := json.Marshal(portfolio)
	if err != nil {
		log.Fatalf("✗ Failed to marshal portfolio JSON: %v\n", err)
	}

	// 3. Resolve Faculty ID
	prodiLower := strings.ToLower(strings.TrimSpace(target.ProdiName))
	var facultyID string
	
	// Query faculties to find matching keywords
	var query string
	var args []interface{}
	
	if strings.Contains(prodiLower, "bisnis") || strings.Contains(prodiLower, "ekonomi") || strings.Contains(prodiLower, "manajemen") || strings.Contains(prodiLower, "digital") {
		query = "SELECT id FROM faculties WHERE active = true AND (LOWER(name) LIKE '%bisnis%' OR LOWER(name) LIKE '%ekonomi%' OR LOWER(name) LIKE '%manajemen%') ORDER BY sort_order LIMIT 1"
	} else if strings.Contains(prodiLower, "farmasi") {
		query = "SELECT id FROM faculties WHERE active = true AND LOWER(name) LIKE '%farmasi%' ORDER BY sort_order LIMIT 1"
	} else {
		query = "SELECT id FROM faculties WHERE active = true AND (LOWER(name) LIKE '%kesehatan%' OR LOWER(name) LIKE '%keperawatan%' OR LOWER(name) LIKE '%fikes%') ORDER BY sort_order LIMIT 1"
	}

	err = db.QueryRowContext(ctx, query, args...).Scan(&facultyID)
	if err != nil {
		// Fallback to first active faculty
		_ = db.QueryRowContext(ctx, "SELECT id FROM faculties WHERE active = true ORDER BY sort_order LIMIT 1").Scan(&facultyID)
	}

	if facultyID == "" {
		log.Fatal("✗ No active faculty found in database to attach the lecturer")
	}

	// Format Education text: e.g. "S2 | Mengajar:15 | Publikasi:2"
	parts := []string{}
	if highestEdu != "" {
		parts = append(parts, highestEdu)
	}
	if teachingCount > 0 {
		parts = append(parts, fmt.Sprintf("Mengajar:%d", teachingCount))
	}
	if pubCount > 0 {
		parts = append(parts, fmt.Sprintf("Publikasi:%d", pubCount))
	}
	if pengCount > 0 {
		parts = append(parts, fmt.Sprintf("Pengabdian:%d", pengCount))
	}
	formattedEdu := strings.Join(parts, " | ")

	slug := makeSlug(target.Name)

	position := academicPosition
	if position == "" {
		position = "Dosen"
	}

	// 4. Save/Upsert into faculty_lecturers
	fmt.Println("Upserting lecturer into database...")
	
	insertSQL := `
		INSERT INTO faculty_lecturers (faculty_id, slug, name, position, education, expertise, photo_url, leadership_group, source, active, sort_order, pddikti_id, pddikti_uuid, portfolio)
		VALUES ($1, $2, $3, $4, $5, $6, '', '', 'pddikti', true, 0, NULLIF($7, ''), NULLIF($8, ''), $9::jsonb)
		ON CONFLICT (pddikti_uuid) DO UPDATE SET
			faculty_id = EXCLUDED.faculty_id,
			name = EXCLUDED.name,
			position = EXCLUDED.position,
			education = EXCLUDED.education,
			expertise = EXCLUDED.expertise,
			source = EXCLUDED.source,
			active = true,
			updated_at = NOW(),
			pddikti_id = EXCLUDED.pddikti_id,
			portfolio = EXCLUDED.portfolio
		RETURNING id
	`

	var insertedID string
	err = db.QueryRowContext(ctx, insertSQL, facultyID, slug, target.Name, position, formattedEdu, target.ProdiName, target.NIDN, uuid, string(portfolioJSON)).Scan(&insertedID)
	if err != nil {
		log.Fatalf("✗ Failed to insert/update lecturer: %v\n", err)
	}

	// Ensure unique slug if needed
	uniqueSlug := slug
	var exists bool
	err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM faculty_lecturers WHERE slug = $1 AND id != $2)", uniqueSlug, insertedID).Scan(&exists)
	if err == nil && exists {
		uniqueSlug = fmt.Sprintf("%s-%s", slug, insertedID[:6])
		_, _ = db.ExecContext(ctx, "UPDATE faculty_lecturers SET slug = $1 WHERE id = $2", uniqueSlug, insertedID)
	}

	fmt.Printf("\n=== IMPORT SUCCESSFUL ===\n")
	fmt.Printf("ID: %s\n", insertedID)
	fmt.Printf("Name: %s\n", target.Name)
	fmt.Printf("NIDN: %s\n", target.NIDN)
	fmt.Printf("Slug: %s\n", uniqueSlug)
	fmt.Printf("Prodi: %s\n", target.ProdiName)
	fmt.Printf("Faculty ID: %s\n", facultyID)
}
