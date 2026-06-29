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
		dbname = "kampuspro"
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
	fmt.Println("✓ Connected to database")

	// Load dosen_data.json
	dosenMap := make(map[string]string)
	dosenBytes, err := os.ReadFile("dosen_data.json")
	if err != nil {
		dosenBytes, err = os.ReadFile("backend/dosen_data.json")
		if err != nil {
			dosenBytes, err = os.ReadFile("../dosen_data.json")
			if err != nil {
				dosenBytes, err = os.ReadFile("../../dosen_data.json")
				if err != nil {
					log.Fatal("Failed to read dosen_data.json: ", err)
				}
			}
		}
	}

	type DosenJSONItem struct {
		Reg  string `json:"reg"`
		Nama string `json:"nama"`
	}
	var items []DosenJSONItem
	if err := json.Unmarshal(dosenBytes, &items); err != nil {
		log.Fatal("Failed to parse dosen_data.json: ", err)
	}

	for _, item := range items {
		clean := cleanName(item.Nama)
		if item.Reg != "" && clean != "" {
			dosenMap[clean] = item.Reg
		}
	}
	fmt.Printf("✓ Loaded %d lecturers from dosen_data.json\n", len(dosenMap))

	// Fetch all lecturers from database
	rows, err := db.QueryContext(context.Background(), "SELECT id, name, COALESCE(portfolio::text, '{}') FROM faculty_lecturers")
	if err != nil {
		log.Fatal("Query failed: ", err)
	}
	defer rows.Close()

	type DbLecturer struct {
		ID        string
		Name      string
		Portfolio string
	}
	var dbLecturers []DbLecturer
	for rows.Next() {
		var l DbLecturer
		if err := rows.Scan(&l.ID, &l.Name, &l.Portfolio); err != nil {
			log.Fatal("Scan failed: ", err)
		}
		dbLecturers = append(dbLecturers, l)
	}
	fmt.Printf("✓ Loaded %d lecturers from database\n", len(dbLecturers))

	updatedCount := 0
	for _, l := range dbLecturers {
		clean := cleanName(l.Name)
		nidn, found := dosenMap[clean]
		if !found {
			continue
		}

		var portfolio map[string]interface{}
		if err := json.Unmarshal([]byte(l.Portfolio), &portfolio); err != nil {
			portfolio = make(map[string]interface{})
		}

		profileRaw, hasProfile := portfolio["profile"]
		var profile map[string]interface{}
		if hasProfile {
			if m, ok := profileRaw.(map[string]interface{}); ok {
				profile = m
			} else {
				profile = make(map[string]interface{})
			}
		} else {
			profile = make(map[string]interface{})
		}

		currentNidn, _ := profile["nidn"].(string)
		if currentNidn == nidn {
			continue // Already correct
		}

		profile["nidn"] = nidn
		portfolio["profile"] = profile

		portfolioJSON, err := json.Marshal(portfolio)
		if err != nil {
			log.Printf("Failed to marshal portfolio for %s: %v\n", l.Name, err)
			continue
		}

		_, err = db.ExecContext(context.Background(), "UPDATE faculty_lecturers SET portfolio = $1, updated_at = NOW() WHERE id = $2", string(portfolioJSON), l.ID)
		if err != nil {
			log.Printf("Failed to update database for %s: %v\n", l.Name, err)
		} else {
			fmt.Printf("Updated portfolio for %s -> NIDN set to %s\n", l.Name, nidn)
			updatedCount++
		}
	}

	fmt.Printf("\n=== MIGRATION COMPLETED ===\n")
	fmt.Printf("Total lecturers updated: %d\n", updatedCount)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
