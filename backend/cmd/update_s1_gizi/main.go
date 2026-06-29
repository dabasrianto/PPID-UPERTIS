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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	// Load env
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

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	fmt.Println("✓ Connected to database")

	visi := `Menjadi Program Studi Sarjana  yang menguasai IPTEKS  inovatif di bidang gizi dan berkarakter tangguh serta berdaya saing di dunia kerja dalam mengembangkan pangan fungsional`

	misi := `1. Menyelenggarakan pendidikan akademik dengan menguasai ipteks inovatif dalam mengembangkan pangan fungsional, yang berkarakter tangguh
2. Melakukan penelitian yang inovatif untuk mengembangkan ipteks pangan fungsional dalam rangka meningkatkan daya saing
3. Melaksanakan pengabdian kepada masyarakat untuk kesejahteraan masyarakat
4. Menjalankan Tri Dharma Perguruan Tinggi Yang Berorientasi Pada Mutu Yang Berkelanjutan Sehingga Mampu Bersaing di  dunia kerja`

	tujuan := `1. Menghasilkan lulusan yang memiliki kemampuan akademik dan menguasai IPTEKS inovatif
2. Menghasilkan lulusan  yang  mampu melakukan penelitian yang inovatif dalam bidang gizi untuk menghasilkan pangan fungsional
3. Menghasilkan lulusan yang mampu mengimplementasi pengetahuan dan keterampilan dalam bidang  gizi kepada masyarakat.
4. Menghasilkan  lulusan berorientasi  Mutu Yang Berkelanjutan Sehingga Mampu Bersaing di  dunia kerja`

	gelar := `Sarjana Gizi (S.Gz) adalah gelar akademis yang diberikan kepada mahasiswa yang telah menyelesaikan program studi tingkat strata satu (S1) di bidang ilmu gizi. Lulusan S1 Gizi memiliki pemahaman mendalam tentang hubungan antara makanan, nutrisi, dan kesehatan manusia, serta dibekali kompetensi untuk memimpin pelayanan asuhan gizi, asuhan dietetik, konseling gizi, kepemimpinan gizi di masyarakat, serta merencanakan dan mengembangkan bisnis makanan fungsional secara profesional.`

	kompetensi := `Profil Lulusan Program Studi Sarjana (S1) Gizi:

1. Manajer dan Care Provider Pelayanan Gizi
Sarjana Gizi yang mampu merencanakan, mengelola, mengembangkan, mengawasi, mengevaluasi dan memiliki keterampilan dalam memilih alternatif atau strategi yang efektif dan efisien dalam memimpin dan menyelenggarakan pelayanan gizi promotif, preventif, kuratif dan rehabilitatif berdasarkan proses asuhan gizi secara profesional.

2. Supervisor Pelayanan Gizi
Sarjana gizi yang memiliki keterampilan untuk mengimplementasikan kekuasaan dan wewenang dalam pemberian upaya perbaikan atau peningkatan gizi dan dietetik masyarakat/kelompok/individu dan pemberian umpan balik melalui rangkaian kegiatan pengkajian, diagnosis, intervensi, monitoring dan evaluasi gizi termasuk kegiatan pendidikan dan pelatihan gizi untuk perubahan pengetahuan, sikap, keterampilan atau perilaku.

3. Inspirator Gizi di Masyarakat (Community Leader)
Sarjana gizi yang mampu menguasai, menanamkan, memanfaatkan ilmu pengetahuan dalam bidang gizi untuk merespon, menggali potensi dan memberikan solusi masalah gizi di masyarakat, dunia usaha dan industri.

4. Peneliti Pembelajar
Sarjana gizi yang mematuhi kaidah-kaidah ilmiah dalam merencanakan, melaksanakan, memanfaatkan, dan mendeseminasikan hasil penelitian sesuai dengan bidang gizi maupun kesiapan untuk belajar pada bidang keahlian yang lain serta memiliki kemampuan menyampaikan gagasannya baik lisan maupun tulisan dalam perspektif global.

5. Entrepreneur
Sarjana gizi yang mampu melakukan perencanaan dan pengembangan bisnis berorientasi pada tindakan, berpikir simpel, mencari dan mengejar peluang usaha yang terbaik dengan disiplin tinggi, mampu mengambil keputusan secara tepat dan berwawasan.`

	strukturOrganisasi := `[
		{"name": "Wilda Laila, SKM, M.Biomed, Dietisien", "role": "Ketua Program Studi", "level": 1, "niidn": "hu-vIbMEh6lc1jKM9B4jrVQ_vmHi_YxJc0luDdkX4EvaQl7BAVrFV9jE6eHOClgRWapnbg=="},
		{"name": "Tika Dwita Adfar, S.ST, M.Biomed, Dietisien", "role": "Sekretaris Prodi", "level": 2, "niidn": "KhQzpBqYzBMXsMpZeeuTtQ3Hu5BYfl5zafNqsqqfDBC06tyLQ-r00MNOqQcIixUDBZrARA=="},
		{"name": "Dezi Ilham, S.Pd, M.Biomed", "role": "Penilaian & Kemahasiswaan", "level": 3, "niidn": "8yj4hBD5T3zBe2K9LT3IPKbANz1uWeK_b1h7IdjdKNtLAcRTRfPdUgtghoTl9zoGl48z-A=="},
		{"name": "Dwi Syaputri Yanti, M.Gizi, S.Tr Gz", "role": "Penilaian & Kemahasiswaan", "level": 3, "niidn": ""},
		{"name": "Erina Masri, M.Biomed, S.K.M.", "role": "Dosen, Tendik, Penelitian dan PKM", "level": 3, "niidn": "ZPsu5FwJy4hjRBrVZqKGTGMNGJk3qqJY_VDARW5f8Absz_Jx6j6pLzAA_1WwYimVKbVElg=="},
		{"name": "Rahmita Yanti, S.KM, M.Kes", "role": "Dosen, Tendik, Penelitian dan PKM", "level": 3, "niidn": "o7lhFPxm1bIya8A22WXl9-Q9fe8NqEGK2NzCDKgJTMZt2hUMShUu70EySvl7td3x2FvRyg=="},
		{"name": "Harleni, S.Pd, M.Pd", "role": "Dosen, Tendik, Penelitian dan PKM", "level": 3, "niidn": "TtSDSF39QRjXj_SNUZM-6FolS0rlRa2p_tx6hj9xGP4_636vzcQ-EScwdfEcmmX-73D6Qg=="},
		{"name": "Nurhamidah, S.KM, M.Biomed", "role": "Sarana Prasarana Pendidikan dan Keuangan", "level": 3, "niidn": "jIEoMFPJxo561Dy7xBYn0Nj-GjBUHSzPqNkWQP30jKrJGH1kyf8-2CmCR8CG4cD3YXkw6g=="},
		{"name": "Widia Dara, SP, MP", "role": "Sarana Prasarana Pendidikan dan Keuangan", "level": 3, "niidn": "9CEJhvtrbGXXvrxAjRG6HL466pUMJHjZlzWS7wak9wcEon3Gv5cuTdKSUpkKBiJNWE77pg=="},
		{"name": "Yeyep Natrio, S.S., M.Hum", "role": "Sarana Prasarana Pendidikan dan Keuangan", "level": 3, "niidn": "RRgXVuk7WZPJiuX5Imp9ZPjbdwdtfBu-286OBoVGLN_yEFHMjI4TpW7qphG1ZQTnHQBV8Q=="}
	]`

	syllabus := `[
		{ "semester": 1, "code": "IEU-103", "name": "Pancasila", "sks": 2, "type": "Wajib" },
		{ "semester": 1, "code": "GZI-261", "name": "Anatomi Manusia", "sks": 3, "type": "Wajib" },
		{ "semester": 1, "code": "GZI-131", "name": "Dasar manajemen", "sks": 2, "type": "Wajib" },
		{ "semester": 1, "code": "GZI-116", "name": "Ilmu Gizi Dasar", "sks": 3, "type": "Wajib" },
		{ "semester": 1, "code": "GZI-121", "name": "Kimia Organik", "sks": 2, "type": "Wajib" },
		{ "semester": 1, "code": "GZI-114", "name": "Pengantar Biologi Manusia", "sks": 2, "type": "Wajib" },
		{ "semester": 1, "code": "GZI-131", "name": "Pengantar Psikologi", "sks": 2, "type": "Wajib" },
		{ "semester": 1, "code": "GZI-112", "name": "Bahasa Indonesia", "sks": 2, "type": "Wajib" },
		{ "semester": 1, "code": "IEU-102", "name": "Kewarganegaraan", "sks": 2, "type": "Wajib" },

		{ "semester": 2, "code": "GZI-104", "name": "Character Building", "sks": 2, "type": "Wajib" },
		{ "semester": 2, "code": "GZI-122", "name": "Kimia Anorganik", "sks": 2, "type": "Wajib" },
		{ "semester": 2, "code": "GZI-262", "name": "Fisiologi Manusia", "sks": 3, "type": "Wajib" },
		{ "semester": 2, "code": "GZI-211", "name": "Dasar Biokimia Gizi", "sks": 3, "type": "Wajib" },
		{ "semester": 2, "code": "GZI-212", "name": "Gizi Dalam Daur Kehidupan", "sks": 3, "type": "Wajib" },
		{ "semester": 2, "code": "GZI-231", "name": "Dasar Kulinari", "sks": 3, "type": "Wajib" },
		{ "semester": 2, "code": "GZI-141", "name": "Pengantar Komunikasi", "sks": 2, "type": "Wajib" },
		{ "semester": 2, "code": "GZI-102", "name": "Agama", "sks": 2, "type": "Wajib" },
		{ "semester": 2, "code": "GZI-105", "name": "Bahasa Inggris", "sks": 3, "type": "Wajib" },
		{ "semester": 2, "code": "GZI-151", "name": "Matematika", "sks": 2, "type": "Wajib" },

		{ "semester": 3, "code": "GZI-361", "name": "Patofisiologi Penyakit Menular", "sks": 3, "type": "Wajib" },
		{ "semester": 3, "code": "GZI-311", "name": "Metabolisme Energi dan Gizi Makro", "sks": 3, "type": "Wajib" },
		{ "semester": 3, "code": "GZI-312", "name": "Metabolisme Zat Gizi Mikro", "sks": 3, "type": "Wajib" },
		{ "semester": 3, "code": "GZI-331", "name": "Ilmu Bahan Pangan", "sks": 3, "type": "Wajib" },
		{ "semester": 3, "code": "GZI-341", "name": "Pengantar Antropologi", "sks": 2, "type": "Wajib" },
		{ "semester": 3, "code": "GZI-342", "name": "Pengantar Sosiologi", "sks": 2, "type": "Wajib" },

		{ "semester": 4, "code": "GZI-461", "name": "Patofisiologi Penyakit Tidak Menular", "sks": 3, "type": "Wajib" },
		{ "semester": 4, "code": "GZI-411", "name": "Penilaian Status Gizi", "sks": 3, "type": "Wajib" },
		{ "semester": 4, "code": "GZI-412", "name": "Penilaian Konsumsi Pangan", "sks": 3, "type": "Wajib" },
		{ "semester": 4, "code": "GZI-431", "name": "Analisa Zat Gizi Pangan", "sks": 3, "type": "Wajib" },
		{ "semester": 4, "code": "GZI-441", "name": "Pendidikan Gizi", "sks": 2, "type": "Wajib" },
		{ "semester": 4, "code": "GZI-402", "name": "Bioetika", "sks": 2, "type": "Wajib" },
		{ "semester": 4, "code": "GZI-451", "name": "Statistik Dasar", "sks": 2, "type": "Wajib" },
		{ "semester": 4, "code": "GZI-481", "name": "Mikrobiologi Pangan", "sks": 2, "type": "Pilihan" },

		{ "semester": 5, "code": "GZI-303", "name": "Kewirausahaan", "sks": 3, "type": "Wajib" },
		{ "semester": 5, "code": "GZI-511", "name": "Dietetika Penyakit Infeksi dan Defisier", "sks": 3, "type": "Wajib" },
		{ "semester": 5, "code": "GZI-512", "name": "Gizi Olahraga", "sks": 2, "type": "Wajib" },
		{ "semester": 5, "code": "GZI-541", "name": "Konsultasi Gizi", "sks": 3, "type": "Wajib" },
		{ "semester": 5, "code": "GZI-531", "name": "Manajemen Sistem Penyelenggaraan Makanan Institusi (MSPMI)", "sks": 3, "type": "Wajib" },
		{ "semester": 5, "code": "GZI-551", "name": "Epidemiologi Gizi", "sks": 3, "type": "Wajib" },
		{ "semester": 5, "code": "GZI-581", "name": "Statistik Lanjut", "sks": 3, "type": "Pilihan" },
		{ "semester": 5, "code": "GZI-582", "name": "Teknologi Pangan", "sks": 3, "type": "Pilihan" },
		{ "semester": 5, "code": "GZI-583", "name": "Pangan Fungsional", "sks": 3, "type": "Pilihan" },

		{ "semester": 6, "code": "GZI-611", "name": "Dietetika Penyakit Tidak Menular", "sks": 3, "type": "Wajib" },
		{ "semester": 6, "code": "GZI-631", "name": "Manajemen Program Gizi", "sks": 2, "type": "Wajib" },
		{ "semester": 6, "code": "GZI-651", "name": "Metode Penelitian Gizi", "sks": 3, "type": "Wajib" },
		{ "semester": 6, "code": "GZI-681", "name": "Pengembangan Media Komunikasi", "sks": 2, "type": "Pilihan" },
		{ "semester": 6, "code": "GZI-682", "name": "Ekonomi Pangan dan Gizi", "sks": 2, "type": "Pilihan" },
		{ "semester": 6, "code": "GZI-683", "name": "Analisa Data Pangan dan Gizi", "sks": 2, "type": "Pilihan" },
		{ "semester": 6, "code": "GZI-684", "name": "Perencanaan Pangan dan Gizi", "sks": 3, "type": "Pilihan" },

		{ "semester": 7, "code": "GZI-304", "name": "PMPKL", "sks": 4, "type": "Wajib" },
		{ "semester": 7, "code": "GZI-791", "name": "Kepaniteraan Dietetik", "sks": 4, "type": "Wajib" },
		{ "semester": 7, "code": "GZI-792", "name": "Kepaniteraan MSPMI", "sks": 3, "type": "Wajib" },
		{ "semester": 7, "code": "GZI-781", "name": "Proposal Skripsi", "sks": 2, "type": "Pilihan" },
		{ "semester": 7, "code": "GZI-782", "name": "Pengawasan Mutu Pangan", "sks": 3, "type": "Pilihan" },
		{ "semester": 7, "code": "GZI-783", "name": "Perkembangan Gizi Terkini", "sks": 2, "type": "Pilihan" },

		{ "semester": 8, "code": "GZI-891", "name": "Skripsi", "sks": 4, "type": "Wajib" }
	]`

	ctx := context.Background()

	// Update S1 Gizi program profile
	res, err := db.ExecContext(ctx, `
		UPDATE faculty_programs
		SET visi = $1,
		    misi = $2,
		    tujuan = $3,
		    gelar_akademik = $4,
		    kompetensi_lulusan = $5,
		    struktur_organisasi = $6,
		    syllabus = $7
		WHERE LOWER(name) LIKE '%s1 gizi%' OR LOWER(slug) = 's1-gizi'
	`, visi, misi, tujuan, gelar, kompetensi, strukturOrganisasi, syllabus)

	if err != nil {
		log.Fatal("Failed to update S1 Gizi program profile: ", err)
	}

	rows, _ := res.RowsAffected()
	fmt.Printf("✓ Updated S1 Gizi program details (%d rows)\n", rows)

	// Update expertise for shared lecturers
	sharedLecturers := []string{
		"Alya Misdhal Rini",
		"Maria Nova",
		"Tika Dwita Adfar",
		"Sepni Asmira",
		"Yensasnidar",
		"Risya Ahriyasna",
	}

	for _, l := range sharedLecturers {
		res2, err := db.ExecContext(ctx, `
			UPDATE faculty_lecturers
			SET expertise = 'Gizi (D3), Gizi (S1)'
			WHERE UPPER(name) LIKE '%' || UPPER($1) || '%'
		`, l)
		if err != nil {
			fmt.Printf("⚠ Failed to update expertise for %s: %v\n", l, err)
		} else {
			cnt, _ := res2.RowsAffected()
			fmt.Printf("✓ Updated expertise for %s (%d rows)\n", l, cnt)
		}
	}

	fmt.Println("Database migration for S1 Gizi finished successfully!")
}
