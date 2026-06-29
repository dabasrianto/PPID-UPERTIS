package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=kampuspro sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var id, title, subtitle, slug, content string
	var published bool
	err = db.QueryRowContext(context.Background(),
		"SELECT id, title, subtitle, slug, content, published FROM pages WHERE slug=$1",
		"profil-lppm").Scan(&id, &title, &subtitle, &slug, &content, &published)

	if err != nil {
		log.Fatal("Query error: ", err)
	}

	fmt.Printf("ID: %s\nTitle: %s\nSubtitle: %s\nSlug: %s\nPublished: %t\nContent Length: %d\n", id, title, subtitle, slug, published, len(content))
}
