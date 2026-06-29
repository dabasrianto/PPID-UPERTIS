-- PostgreSQL Schema for KampusPro
-- Run this SQL in your PostgreSQL database

-- No extensions needed: using built-in gen_random_uuid() (PostgreSQL 13+)

-- Users table (for authentication)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- User roles table
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Faculty admins (which faculty a faculty_admin can manage)
CREATE TABLE IF NOT EXISTS faculty_admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    faculty_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Hero slides
CREATE TABLE IF NOT EXISTS hero_slides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    eyebrow VARCHAR(255),
    description TEXT,
    image_url VARCHAR(500),
    cta_primary_label VARCHAR(100),
    cta_primary_href VARCHAR(255),
    cta_secondary_label VARCHAR(100),
    cta_secondary_href VARCHAR(255),
    video_url TEXT,
    active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Faculties
CREATE TABLE IF NOT EXISTS faculties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(10) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    vision TEXT,
    mission TEXT,
    about_content TEXT,
    hero_title VARCHAR(255),
    hero_eyebrow VARCHAR(255),
    hero_description TEXT,
    cover_image_url VARCHAR(500),
    cover_images JSONB DEFAULT '[]'::jsonb,
    accent VARCHAR(50) DEFAULT 'navy',
    active BOOLEAN DEFAULT true,
    programs INTEGER DEFAULT 0,
    sort_order INTEGER DEFAULT 0,
    facilities JSONB,
    contact_info JSONB,
    tujuan TEXT DEFAULT '',
    struktur_organisasi_url VARCHAR(500) DEFAULT '',
    kerjasama JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Faculty programs (study programs)
