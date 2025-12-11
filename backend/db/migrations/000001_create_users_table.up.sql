CREATE EXTENSION pg_trgm;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    public_id VARCHAR(12) NOT NULL,
    password_hash TEXT NOT NULL,
    bio TEXT NOT NULL DEFAULT 'Add bio',
    profile_image TEXT NOT NULL,
    status VARCHAR NOT NULL CHECK(status IN ('unverified', 'verified')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE OR REPLACE FUNCTION generate_short_id(length INT DEFAULT 8)
RETURNS TEXT AS $$
DECLARE
    chars TEXT := '23456789ABCDEFGHJKMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz';
    result TEXT := '';
    bytes BYTEA;
    i INT;
    chars_length INT;
BEGIN
    chars_length := length(chars);
    bytes := gen_random_bytes(length);
    
    FOR i IN 0..length-1 LOOP
        result := result || substr(chars, (get_byte(bytes, i) % chars_length) + 1, 1);
    END LOOP;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION set_user_public_id()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.public_id IS NULL OR NEW.public_id = '' THEN
        NEW.public_id := 'USR-' || generate_short_id(8);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_user_public_id
    BEFORE INSERT ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_user_public_id();

CREATE INDEX idx_userpublicid_gin ON users USING GIN (public_id gin_trgm_ops);
CREATE INDEX idx_username_gin ON users USING GIN (name gin_trgm_ops);
CREATE INDEX idx_useremail_gin ON users USING GIN (email gin_trgm_ops);