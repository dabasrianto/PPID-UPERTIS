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
	_ = godotenv.Load("backend/.env")
	_ = godotenv.Load(".env")

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
		log.Fatal("Failed to connect to db: ", err)
	}
	defer db.Close()

	ctx := context.Background()

	fmt.Println("=== Checking Downloads Categories in Local Database ===")
	rows, err := db.QueryContext(ctx, "SELECT id, title, category, active FROM downloads")
	if err != nil {
		log.Fatal("Query failed: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, title, category string
		var active bool
		rows.Scan(&id, &title, &category, &active)
		fmt.Printf("Download: ID=%s, Title=%q, Category=%q, Active=%v\n", id, title, category, active)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
