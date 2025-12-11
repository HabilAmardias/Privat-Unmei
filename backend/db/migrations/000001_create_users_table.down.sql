DROP INDEX IF EXISTS idx_username_gin;
DROP INDEX IF EXISTS idx_useremail_gin;
DROP INDEX IF EXISTS idx_userpublicid_gin;
DROP TRIGGER IF EXISTS trigger_set_user_public_id ON users;
DROP FUNCTION IF EXISTS set_user_public_id();
DROP FUNCTION IF EXISTS generate_short_id(INT);
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS pg_trgm;
DROP EXTENSION IF EXISTS pgcrypto;