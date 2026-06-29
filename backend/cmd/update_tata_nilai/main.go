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
	// Load .env from current directory or parent directory
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("backend/.env")
		if err != nil {
			log.Println("Warning: No .env file found, using defaults/environment variables")
		}
	}

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

	ctx := context.Background()

	// JSON structure containing the updated 7 values matching the UPERTIS flyer
	newValueJSON := `{
  "intro": "Tujuh Tata Nilai UPERTIS merupakan fondasi karakter yang menjiwai seluruh aktivitas akademik dan non-akademik di lingkungan Universitas Perintis Indonesia.",
  "values": [
    {
      "name": "Unggul",
      "description": "Bahwa lulusan memiliki kemampuan berdaya saing pada taraf ASEAN dan sukses dalam karirnya.",
      "icon": "🏆"
    },
    {
      "name": "Profesional",
      "description": "Merupakan kemampuan lulusan dalam melakukan unjuk kerja dengan menggunakan konsep, teori, metode, bahan, dan/atau instrument, yang diperoleh melalui pengalaman pembelajaran.",
      "icon": "💼"
    },
    {
      "name": "Etika",
      "description": "Bahwa lulusan memiliki nilai, norma dan etika akademik.",
      "icon": "⚖️"
    },
    {
      "name": "Rasional",
      "description": "Bahwa lulusan memiliki sikap yang rasional untuk bertindak dalam berbagai lapisan masyarakat.",
      "icon": "🧠"
    },
    {
      "name": "Tangguh",
      "description": "Bahwa lulusan memiliki karakter yang mampu beradaptasi dalam berbagai situasi dan kondisi, ulet dalam berbagai usaha dan upaya, serta dinamis terhadap perubahan lingkungan strategis",
      "icon": "💪"
    },
    {
      "name": "Integritas",
      "description": "Bahwa lulusan memiliki nilai kejujuran ilmiah dalam bertindak sehingga karyanya dapat dipercaya.",
      "icon": "💎"
    },
    {
      "name": "Solidaritas",
      "description": "Bahwa lulusan menghargai perbedaan dan keragaman ras, suku, agama, budaya, dan status sosial dalam kehidupan bermasyarakat.",
      "icon": "👥"
    }
  ]
}`

	_, err = db.ExecContext(ctx, `
		INSERT INTO site_settings (key, value, updated_at)
		VALUES ('page_tata_nilai', $1, NOW())
		ON CONFLICT (key) 
		DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
	`, newValueJSON)
	if err != nil {
		log.Fatal("Failed to update page_tata_nilai in database:", err)
	}
	fmt.Println("Successfully updated page_tata_nilai settings with the UPERTIS acronym and correct flyer values!")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
