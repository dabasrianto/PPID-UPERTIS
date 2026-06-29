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
	ctx := context.Background()

	// 1. Check duplicate names (same name in same faculty)
	fmt.Println("=== DUPLIKAT NAMA (nama sama di fakultas sama) ===")
	rows, err := db.QueryContext(ctx, `
		SELECT name, faculty_id, COUNT(*) as cnt
		FROM faculty_lecturers
		GROUP BY UPPER(name), name, faculty_id
		HAVING COUNT(*) > 1
		ORDER BY cnt DESC, name
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	dupNameCount := 0
	for rows.Next() {
		var name, facultyID string
		var cnt int
		rows.Scan(&name, &facultyID, &cnt)
		fmt.Printf("  ⚠ %s (faculty: %s...) -> %d records\n", name, facultyID[:8], cnt)
		dupNameCount++
	}
	if dupNameCount == 0 {
		fmt.Println("  ✓ Tidak ada duplikat nama dalam fakultas yang sama")
	}

	// 2. Check duplicate names across all faculties
	fmt.Println("\n=== DUPLIKAT NAMA (nama sama lintas fakultas) ===")
	rows2, err := db.QueryContext(ctx, `
		SELECT UPPER(name) as uname, COUNT(*) as cnt, COUNT(DISTINCT faculty_id) as fac_cnt
		FROM faculty_lecturers
		GROUP BY UPPER(name)
		HAVING COUNT(*) > 1
		ORDER BY cnt DESC, uname
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows2.Close()
	dupCrossCount := 0
	for rows2.Next() {
		var name string
		var cnt, facCnt int
		rows2.Scan(&name, &cnt, &facCnt)
		fmt.Printf("  ⚠ %s -> %d records across %d faculties\n", name, cnt, facCnt)
		dupCrossCount++
	}
	if dupCrossCount == 0 {
		fmt.Println("  ✓ Tidak ada duplikat nama lintas fakultas")
	}

	// 3. Check duplicate pddikti_id
	fmt.Println("\n=== DUPLIKAT PDDIKTI_ID ===")
	rows3, err := db.QueryContext(ctx, `
		SELECT pddikti_id, COUNT(*) as cnt
		FROM faculty_lecturers
		WHERE pddikti_id IS NOT NULL AND pddikti_id != ''
		GROUP BY pddikti_id
		HAVING COUNT(*) > 1
		ORDER BY cnt DESC
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows3.Close()
	dupPddiktiCount := 0
	for rows3.Next() {
		var pddiktiID string
		var cnt int
		rows3.Scan(&pddiktiID, &cnt)
		fmt.Printf("  ⚠ pddikti_id=%s -> %d records\n", pddiktiID, cnt)
		dupPddiktiCount++
	}
	if dupPddiktiCount == 0 {
		fmt.Println("  ✓ Tidak ada duplikat pddikti_id")
	}

	// 4. Check duplicate slugs
	fmt.Println("\n=== DUPLIKAT SLUG ===")
	rows4, err := db.QueryContext(ctx, `
		SELECT slug, COUNT(*) as cnt
		FROM faculty_lecturers
		WHERE slug IS NOT NULL AND slug != ''
		GROUP BY slug
		HAVING COUNT(*) > 1
		ORDER BY cnt DESC
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows4.Close()
	dupSlugCount := 0
	for rows4.Next() {
		var slug string
		var cnt int
		rows4.Scan(&slug, &cnt)
		fmt.Printf("  ⚠ slug=%s -> %d records\n", slug, cnt)
		dupSlugCount++
	}
	if dupSlugCount == 0 {
		fmt.Println("  ✓ Tidak ada duplikat slug")
	}

	// 5. Show all records with detail for manual check
	fmt.Println("\n=== RINGKASAN ===")
	var total int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM faculty_lecturers").Scan(&total)
	fmt.Printf("Total dosen: %d\n", total)

	var nullPddikti int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM faculty_lecturers WHERE pddikti_id IS NULL OR pddikti_id = ''").Scan(&nullPddikti)
	fmt.Printf("Dosen tanpa pddikti_id: %d\n", nullPddikti)

	// By source
	fmt.Println("\nPer source:")
	rows5, _ := db.QueryContext(ctx, `
		SELECT COALESCE(source, 'manual'), COUNT(*) FROM faculty_lecturers GROUP BY source ORDER BY COUNT(*) DESC
	`)
	defer rows5.Close()
	for rows5.Next() {
		var src string
		var cnt int
		rows5.Scan(&src, &cnt)
		fmt.Printf("  - %s: %d\n", src, cnt)
	}

	// By faculty
	fmt.Println("\nPer fakultas:")
	rows6, _ := db.QueryContext(ctx, `
		SELECT f.name, COUNT(l.id)
		FROM faculty_lecturers l
		JOIN faculties f ON f.id = l.faculty_id
		GROUP BY f.name ORDER BY f.name
	`)
	defer rows6.Close()
	for rows6.Next() {
		var name string
		var cnt int
		rows6.Scan(&name, &cnt)
		fmt.Printf("  - %s: %d dosen\n", name, cnt)
	}
}
