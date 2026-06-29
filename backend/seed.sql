-- Seed script for KampusPro
-- Run this script using: psql -U postgres -d kampuspro -f seed.sql

-- Helper to make queries simpler
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Hero Slides
DO $$
BEGIN
    FOR i IN 1..10 LOOP
        INSERT INTO hero_slides (title, eyebrow, description, image_url, cta_primary_label, cta_primary_href, cta_secondary_label, cta_secondary_href, active, sort_order)
        VALUES (
            'Selamat Datang di KampusPro ' || i,
            'Penerimaan Mahasiswa Baru ' || (2024 + i),
            'Wujudkan masa depan gemilang bersama universitas terkemuka dengan fasilitas modern dan kurikulum industri.',
            'https://images.unsplash.com/photo-1541339907198-e08756dedf3f?auto=format&fit=crop&q=80',
            'Daftar Sekarang', '/pmb',
            'Pelajari Lebih Lanjut', '/about',
            true, i
        );
    END LOOP;
END $$;

-- 2. Faculties
DO $$
DECLARE
    f_id UUID;
    fac_names TEXT[] := ARRAY['Fakultas Ilmu Komputer', 'Fakultas Kedokteran', 'Fakultas Ekonomi', 'Fakultas Teknik', 'Fakultas Hukum', 'Fakultas Ilmu Budaya', 'Fakultas Psikologi', 'Fakultas MIPA', 'Fakultas Pertanian', 'Fakultas Komunikasi'];
    fac_codes TEXT[] := ARRAY['FIK', 'FK', 'FE', 'FT', 'FH', 'FIB', 'FPSi', 'FMIPA', 'FP', 'FIKOM'];
    suffix TEXT;
BEGIN
    FOR i IN 1..10 LOOP
        suffix := substr(md5(random()::text), 1, 4);
        INSERT INTO faculties (code, name, slug, description, vision, mission, about_content, hero_title, hero_eyebrow, hero_description, cover_image_url, accent, active, programs, sort_order)
        VALUES (
            fac_codes[i] || '_' || suffix,
            fac_names[i] || ' ' || suffix,
            LOWER(REPLACE(fac_names[i], ' ', '-')) || '-' || suffix,
            'Fakultas terbaik dengan fasilitas lengkap.',
            'Menjadi fakultas unggul bertaraf internasional',
            '1. Menyelenggarakan pendidikan berkualitas\n2. Melakukan penelitian inovatif',
            'Sejarah dan profil lengkap fakultas.',
            'Mari Bergabung di ' || fac_codes[i],
            'Pilihan Cerdas',
            'Fakultas dengan lulusan yang sangat diminati di dunia kerja.',
            'https://images.unsplash.com/photo-1562774053-701939374585?auto=format&fit=crop&q=80',
            'navy', true, 3, i
        ) RETURNING id INTO f_id;

        -- 3. Programs (3 per faculty)
        FOR j IN 1..3 LOOP
            INSERT INTO faculty_programs (faculty_id, name, level, description, accreditation, duration_years, active, sort_order)
            VALUES (
                f_id,
                'Program Studi ' || fac_codes[i] || ' ' || j,
                'S1', 'Program sarjana dengan kurikulum industri',
                'A', 4, true, j
            );
        END LOOP;

        -- 4. Lecturers (3 per faculty)
        FOR j IN 1..3 LOOP
            INSERT INTO faculty_lecturers (faculty_id, name, position, education, expertise, photo_url, active, sort_order)
            VALUES (
                f_id,
                'Dr. Dosen ' || fac_codes[i] || ' ' || j,
                'Lektor Kepala', 'S3 Universitas Terbaik', 'Spesialis ' || fac_names[i],
                'https://ui-avatars.com/api/?name=Dosen+' || j || '&background=random',
                true, j
            );
        END LOOP;

        -- 7. Blog Posts (2 per faculty)
        FOR j IN 1..2 LOOP
            INSERT INTO blog_posts (title, slug, excerpt, content, cover_image_url, author_name, category, tags, faculty_id, status, published_at, seo_title, seo_description)
            VALUES (
                'Artikel Menarik dari ' || fac_codes[i] || ' ' || j || ' ' || suffix,
                LOWER(fac_codes[i]) || '-artikel-' || j || '-' || suffix,
                'Sebuah ulasan singkat tentang inovasi di fakultas.',
                'Isi konten yang panjang dan mendetail mengenai topik yang sedang dibahas di lingkungan kampus...',
                'https://images.unsplash.com/photo-1454165804606-c3d57bc86b40?auto=format&fit=crop&q=80',
                'Admin Kampus', 'Pendidikan', ARRAY['Kampus', 'Inovasi'],
                f_id, 'published', NOW() - (j || ' days')::INTERVAL,
                'Artikel ' || fac_codes[i], 'Meta deskripsi artikel'
            );
        END LOOP;
    END LOOP;
END $$;

-- 5. News
DO $$
BEGIN
    FOR i IN 1..10 LOOP
        INSERT INTO news (title, excerpt, content, image_url, category, date, featured, active)
        VALUES (
            'Berita Terkini Kampus - Edisi ' || i,
            'Kampus kita kembali menorehkan prestasi gemilang di tingkat nasional.',
            'Isi berita lengkap yang memuat rincian tentang pencapaian mahasiswa dan dosen...',
            'https://images.unsplash.com/photo-1523050854058-8df90110c9f1?auto=format&fit=crop&q=80',
            'Prestasi', NOW() - (i || ' days')::INTERVAL, (i % 3 = 0), true
        );
    END LOOP;
END $$;

-- 6. Testimonials
DO $$
BEGIN
    FOR i IN 1..10 LOOP
        INSERT INTO testimonials (name, role, quote, image_url, year, active, sort_order)
        VALUES (
            'Alumni Sukses ' || i,
            'CEO Startup',
            'Kuliah di sini adalah keputusan terbaik yang pernah saya buat. Fasilitas dan dosennya sangat mendukung.',
            'https://ui-avatars.com/api/?name=Alumni+' || i || '&background=random',
            '202' || (i % 5), true, i
        );
    END LOOP;
END $$;

-- 8. Pages
DO $$
DECLARE
    suffix TEXT;
BEGIN
    FOR i IN 1..10 LOOP
        suffix := substr(md5(random()::text), 1, 4);
        INSERT INTO pages (title, subtitle, slug, content, cover_image_url, published, sort_order, seo_title, seo_description)
        VALUES (
            'Halaman Informasi ' || i || ' ' || suffix,
            'Penjelasan detail halaman ' || i,
            'halaman-info-' || i || '-' || suffix,
            '<h2>Selamat Datang di Halaman ' || i || '</h2><p>Ini adalah halaman statis yang dibuat secara dinamis untuk keperluan informasi publik.</p>',
            'https://images.unsplash.com/photo-1497366216548-37526070297c?auto=format&fit=crop&q=80',
            true, i,
            'SEO Title Halaman ' || i, 'Deskripsi lengkap untuk SEO halaman ' || i
        );
    END LOOP;
END $$;
