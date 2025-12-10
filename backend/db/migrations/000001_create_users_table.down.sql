DROP INDEX IF EXISTS idx_username_gin, idx_useremail_gin, idx_userpublicid_gin;
DROP TRIGGER IF EXISTS trigger_set_user_public_id;
DROP FUNCTION IF EXISTS set_user_public_id;
DROP FUNCTION IF EXISTS generate_short_id;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS pg_trgm;
DROP EXTENSION IF EXISTS pgcrypto;

