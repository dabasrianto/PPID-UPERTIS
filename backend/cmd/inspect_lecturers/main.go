package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("backend/.env")
		if err != nil {
			err = godotenv.Load("../.env")
			if err != nil {
				log.Println("Warning: No .env file found, using defaults")
			}
		}
	}

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

	fmt.Println("--- Lecturer Status Verification ---")
	var total, active, inactive int
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM faculty_lecturers").Scan(&total)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM faculty_lecturers WHERE active = true").Scan(&active)
	db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM faculty_lecturers WHERE active = false").Scan(&inactive)
	
	fmt.Printf("Total Lecturers: %d (Expected: 127)\n", total)
	fmt.Printf("Active Lecturers: %d (Expected: 118)\n", active)
	fmt.Printf("Inactive/Deactivated: %d (Expected: 9)\n\n", inactive)

	fmt.Println("--- Inactive Lecturers ---")
	rows, err := db.QueryContext(context.Background(), "SELECT name, pddikti_id FROM faculty_lecturers WHERE active = false ORDER BY name")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name, nidn string
			rows.Scan(&name, &nidn)
			fmt.Printf("- %s (NIDN: %s)\n", name, nidn)
		}
	}

	fmt.Println("\n--- Newly Inserted (rekap_sync) ---")
	rows2, err := db.QueryContext(context.Background(), "SELECT name, pddikti_id, expertise FROM faculty_lecturers WHERE source = 'rekap_sync' AND active = true ORDER BY created_at DESC LIMIT 10")
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var name, nidn, expertise string
			rows2.Scan(&name, &nidn, &expertise)
			fmt.Printf("- %s (NIDN: %s, Prodi: %s)\n", name, nidn, expertise)
		}
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
