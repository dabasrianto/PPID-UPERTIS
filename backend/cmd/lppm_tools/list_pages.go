//go:build ignore

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
	godotenv.Load(".env")

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "kampuspro")
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

	rows, err := db.QueryContext(context.Background(), "SELECT slug, title FROM pages")
	if err != nil {
		log.Fatal("Failed to query pages table: ", err)
	}
	defer rows.Close()

	outPath := `c:\Users\msii\.gemini\antigravity-ide\brain\467d675d-b553-4e32-8c81-328e54778d39\scratch\pages_in_db.txt`
	f, err := os.Create(outPath)
	if err != nil {
		log.Fatal("Failed to create output file: ", err)
	}
	defer f.Close()

	fmt.Fprintln(f, "Existing Pages in DB:")
	fmt.Println("Existing Pages in DB:")
	for rows.Next() {
		var slug, title string
		if err := rows.Scan(&slug, &title); err == nil {
			fmt.Fprintf(f, "- %s: %s\n", slug, title)
			fmt.Printf("- %s: %s\n", slug, title)
		}
	}
	fmt.Println("Done writing to file.")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
