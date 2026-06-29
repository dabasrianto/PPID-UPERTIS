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
	// Load .env from current directory or backend directory
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../.env")
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

	// Menu JSON structure containing all 12 items under Profil group
	newMenuJSON := `[
  {"group": "Profil", "items": [
    {"href": "/halaman/sambutan-rektor", "label": "Sambutan Rektor"},
    {"href": "/halaman/visi-misi", "label": "Visi & Misi"},
    {"href": "/halaman/video-tour", "label": "Video Tour Kampus"},
    {"href": "/halaman/sejarah", "label": "Perjalanan Sejarah"},
    {"href": "/halaman/pimpinan-manajemen", "label": "Pimpinan & Manajemen"},
    {"href": "/halaman/renstra", "label": "Rencana Strategis"},
    {"href": "/halaman/kontak", "label": "Kontak & Lokasi"},
    {"href": "/halaman/panduan-lambang", "label": "Panduan Lambang"},
    {"href": "/halaman/pedoman-media-sosial", "label": "Pedoman Media Sosial"},
    {"href": "/halaman/tata-nilai", "label": "7 Tata Nilai UPERTIS"},
    {"href": "/halaman/fasilitas", "label": "Fasilitas Kampus"},
    {"href": "/halaman/tugas-fungsi", "label": "Tugas & Fungsi"}
  ]},
  {"group": "Pendidikan", "href": "/fakultas", "type": "faculties", "items": []},
  {"group": "Kalender Akademik", "items": [
    {
      "label": "4.1 Penelitian",
      "href": "/penelitian",
      "items": [
        {"label": "Pusat Penelitian", "href": "/penelitian"},
        {
          "label": "Profil LPPM",
          "href": "/halaman/profil-lppm",
          "items": [
            {"label": "LPPM", "href": "/halaman/profil-lppm"},
            {"label": "Jurnal UPERTIS", "href": "https://jurnal.upertis.ac.id/"},
            {"label": "Pusat Studi", "href": "/penelitian"},
            {"label": "HKI & Paten", "href": "/penelitian"}
          ]
        },
        {"label": "Laboratorium Penelitian", "href": "/halaman/laboratorium"},
        {"label": "Rencana Penelitian", "href": "/penelitian"}
      ]
    },
    {
      "label": "4.2 Pengabdian",
      "href": "/pengabdian",
      "items": [
        {"label": "Pusat Pengabdian", "href": "/pengabdian"},
        {"label": "Rencana Pengabdian", "href": "/pengabdian"}
      ]
    }
  ]},
  {"group": "PMB", "items": [
    {"href": "/pmb#informasi", "label": "Informasi Pendaftaran"},
    {"href": "/pmb#jalur", "label": "Jalur Masuk"},
    {"href": "/pmb#biaya", "label": "Biaya Kuliah"},
    {"href": "/pmb#beasiswa", "label": "Beasiswa"},
    {"href": "/pmb#alur", "label": "Alur Pendaftaran"},
    {"href": "/pmb#penting", "label": "Informasi Penting"},
    {"href": "/pmb/daftar", "label": "Daftar Sekarang"},
    {"href": "/pmb#testimoni", "label": "Testimoni Alumni"},
    {"href": "/pmb/minat-bakat", "label": "Game Tes Minat Bakat"}
  ]},
  {"group": "Layanan", "items": [
    {
      "label": "Kemahasiswaan",
      "href": "/layanan/kemahasiswaan",
      "items": [
        {"label": "Organisasi Mahasiswa (BEM, DLM, HIMA)", "href": "/layanan/kemahasiswaan#organisasi"},
        {"label": "UKM (Seni, Olahraga, Kerohanian)", "href": "/layanan/kemahasiswaan#ukm"},
        {"label": "Prestasi Mahasiswa", "href": "/layanan/kemahasiswaan#prestasi"},
        {"label": "Bimbingan Konseling", "href": "/layanan/kemahasiswaan#konseling"}
      ]
    },
    {
      "label": "Alumni",
      "href": "/layanan/alumni",
      "items": [
        {"label": "Profil Tracer Study", "href": "/layanan/alumni#tracer"},
        {"label": "Karir & Alumni", "href": "/layanan/alumni#karir"}
      ]
    },
    {
      "label": "Humas dan Kerjasama",
      "href": "/layanan/humas-dan-kerjasama",
      "items": [
        {"label": "Profil", "href": "/layanan/humas-dan-kerjasama#profil"},
        {"label": "Mitra", "href": "/layanan/humas-dan-kerjasama#mitra"},
        {"label": "Layanan Kerjasama", "href": "/layanan/humas-dan-kerjasama#layanan"},
        {"label": "Website Kerjasama", "href": "https://kerjasama.upertis.ac.id/"}
      ]
    },
    {
      "label": "Penjaminan Mutu",
      "href": "/layanan/penjaminan-mutu",
      "items": [
        {"label": "Profil", "href": "/layanan/penjaminan-mutu#profil"},
        {"label": "Laporan Penjaminan Mutu", "href": "/layanan/penjaminan-mutu#laporan"},
        {"label": "Website P2AMIA", "href": "http://p2amia.upertis.ac.id/"}
      ]
    }
  ]},
  {"group": "Download", "items": [
    {"href": "/media", "label": "Galeri"},
    {"href": "/download", "label": "Download"},
    {"href": "/halaman/live-sosmed", "label": "Live & Sosial Media"}
  ]}
]`

	_, err = db.ExecContext(ctx, "UPDATE site_settings SET value = $1 WHERE key = 'menu'", newMenuJSON)
	if err != nil {
		log.Fatal("Failed to sync menu setting in database:", err)
	}
	fmt.Println("Successfully synced database menu settings with the complete 12 Profil submenus!")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
