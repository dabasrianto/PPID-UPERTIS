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
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()

	// 1. Get Afri Andika details
	var facID, facName string
	err = db.QueryRowContext(ctx, `
		SELECT l.faculty_id, f.name 
		FROM faculty_lecturers l 
		JOIN faculties f ON f.id = l.faculty_id 
		WHERE l.name ILIKE '%AFRI ANDIKA%'
	`).Scan(&facID, &facName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Afri Andika's Faculty in Database: %s (ID: %s)\n", facName, facID)

	// 2. Print all faculties in database
	fmt.Println("\n=== Faculties in Database ===")
	rows, err := db.QueryContext(ctx, "SELECT id, code, name FROM faculties")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id, code, name string
		rows.Scan(&id, &code, &name)
		fmt.Printf("- %s (%s): %s\n", code, name, id)
	}
	rows.Close()
}
