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
		log.Fatal("Failed to connect to db: ", err)
	}
	defer db.Close()

	ctx := context.Background()

	// 1. Query all faculties
	fmt.Println("=== Checking Faculties in Production Database ===")
	facRows, err := db.QueryContext(ctx, "SELECT id, code, name FROM faculties")
	if err != nil {
		log.Fatal("Failed to query faculties: ", err)
	}
	facultyMap := make(map[string]string)
	for facRows.Next() {
		var id, code, name string
		facRows.Scan(&id, &code, &name)
		fmt.Printf("Faculty: ID=%s, Code=%s, Name=%s\n", id, code, name)
		facultyMap[code] = id
	}
	facRows.Close()

	// 2. Try to search for AFRI ANDIKA in production database
	fmt.Println("\n=== Checking for existing AFRI ANDIKA in production database ===")
	rows, err := db.QueryContext(ctx, "SELECT id, faculty_id, name, slug, pddikti_id, source FROM faculty_lecturers WHERE name ILIKE '%AFRI ANDIKA%'")
	if err != nil {
		log.Fatal("Search query failed: ", err)
	}
	count := 0
	for rows.Next() {
		count++
		var id, facultyID, name, slug, pddiktiID, source string
		rows.Scan(&id, &facultyID, &name, &slug, &pddiktiID, &source)
		fmt.Printf("Found Lecturer: ID=%s, FacultyID=%s, Name=%s, Slug=%s, PDDIKTI_ID=%s, Source=%s\n", id, facultyID, name, slug, pddiktiID, source)
	}
	rows.Close()
	if count == 0 {
		fmt.Println("No records matching 'AFRI ANDIKA' found in database.")
	}

	// 3. Test the lookup query exactly like in import_dosen
	fmt.Println("\n=== Testing Lookup Query like in import_dosen ===")
	febID := facultyMap["FEB"]
	if febID == "" {
		fmt.Println("WARNING: FEB faculty not found in database!")
	} else {
		var existingID string
		err = db.QueryRowContext(ctx,
			"SELECT id FROM faculty_lecturers WHERE (UPPER(name) = UPPER($1) OR UPPER(name) LIKE UPPER($1) || ',%') AND faculty_id = $2",
			"AFRI ANDIKA", febID).Scan(&existingID)
		if err != nil {
			fmt.Printf("Lookup returned error: %v (IsErrNoRows=%v)\n", err, err == sql.ErrNoRows)
		} else {
			fmt.Printf("Lookup found existing ID: %s\n", existingID)
		}
	}

	// 4. Test the insert query in a rollback transaction with dynamic faculty ID
	fmt.Println("\n=== Testing Raw Insert with Dynamic FEB Faculty ID ===")
	if febID == "" {
		fmt.Println("Cannot test insert without FEB ID")
		return
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal("Failed to begin tx: ", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6, 'excel_import', $7, true, 0)`,
		febID, "AFRI ANDIKA, S.M,M.M", "afri-andika-test", "Dosen", "S2, S.M,M.M", "Bisnis Digital (S1)", nil)
	
	if err != nil {
		fmt.Printf(">>> PostgreSQL INSERT Error: %v\n", err)
	} else {
		fmt.Println(">>> Insert succeeded in transaction simulation using FEB Faculty ID!")
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
