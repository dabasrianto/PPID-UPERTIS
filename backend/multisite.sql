-- Multi-Site Migration: Add faculty_id to content tables
-- Run on production: /www/server/pgsql/bin/psql -U db_up3rt15z -d db_up3rt15z -f backend/multisite.sql

-- 1. Add faculty_id to posts
ALTER TABLE posts ADD COLUMN IF NOT EXISTS faculty_id UUID REFERENCES faculties(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_posts_faculty_id ON posts(faculty_id);

-- 2. Add faculty_id to campus_events
ALTER TABLE campus_events ADD COLUMN IF NOT EXISTS faculty_id UUID REFERENCES faculties(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_events_faculty_id ON campus_events(faculty_id);

-- 3. Add faculty_id to news (legacy)
ALTER TABLE news ADD COLUMN IF NOT EXISTS faculty_id UUID REFERENCES faculties(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_news_faculty_id ON news(faculty_id);

-- 4. Add site_id to site_settings for per-site branding (Phase 4)
-- For now, keep single-key approach. Will migrate later if needed.

-- Verify
SELECT 'posts.faculty_id' AS col, COUNT(*) FROM posts WHERE faculty_id IS NOT NULL
UNION ALL
SELECT 'events.faculty_id', COUNT(*) FROM campus_events WHERE faculty_id IS NOT NULL
UNION ALL
SELECT 'news.faculty_id', COUNT(*) FROM news WHERE faculty_id IS NOT NULL;
