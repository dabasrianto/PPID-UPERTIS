package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func initDB() {
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

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}

	// Connection pool configuration
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Database connected successfully")

	// Enable slow query logging (PostgreSQL side — log queries > 500ms)
	_, _ = db.ExecContext(context.Background(), `SET log_min_duration_statement = 500`)
	// Ensure required tables exist
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS site_settings (
			key TEXT PRIMARY KEY,
			value JSONB NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatal("Failed to create site_settings table: ", err)
	}

	// Add status column to users if not exists
	_, err = db.ExecContext(context.Background(), `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='status') THEN
				ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active';
			END IF;
		END $$;
	`)
	if err != nil {
		log.Println("Warning: Failed to add status column to users table:", err)
	}

	// Add phone_whatsapp column to users if not exists
	_, err = db.ExecContext(context.Background(), `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='phone_whatsapp') THEN
				ALTER TABLE users ADD COLUMN phone_whatsapp VARCHAR(50) DEFAULT '';
			END IF;
		END $$;
	`)
	if err != nil {
		log.Println("Warning: Failed to add phone_whatsapp column to users table:", err)
	}

	// Auto-migrate column map_coordinates for campus_events
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE campus_events ADD COLUMN IF NOT EXISTS map_coordinates VARCHAR(100)`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_lecturers ADD COLUMN IF NOT EXISTS gelar VARCHAR(100) DEFAULT ''`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_lecturers ADD COLUMN IF NOT EXISTS leadership_group VARCHAR(50)`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_lecturers ADD COLUMN IF NOT EXISTS source VARCHAR(50)`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_lecturers ADD COLUMN IF NOT EXISTS pddikti_id TEXT`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_lecturers ADD COLUMN IF NOT EXISTS pddikti_uuid TEXT`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_lecturers ADD COLUMN IF NOT EXISTS portfolio JSONB DEFAULT '{}'::jsonb`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_lecturers ADD COLUMN IF NOT EXISTS slug VARCHAR(320)`)
	_, _ = db.ExecContext(context.Background(), `CREATE UNIQUE INDEX IF NOT EXISTS idx_faculty_lecturers_slug ON faculty_lecturers(slug)`)
	
	// Move UUID to pddikti_uuid column
	_, _ = db.ExecContext(context.Background(), `
		UPDATE faculty_lecturers 
		SET pddikti_uuid = pddikti_id
		WHERE LENGTH(pddikti_id) > 20 AND pddikti_uuid IS NULL
	`)
	_, _ = db.ExecContext(context.Background(), `
		UPDATE faculty_lecturers 
		SET pddikti_id = NULL 
		WHERE LENGTH(pddikti_id) > 20
	`)

	// Restore pddikti_id (NIDN) from portfolio profile
	_, _ = db.ExecContext(context.Background(), `
		UPDATE faculty_lecturers 
		SET pddikti_id = portfolio -> 'profile' ->> 'nidn'
		WHERE (pddikti_id IS NULL OR LENGTH(pddikti_id) <> 10)
		  AND portfolio -> 'profile' ->> 'nidn' IS NOT NULL
		  AND portfolio -> 'profile' ->> 'nidn' <> ''
		  AND LENGTH(portfolio -> 'profile' ->> 'nidn') = 10
	`)

	// Run Go helper migration from JSON
	migrateNIDNsFromJSON(db)
	syncGelarFromJSON(db)


	// Final cleanup to ensure NO invalid values are in pddikti_id before indexing
	_, _ = db.ExecContext(context.Background(), `UPDATE faculty_lecturers SET pddikti_id = NULL WHERE pddikti_id = '' OR pddikti_id = '-' OR (pddikti_id IS NOT NULL AND LENGTH(pddikti_id) <> 10)`)

	// Clean up duplicates for unique index
	_, _ = db.ExecContext(context.Background(), `
		DELETE FROM faculty_lecturers a
		USING faculty_lecturers b
		WHERE a.pddikti_id IS NOT NULL
		  AND a.pddikti_id <> ''
		  AND a.pddikti_id = b.pddikti_id
		  AND a.ctid > b.ctid
	`)
	_, _ = db.ExecContext(context.Background(), `CREATE UNIQUE INDEX IF NOT EXISTS idx_faculty_lecturers_pddikti_id ON faculty_lecturers(pddikti_id) WHERE pddikti_id IS NOT NULL AND pddikti_id <> ''`)
	_, _ = db.ExecContext(context.Background(), `CREATE UNIQUE INDEX IF NOT EXISTS idx_faculty_lecturers_pddikti_uuid ON faculty_lecturers(pddikti_uuid)`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_lecturers ADD COLUMN IF NOT EXISTS scholar_data JSONB DEFAULT '{}'::jsonb`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_lecturers ADD COLUMN IF NOT EXISTS scholar_fetched_at TIMESTAMP`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_lecturers ADD COLUMN IF NOT EXISTS scholar_id TEXT`)
	_, _ = db.ExecContext(context.Background(), `
		UPDATE faculty_lecturers
		SET slug = LOWER(REGEXP_REPLACE(COALESCE(name,''), '[^a-zA-Z0-9]+', '-', 'g')) || '-' || LEFT(REPLACE(id::text,'-',''),8)
		WHERE COALESCE(slug,'') = ''
	`)
	_, _ = db.ExecContext(context.Background(), `
		UPDATE faculty_lecturers
		SET source = 'pddikti',
		    leadership_group = ''
		WHERE COALESCE(source,'') = '' AND COALESCE(leadership_group,'') = 'pddikti'
	`)

	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_programs ADD COLUMN IF NOT EXISTS slug VARCHAR(255)`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculties ADD COLUMN IF NOT EXISTS tujuan TEXT DEFAULT ''`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculties ADD COLUMN IF NOT EXISTS struktur_organisasi_url VARCHAR(500) DEFAULT ''`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculties ADD COLUMN IF NOT EXISTS kerjasama JSONB DEFAULT '[]'::jsonb`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculties ADD COLUMN IF NOT EXISTS about_image_url VARCHAR(500) DEFAULT ''`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculties ADD COLUMN IF NOT EXISTS about_images JSONB DEFAULT '[]'::jsonb`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculties ADD COLUMN IF NOT EXISTS facility_images JSONB DEFAULT '[]'::jsonb`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculties ADD COLUMN IF NOT EXISTS cover_images JSONB DEFAULT '[]'::jsonb`)
	_, _ = db.ExecContext(context.Background(), `
		UPDATE faculty_programs
		SET slug = TRIM(BOTH '-' FROM LOWER(REGEXP_REPLACE(
			CASE 
				WHEN LOWER(name) LIKE LOWER(COALESCE(level,'')) || '%' THEN name
				ELSE COALESCE(level,'') || '-' || name
			END, '[^a-zA-Z0-9]+', '-', 'g')))
	`)
	_, _ = db.ExecContext(context.Background(), `CREATE UNIQUE INDEX IF NOT EXISTS idx_faculty_programs_slug ON faculty_programs(slug)`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_programs ADD COLUMN IF NOT EXISTS syllabus JSONB DEFAULT '[]'::jsonb`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_programs ADD COLUMN IF NOT EXISTS career_paths JSONB DEFAULT '[]'::jsonb`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_programs ADD COLUMN IF NOT EXISTS visi TEXT DEFAULT ''`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_programs ADD COLUMN IF NOT EXISTS misi TEXT DEFAULT ''`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_programs ADD COLUMN IF NOT EXISTS tujuan TEXT DEFAULT ''`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_programs ADD COLUMN IF NOT EXISTS gelar_akademik TEXT DEFAULT ''`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_programs ALTER COLUMN gelar_akademik TYPE TEXT`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_programs ADD COLUMN IF NOT EXISTS kompetensi_lulusan TEXT DEFAULT ''`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_programs ADD COLUMN IF NOT EXISTS fasilitas_laboratorium TEXT DEFAULT ''`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_programs ADD COLUMN IF NOT EXISTS struktur_organisasi JSONB DEFAULT '[]'::jsonb`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_programs ADD COLUMN IF NOT EXISTS fasilitas_laboratorium_image VARCHAR(500) DEFAULT ''`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_programs ADD COLUMN IF NOT EXISTS fasilitas_laboratorium_images JSONB DEFAULT '[]'::jsonb`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_programs ADD COLUMN IF NOT EXISTS cover_image_url VARCHAR(500) DEFAULT ''`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE faculty_programs ADD COLUMN IF NOT EXISTS card_bg_color VARCHAR(50) DEFAULT ''`)



	// Migrations for direct WhatsApp replies
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE contact_messages ADD COLUMN IF NOT EXISTS reply_message TEXT`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE contact_messages ADD COLUMN IF NOT EXISTS replied_at TIMESTAMP`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE contact_messages ADD COLUMN IF NOT EXISTS is_read_user BOOLEAN DEFAULT true`)

	_, _ = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS contact_message_replies (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			message_id UUID REFERENCES contact_messages(id) ON DELETE CASCADE,
			sender_type VARCHAR(50) NOT NULL,
			sender_name VARCHAR(255) NOT NULL,
			message TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)

	// Create downloads table
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS downloads (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title VARCHAR(255) NOT NULL,
			description TEXT,
			file_url VARCHAR(500) NOT NULL,
			category VARCHAR(100) NOT NULL DEFAULT 'umum',
			active BOOLEAN DEFAULT true,
			downloads_count INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Println("Warning: Failed to create downloads table:", err)
	}
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE downloads ALTER COLUMN file_url TYPE TEXT`)

	// Create gallery table
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS gallery (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title VARCHAR(255) NOT NULL,
			description TEXT,
			media_url VARCHAR(500) NOT NULL,
			media_type VARCHAR(50) NOT NULL, -- 'image' or 'video'
			active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Println("Warning: Failed to create gallery table:", err)
	}

	// Create tool_usages table
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS tool_usages (
			slug VARCHAR(100) PRIMARY KEY,
			count INTEGER DEFAULT 0,
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Println("Warning: Failed to create tool_usages table:", err)
	}

	// Create visitor_country_counts table
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS visitor_country_counts (
			country_code VARCHAR(10) PRIMARY KEY,
			count INTEGER DEFAULT 0,
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Println("Warning: Failed to create visitor_country_counts table:", err)
	}


	// Create popup_banners table
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS popup_banners (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title VARCHAR(255) NOT NULL,
			description TEXT,
			image_url VARCHAR(500),
			link_text VARCHAR(100),
			link_url VARCHAR(255),
			start_date TIMESTAMP,
			end_date TIMESTAMP,
			active BOOLEAN DEFAULT false,
			sort_order INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Println("Warning: Failed to create popup_banners table:", err)
	}
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE popup_banners ADD COLUMN IF NOT EXISTS show_image_only BOOLEAN DEFAULT false`)

	// Create faqs table
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS faqs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			question VARCHAR(500) NOT NULL,
			answer TEXT NOT NULL,
			category VARCHAR(100) DEFAULT 'Umum',
			sort_order INTEGER DEFAULT 0,
			active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Println("Warning: Failed to create faqs table:", err)
	}
	_, _ = db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_faqs_active_sort ON faqs (active, sort_order ASC)`)
	// Seed default FAQ entries if table is empty or has duplicates
	var faqCount int
	_ = db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM faqs`).Scan(&faqCount)
	var duplicateCount int
	_ = db.QueryRowContext(context.Background(), `SELECT COUNT(*) - COUNT(DISTINCT question) FROM faqs`).Scan(&duplicateCount)
	if faqCount < 6 || duplicateCount > 0 {
		_, _ = db.ExecContext(context.Background(), `TRUNCATE TABLE faqs`)
		_, _ = db.ExecContext(context.Background(), `
			INSERT INTO faqs (question, answer, category, sort_order, active) VALUES
			('Apa saja Fakultas dan Program Studi yang ada di UPERTIS?', 'UPERTIS memiliki 3 Fakultas utama: Fakultas Farmasi (D3, S1, Profesi Apoteker), Fakultas Ilmu Kesehatan (D3, S1, Profesi Ners & Bidan), dan Fakultas Ekonomi Bisnis & Ilmu Sosial (S1 Manajemen, S1 Akuntansi). Anda dapat melihat daftar lengkap program studi di halaman Pendidikan.', 'Akademik', 1, true),
			('Di mana lokasi kampus UPERTIS?', 'Universitas Perintis Indonesia (UPERTIS) memiliki kampus di kota Padang (Jl. Adinegoro Simp. Kalumpang Lubuk Buaya, Padang) dan Bukittinggi. Informasi peta dan rincian kontak dapat diakses melalui menu Kontak & Lokasi.', 'Umum', 2, true),
			('Bagaimana prosedur pendaftaran bagi Calon Mahasiswa Baru (PMB)?', 'Pendaftaran mahasiswa baru dapat dilakukan secara online melalui portal PMB UPERTIS di pmb.online.upertis.ac.id atau melalui menu PMB di website ini. Calon pendaftar perlu mengisi formulir, mengunggah ijazah/nilai rapor, pas foto, serta bukti pembayaran pendaftaran.', 'PMB', 3, true),
			('Apakah ada program beasiswa yang ditawarkan di UPERTIS?', 'Ya, UPERTIS menyediakan berbagai beasiswa, termasuk KIP-Kuliah (Kemendikbudristek), beasiswa prestasi akademik/non-akademik, beasiswa yayasan bagi keluarga kurang mampu, serta beasiswa dari mitra industri. Informasi pengajuan dapat ditanyakan langsung ke bagian Kemahasiswaan.', 'Beasiswa', 4, true),
			('Bagaimana status akreditasi Universitas Perintis Indonesia?', 'Institusi Universitas Perintis Indonesia (UPERTIS) beserta program studi di bawahnya telah terakreditasi oleh BAN-PT dan LAM-PTKes (untuk rumpun kesehatan). Mayoritas program studi kesehatan kami memiliki akreditasi Baik Sekali hingga Unggul.', 'Akademik', 5, true),
			('Apakah UPERTIS menyediakan asrama atau fasilitas tempat tinggal mahasiswa?', 'Di sekitar lokasi kampus Padang maupun Bukittinggi terdapat banyak pilihan kost dan kontrakan mahasiswa. Pihak kemahasiswaan menyediakan daftar rekomendasi tempat tinggal yang aman dan terjangkau di sekitar lingkungan kampus.', 'Fasilitas', 6, true)
		`)
	}


	// Create pmb_batches table
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS pmb_batches (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			academic_year VARCHAR(50),
			start_date DATE,
			end_date DATE,
			registration_fee NUMERIC DEFAULT 250000,
			is_active BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Println("Warning: Failed to create pmb_batches table:", err)
	}

	// Create pmb_candidates table
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS pmb_candidates (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			batch_id UUID REFERENCES pmb_batches(id) ON DELETE SET NULL,
			registration_number VARCHAR(100),
			full_name VARCHAR(255) NOT NULL,
			phone_whatsapp VARCHAR(50),
			nisn VARCHAR(50),
			school_origin VARCHAR(255),
			first_choice_program_id UUID,
			second_choice_program_id UUID,
			status VARCHAR(50) DEFAULT 'DRAFT',
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Println("Warning: Failed to create pmb_candidates table:", err)
	}

	// Create pmb_documents table
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS pmb_documents (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			candidate_id UUID REFERENCES pmb_candidates(id) ON DELETE CASCADE,
			document_type VARCHAR(50) NOT NULL,
			file_url VARCHAR(500) NOT NULL,
			status VARCHAR(50) DEFAULT 'PENDING',
			admin_notes TEXT,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Println("Warning: Failed to create pmb_documents table:", err)
	}

	// Add unique constraint on pmb_documents (candidate_id, document_type)
	_, _ = db.ExecContext(context.Background(), `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'pmb_documents_candidate_doc_unique') THEN
				ALTER TABLE pmb_documents ADD CONSTRAINT pmb_documents_candidate_doc_unique UNIQUE (candidate_id, document_type);
			END IF;
		END $$;
	`)

	// Create blog_post_views table
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS blog_post_views (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
			viewed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Println("Warning: Failed to create blog_post_views table:", err)
	}

	// Create pmb_payments table
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS pmb_payments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			candidate_id UUID REFERENCES pmb_candidates(id) ON DELETE CASCADE,
			payment_type VARCHAR(50),
			amount NUMERIC,
			proof_image_url VARCHAR(500),
			status VARCHAR(50) DEFAULT 'PENDING',
			verified_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Println("Warning: Failed to create pmb_payments table:", err)
	}
	db.ExecContext(context.Background(), `ALTER TABLE pmb_payments ADD COLUMN IF NOT EXISTS verified_at TIMESTAMP`)
	// Create otps table
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS otps (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			phone VARCHAR(20) NOT NULL,
			code VARCHAR(10) NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			is_used BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Println("Warning: Failed to create otps table:", err)
	}

	// Create pages table
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS pages (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title VARCHAR(255) NOT NULL,
			subtitle VARCHAR(255),
			slug VARCHAR(255) UNIQUE NOT NULL,
			content TEXT,
			cover_image_url VARCHAR(500),
			published BOOLEAN DEFAULT false,
			sort_order INTEGER DEFAULT 0,
			seo_title VARCHAR(255),
			seo_description TEXT,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Println("Warning: Failed to create pages table:", err)
	}

	// Create permohonan_informasi table
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS permohonan_informasi (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			ticket_number VARCHAR(100) UNIQUE NOT NULL,
			applicant_type VARCHAR(50) NOT NULL,
			name VARCHAR(255) NOT NULL,
			identity_number VARCHAR(100) NOT NULL,
			email VARCHAR(255) NOT NULL,
			phone VARCHAR(100) NOT NULL,
			address TEXT NOT NULL,
			details TEXT NOT NULL,
			purpose TEXT NOT NULL,
			obtain_method VARCHAR(100) NOT NULL,
			delivery_method VARCHAR(100) NOT NULL,
			attachment_url VARCHAR(500),
			status VARCHAR(50) DEFAULT 'pending',
			admin_response TEXT,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Println("Warning: Failed to create permohonan_informasi table:", err)
	}

	// Auto-migrate columns for pages table
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE pages ADD COLUMN IF NOT EXISTS subtitle VARCHAR(255)`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE pages ADD COLUMN IF NOT EXISTS cover_image_url VARCHAR(500)`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE pages ADD COLUMN IF NOT EXISTS sort_order INTEGER DEFAULT 0`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE pages ADD COLUMN IF NOT EXISTS seo_title VARCHAR(255)`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE pages ADD COLUMN IF NOT EXISTS seo_description TEXT`)

	// Seed default PPID Info Publik pages if they do not exist
	berkalaJSON := `{
  "intro": "<p class=\"text-xs lg:text-sm text-slate-500 font-medium leading-relaxed mb-6\">Berikut adalah daftar Informasi Publik Berkala Universitas Perintis Indonesia yang wajib disediakan dan diumumkan secara berkala sekurang-kurangnya 6 bulan sekali sesuai ketentuan UU KIP.</p>",
  "sections": [
    {
      "text": "<h3 class=\"text-xs font-extrabold text-[#002147] uppercase tracking-wider mb-2\">1. Informasi Badan Publik UPERTIS</h3><p class=\"text-[11px] text-slate-400 mb-4\">Dokumen pendirian, kedudukan hukum, rencana jangka panjang, pimpinan, dan tata pamong universitas:</p><div class=\"grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-2\"><a href=\"/halaman/profil\" class=\"flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group\"><div class=\"p-2.5 bg-white text-blue-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap="round" stroke-linejoin="round\" stroke-width=\"2\" d=\"M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4\"></path></svg></div><div><span class=\"text-xs font-bold text-slate-800 block\">Profil & Sejarah</span><span class=\"text-[10px] text-slate-400 block mt-0.5\">Kedudukan & Domisili UPERTIS</span></div></a><a href=\"/halaman/visi-misi\" class=\"flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group\"><div class=\"p-2.5 bg-white text-indigo-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap="round" stroke-linejoin="round\" stroke-width=\"2\" d=\"M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z\"></path></svg></div><div><span class=\"text-xs font-bold text-slate-800 block\">Visi & Misi PPID</span><span class=\"text-[10px] text-slate-400 block mt-0.5\">Arah gerak & komitmen layanan</span></div></a><a href=\"/halaman/tugas-dan-fungsi\" class=\"flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group\"><div class=\"p-2.5 bg-white text-amber-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap="round" stroke-linejoin="round\" stroke-width=\"2\" d=\"M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01\"></path></svg></div><div><span class=\"text-xs font-bold text-slate-800 block\">Tugas & Fungsi</span><span class=\"text-[10px] text-slate-400 block mt-0.5\">Uraian tugas pokok pengelola</span></div></a><a href=\"/halaman/struktur-organisasi-2\" class=\"flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group\"><div class=\"p-2.5 bg-white text-rose-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap="round" stroke-linejoin="round\" stroke-width=\"2\" d=\"M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z\"></path></svg></div><div><span class=\"text-xs font-bold text-slate-800 block\">Struktur Organisasi</span><span class=\"text-[10px] text-slate-400 block mt-0.5\">Bagan tata koordinasi PPID</span></div></a><a href=\"/halaman/pimpinan\" class=\"flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group\"><div class=\"p-2.5 bg-white text-emerald-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap="round" stroke-linejoin="round\" stroke-width=\"2\" d=\"M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z\"></path></svg></div><div><span class=\"text-xs font-bold text-slate-800 block\">Profil Pimpinan</span><span class=\"text-[10px] text-slate-400 block mt-0.5\">Jajaran rektorat & struktural</span></div></a><a href=\"/download\" class=\"flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group\"><div class=\"p-2.5 bg-white text-teal-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap="round" stroke-linejoin="round\" stroke-width=\"2\" d=\"M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253\"></path></svg></div><div><span class=\"text-xs font-bold text-slate-800 block\">Regulasi & Pedoman</span><span class=\"text-[10px] text-slate-400 block mt-0.5\">Statuta, Renstra, RIP, Akademik</span></div></a></div>",
      "imagePosition": "right"
    },
    {
      "text": "<h3 class=\"text-xs font-extrabold text-[#002147] uppercase tracking-wider mb-2\">2. Kebijakan Keterbukaan & Kinerja</h3><p class=\"text-[11px] text-slate-400 mb-4\">Alur permohonan informasi, laporan pertanggungjawaban keuangan, dan data kemitraan:</p><div class=\"grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-2\"><a href=\"/halaman/jadwal-layanan-informasi\" class=\"flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group\"><div class=\"p-2.5 bg-white text-purple-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap="round" stroke-linejoin="round\" stroke-width=\"2\" d=\"M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z\"></path></svg></div><div><span class=\"text-xs font-bold text-slate-800 block\">Tata Cara Layanan</span><span class=\"text-[10px] text-slate-400 block mt-0.5\">Formulir, keberatan, & sengketa</span></div></a><a href=\"/download\" class=\"flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group\"><div class=\"p-2.5 bg-white text-orange-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap="round" stroke-linejoin="round\" stroke-width=\"2\" d=\"M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z\"></path></svg></div><div><span class=\"text-xs font-bold text-slate-800 block\">Keuangan & Kinerja</span><span class=\"text-[10px] text-slate-400 block mt-0.5\">RKAT, LKT, Laporan Layanan</span></div></a><a href=\"/download\" class=\"flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group\"><div class=\"p-2.5 bg-white text-pink-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap="round" stroke-linejoin="round\" stroke-width=\"2\" d=\"M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1\"></path></svg></div><div><span class=\"text-xs font-bold text-slate-800 block\">Kerja Sama & Akreditasi</span><span class=\"text-[10px] text-slate-400 block mt-0.5\">Sertifikat prodi & data kemitraan</span></div></a></div>",
      "imagePosition": "right"
    }
  ]
}`

	setiapSaatJSON := `{
  "intro": "<p class=\"text-xs lg:text-sm text-slate-500 font-medium leading-relaxed mb-6\">Berikut adalah daftar Informasi Publik Tersedia Setiap Saat Universitas Perintis Indonesia yang wajib disediakan dan siap diberikan kepada pemohon informasi saat dibutuhkan sesuai dengan ketentuan UU KIP.</p>",
  "sections": [
    {
      "text": "<h3 class=\"text-xs font-extrabold text-[#002147] uppercase tracking-wider mb-2\">1. Dokumen Resmi Perguruan Tinggi</h3><p class=\"text-[11px] text-slate-400 mb-4\">Unduh dokumen statuta, rencana pembangunan jangka menengah, organisasi tata kelola, dan regulasi akademik:</p><div class=\"grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-2\"><a href=\"/download\" class=\"flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group\"><div class=\"p-2.5 bg-white text-blue-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap="round" stroke-linejoin="round\" stroke-width=\"2\" d=\"M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253\"></path></svg></div><div><span class=\"text-xs font-bold text-slate-800 block\">Statuta UPERTIS</span><span class=\"text-[10px] text-slate-400 block mt-0.5\">Landasan hukum operasional</span></div></a><a href=\"/download\" class=\"flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group\"><div class=\"p-2.5 bg-white text-indigo-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap="round" stroke-linejoin="round\" stroke-width=\"2\" d=\"M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2\"></path></svg></div><div><span class=\"text-xs font-bold text-slate-800 block\">Rencana Strategis</span><span class=\"text-[10px] text-slate-400 block mt-0.5\">Rencana Kerja UPERTIS</span></div></a><a href=\"/download\" class=\"flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group\"><div class=\"p-2.5 bg-white text-emerald-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap="round" stroke-linejoin="round\" stroke-width=\"2\" d=\"M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 4a2 2 0 002-2v-1a2 2 0 00-2-2m-2 4h.01m-2-4H5m14 8H5m14 4H5\"></path></svg></div><div><span class=\"text-xs font-bold text-slate-800 block\">OTK & RIP UPERTIS</span><span class=\"text-[10px] text-slate-400 block mt-0.5\">Rencana pengembangan</span></div></a></div>",
      "imagePosition": "right"
    }
  ]
}`

	sertaMertaJSON := `{
  "intro": "<p class=\"text-xs lg:text-sm text-slate-500 font-medium leading-relaxed mb-6\">Berikut adalah daftar Informasi Serta Merta Universitas Perintis Indonesia yang wajib diumumkan secara serta merta menyangkut hajat hidup orang banyak dan ketertiban umum.</p>",
  "sections": [
    {
      "text": "<h3 class=\"text-xs font-extrabold text-[#002147] uppercase tracking-wider mb-2\">1. Sistem Peringatan Bencana & Keselamatan</h3><p class=\"text-[11px] text-slate-400 mb-4\">Panduan mitigasi, penanggulangan keadaan darurat, dan pedoman keselamatan kerja sivitas akademika:</p><div class=\"grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-2\"><a href=\"/download\" class=\"flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group\"><div class=\"p-2.5 bg-white text-rose-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap="round" stroke-linejoin="round\" stroke-width=\"2\" d=\"M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z\"></path></svg></div><div><span class=\"text-xs font-bold text-slate-800 block\">SOP Darurat Kampus</span><span class=\"text-[10px] text-slate-400 block mt-0.5\">Langkah mitigasi & evakuasi</span></div></a><a href=\"/download\" class=\"flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group\"><div class=\"p-2.5 bg-white text-emerald-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap="round" stroke-linejoin="round\" stroke-width=\"2\" d=\"M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z\"></path></svg></div><div><span class=\"text-xs font-bold text-slate-800 block\">Pedoman K3 UPERTIS</span><span class=\"text-[10px] text-slate-400 block mt-0.5\">Keselamatan sivitas akademika</span></div></a></div>",
      "imagePosition": "right"
    }
  ]
}`

	zonaJSON := `{
  "intro": "<p class=\"text-xs lg:text-sm text-slate-500 font-medium leading-relaxed mb-6\">Pembangunan Zona Integrasi di lingkungan Universitas Perintis Indonesia merupakan wujud komitmen nyata mewujudkan Wilayah Bebas Korupsi (WBK) dan Wilayah Birokrasi Bersih Melayani (WBBM).</p>",
  "sections": [
    {
      "text": "<h3 class=\"text-xs font-extrabold text-[#002147] uppercase tracking-wider mb-2\">1. Pilar Pembangunan Zona Integrasi</h3><p class=\"text-[11px] text-slate-400 mb-4\">Dokumen rencana kerja, lembar kerja evaluasi, dan rekapitulasi data dukung 6 area perubahan:</p><div class=\"grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-2\"><a href=\"/download\" class=\"flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group\"><div class=\"p-2.5 bg-white text-blue-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap="round" stroke-linejoin="round\" stroke-width=\"2\" d=\"M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z\"></path></svg></div><div><span class=\"text-xs font-bold text-slate-800 block\">Dokumen ZI UPERTIS</span><span class=\"text-[10px] text-slate-400 block mt-0.5\">Rencana Kerja & LKE lengkap</span></div></a></div>",
      "imagePosition": "right"
    }
  ]
}`

	_, _ = db.ExecContext(context.Background(), `
		INSERT INTO pages (title, subtitle, slug, content, published, sort_order) 
		VALUES 
		(
			'Informasi Publik Berkala', 
			'Daftar Informasi Publik Universitas Perintis Indonesia yang Wajib Disediakan dan Diumumkan secara Berkala', 
			'informasi-publik-berkala', 
			$1, 
			true, 1
		),
		(
			'Informasi Tersedia Setiap Saat', 
			'Daftar Informasi Publik Universitas Perintis Indonesia yang Wajib Tersedia Setiap Saat', 
			'informasi-tersedia-setiap-saat', 
			$2, 
			true, 2
		),
		(
			'Informasi Serta Merta', 
			'Informasi yang Wajib Diumumkan secara Serta Merta oleh Universitas Perintis Indonesia', 
			'info-serta-merta', 
			$3, 
			true, 3
		),
		(
			'Zona Integrasi (ZI)', 
			'Pembangunan Zona Integrasi Menuju Wilayah Bebas dari Korupsi (WBK) & Wilayah Birokrasi Bersih dan Melayani (WBBM) UPERTIS', 
			'zona-integrasi', 
			$4, 
			true, 4
		),
		(
			'Pengajuan Keberatan Informasi', 
			'Mekanisme pengajuan keberatan atas layanan informasi publik PPID UPERTIS', 
			'keberatan-informasi', 
			'{"intro": "<p>Pemohon Informasi Publik dapat mengajukan keberatan secara tertulis kepada Atasan PPID UPERTIS jika permohonan informasi ditolak, tidak ditanggapi secara tepat waktu, atau biaya yang dikenakan terlalu tinggi. Silakan unduh formulir di bawah ini, isi dengan lengkap, lalu serahkan kepada petugas desk layanan PPID.</p>", "docs": [{"title": "Formulir Pengajuan Keberatan Informasi UPERTIS", "description": "Formulir cetak resmi pengajuan keberatan layanan informasi publik.", "file_url": "/uploads/formulir-keberatan-ppid.pdf"}]}', 
			true, 5
		)
		ON CONFLICT (slug) DO NOTHING
	`, berkalaJSON, setiapSaatJSON, sertaMertaJSON, zonaJSON)

	// Automatically upgrade existing legacy plain-text seed placeholders to premium JSON layout
	_, _ = db.ExecContext(context.Background(), `
		UPDATE pages 
		SET content = $1 
		WHERE slug = 'informasi-publik-berkala' 
		  AND (content LIKE '<h2>Informasi Publik Berkala</h2>%' OR content IS NULL OR content = '')
	`, berkalaJSON)

	_, _ = db.ExecContext(context.Background(), `
		UPDATE pages 
		SET content = $1 
		WHERE slug = 'informasi-tersedia-setiap-saat' 
		  AND (content LIKE '<h2>Informasi Tersedia Setiap Saat</h2>%' OR content IS NULL OR content = '')
	`, setiapSaatJSON)

	_, _ = db.ExecContext(context.Background(), `
		UPDATE pages 
		SET content = $1 
		WHERE slug = 'info-serta-merta' 
		  AND (content LIKE '<h2>Informasi Serta Merta</h2>%' OR content IS NULL OR content = '')
	`, sertaMertaJSON)

	_, _ = db.ExecContext(context.Background(), `
		UPDATE pages 
		SET content = $1 
		WHERE slug = 'zona-integrasi' 
		  AND (content LIKE '<h2>Zona Integrasi (ZI)</h2>%' OR content IS NULL OR content = '')
	`, zonaJSON)

	// Auto-migrate campus_events
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE campus_events ADD COLUMN IF NOT EXISTS category VARCHAR(100)`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE campus_events ADD COLUMN IF NOT EXISTS registration_url VARCHAR(500)`)
	_, _ = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS event_categories (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(100) NOT NULL UNIQUE,
			type VARCHAR(50) DEFAULT 'all',
			active BOOLEAN DEFAULT true
		)
	`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE event_categories ADD COLUMN IF NOT EXISTS type VARCHAR(50) DEFAULT 'all'`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE event_categories ADD COLUMN IF NOT EXISTS active BOOLEAN DEFAULT true`)

	_, _ = db.ExecContext(context.Background(), `ALTER TABLE campus_events ADD COLUMN IF NOT EXISTS is_internal_registration BOOLEAN DEFAULT false`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE campus_events ADD COLUMN IF NOT EXISTS capacity INTEGER DEFAULT 0`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE campus_events ADD COLUMN IF NOT EXISTS speakers JSONB DEFAULT '[]'::jsonb`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE campus_events ADD COLUMN IF NOT EXISTS wa_message_template TEXT`)
	_, _ = db.ExecContext(context.Background(), `ALTER TABLE campus_events ADD COLUMN IF NOT EXISTS images JSONB DEFAULT '[]'::jsonb`)

	// Auto-migrate event_registrations table if missing
	_, _ = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS event_registrations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			event_id UUID REFERENCES campus_events(id) ON DELETE CASCADE,
			full_name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			whatsapp VARCHAR(50) NOT NULL,
			status VARCHAR(50) DEFAULT 'registered',
			registered_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(event_id, email),
			UNIQUE(event_id, whatsapp)
		)
	`)

	// Insert default categories
	_, _ = db.ExecContext(context.Background(), `
		INSERT INTO event_categories (name, type) VALUES 
		('Berita Kampus', 'all'),
		('Pengumuman', 'all'),
		('Event Akademik', 'all'),
		('Prestasi', 'news'),
		('Wisuda', 'news'),
		('Seminar', 'event'),
		('Workshop', 'event'),
		('Riset', 'blog'),
		('Umum', 'all')
		ON CONFLICT (name) DO NOTHING
	`)

	// Ensure critical indexes exist for query performance
	indexStatements := []string{
		// News table: add slug, seo, author columns
		`ALTER TABLE news ADD COLUMN IF NOT EXISTS slug VARCHAR(255) UNIQUE`,
		`ALTER TABLE news ADD COLUMN IF NOT EXISTS seo_title VARCHAR(255)`,
		`ALTER TABLE news ADD COLUMN IF NOT EXISTS seo_description TEXT`,
		`ALTER TABLE news ADD COLUMN IF NOT EXISTS author_name VARCHAR(255) DEFAULT 'Admin'`,
		`UPDATE news SET slug = LOWER(REPLACE(REPLACE(REPLACE(title, ' ', '-'), '.', ''), ',', '')) WHERE slug IS NULL`,
		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_news_slug ON news(slug)`,
		`CREATE INDEX IF NOT EXISTS idx_blog_posts_slug_lower ON blog_posts (LOWER(slug))`,
		`CREATE INDEX IF NOT EXISTS idx_blog_posts_status ON blog_posts (status)`,
		`CREATE INDEX IF NOT EXISTS idx_blog_posts_published_at ON blog_posts (published_at DESC NULLS LAST)`,
		`CREATE INDEX IF NOT EXISTS idx_pages_slug ON pages (slug)`,
		`CREATE INDEX IF NOT EXISTS idx_pages_published ON pages (published)`,
		`CREATE INDEX IF NOT EXISTS idx_faculties_slug ON faculties (slug)`,
		`CREATE INDEX IF NOT EXISTS idx_faculties_active ON faculties (active)`,
		`CREATE INDEX IF NOT EXISTS idx_news_active_date ON news (active, date DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_campus_events_active_date ON campus_events (active, event_date ASC)`,
		`CREATE INDEX IF NOT EXISTS idx_campus_events_slug ON campus_events (slug)`,
		`CREATE INDEX IF NOT EXISTS idx_contact_messages_read ON contact_messages (is_read, created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users (email)`,
		`CREATE INDEX IF NOT EXISTS idx_pmb_candidates_user_id ON pmb_candidates (user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_pmb_candidates_status ON pmb_candidates (status)`,
		`CREATE INDEX IF NOT EXISTS idx_event_registrations_event ON event_registrations (event_id)`,
		`CREATE INDEX IF NOT EXISTS idx_otps_phone_code ON otps (phone, code, expires_at)`,
		// Testimonials table: add video/link columns
		`ALTER TABLE testimonials ADD COLUMN IF NOT EXISTS video_url VARCHAR(500)`,
		`ALTER TABLE testimonials ADD COLUMN IF NOT EXISTS youtube_url VARCHAR(500)`,
		`ALTER TABLE testimonials ADD COLUMN IF NOT EXISTS tiktok_url VARCHAR(500)`,
		`ALTER TABLE testimonials ADD COLUMN IF NOT EXISTS instagram_url VARCHAR(500)`,
	}
	for _, stmt := range indexStatements {
		if _, err := db.ExecContext(context.Background(), stmt); err != nil {
			log.Printf("Warning: index creation failed: %v", err)
		}
	}

	ensureDefaultPages()
	ensureDefaultDownloads()
	ensureDefaultSettings()
	migratePostsTable()
	migrateAuthOTPTables()
	migrateCommentsTable()
	migrateMediaTable()
}

func ensureDefaultPages() {
	defaultPages := []struct {
		slug     string
		title    string
		subtitle string
	}{
		{"sejarah", "Perjalanan Sejarah", "Cerita sekilas mengenai UPERTIS dari masa ke masa"},
		{"pimpinan", "Pimpinan & Manajemen", "Mengenal jajaran pimpinan Universitas, Fakultas, dan Unit"},
		{"renstra", "Rencana Strategis", "Rencana pengembangan Universitas Perintis Indonesia"},
		{"kontak", "Kontak & Lokasi Kampus", "Informasi kontak resmi dan lokasi kampus UPERTIS"},
		{"lambang", "Panduan Lambang", "Filosofi dan panduan penggunaan lambang Universitas"},
		{"medsos", "Pedoman Media Sosial", "Panduan resmi komunikasi di media sosial"},
		{"tata-nilai", "7 Tata Nilai UPERTIS", "Internalisasi nilai-nilai utama sivitas akademika"},
		{"fasilitas", "Fasilitas Kampus", "Eksplorasi sarana dan prasarana modern kami"},
		{"tugas-fungsi", "Tugas & Fungsi", "Uraian tugas pokok dan fungsi unit kerja"},
		{"tugas-dan-fungsi", "Tugas & Fungsi PPID", "Kedudukan hukum, wewenang operasional, dan komitmen pelayanan PPID UPERTIS"},
		{"regulasi", "Regulasi KIP", "Landasan hukum dan dasar regulasi Keterbukaan Informasi Publik"},
		{"profil", "Profil Kampus", "Profil Lengkap Universitas Perintis Indonesia"},
		{"profil-lppm", "Profil LPPM", "Profil Lembaga Penelitian dan Pengabdian Masyarakat UPERTIS"},
		{"permohonan-informasi", "Formulir Permohonan Informasi Publik", "Layanan pengajuan dokumen dan informasi publik secara online"},
		{"visi-misi", "Visi & Misi PPID", "Visi, misi, dan komitmen pelayanan informasi publik"},
		{"maklumat", "Maklumat Pelayanan PPID", "Maklumat komitmen resmi standar pelayanan informasi publik"},
		{"struktur-organisasi-2", "Struktur Organisasi PPID", "Bagan dan susunan keanggotaan PPID UPERTIS"},
	}

	for _, p := range defaultPages {
		_, err := db.ExecContext(context.Background(),
			`INSERT INTO pages (slug, title, subtitle, content, published)
			 VALUES ($1, $2, $3, $4, true)
			 ON CONFLICT (slug) DO NOTHING`,
			p.slug, p.title, p.subtitle, "## "+p.title+"\n\nKonten untuk halaman ini belum diisi oleh admin. Silakan perbarui melalui dashboard admin.",
		)
		if err != nil {
			log.Printf("Warning: Failed to ensure page %s: %v", p.slug, err)
		}
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func ensureSuperAdmin(email, password, fullName string) error {
	hashed, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	_, err = db.ExecContext(context.Background(),
		`INSERT INTO users (email, password_hash, full_name, role)
		VALUES ($1, $2, $3, 'admin')
		ON CONFLICT (email) DO NOTHING`,
		email, hashed, fullName)
	return err
}

func ensureDefaultSettings() {
	// Seed campus_highlights
	highlightsJSON := `{
		"title": "Keunggulan Kampus",
		"cards": [
			{
				"subtitle": "FASILITAS MODERN",
				"title": "Laboratorium Klinik Terpadu & CBT Center Nasional",
				"link_text": "Eksplorasi Fasilitas",
				"link_url": "/halaman/fasilitas",
				"image_url": "/kampus-unpri.jpeg"
			},
			{
				"subtitle": "KARAKTER & ETIKA",
				"title": "Internalisasi 7 Tata Nilai Utama UPERTIS",
				"link_text": "Lihat Selengkapnya",
				"link_url": "/halaman/tata-nilai",
				"image_url": "",
				"items": ["Unggul", "Profesional", "Integritas", "Tangguh"]
			},
			{
				"subtitle": "MOTTO PENDIDIKAN",
				"title": "Your Dream Is Our Mission — Mimpi Anda Adalah Misi Kami",
				"link_text": "Visi & Rencana Strategis",
				"link_url": "/halaman/visi-misi",
				"image_url": ""
			}
		]
	}`

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO site_settings (key, value, updated_at)
		 VALUES ('campus_highlights', $1, NOW())
		 ON CONFLICT (key) DO NOTHING`,
		highlightsJSON,
	)
	if err != nil {
		log.Printf("Warning: Failed to ensure default campus highlights: %v", err)
	}

	// Seed page_faq settings
	faqJSON := `{
		"badge": "Frequently Asked Questions",
		"title_1": "Got Questions?",
		"title_2": "We've Got Answers",
		"subtitle": "Temukan jawaban atas pertanyaan yang paling sering diajukan seputar Universitas Perintis Indonesia.",
		"support_title": "Still have questions?",
		"support_desc": "Belum menemukan jawaban yang Anda cari? Tim dukungan kami siap membantu Anda.",
		"support_button": "Contact Support"
	}`

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO site_settings (key, value, updated_at)
		 VALUES ('page_faq', $1, NOW())
		 ON CONFLICT (key) DO NOTHING`,
		faqJSON,
	)
	if err != nil {
		log.Printf("Warning: Failed to ensure default page_faq settings: %v", err)
	}
}

func migrateNIDNsFromJSON(db *sql.DB) {
	dosenBytes, err := os.ReadFile("dosen_data.json")
	if err != nil {
		dosenBytes, err = os.ReadFile("backend/dosen_data.json")
	}
	if err != nil {
		log.Println("Warning: failed to read dosen_data.json during migration:", err)
		return
	}

	type DosenJSONItem struct {
		Reg  string `json:"reg"`
		Nama string `json:"nama"`
	}
	var items []DosenJSONItem
	if err := json.Unmarshal(dosenBytes, &items); err != nil {
		log.Println("Warning: failed to parse dosen_data.json during migration:", err)
		return
	}

	dosenMap := make(map[string]string)
	for _, item := range items {
		if item.Reg != "" && len(item.Reg) == 10 {
			clean := cleanNameForMigration(item.Nama)
			if clean != "" {
				dosenMap[clean] = item.Reg
			}
		}
	}

	type LecturerItem struct {
		ID   string
		Name string
	}
	var lecturersToMigrate []LecturerItem

	rows, err := db.Query("SELECT id, name FROM faculty_lecturers WHERE pddikti_id IS NULL OR LENGTH(pddikti_id) > 20")
	if err == nil {
		for rows.Next() {
			var id, name string
			if rows.Scan(&id, &name) == nil {
				lecturersToMigrate = append(lecturersToMigrate, LecturerItem{ID: id, Name: name})
			}
		}
		rows.Close()
	}

	for _, l := range lecturersToMigrate {
		clean := cleanNameForMigration(l.Name)
		if reg, found := dosenMap[clean]; found {
			_, _ = db.Exec("UPDATE faculty_lecturers SET pddikti_id = $1 WHERE id = $2", reg, l.ID)
		}
	}
}

func cleanNameForMigration(name string) string {
	name = strings.ToLower(name)
	prefixes := []string{"dr.", "prof.", "apt.", "ns.", "drs.", "dra.", "hj.", "h.", "dr", "prof", "apt", "ns", "drs", "dra"}
	for _, p := range prefixes {
		if strings.HasPrefix(name, p+" ") {
			name = strings.TrimPrefix(name, p+" ")
		}
	}
	if idx := strings.Index(name, ","); idx != -1 {
		name = name[:idx]
	}
	var b strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return strings.TrimSpace(b.String())
}

func syncGelarFromJSON(db *sql.DB) {
	dosenBytes, err := os.ReadFile("dosen_data.json")
	if err != nil {
		dosenBytes, err = os.ReadFile("backend/dosen_data.json")
	}
	if err != nil {
		log.Println("Warning: failed to read dosen_data.json during gelar migration:", err)
		return
	}

	type DosenJSONItem struct {
		Reg   string `json:"reg"`
		Nama  string `json:"nama"`
		Gelar string `json:"gelar"`
	}
	var items []DosenJSONItem
	if err := json.Unmarshal(dosenBytes, &items); err != nil {
		log.Println("Warning: failed to parse dosen_data.json during gelar migration:", err)
		return
	}

	gelarMapByNidn := make(map[string]string)
	gelarMapByName := make(map[string]string)
	for _, item := range items {
		g := strings.TrimSpace(item.Gelar)
		if g != "" {
			if item.Reg != "" && len(item.Reg) == 10 {
				gelarMapByNidn[item.Reg] = g
			}
			clean := cleanNameForMigration(item.Nama)
			if clean != "" {
				gelarMapByName[clean] = g
			}
		}
	}

	type LecturerGelarItem struct {
		ID        string
		Name      string
		PddiktiID string
	}
	var lecturers []LecturerGelarItem

	rows, err := db.Query("SELECT id, name, COALESCE(pddikti_id, '') FROM faculty_lecturers")
	if err == nil {
		for rows.Next() {
			var id, name, pddiktiId string
			if rows.Scan(&id, &name, &pddiktiId) == nil {
				lecturers = append(lecturers, LecturerGelarItem{ID: id, Name: name, PddiktiID: pddiktiId})
			}
		}
		rows.Close()
	}

	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	updatedCount := 0
	for _, l := range lecturers {
		var gelar string
		found := false
		if l.PddiktiID != "" && len(l.PddiktiID) == 10 {
			if g, exists := gelarMapByNidn[l.PddiktiID]; exists {
				gelar = g
				found = true
			}
		}
		if !found {
			clean := cleanNameForMigration(l.Name)
			if g, exists := gelarMapByName[clean]; exists {
				gelar = g
				found = true
			}
		}

		if found && gelar != "" {
			_, err = tx.Exec("UPDATE faculty_lecturers SET gelar = $1 WHERE id = $2", gelar, l.ID)
			if err == nil {
				updatedCount++
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println("Warning: failed to commit gelar updates:", err)
	} else {
		log.Printf("✓ Successfully synchronized %d academic titles (gelar) from JSON\n", updatedCount)
	}
}

func ensureDefaultDownloads() {
	defaultDownloads := []struct {
		title       string
		description string
		fileURL     string
		category    string
	}{
		{"Tata Cara keberatan Informasi Publik", "Panduan tata cara mengajukan keberatan atas layanan informasi publik.", "https://ppid.upertis.ac.id/?wpdmdl=2869", "umum"},
		{"Tata Cara pengajuan penyelesaian sengketa Informasi Publik", "Panduan tata cara pengajuan penyelesaian sengketa informasi publik.", "https://ppid.upertis.ac.id/?wpdmdl=2835", "umum"},
		{"Tata Cara Pengajuan Permohonan Informasi Publik", "Panduan tata cara permohonan informasi publik secara online/offline.", "https://ppid.upertis.ac.id/?wpdmdl=2833", "umum"},
		{"DATA AKREDITASI PRODI (Informasi Setiap Saat)", "Sertifikat dan dokumen akreditasi seluruh program studi di lingkungan UPERTIS.", "https://ppid.upertis.ac.id/?wpdmdl=2831", "ppid-setiap-saat"},
		{"Daftar Informasi Dikecualikan (Informasi Setiap Saat)", "Daftar informasi publik UPERTIS yang dikecualikan berdasarkan uji konsekuensi.", "https://ppid.upertis.ac.id/?wpdmdl=2829", "ppid-dikecualikan"},
		{"Daftar Informasi Publik Upertis (informasi Setiap Saat)", "Daftar informasi publik (DIP) UPERTIS yang wajib disediakan setiap saat.", "https://ppid.upertis.ac.id/?wpdmdl=2827", "ppid-setiap-saat"},
		{"CAPAIAN UPERTIS BERDASARKAN 9 PILAR", "Dokumen laporan capaian Universitas Perintis Indonesia berdasarkan instrumen 9 pilar.", "https://ppid.upertis.ac.id/?wpdmdl=2825", "umum"},
		{"Laporan Layanan Informasi 2025", "Laporan tahunan kinerja pelayanan informasi publik PPID UPERTIS tahun 2025.", "https://ppid.upertis.ac.id/?wpdmdl=2823", "ppid-berkala"},
		{"Laporan Layanan Informasi 2024", "Laporan tahunan kinerja pelayanan informasi publik PPID UPERTIS tahun 2024.", "https://ppid.upertis.ac.id/?wpdmdl=2821", "ppid-berkala"},
		{"OTK UPERTIS", "Dokumen resmi Organisasi dan Tata Kerja (OTK) Universitas Perintis Indonesia.", "https://ppid.upertis.ac.id/?wpdmdl=2819", "ppid-setiap-saat"},
		{"RENSTRA UPERTIS", "Dokumen Rencana Strategis (Renstra) Universitas Perintis Indonesia.", "https://ppid.upertis.ac.id/?wpdmdl=2817", "ppid-setiap-saat"},
		{"RIP UPERTIS edit 9 Des 2022", "Rencana Induk Pengembangan (RIP) Universitas Perintis Indonesia.", "https://ppid.upertis.ac.id/?wpdmdl=2815", "ppid-setiap-saat"},
		{"STATUTA UPERTIS", "Statuta Universitas Perintis Indonesia sebagai landasan hukum operasional.", "https://ppid.upertis.ac.id/?wpdmdl=2813", "ppid-setiap-saat"},
		{"RKAT FIX 2023-2024 UPERTIS", "Rencana Kerja dan Anggaran Tahunan (RKAT) Universitas Perintis Indonesia.", "https://ppid.upertis.ac.id/?wpdmdl=2811", "ppid-berkala"},
		{"REVISI Peraturan Akademik UPERTIS 2024", "Buku pedoman peraturan akademik Universitas Perintis Indonesia revisi terbaru 2024.", "https://ppid.upertis.ac.id/?wpdmdl=2809", "ppid-setiap-saat"},
		{"LAPORAN KINERJA 2022 2023", "Laporan akuntabilitas kinerja Universitas Perintis Indonesia periode 2022/2023.", "https://ppid.upertis.ac.id/?wpdmdl=2807", "ppid-berkala"},
		{"LAPORAN KINERJA 2021 2022", "Laporan akuntabilitas kinerja Universitas Perintis Indonesia periode 2021/2022.", "https://ppid.upertis.ac.id/?wpdmdl=2805", "ppid-berkala"},
	}

	for _, d := range defaultDownloads {
		_, err := db.ExecContext(context.Background(),
			`INSERT INTO downloads (title, description, file_url, category, active)
			 SELECT $1, $2, $3, $4, true
			 WHERE NOT EXISTS (
				 SELECT 1 FROM downloads WHERE file_url = $3
			 )`,
			d.title, d.description, d.fileURL, d.category,
		)
		if err != nil {
			log.Printf("Warning: Failed to seed download %s: %v", d.title, err)
		}
	}
}


