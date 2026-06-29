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

	// 1. Fetch lecturers from database
	rows, err := db.QueryContext(context.Background(), "SELECT id, name, COALESCE(pddikti_id, '') FROM faculty_lecturers")
	if err != nil {
		log.Fatal("Query failed: ", err)
	}
	defer rows.Close()

	var dbLecturers []DbLecturer
	dbMapByName := make(map[string]DbLecturer)
	dbMapByNidn := make(map[string]DbLecturer)

	for rows.Next() {
		var l DbLecturer
		if err := rows.Scan(&l.ID, &l.Name, &l.PddiktiID); err != nil {
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

	// 2. Load rekap JSON
	rekapJSONPath := "backend/rekap_penelitian.json"
	rekapBytes, err := os.ReadFile(rekapJSONPath)
	if err != nil {
		// Fallback
		rekapJSONPath = "rekap_penelitian.json"
		rekapBytes, err = os.ReadFile(rekapJSONPath)
		if err != nil {
			log.Fatal("Failed to read rekap_penelitian.json: ", err)
		}
	}

	var rekapItems []RekapItem
	if err := json.Unmarshal(rekapBytes, &rekapItems); err != nil {
		log.Fatal("Failed to parse JSON: ", err)
	}

	// 3. Perform comparison
	matchedCount := 0
	var unmatchedRekap []RekapItem

	fmt.Printf("=== PERBANDINGAN DOSEN REKAP VS DATABASE ===\n")
	fmt.Printf("Total Dosen di Rekap: %d\n", len(rekapItems))
	fmt.Printf("Total Dosen di Database (faculty_lecturers): %d\n\n", len(dbLecturers))

	for _, item := range rekapItems {
		matchFound := false
		var matchedLecturer DbLecturer

		// Match 1: NIDN
		cleanRekapNidn := strings.TrimSpace(item.Nidn)
		if cleanRekapNidn != "" && cleanRekapNidn != "-" {
			if l, exists := dbMapByNidn[cleanRekapNidn]; exists {
				matchFound = true
				matchedLecturer = l
			}
		}

		// Match 2: Name
		if !matchFound {
			cleanRekapName := cleanName(item.NamaDosen)
			if l, exists := dbMapByName[cleanRekapName]; exists {
				matchFound = true
				matchedLecturer = l
			}
		}

		if matchFound {
			matchedCount++
			_ = matchedLecturer // Used
		} else {
			unmatchedRekap = append(unmatchedRekap, item)
		}
	}

	matchPercentage := float64(matchedCount) / float64(len(rekapItems)) * 100
	fmt.Printf("Dosen Rekap yang COCOK dengan Database: %d / %d (%.2f%%)\n", matchedCount, len(rekapItems), matchPercentage)
	fmt.Printf("Dosen Rekap yang TIDAK COCOK (Tidak ada di database): %d\n\n", len(unmatchedRekap))

	if len(unmatchedRekap) > 0 {
		fmt.Println("--- DOSEN REKAP YANG TIDAK ADA DI DATABASE APP ---")
		for idx, item := range unmatchedRekap {
			fmt.Printf("%d. Nama: %s | NIDN: %s | Prodi: %s\n", idx+1, item.NamaDosen, item.Nidn, item.ProgramStudi)
		}
		fmt.Println()
	}

	// 4. Check if some database lecturers are not in the rekap
	rekapNamesMap := make(map[string]bool)
	rekapNidnMap := make(map[string]bool)
	for _, item := range rekapItems {
		cleanN := cleanName(item.NamaDosen)
		if cleanN != "" {
			rekapNamesMap[cleanN] = true
		}
		cleanNidn := strings.TrimSpace(item.Nidn)
		if cleanNidn != "" && cleanNidn != "-" {
			rekapNidnMap[cleanNidn] = true
		}
	}

	var unmatchedDb []DbLecturer
	for _, l := range dbLecturers {
		matchFound := false
		cleanNidn := strings.TrimSpace(l.PddiktiID)
		if cleanNidn != "" {
			if rekapNidnMap[cleanNidn] {
				matchFound = true
			}
		}
		if !matchFound {
			cleanN := cleanName(l.Name)
			if rekapNamesMap[cleanN] {
				matchFound = true
			}
		}
		if !matchFound {
			unmatchedDb = append(unmatchedDb, l)
		}
	}

	fmt.Printf("Dosen di Database yang TIDAK ADA di Rekap: %d / %d\n", len(unmatchedDb), len(dbLecturers))
	if len(unmatchedDb) > 0 {
		fmt.Println("--- CONTOH DOSEN DATABASE YANG TIDAK ADA DI REKAP (Max 10) ---")
		limit := 10
		if len(unmatchedDb) < limit {
			limit = len(unmatchedDb)
		}
		for i := 0; i < limit; i++ {
			l := unmatchedDb[i]
			fmt.Printf("%d. Nama: %s | NIDN/PDDIKTI: %s\n", i+1, l.Name, l.PddiktiID)
		}
	}
}
