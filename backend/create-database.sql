-- Run this against the default postgres database before running schema.sql.
-- Example:
-- psql -U postgres -d postgres -f backend/create-database.sql

SELECT 'CREATE DATABASE kampuspro'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'kampuspro')\gexec
