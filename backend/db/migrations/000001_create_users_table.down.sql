DROP INDEX IF EXISTS idx_username_gin, idx_useremail_gin;

DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS pg_trgm;