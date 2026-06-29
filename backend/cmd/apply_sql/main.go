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
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run backend/cmd/apply_sql/main.go <path_to_sql_file>")
	}
	sqlPath := os.Args[1]

	// Find and load .env file
	loadEnv()

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

	content, err := os.ReadFile(sqlPath)
	if err != nil {
		log.Fatalf("Failed to read SQL file %s: %v", sqlPath, err)
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Fatal("Failed to begin transaction: ", err)
	}

	fmt.Printf("Executing SQL from %s on database %s...\n", sqlPath, dbname)
	_, err = tx.ExecContext(context.Background(), string(content))
	if err != nil {
		tx.Rollback()
		log.Fatal("Failed to execute SQL: ", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("Failed to commit transaction: ", err)
	}

	fmt.Println("Successfully applied SQL updates!")
}

func loadEnv() {
	// Try multiple locations to find .env file
	paths := []string{
		".env",
		"backend/.env",
		"../.env",
		"../../.env",
		"../../../.env",
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			if err := godotenv.Load(p); err == nil {
				fmt.Printf("Loaded environment from %s\n", p)
				return
			}
		}
	}
	log.Println("Warning: No .env file found, using system environment variables")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
