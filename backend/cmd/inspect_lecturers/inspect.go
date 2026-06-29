//go:build ignore

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

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

	rows, err := db.QueryContext(context.Background(), "SELECT id, name, active, pddikti_id, portfolio FROM faculty_lecturers ORDER BY name")
	if err != nil {
		log.Fatal("Failed to query: ", err)
	}
	defer rows.Close()

	fmt.Printf("%-35s | %-6s | %-12s | %s\n", "NAME", "ACTIVE", "PDDIKTI_ID", "PORTFOLIO")
	fmt.Println(strings.Repeat("-", 75))
	for rows.Next() {
		var id, name string
		var active bool
		var pddiktiId sql.NullString
		var portfolio sql.NullString
		
		err := rows.Scan(&id, &name, &active, &pddiktiId, &portfolio)
		if err != nil {
			fmt.Printf("SCAN ERROR for row: %v\n", err)
			continue
		}
		
		portStatus := "Populated"
		if !portfolio.Valid || portfolio.String == "" || portfolio.String == "{}" {
			portStatus = "EMPTY"
		}
		
		dispId := "-"
		if pddiktiId.Valid {
			dispId = pddiktiId.String
			if len(dispId) > 15 {
				dispId = dispId[:12] + "..."
			}
		}
		
		fmt.Printf("%-35s | %-6t | %-12s | %s\n", name, active, dispId, portStatus)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
