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

	content := `Lembaga Penelitian dan Pengabdian Masyarakat Universitas Perintis Indonesia (LPPM UPERTIS) adalah unsur pelaksana akademik di lingkungan UPERTIS yang bertugas mengoordinasikan, memantau, dan mengevaluasi pelaksanaan kegiatan penelitian serta pengabdian kepada masyarakat.

Lembaga ini berfungsi sebagai wadah untuk mendorong para dosen dan civitas akademika dalam mengembangkan IPTEKS inovatif yang tangguh dan kompetitif di tingkat nasional.

---

## Visi & Misi LPPM UPERTIS

### Visi
> Menjadi pusat lembaga kegiatan penelitian dan pengabdian masyarakat yang **Unggul, Inovatif serta Mampu Bersaing secara Nasional dan Internasional** khususnya dalam bidang kesehatan dan bencana, obat-obatan, ilmu komunikasi, dan bisnis digital.

### Misi
1. **Mengembangkan relevansi** penelitian dan pengabdian kepada masyarakat untuk meningkatkan mutu pendidikan.
2. **Mengembangkan penelitian dan pengabdian** kepada masyarakat berbasis kesehatan dan bencana, obat-obatan, ilmu komunikasi dan bisnis digital.
3. **Mengembangkan publikasi dan luaran** penelitian dalam bentuk jurnal nasional, internasional, perolehan Hak Kekayaan Intelektual (HKI), buku teks, serta seminar nasional dan internasional.
4. **Mengembangkan hilirisasi** penelitian dan pengabdian sebagai produk komersial.

---

## Struktur Organisasi LPPM

Berikut adalah susunan kepengurusan Lembaga Penelitian dan Pengabdian kepada Masyarakat (LPPM) UPERTIS:

| Jabatan | Nama Pejabat |
| :--- | :--- |
| **Ketua LPPM** | Ns. Falerisiska Yunere, M.Kep |
| **Pj. Penelitian** | Sri Indrayati, M.Si |
| **Pj. Pengabdian Masyarakat** | Ns. Falerisiska Yunere, M.Kep |
| **Publikasi & Kegiatan Ilmiah** | Feni Wartisa, S.Sit, M.KM / Epi Supriadi, M.Si |
| **HKI, Inovasi & KATSINOV** | Ria Afrianti, M.Si, Apt |
| **Komisi Etik Penelitian** | Defrimal, S.Kp, M.Biomed, PAK |

---

## Uraian Wewenang, Tanggung Jawab, dan Tugas

### 1. Ketua LPPM
* **Nama**: Ns. Falerisiska Yunere, M.Kep
* **Wewenang**:
  - Menyusun program kerja Lembaga Penelitian dan Pengabdian kepada masyarakat sebagai pedoman pelaksanaan tugas.
  - Menyusun Rencana Induk Penelitian (RIP) dan Rencana Strategis PkM berdasarkan Road Map Penelitian dan mengembangkan payung penelitian dan PkM berbasis IPTEKS serta menentukan arah Penelitian dan PkM.
  - Menyusun ketentuan, aturan, SOP, dan standar yang berkaitan dengan penelitian, pengabdian, HKI, ilmiah, serta etik penelitian berkoordinasi dengan Pihak terkait khususnya LPMI dan Rektor.
  - Mengorganisasikan Koordinator Penelitian, PkM, Publikasi, HKI, dan Komite Etik dalam melaksanakan tugas agar terjalin kerjasama yang baik.
  - Melaksanakan monitoring sasaran mutu dan rencana mutu LPPM.
  - Membina anggota di lingkungan LPPM untuk meningkatkan kemampuan dan disiplin kerja.
  - Menetapkan rumusan informasi hasil Penelitian dan PkM berdasarkan ketentuan yang berlaku untuk diketahui oleh masyarakat.
  - Berkoordinasi dengan Wakil Rektor II dalam menetapkan rumusan naskah kerjasama Penelitian dan PkM dengan instansi terkait di luar sekolah/kampus sebagai pedoman kerja.
  - Menyusun laporan kinerja lembaga Penelitian dan PkM sesuai dengan hasil yang telah dicapai sebagai pertanggungjawaban pelaksanaan tugas.
  - Menghimpun dan mengkaji pedoman pendataan dan sistem informasi di bidang Penelitian dan Pengabdian kepada Masyarakat.
  - Melakukan kegiatan pengembangan dosen berkaitan penelitian, publikasi, HKI, inovasi, dan pengabdian masyarakat.
  - Melakukan monitoring berkaitan pelaksanaan penelitian, pengabdian, publikasi, HKI, Inovasi dan Etik Penelitian.
  - Mencari mitra penelitian dan pengabdian masyarakat dengan pihak ketiga baik dari segi materi, materil, atau sarana prasarana.
  - Mencari hibah penelitian dan pengabdian masyarakat atau mendapatkan dari pihak ketiga.
* **Tanggung Jawab**:
  - Bertanggung jawab terhadap mekanisme dan pengelolaan dana hibah penelitian Yayasan untuk dosen.
  - Memberikan laporan kepada Wakil Rektor II setiap akhir semester tentang penelitian dan Pengabdian Masyarakat dosen.
  - Ketua bersama penanggung jawab penelitian, bertanggung jawab untuk mengawal setiap pengajuan proposal penelitian dosen, harus sesuai dengan penciri Program Studi dan Road Map penelitian Fakultas dan Road Map penelitian Program Studi.
  - Bersama dengan Wakil Rektor II UPERTIS melaksanakan menilai reward dan punishment terhadap dosen yang melakukan / yang tidak melakukan penelitian, pengabdian dan publikasi.
  - Melaksanakan tugas lain yang diberikan oleh Pimpinan.
* **Tugas Pokok**:
  - Menyusun RIP Penelitian dan Renstra Pengabmas UPERTIS.
  - Menyusun Rencana kerja tahunan kegiatan penelitian dan pengabdian kepada masyarakat.
  - Meningkatkan Luaran penelitian dan pengabdian kepada masyarakat berupa Jurnal Nasional maupun Internasional.
  - Meningkatkan produk inovasi dan HKI.
  - Menyelenggarakan penelitian yang berlandaskan etika penelitian.
  - Menyusun Road Map Penelitian tingkat Universitas UPERTIS.
  - Membantu Penyusunan Road Map Penelitian tingkat Fakultas dan tingkat Prodi.

---

### 2. Bagian Penelitian
* **Nama**: Sri Indrayati, M.Si
* **Wewenang**:
  - Mempersiapkan kegiatan untuk evaluasi proposal dan monitoring pelaksanaan kegiatan penelitian, seminar hasil penelitian/Deseminasi hasil penelitian.
  - Melakukan monitoring baik langsung maupun pada laman simlibtamas berkaitan dosen yang mendapatkan hibah kementerian atau hibah eksternal lainnya.
  - Melakukan kegiatan diseminasi hasil penelitian baik internal dan eksternal.
  - Melakukan pendokumentasian semua berkaitan penelitian dosen (proposal, laporan kemajuan, laporan akhir, laporan keuangan, luaran penelitian).
  - Membuat kontrak penelitian, surat tugas dan distribusi serta dokumentasi.
  - Melaksanakan penyimpanan dokumen dan surat-surat di bidang penelitian.
  - Memberikan layanan dan informasi penelitian.
  - Melakukan verifikasi data dosen berkaitan Sinta dan simlibtamas.
  - Melakukan koordinasi dengan Prodi berkaitan penelitian mahasiswa baik secara mandiri maupun bekerjasama dengan dosen.
  - Membantu Biro kemahasiswaan untuk melakukan sosialisasi, pembinaan, bimbingan pembuatan proposal penelitian dalam rangka mengikutsertakan dalam kegiatan lomba kewirausahaan atau lomba karya tulis ilmiah untuk tingkat mahasiswa sesuai kebutuhan biro kemahasiswaan.
  - Bersama dengan ketua LPPM, Biro Kemahasiswaan mencari hibah penelitian untuk mahasiswa dari pihak ketiga.
  - Bersama ketua LPPM untuk melakukan monitoring capaian Luaran penelitian dosen di masing fakultas / prodi yang harus dicapai.
  - Melakukan koordinasi dengan ketua prodi, departemen keilmuan berkaitan penelitian mahasiswa (skripsi khususnya berkaitan publikasi mahasiswa dan proses plagiarism).
  - Membuat rekapan penelitian yang dilakukan dosen persemester dan laporan kinerja pertahun.
  - Melakukan kegiatan pengembangan kemampuan mahasiswa dan dosen dalam penyusunan proposal penelitian, pengabdian masyarakat, publikasi, kewirausahaan dan karya ilmiah lainnya.
  - Melaksanakan tugas lain yang diberikan oleh Pimpinan.
* **Tanggung Jawab**:
  - Ikut Bertanggung jawab menyusun dan menentukan road map penelitian tingkat Universitas, Fakultas dan Program Studi bersama Ketua LPPM, Fakultas dan Prodi.
  - Berkoordinasi dengan Penanggung Jawab Penelitian di Prodi dan Fakultas dalam melaksanakan tugas agar terjalin kerjasama yang baik.
  - Melaporkan kinerja penelitian pada simlibtamas.
  - Bertanggung jawab berkaitan sinta, kinerja sinta dan citasi dosen.
  - Bertanggung jawab dalam pengelolaan penelitian dosen di tingkat universitas dengan bekerjasama dengan ketua prodi/Pj Penelitian dan Pengabdian di Prodi untuk pencapaian jumlah penelitian 1 dosen per tahun.
* **Tugas Pokok**:
  - Menyusun Rencana kerja tahunan kegiatan penelitian.
  - Melakukan pengelolaan berkaitan kegiatan penelitian.
  - Menyelenggarakan penelitian yang berlandaskan etika penelitian.

---

### 3. Bagian Pengabdian Masyarakat (PkM)
* **Nama**: Ns. Falerisiska Yunere, M.Kep
* **Wewenang**:
  - Menyusun rencana dan program kerja PkM sebagai pedoman pelaksanaan tugas.
  - Berkoordinasi dengan Penanggung Jawab di Program Studi dan Fakultas dalam melaksanakan tugas agar terjalin kerjasama yang baik.
  - Melakukan pelaporan kinerja pengabdian pada simlibtamas.
  - Mempersiapkan kegiatan untuk evaluasi proposal dan monitoring pelaksanaan kegiatan pengabdian, seminar hasil pengabdian/Deseminasi hasil pengabdian.
  - Melakukan monitoring baik langsung maupun pada laman simlibtamas berkaitan dosen yang mendapatkan hibah kementerian atau hibah eksternal lainnya.
  - Melakukan pendokumentasian semua berkaitan pengabdian dosen (proposal, laporan kemajuan, laporan akhir, laporan keuangan, luaran pengabdian).
  - Membuat kontrak pengabdian, surat tugas dan distribusi serta dokumentasi.
  - Melakukan kegiatan diseminasi pengabdian baik internal maupun eksternal.
  - Melaksanakan penyimpanan dokumen dan surat-surat di bidang pengabdian kepada masyarakat.
  - Berkoordinasi dengan panitia pelaksana PMPKL/KKN untuk pelaksanaan kegiatan dan pengarsipan.
  - Memberikan layanan dan informasi pengabdian masyarakat.
  - Melakukan pengelolaan tingkat universitas kegiatan Pengabdian Masyarakat mahasiswa baik secara mandiri maupun bekerjasama dengan dosen dengan berkoordinasi dengan PJ Penelitian dan Pengabdian di Prodi.
  - Membantu Biro Kemahasiswaan untuk melakukan sosialisasi, pembinaan, bimbingan pembuatan proposal Pengabdian Masyarakat dalam rangka keikutsertaan mahasiswa dalam kegiatan lomba proposal PKM atau lomba karya tulis ilmiah untuk tingkat mahasiswa sesuai kebutuhan biro kemahasiswaan.
  - Bersama dengan ketua LPPM mencari hibah Pengabdian Masyarakat untuk mahasiswa dari pihak ketiga.
  - Bersama ketua LPPM Bertanggung jawab untuk pemantauan pencapaian luaran pengabdian masyarakat yang harus dicapai dosen di prodi / fakultas.
  - Membantu Biro Kemahasiswaan dalam meningkatkan jumlah pengabdian masyarakat 1 himpunan mahasiswa (HIMA) atau gabungan 1 kali per tahun pada saat kegiatan kemah bakti mahasiswa sesuai kebutuhan dari Biro Kemahasiswaan.
  - Membuat rekapan pengabdian yang dilakukan dosen persemester dan laporan kinerja pertahun.
  - Melaksanakan tugas lain yang diberikan oleh Pimpinan.
* **Tanggung Jawab**:
  - Ikut bertanggung jawab menyusun dan menentukan road map pengabdian masyarakat tingkat Universitas, dan membantu Fakultas dan Program Studi dalam penyusunan road map pengabdiannya.
  - Bertanggung jawab dalam pelaksanaan PMPKL (KKN) tingkat Universitas 1 kali per tahun.
  - Bertanggung jawab untuk dalam pengelolaan pengabdian masyarakat di prodi (semua dosen terlibat pengabdian/semester/minimal 2 pengabdian/prodi/semester).
* **Tugas Pokok**:
  - Menyusun Rencana kerja tahunan kegiatan pengabdian.
  - Melakukan pengelolaan berkaitan kegiatan pengabdian.

---

### 4. Bagian Publikasi dan Kegiatan Ilmiah
* **Nama**: Feni Wartisa, S.Sit, M.KM / Epi Supriadi, M.Si
* **Wewenang**:
  - Menyusun rencana dan program kerja sebagai pedoman pelaksanaan tugas.
  - Melakukan pengelolaan dan pemantauan terhadap jurnal yang ada.
  - Mempersiapkan pengaturan penyebarluasan hasil Penelitian dan PkM melalui kegiatan pada Jurnal, Seminar Penelitian dan Pengabdian Masyarakat.
  - Melakukan pengelolaan data artikel/publikasi jurnal penelitian dan PkM.
  - Menyimpan, memelihara dokumen dan surat yang berhubungan dengan data artikel jurnal hasil penelitian dan PkM baik lokal, nasional atau internasional.
  - Melakukan upaya peningkatan jurnal baik secara kuantitas dan kualitas.
  - Melakukan kegiatan seminar ilmiah lingkup UPERTIS.
* **Tanggung Jawab**:
  - Bertanggung jawab dalam Meningkatkan akreditasi dan penambahan jurnal penelitian dan pengabdian pada UPERTIS.
  - Bertanggung jawab dalam Melakukan upaya peningkatan publikasi dosen, citasi dosen dan berusaha untuk meningkatkan jumlah artikel penelitian dosen yang dipublikasikan di Jurnal Terindeks (indeks scopus) dan jurnal-jurnal Internasional lainnya.
  - Bertanggung jawab dalam Membuat rekapan publikasi yang dilakukan dosen persemester dan laporan kinerja pertahun.
  - Melaksanakan tugas lain yang diberikan oleh Pimpinan.
* **Tugas Pokok**:
  - Menyusun Rencana kerja tahunan kegiatan publikasi ilmiah.
  - Melakukan pengelolaan berkaitan kegiatan publikasi ilmiah dan jurnal.

---

### 5. Bagian HKI, Inovasi, dan KATSINOV
* **Nama**: Ria Afrianti, M.Si, Apt
* **Wewenang**:
  - Menyusun rencana dan program kerja sebagai pedoman pelaksanaan tugas.
  - Menerima naskah dan bahan berkaitan HKI.
  - Melakukan administrasi berkaitan pengurusan HKI dosen dan mahasiswa.
  - Memproses naskah/bahan/ untuk penerbitan HKI.
  - Melaporkan kinerja berkaitan HKI, Inovasi pada Kinerja Inovasi.
  - Melakukan upaya peningkatan HKI, inovasi dan produk lainnya dosen dan mahasiswa.
  - Membantu pengembangan karya inovasi dosen, mahasiswa.
* **Tanggung Jawab**:
  - Melakukan pengisian kinerja KATSINOV.
  - Bertanggung jawab untuk meningkatkan jumlah HKI 2 pertahun/prodi.
  - Membuat rekapan karya hasil penelitian, pengabdian dan lainnya dan Membuat laporan kinerja.
  - Melaksanakan tugas lain yang diberikan oleh Pimpinan.
* **Tugas Pokok**:
  - Menyusun Rencana kerja tahunan kegiatan HKI, Inovasi.
  - Melakukan pengelolaan berkaitan kegiatan dalam menghasilkan karya HKI, Inovasi.

---

### 6. Komisi Etik Penelitian
* **Nama**: Defrimal, S.Kp, M.Biomed, PAK
* **Wewenang**:
  - Melakukan kajian etik protokol penelitian kesehatan yang mengikutsertakan manusia dan/atau menggunakan hewan percobaan sebagai subyek penelitian beserta tim.
  - Memberikan persetujuan etik (ethical clearance) terhadap protokol penelitian.
  - Melakukan monitoring dan evaluasi terhadap pelaksanaan penelitian yang telah memperoleh persetujuan etik beserta tim.
  - Melakukan sosialisasi pedoman etik sesuai standar dan pedoman WHO.
  - Mengusulkan pemberhentian pelaksanaan penelitian kesehatan terhadap penelitian yang menyimpang/tidak sesuai protokol yang telah diberikan persetujuan etik.
  - Mengajukan kajian ulang protokol penelitian kesehatan dari institusi/ lembaga penelitian lainnya yang bersengketa dengan peneliti.
* **Tanggung Jawab**:
  - Membuat rekapan izin etik yang dilakukan persemester dan laporan kegiatan Komisi Etik.
  - Melakukan sosialisasi kode etik kepada dosen dan mahasiswa di UPERTIS.
  - Melaksanakan tugas lain yang diberikan oleh Pimpinan.
* **Tugas Pokok**:
  - Menyusun Rencana kerja tahunan kegiatan berkaitan Etik Penelitian.
  - Melakukan pengelolaan berkaitan kegiatan etik penelitian di UPERTIS.`


	_, err = db.ExecContext(context.Background(), `
		INSERT INTO pages (slug, title, subtitle, content, published, sort_order, seo_title, seo_description, created_at, updated_at)
		VALUES ('profil-lppm', 'Profil LPPM', 'Lembaga Penelitian dan Pengabdian Masyarakat UPERTIS', $1, true, 10, 'Profil LPPM — Universitas Perintis Indonesia', 'Lembaga Penelitian dan Pengabdian Masyarakat Universitas Perintis Indonesia', NOW(), NOW())
		ON CONFLICT (slug)
		DO UPDATE SET title='Profil LPPM', subtitle='Lembaga Penelitian dan Pengabdian Masyarakat UPERTIS', content=$1, published=true, updated_at=NOW()
	`, content)
	if err != nil {
		log.Fatal("Failed to insert/update profil-lppm page: ", err)
	}

	fmt.Println("Successfully inserted/updated Profil LPPM page in the database!")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
