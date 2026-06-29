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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	_ = godotenv.Load(".env")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "")
	if dbname == "" { log.Fatal("DB_NAME required") }
	sslmode := getEnv("DB_SSLMODE", "disable")

	db, _ := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode))
	defer db.Close()
	ctx := context.Background()

	fmt.Println("=== Program Studi (faculty_programs) ===")
	rows, _ := db.QueryContext(ctx, `
		SELECT fp.id, fp.name, fp.level, fp.slug, f.name as faculty_name
		FROM faculty_programs fp
		JOIN faculties f ON f.id = fp.faculty_id
		WHERE fp.active = true
		ORDER BY f.name, fp.name
	`)
	defer rows.Close()
	for rows.Next() {
		var id, name, level, slug, fac string
		rows.Scan(&id, &name, &level, &slug, &fac)
		fmt.Printf("  [%s] %s (%s) -> slug: %s, fac: %s\n", id[:8], name, level, slug, fac)
	}

	fmt.Println("\n=== Expertise Values di faculty_lecturers ===")
	rows2, _ := db.QueryContext(ctx, `
		SELECT COALESCE(expertise,'(kosong)'), COUNT(*) 
		FROM faculty_lecturers 
		GROUP BY expertise 
		ORDER BY expertise
	`)
	defer rows2.Close()
	for rows2.Next() {
		var exp string
		var cnt int
		rows2.Scan(&exp, &cnt)
		fmt.Printf("  %s: %d dosen\n", exp, cnt)
	}
}
