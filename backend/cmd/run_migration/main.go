package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env
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

	fmt.Printf("Connecting to database: host=%s, port=%s, user=%s, dbname=%s\n", host, port, user, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database connection: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database: ", err)
	}
	fmt.Println("Connected to database successfully.")

	// Read sync_data.sql
	sqlFilePaths := []string{
		"sync_data.sql",
		"backend/sync_data.sql",
		"../backend/sync_data.sql",
	}

	var sqlContent []byte
	var loadedPath string
	for _, path := range sqlFilePaths {
		content, err := ioutil.ReadFile(path)
		if err == nil {
			sqlContent = content
			loadedPath = path
			break
		}
	}

	if len(sqlContent) == 0 {
		log.Fatal("Could not find sync_data.sql in standard paths")
	}

	fmt.Printf("Loaded sql file from %s (%d bytes). Running migration...\n", loadedPath, len(sqlContent))

	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	_, err = db.ExecContext(ctx, string(sqlContent))
	if err != nil {
		log.Fatal("Failed to execute SQL migration: ", err)
	}

	fmt.Printf("Migration completed successfully in %v!\n", time.Since(startTime))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
