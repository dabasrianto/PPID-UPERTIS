package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("backend/.env")

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
		log.Fatal("Failed to connect to local db: ", err)
	}
	defer db.Close()

	ctx := context.Background()
	var sqlBuilder strings.Builder

	sqlBuilder.WriteString("-- Auto-generated PPID Sync Script\n")
	sqlBuilder.WriteString("-- Generated on: " + time.Now().Format(time.RFC1123) + "\n\n")
	sqlBuilder.WriteString("BEGIN;\n\n")

	// 1. Export downloads
	sqlBuilder.WriteString("-- ==========================================\n")
	sqlBuilder.WriteString("-- DOWNLOADS TABLE\n")
	sqlBuilder.WriteString("-- ==========================================\n")
	sqlBuilder.WriteString("TRUNCATE TABLE downloads CASCADE;\n\n")

	downRows, err := db.QueryContext(ctx, "SELECT id, title, description, file_url, category, active, downloads_count, created_at, updated_at FROM downloads")
	if err != nil {
		log.Fatal("Failed to query downloads: ", err)
	}
	defer downRows.Close()

	for downRows.Next() {
		var id, title, description, fileURL, category string
		var active bool
		var downloadsCount int
		var createdAt, updatedAt time.Time

		err := downRows.Scan(&id, &title, &description, &fileURL, &category, &active, &downloadsCount, &createdAt, &updatedAt)
		if err != nil {
			log.Fatal("Failed to scan download row: ", err)
		}

		sqlBuilder.WriteString(fmt.Sprintf(
			"INSERT INTO downloads (id, title, description, file_url, category, active, downloads_count, created_at, updated_at) VALUES (%s, %s, %s, %s, %s, %t, %d, '%s', '%s');\n",
			quoteString(id), quoteString(title), quoteString(description), quoteString(fileURL), quoteString(category),
			active, downloadsCount, createdAt.Format("2006-01-02 15:04:05.999999"), updatedAt.Format("2006-01-02 15:04:05.999999"),
		))
	}

	// 2. Export pages
	sqlBuilder.WriteString("\n-- ==========================================\n")
	sqlBuilder.WriteString("-- PAGES TABLE\n")
	sqlBuilder.WriteString("-- ==========================================\n")
	sqlBuilder.WriteString("TRUNCATE TABLE pages CASCADE;\n\n")

	pageRows, err := db.QueryContext(ctx, "SELECT id, title, subtitle, slug, content, cover_image_url, published, sort_order, seo_title, seo_description, created_at, updated_at FROM pages")
	if err != nil {
		log.Fatal("Failed to query pages: ", err)
	}
	defer pageRows.Close()

	for pageRows.Next() {
		var id, title, subtitle, slug, content, coverImageURL, seoTitle, seoDescription string
		var published bool
		var sortOrder int
		var createdAt, updatedAt time.Time

		var subVal, covVal, seoTVal, seoDVal sql.NullString
		err := pageRows.Scan(&id, &title, &subVal, &slug, &content, &covVal, &published, &sortOrder, &seoTVal, &seoDVal, &createdAt, &updatedAt)
		if err != nil {
			log.Fatal("Failed to scan page row: ", err)
		}

		if subVal.Valid {
			subtitle = subVal.String
		}
		if covVal.Valid {
			coverImageURL = covVal.String
		}
		if seoTVal.Valid {
			seoTitle = seoTVal.String
		}
		if seoDVal.Valid {
			seoDescription = seoDVal.String
		}

		sqlBuilder.WriteString(fmt.Sprintf(
			"INSERT INTO pages (id, title, subtitle, slug, content, cover_image_url, published, sort_order, seo_title, seo_description, created_at, updated_at) VALUES (%s, %s, %s, %s, %s, %s, %t, %d, %s, %s, '%s', '%s');\n",
			quoteString(id), quoteString(title), quoteOrNull(subtitle), quoteString(slug), quoteString(content),
			quoteOrNull(coverImageURL), published, sortOrder, quoteOrNull(seoTitle), quoteOrNull(seoDescription),
			createdAt.Format("2006-01-02 15:04:05.999999"), updatedAt.Format("2006-01-02 15:04:05.999999"),
		))
	}

	// 3. Export site_settings
	sqlBuilder.WriteString("\n-- ==========================================\n")
	sqlBuilder.WriteString("-- SITE SETTINGS TABLE\n")
	sqlBuilder.WriteString("-- ==========================================\n")

	settingRows, err := db.QueryContext(ctx, "SELECT key, value, updated_at FROM site_settings")
	if err != nil {
		log.Fatal("Failed to query site_settings: ", err)
	}
	defer settingRows.Close()

	for settingRows.Next() {
		var key string
		var valueRaw []byte
		var updatedAt time.Time

		err := settingRows.Scan(&key, &valueRaw, &updatedAt)
		if err != nil {
			log.Fatal("Failed to scan site_setting row: ", err)
		}

		valueStr := string(valueRaw)

		sqlBuilder.WriteString(fmt.Sprintf(
			"INSERT INTO site_settings (key, value, updated_at) VALUES (%s, %s, '%s') ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = EXCLUDED.updated_at;\n",
			quoteString(key), quoteString(valueStr), updatedAt.Format("2006-01-02 15:04:05.999999"),
		))
	}

	sqlBuilder.WriteString("\nCOMMIT;\n")

	outputPath := "sync_data.sql"
	err = ioutil.WriteFile(outputPath, []byte(sqlBuilder.String()), 0644)
	if err != nil {
		log.Fatal("Failed to write SQL file: ", err)
	}

	fmt.Println("Successfully generated clean, self-replacing sync_data.sql file with pages, downloads, and site_settings tables!")
}

func quoteString(s string) string {
	escaped := strings.ReplaceAll(s, "'", "''")
	return fmt.Sprintf("'%s'", escaped)
}

func quoteOrNull(s string) string {
	if s == "" {
		return "NULL"
	}
	return quoteString(s)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
