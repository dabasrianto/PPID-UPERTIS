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
	"sync"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Lecturer struct {
	ID        string
	Name      string
	Slug      string
	NIDN      string
	UUID      string
	Portfolio string
}

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

func main() {
	// Parse CLI flags
	slugFlag := flag.String("slug", "", "Sync only a specific lecturer by their slug")
	allFlag := flag.Bool("all", false, "Sync all lecturers who have empty portfolios")
	forceFlag := flag.Bool("force", false, "Force sync even if portfolio is already populated")
	flag.Parse()

	if *slugFlag == "" && !*allFlag {
		fmt.Println("Usage:")
		fmt.Println("  go run ./cmd/sync_pddikti_portfolios/main.go -slug <lecturer-slug>")
		fmt.Println("  go run ./cmd/sync_pddikti_portfolios/main.go -all")
		fmt.Println("  Add -force to re-sync already populated portfolios")
		return
	}

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

	// Load dosen_data.json for NIDN matching
	dosenMap := make(map[string]string)
	dosenBytes, err := os.ReadFile("dosen_data.json")
	if err != nil {
		dosenBytes, err = os.ReadFile("backend/dosen_data.json")
	}
	if err == nil {
		type DosenJSONItem struct {
			Reg  string `json:"reg"`
			Nama string `json:"nama"`
		}
		var items []DosenJSONItem
		if err := json.Unmarshal(dosenBytes, &items); err == nil {
			for _, item := range items {
				clean := cleanName(item.Nama)
				if item.Reg != "" && clean != "" {
					dosenMap[clean] = item.Reg
				}
			}
			fmt.Printf("✓ Loaded %d lecturers from dosen_data.json\n", len(dosenMap))
		} else {
			fmt.Printf("⚠ Failed to parse dosen_data.json: %v\n", err)
		}
	} else {
		fmt.Printf("⚠ Failed to read dosen_data.json: %v\n", err)
	}

	// Query lecturers
	var query string
	var args []interface{}
	if *slugFlag != "" {
		query = "SELECT id, name, slug, COALESCE(pddikti_id, ''), COALESCE(pddikti_uuid, ''), portfolio::text FROM faculty_lecturers WHERE slug = $1"
		args = append(args, *slugFlag)
	} else {
		if *forceFlag {
			query = "SELECT id, name, slug, COALESCE(pddikti_id, ''), COALESCE(pddikti_uuid, ''), portfolio::text FROM faculty_lecturers"
		} else {
			query = "SELECT id, name, slug, COALESCE(pddikti_id, ''), COALESCE(pddikti_uuid, ''), portfolio::text FROM faculty_lecturers WHERE portfolio IS NULL OR portfolio::text = '{}'"
		}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal("Query failed: ", err)
	}
	defer rows.Close()

	var lecturers []Lecturer
	for rows.Next() {
		var l Lecturer
		if err := rows.Scan(&l.ID, &l.Name, &l.Slug, &l.NIDN, &l.UUID, &l.Portfolio); err != nil {
			log.Fatal(err)
		}
		lecturers = append(lecturers, l)
	}
	fmt.Printf("✓ Loaded %d lecturers to sync\n", len(lecturers))

	nidnReg := regexp.MustCompile(`^\d{10}$`)

	var wg sync.WaitGroup
	sem := make(chan struct{}, 3) // Concurrency throttle

	for _, l := range lecturers {
		wg.Add(1)
		go func(lect Lecturer) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			fmt.Printf("[%s] Starting sync for: %s\n", lect.Slug, lect.Name)

			uuid := lect.UUID
			isNidn := nidnReg.MatchString(lect.NIDN)
			var resolvedNidn string

			if uuid == "" && lect.NIDN != "" && isNidn {
				fmt.Printf("[%s] NIDN '%s' detected but no UUID in DB. Resolving to PDDIKTI UUID...\n", lect.Slug, lect.NIDN)
				searchURL := "https://api-pddikti.kemdiktisaintek.go.id/pencarian/dosen/" + url.PathEscape(lect.NIDN)
				body, status, err := pddiktiGET(ctx, searchURL)
				if err == nil && status == 200 {
					var searchItems []pddiktiDosenSearchItem
					if err := json.Unmarshal(body, &searchItems); err == nil {
						for _, it := range searchItems {
							ptLower := strings.ToLower(it.PTName)
							if strings.Contains(ptLower, "perintis") || strings.Contains(ptLower, "upertis") {
								uuid = it.ID
								resolvedNidn = it.NIDN
								fmt.Printf("[%s] ✓ Resolved NIDN %s to UUID: %s\n", lect.Slug, lect.NIDN, uuid)
								break
							}
						}
					}
				}
			}

			if uuid == "" {
				// No UUID! Let's try name search
				cleanNameStr := lect.Name
				if idx := strings.Index(cleanNameStr, ","); idx != -1 {
					cleanNameStr = cleanNameStr[:idx]
				}
				cleanNameStr = strings.TrimSpace(cleanNameStr)
				
				resolvedNidn = dosenMap[cleanName(lect.Name)]
				
				fmt.Printf("[%s] Empty identifier in DB. Searching PDDIKTI for name: %s...\n", lect.Slug, cleanNameStr)
				searchURL := "https://api-pddikti.kemdiktisaintek.go.id/pencarian/dosen/" + url.PathEscape(cleanNameStr)
				if body, status, err := pddiktiGET(ctx, searchURL); err == nil && status == 200 {
					var nameItems []pddiktiDosenSearchItem
					if err := json.Unmarshal(body, &nameItems); err == nil {
						for _, it := range nameItems {
							ptLower := strings.ToLower(it.PTName)
							if strings.Contains(ptLower, "perintis") || strings.Contains(ptLower, "upertis") {
								uuid = it.ID
								resolvedNidn = it.NIDN
								fmt.Printf("[%s] ✓ Resolved empty identifier to UUID: %s, NIDN: %s\n", lect.Slug, uuid, resolvedNidn)
								break
							}
						}
					}
				}
			} else {
				// Already a UUID
				if resolvedNidn == "" {
					resolvedNidn = lect.NIDN
				}
				if resolvedNidn == "" {
					resolvedNidn = dosenMap[cleanName(lect.Name)]
				}
				if resolvedNidn == "" {
					cleanNameStr := lect.Name
					if idx := strings.Index(cleanNameStr, ","); idx != -1 {
						cleanNameStr = cleanNameStr[:idx]
					}
					cleanNameStr = strings.TrimSpace(cleanNameStr)
					fmt.Printf("[%s] UUID set but NIDN not in map. Searching PDDIKTI for NIDN...\n", lect.Slug)
					searchURL := "https://api-pddikti.kemdiktisaintek.go.id/pencarian/dosen/" + url.PathEscape(cleanNameStr)
					if body, status, err := pddiktiGET(ctx, searchURL); err == nil && status == 200 {
						var nameItems []pddiktiDosenSearchItem
						if err := json.Unmarshal(body, &nameItems); err == nil {
							for _, it := range nameItems {
								ptLower := strings.ToLower(it.PTName)
								if (strings.Contains(ptLower, "perintis") || strings.Contains(ptLower, "upertis")) && (it.ID == uuid) {
									resolvedNidn = it.NIDN
									fmt.Printf("[%s] ✓ Found NIDN via name search: %s\n", lect.Slug, resolvedNidn)
									break
								}
							}
						}
					}
				}
			}

			if uuid == "" || len(uuid) < 20 {
				fmt.Printf("[%s] ✗ Invalid UUID key: %s\n", lect.Slug, uuid)
				return
			}

			// We have a UUID! Scrape the portfolio parts.
			portfolio := map[string]interface{}{}
			var teachingCount, pubCount, pengCount int
			var highestEdu, academicPosition string

			// 1. Profile
			profileEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/profile/" + url.PathEscape(uuid)
			if body, status, err := pddiktiGET(ctx, profileEndpoint); err == nil && status == 200 {
				var raw interface{}
				json.Unmarshal(body, &raw)

				nidnVal := ""
				if isNidn {
					nidnVal = lect.NIDN
				} else if resolvedNidn != "" {
					nidnVal = resolvedNidn
				} else {
					nidnVal = dosenMap[cleanName(lect.Name)]
				}

				if m, ok := raw.(map[string]interface{}); ok && nidnVal != "" {
					m["nidn"] = nidnVal
				}
				portfolio["profile"] = raw

				var m map[string]interface{}
				json.Unmarshal(body, &m)
				if ja, ok := m["jabatan_akademik"].(string); ok && ja != "" {
					academicPosition = ja
				}
				if pt, ok := m["pendidikan_tertinggi"].(string); ok && pt != "" {
					highestEdu = pt
				}
			} else {
				fmt.Printf("[%s] ✗ Profile endpoint failed: status=%d, err=%v\n", lect.Slug, status, err)
				return // If profile fails, skip the rest
			}

			// 2. Teaching History
			teachEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/teaching-history/" + url.PathEscape(uuid)
			if body, status, err := pddiktiGET(ctx, teachEndpoint); err == nil && status == 200 {
				var raw interface{}
				json.Unmarshal(body, &raw)
				portfolio["teaching_history"] = raw
				teachingCount = pddiktiJSONArrayLen(body)
			}

			// 3. Study History
			studyEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/study-history/" + url.PathEscape(uuid)
			if body, status, err := pddiktiGET(ctx, studyEndpoint); err == nil && status == 200 {
				var raw interface{}
				json.Unmarshal(body, &raw)
				portfolio["study_history"] = raw
			}

			// 4. Penelitian
			penelitianEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/penelitian/" + url.PathEscape(uuid)
			if body, status, err := pddiktiGET(ctx, penelitianEndpoint); err == nil && status == 200 {
				var raw interface{}
				json.Unmarshal(body, &raw)
				portfolio["penelitian"] = raw
				pubCount = pddiktiJSONArrayLen(body)
			}

			// 5. Pengabdian
			pengabdianEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/pengabdian/" + url.PathEscape(uuid)
			if body, status, err := pddiktiGET(ctx, pengabdianEndpoint); err == nil && status == 200 {
				var raw interface{}
				json.Unmarshal(body, &raw)
				portfolio["pengabdian"] = raw
				pengCount = pddiktiJSONArrayLen(body)
			}

			// 6. Karya
			karyaEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/karya/" + url.PathEscape(uuid)
			if body, status, err := pddiktiGET(ctx, karyaEndpoint); err == nil && status == 200 {
				var raw interface{}
				json.Unmarshal(body, &raw)
				portfolio["karya"] = raw
			}

			// 7. Paten
			patenEndpoint := "https://api-pddikti.kemdiktisaintek.go.id/dosen/portofolio/paten/" + url.PathEscape(uuid)
			if body, status, err := pddiktiGET(ctx, patenEndpoint); err == nil && status == 200 {
				var raw interface{}
				json.Unmarshal(body, &raw)
				portfolio["paten"] = raw
			}

			// Marshal portfolio
			portfolioJSON, err := json.Marshal(portfolio)
			if err != nil {
				fmt.Printf("[%s] ✗ Failed to marshal portfolio map: %v\n", lect.Slug, err)
				return
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

			// Update database record
			// We only update position if it's currently a standard position (or empty/dosen) and not leadership positions like rektor/dekan/kaprodi
			var updateQuery string
			var updateArgs []interface{}

			if academicPosition != "" {
				updateQuery = `
					UPDATE faculty_lecturers SET
						pddikti_uuid = $1,
						portfolio = $2,
						education = $3,
						position = CASE 
							WHEN LOWER(position) LIKE '%rektor%' 
							  OR LOWER(position) LIKE '%dekan%' 
							  OR LOWER(position) LIKE '%kaprodi%' 
							  OR LOWER(position) LIKE '%prodi%' 
							  OR LOWER(position) LIKE '%yayasan%' 
							  OR LOWER(position) LIKE '%founder%' 
							  THEN position
							ELSE $4
						END,
						source = 'pddikti',
						pddikti_id = NULLIF($5, ''),
						updated_at = NOW()
					WHERE id = $6
				`
				updateArgs = []interface{}{uuid, string(portfolioJSON), formattedEdu, academicPosition, resolvedNidn, lect.ID}
			} else {
				updateQuery = `
					UPDATE faculty_lecturers SET
						pddikti_uuid = $1,
						portfolio = $2,
						education = $3,
						source = 'pddikti',
						pddikti_id = NULLIF($4, ''),
						updated_at = NOW()
					WHERE id = $5
				`
				updateArgs = []interface{}{uuid, string(portfolioJSON), formattedEdu, resolvedNidn, lect.ID}
			}

			_, err = db.Exec(updateQuery, updateArgs...)
			if err != nil {
				fmt.Printf("[%s] ✗ Failed to update database row: %v\n", lect.Slug, err)
			} else {
				fmt.Printf("[%s] ✓ Portfolio sync successful! (T:%d P:%d C:%d)\n", lect.Slug, teachingCount, pubCount, pengCount)
			}

			// Add a short delay to space out API requests
			time.Sleep(300 * time.Millisecond)

		}(l)
	}

	wg.Wait()
	fmt.Println("\n=== Portfolio Sync Completed ===")
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
