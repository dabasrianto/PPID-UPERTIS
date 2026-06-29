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

	ctx := context.Background()

	// Perform regex update on the JSON string value of key 'menu' in site_settings
	res, err := db.ExecContext(ctx, `
		UPDATE site_settings 
		SET value = regexp_replace(value::text, '"/layanans"', '"/layanan"', 'g')::jsonb 
		WHERE key = 'menu'
	`)
	if err != nil {
		log.Fatal("Migration query failed: ", err)
	}

	rows, _ := res.RowsAffected()
	fmt.Printf("Migration completed successfully! Affected rows: %d\n", rows)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
