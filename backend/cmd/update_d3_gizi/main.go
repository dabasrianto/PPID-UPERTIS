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
	// Load .env
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("backend/.env")
		if err != nil {
			err = godotenv.Load("../.env")
			if err != nil {
				log.Println("Warning: No .env file found, using defaults")
			}
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

	visi := `Menjadi Program Studi Diploma yang menguasai IPTEKS Inovatif, kolaboratif, mandiri di bidang gizi dan berkarakter tangguh serta berdaya saing di dunia kerja dalam pengembangan produk pangan`

	misi := `1. Menyelenggarakan Pendidikan dan Proses Pembelajaran untuk Menghasilkan Sumber daya manusia di Bidang Gizi yang terampil dalam pengembangan produk pangan.
2. Menyelenggarakan Penelitian Ilmiah dalam mengembangkan Ilmu Pengetahuan dan teknologi di bidang Gizi yang terampil dalam pengembangan produk pangan.
3. Menyelenggarakan Pengabdian Masyarakat Dalam Upaya Pemberdayaan Masyarakat Berbasis Hasil Penelitian
4. Menjalin kerjasama secara aktif dengan organisasi profesi gizi, institusi pemerintah, swasta, nasional dan internasional dalam menunjang Tridharma Perguruan Tinggi
5. Menjalankan Sistim Tri Dharma Pergutuan Tinggi Yang Berorientasi Pada Mutu Yang Berkelanjutan Sehingga Mampu Bersaing dunia kerja`

	tujuan := `1. Menghasilkan lulusan yang unggul, inovatif ,berkarakter tangguh dan mandiri yang mampu menerapkan ilmu pengetahuan dan teknologi yang terampil dalam pengembangan produk pangan yang mampu bersaing di dunia kerja.
2. Menghasilkan penelitian ilmiah yang berbasis unggul, inovatif, kolaboratif dan berintegrasi terampil dalam pengembangan produk pangan yang dapat menjadi rujukan nasional dan internasioanal
3. Mewujudkan kegiatan pengabdian masyarakat dibidang gizi yang terampil dalam pengembangan produk pangan serta memperdayakan masyarakat secara inovatif agar mampu menyelesaikan masalah secara mandiri dan berkelanjutan
4. Mengembangkan informasi yang diperoleh dari asosiasi dan profesi lain untuk menuju tata kelola Program Studi yang lebih baik.
5. Menghasilkan kerjasama dengan lembaga pemerintah dan swasta dalam bidang Tridharma Perguruan Tinggi`

	gelar := `Ahli Madya Gizi (AMd.Gz) adalah seorang yang telah mengikuti dan menyelesaikan Program Pendidikan Diploma III Gizi sesuai aturan yang berlaku, mempunyai tugas, tanggung jawab dan wewenang yang didasari dengan penguasaan ilmu gizi yang memadai, untuk melakukan kegiatan fungsional dalam bidang pelayanan gizi, makanan dan dietetik baik di masyarakat, individu, atau institusi. serta mempunyai kemampuan manajerial dan organisasi yang sesuai kewenangannya tugasnya.`

	kompetensi := `Capaian Pembelajaran (Learning outcome) Program Studi Diploma III Gizi:

CPL 1. Mampu melaksanakan pelayanan asuhan gizi pada daur kehidupan sesuai standar profesi dengan sikap bertanggung jawab dan taat hukum.
CPL 2. Mampu melaksanakan program gizi masyarakat sesuai dengan standar profesi dilandasi kerjasama lintas program dan lintas sektor.
CPL 3. Mampu melaksanakan penyelenggaraan makanan institusi dengan keunggulan pengembangan produk pangan sesuai dengan standar profesi dan menunjukkan sikap tanggung jawab
CPL 4. Mampu melaksanakan promosi gizi untuk pengendalian faktor risiko sesuai standar profesi dengan menunjukkan sikap menghargai keanekaragaman budaya.
CPL 5. Mampu menerapkan literasi IPTEKS gizi dalam pengendalian faktor risiko masalah gizi sesuai standar profesi secara professional.

Profil Lulusan Prodi D III Gizi:
Ahli Madya Gizi (Nutrisionis Terampil), sebagai :

1. Pelaksana Asuhan Gizi
Memberikan pelayanan asuhan gizi dan konseling gizi pada klien sehat dan klien kategori khusus (ibu hamil, ibu menyusui, bayi, balita, anak, remaja, dewasa, lansia) secara individu, kelompok, dan masyarakat dengan menggunakan prosedur skrining gizi dan PAGT.
2. Pelaksana Program Gizi Masyarakat
Memberikan pelayanan program gizi masyarakat secara promotif, preventif, kuratif, dan rehabilitatif sesuai dengan prosedur baku dan mekanisme yang telah ditetapkan, pada individu, maupun kelompok, dan masyarakat, dengan keunggulan promotive dan preventif.
3. Pelaksana Penyelenggaraan Makanan Institusi
Memberikan pelayanan penyelenggaraan makanan institusi untuk pemenuhan kebutuhan gizi dan dietetik dalam kondisi normal maupun darurat yang meliputi matra darat, laut, dan udara, dengan menggunakan prosedur baku yang telah ditetapkan, pada klien individu, kelompok, masyarakat.
4. Edukator dan Komunikator Gizi
Menggali dan bertukar informasi secara verbal dan nonverbal, dengan memanfaatkan berbagai bentuk media, dalam memberikan promosi dan edukasi, gizi, pangan dan kesehatan kepada klien individu, kelompok, masyarakat, menggunakan Ilmu Pengetahuan Teknologi & Seni (IPTEKS) mutakhir, dengan keunggulan promotive dan preventif.`

	ctx := context.Background()
	
	// Let's check if D III Gizi exists
	var exists bool
	err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM faculty_programs WHERE LOWER(name) LIKE '%d3 gizi%' OR LOWER(name) LIKE '%d iii gizi%' OR LOWER(name) LIKE '%gizi%')").Scan(&exists)
	if err != nil {
		log.Fatal("Failed to query program existence:", err)
	}

	if !exists {
		log.Fatal("Program study 'D III Gizi' not found in database. Please add it first.")
	}

	strukturOrganisasi := `[
		{"role": "Ka. Prodi", "name": "Alya Misdhal Rini, S. Gz, M. Biomed", "niidn": "1001017604", "level": 1},
		{"role": "Sekretaris", "name": "Maria Nova, SKM. M. Kes", "niidn": "1023118301", "level": 2},
		{"role": "Pj. Kurikulum, Pengajaran dan Evaluasi", "name": "Sepni Asmira, STP, MP", "niidn": "1024097801", "level": 3},
		{"role": "Pj. PKL, Penelitian, Pengabmas dan KTI", "name": "Tika Dwita Adfar, M.Biomed", "niidn": "1018039001", "level": 3},
		{"role": "PJ. Kemahasiswaan dan Alumni", "name": "Yensasnidar, S.Gz. M.Pd", "niidn": "1016076701", "level": 3}
	]`

	result, err := db.ExecContext(ctx, `
		UPDATE faculty_programs
		SET visi = $1,
		    misi = $2,
		    tujuan = $3,
		    gelar_akademik = $4,
		    kompetensi_lulusan = $5,
		    struktur_organisasi = $6
		WHERE LOWER(name) LIKE '%d3 gizi%' OR LOWER(name) LIKE '%d iii gizi%' OR LOWER(name) LIKE '%gizi%'
	`, visi, misi, tujuan, gelar, kompetensi, strukturOrganisasi)

	if err != nil {
		log.Fatal("Failed to update database record:", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Successfully updated D III Gizi program data in the database (%d row(s) updated)!\n", rowsAffected)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
