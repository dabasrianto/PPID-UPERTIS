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

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if dbname == "" {
		dbname = "kampuspro"
	}
	if host == "" { host = "localhost" }
	if port == "" { port = "5432" }
	if user == "" { user = "postgres" }
	if password == "" { password = "postgres" }
	if sslmode == "" { sslmode = "disable" }

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queries := []string{
		"Elsa", "Shinta", "Merisa", "Elvi", "Dina", "Suryani",
	}

	fmt.Println("=== PENCARIAN DOSEN MIRIP DI DATABASE ===")
	for _, q := range queries {
		fmt.Printf("\nQuery: %s\n", q)
		rows, err := db.QueryContext(context.Background(), 
			"SELECT name, COALESCE(pddikti_id, '') FROM faculty_lecturers WHERE name ILIKE $1", 
			"%"+q+"%")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		found := false
		for rows.Next() {
			var name, nidn string
			rows.Scan(&name, &nidn)
			fmt.Printf("  -> Nama: %s | NIDN: %s\n", name, nidn)
			found = true
		}
		if !found {
			fmt.Println("  (Tidak ditemukan nama mirip)")
		}
	}
}
