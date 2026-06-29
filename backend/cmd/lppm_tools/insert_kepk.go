//go:build ignore

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
	godotenv.Load(".env")

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "kampuspro")
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

	content := `KEPK (Komite Etik Penelitian Kesehatan) Universitas Perintis Indonesia adalah badan resmi yang menilai kelayakan etik proposal penelitian di bidang kesehatan, kedokteran, dan farmasi. Komite ini memegang peranan penting untuk memastikan perlindungan hak asasi dan kesejahteraan subjek penelitian.

---

## Informasi Utama & Kontak KEPK UPERTIS

| Informasi / Kontak | Detail |
| :--- | :--- |
| **No. Registrasi KEPPKN Kemenkes RI** | 0116221371 |
| **Standar Prosedur** | Mengacu pada standar CIOMS-WHO 2016 |
| **Alamat Sekretariat** | Kampus 1 Universitas Perintis Indonesia, Jl. Adinegoro KM.17, Lubuk Buaya, Padang, Sumatera Barat |
| **Email** | [ethics.upertis@gmail.com](mailto:ethics.upertis@gmail.com) |
| **No. Telepon / WhatsApp** | [+62 813-4830-5867](https://wa.me/6281348305867) |

---

## Proses Pengajuan *Ethical Clearance* (Kelayakan Etik)

Proses pengajuan *Ethical Clearance* di KEPK UPERTIS umumnya diperuntukkan bagi mahasiswa atau peneliti yang melibatkan:
* **Subjek Manusia** (penelitian eksperimental, survei klinis, dsb.)
* **Hewan Coba**
* **Data Rekam Medis**

### Berkas yang Harus Dipersiapkan:
1. Protokol Etik Penelitian
2. Formulir Persetujuan (*Informed Consent*)
3. Instrumen Penelitian (kuesioner, lembar observasi, dsb.)

---

## Tautan Pengajuan Resmi

Silakan ajukan kelayakan etik penelitian Anda melalui tautan di bawah ini:

👉 **[KLIK DI SINI UNTUK PENGAJUAN (bit.ly/EthicsUPERTIS)](http://bit.ly/EthicsUPERTIS)**`

	_, err = db.ExecContext(context.Background(), `
		INSERT INTO pages (slug, title, subtitle, content, published, sort_order, seo_title, seo_description, created_at, updated_at)
		VALUES ('komite-etik-penelitian', 'Komite Etik Penelitian Kesehatan', 'Komite Etik Penelitian Kesehatan (KEPK) Universitas Perintis Indonesia', $1, true, 20, 'Komite Etik Penelitian Kesehatan — Universitas Perintis Indonesia', 'Komite Etik Penelitian Kesehatan (KEPK) Universitas Perintis Indonesia', NOW(), NOW())
		ON CONFLICT (slug)
		DO UPDATE SET title='Komite Etik Penelitian Kesehatan', subtitle='Komite Etik Penelitian Kesehatan (KEPK) Universitas Perintis Indonesia', content=$1, published=true, updated_at=NOW()
	`, content)
	if err != nil {
		log.Fatal("Failed to insert/update komite-etik-penelitian page: ", err)
	}

	fmt.Println("Successfully inserted/updated KEPK page in the database!")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
