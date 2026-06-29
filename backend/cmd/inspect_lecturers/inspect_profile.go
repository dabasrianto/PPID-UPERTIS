//go:build ignore

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")
	_ = godotenv.Load("backend/.env")

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "")
	if dbname == "" {
		dbname = "kampuspro"
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

	var name, portfolioStr string
	err = db.QueryRowContext(context.Background(), "SELECT name, COALESCE(portfolio::text, '') FROM faculty_lecturers WHERE portfolio IS NOT NULL AND portfolio::text <> '{}' LIMIT 1").Scan(&name, &portfolioStr)
	if err != nil {
		log.Fatal("Failed to query: ", err)
	}
	fmt.Printf("Lecturer Name: %s\n", name)


	var p map[string]interface{}
	if err := json.Unmarshal([]byte(portfolioStr), &p); err != nil {
		log.Fatal("Unmarshal err: ", err)
	}

	profile := p["profile"]
	profileJSON, _ := json.MarshalIndent(profile, "", "  ")
	fmt.Println(string(profileJSON))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