CREATE TABLE IF NOT EXISTS faculty_programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    faculty_id UUID REFERENCES faculties(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE,
    level VARCHAR(50),
    description TEXT,
    accreditation VARCHAR(50),
    duration_years INTEGER,
    active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    syllabus JSONB DEFAULT '[]'::jsonb,
    career_paths JSONB DEFAULT '[]'::jsonb,
    visi TEXT DEFAULT '',
    misi TEXT DEFAULT '',
    tujuan TEXT DEFAULT '',
    gelar_akademik TEXT DEFAULT '',
    kompetensi_lulusan TEXT DEFAULT '',
    fasilitas_laboratorium TEXT DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Faculty lecturers
CREATE TABLE IF NOT EXISTS faculty_lecturers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    faculty_id UUID REFERENCES faculties(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    gelar VARCHAR(100) DEFAULT '',
    slug VARCHAR(320),
    position VARCHAR(255),
    education TEXT,
    expertise TEXT,
    photo_url VARCHAR(500),
    leadership_group VARCHAR(50),
    source VARCHAR(50),
    pddikti_id TEXT,
    pddikti_uuid TEXT,
    portfolio JSONB DEFAULT '{}'::jsonb,
    active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- News
CREATE TABLE IF NOT EXISTS news (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    excerpt TEXT,
    content TEXT,
    image_url VARCHAR(500),
    category VARCHAR(100),
    date DATE DEFAULT CURRENT_DATE,
    featured BOOLEAN DEFAULT false,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Testimonials
CREATE TABLE IF NOT EXISTS testimonials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    role VARCHAR(255),
    quote TEXT NOT NULL,
    image_url VARCHAR(500),
    year VARCHAR(50),
    video_url VARCHAR(500),
    youtube_url VARCHAR(500),
    tiktok_url VARCHAR(500),
    instagram_url VARCHAR(500),
    active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Blog posts
CREATE TABLE IF NOT EXISTS blog_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    excerpt TEXT,
    content TEXT,
    cover_image_url VARCHAR(500),
    author_name VARCHAR(255),
    category VARCHAR(100),
    tags TEXT[],
    faculty_id UUID REFERENCES faculties(id) ON DELETE SET NULL,
    status VARCHAR(50) DEFAULT 'draft',
    published_at TIMESTAMP,
    seo_title VARCHAR(255),
    seo_description TEXT,
    views INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Pages
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
);

-- Contact messages
CREATE TABLE IF NOT EXISTS contact_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    subject VARCHAR(255),
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Site settings
CREATE TABLE IF NOT EXISTS site_settings (
    key VARCHAR(255) PRIMARY KEY,
    value JSONB,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Event Categories
CREATE TABLE IF NOT EXISTS event_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE
);

-- Campus Events
CREATE TABLE IF NOT EXISTS campus_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    image_url VARCHAR(500),
    event_date DATE NOT NULL,
    start_time TIME DEFAULT '00:00',
    end_time TIME DEFAULT '23:59',
    location VARCHAR(255),
    map_coordinates VARCHAR(100),
    active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    category VARCHAR(100),
    registration_url VARCHAR(500),
    is_internal_registration BOOLEAN DEFAULT false,
    capacity INTEGER DEFAULT 0,
    speakers JSONB DEFAULT '[]'::jsonb,
    wa_message_template TEXT,
    images JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Event Registrations (Pendaftaran Internal)
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
);

-- ==========================================
-- PMB (Penerimaan Mahasiswa Baru) / PPDB
-- ==========================================

-- PMB Batches (Gelombang Pendaftaran)
CREATE TABLE IF NOT EXISTS pmb_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    academic_year VARCHAR(50) NOT NULL,
    start_date DATE,
    end_date DATE,
    registration_fee DECIMAL(12,2) DEFAULT 0,
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- PMB Candidates (Pendaftar)
CREATE TABLE IF NOT EXISTS pmb_candidates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    batch_id UUID REFERENCES pmb_batches(id) ON DELETE RESTRICT,
    registration_number VARCHAR(100) UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    nisn VARCHAR(50),
    phone_whatsapp VARCHAR(50),
    school_origin VARCHAR(255),
    first_choice_program_id UUID REFERENCES faculty_programs(id) ON DELETE SET NULL,
    second_choice_program_id UUID REFERENCES faculty_programs(id) ON DELETE SET NULL,
    status VARCHAR(50) DEFAULT 'DRAFT', -- DRAFT, WAITING_PAYMENT, PAYMENT_VERIFIED, DOCUMENT_VERIFIED, PASSED, FAILED, ENROLLED
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- PMB Documents (Berkas)
CREATE TABLE IF NOT EXISTS pmb_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidate_id UUID REFERENCES pmb_candidates(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL, -- KTP, KK, IJAZAH, PAS_FOTO, RAPOR, SERTIFIKAT
    file_url VARCHAR(500) NOT NULL,
    status VARCHAR(50) DEFAULT 'PENDING', -- PENDING, APPROVED, REJECTED
    admin_notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- PMB Payments (Pembayaran)
CREATE TABLE IF NOT EXISTS pmb_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidate_id UUID REFERENCES pmb_candidates(id) ON DELETE CASCADE,
    payment_type VARCHAR(50) NOT NULL, -- REGISTRATION, TUITION
    amount DECIMAL(12,2) NOT NULL,
    proof_image_url VARCHAR(500) NOT NULL,
    status VARCHAR(50) DEFAULT 'PENDING', -- PENDING, VERIFIED, REJECTED
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_faculties_slug ON faculties(slug);
CREATE INDEX IF NOT EXISTS idx_faculties_active ON faculties(active);
CREATE INDEX IF NOT EXISTS idx_blog_posts_slug ON blog_posts(slug);
CREATE INDEX IF NOT EXISTS idx_blog_posts_status ON blog_posts(status);
CREATE INDEX IF NOT EXISTS idx_pages_slug ON pages(slug);
CREATE INDEX IF NOT EXISTS idx_pages_published ON pages(published);
CREATE INDEX IF NOT EXISTS idx_news_active ON news(active);
CREATE INDEX IF NOT EXISTS idx_news_date ON news(date DESC);
CREATE INDEX IF NOT EXISTS idx_testimonials_active ON testimonials(active);
CREATE UNIQUE INDEX IF NOT EXISTS idx_faculty_programs_slug ON faculty_programs(slug);

-- Default admin user is created by the Go backend on startup (ensureSuperAdmin)

-- Insert sample faculties
INSERT INTO faculties (code, name, slug, description, vision, mission, accent, active, programs, sort_order)
VALUES 
('FK', 'Fakultas Ilmu Kesehatan', 'ilmu-kesehatan', 'Mencetak tenaga kesehatan profesional berlandaskan ilmu dan nilai.', 'Menjadi fakultas kesehatan yang unggul dan bermanfaat', 'Mencetak tenaga kesehatan yang kompeten dan berakhlak', 'red', true, 3, 1),
('FF', 'Fakultas Farmasi', 'farmasi', 'Mengembangkan inovasi farmasi untuk kesehatan masyarakat.', 'Menjadi pusat pendidikan farmasi yang inovatif', 'Menghasilkan lulusan farmasi yang profesional', 'blue', true, 4, 2),
('FEB', 'Fakultas Ekonomi Bisnis & Ilmu Sosial', 'ekonomi-bisnis', 'Mengembangkan talenta bisnis, manajemen, dan digital yang dinamis.', 'Menjadi fakultas ekonomi bisnis & ilmu sosial unggul', 'Mencetak profesional di bidang ekonomi bisnis & ilmu sosial', 'green', true, 4, 3)
ON CONFLICT (code) DO NOTHING;

-- Insert sample hero slides
INSERT INTO hero_slides (title, eyebrow, description, image_url, cta_primary_label, cta_primary_href, cta_secondary_label, cta_secondary_href, active, sort_order)
VALUES 
('Bersama Membentuk Masa Depan Bermakna', 'Universitas Perintis Indonesia', 'Unggul, Professional, Etika, Rasional, Tangguh, Integritas, dan Solidaritas.', '/kampus-unpri.jpeg', 'Jelajahi UPERTIS', '/halaman/tentang', 'Kenali Program Studi', '#program', true, 1)
ON CONFLICT DO NOTHING;

-- Visitor country counts for self-hosted flag counter
CREATE TABLE IF NOT EXISTS visitor_country_counts (
    country_code VARCHAR(10) PRIMARY KEY,
    count INTEGER DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW()
);

