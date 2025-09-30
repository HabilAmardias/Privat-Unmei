CREATE EXTENSION pg_trgm;

CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    bio TEXT NOT NULL DEFAULT 'Add bio',
    profile_image TEXT NOT NULL,
    status VARCHAR NOT NULL CHECK(status IN ('unverified', 'verified')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_username_gin ON users USING GIN (name gin_trgm_ops);
CREATE INDEX idx_useremail_gin ON users USING GIN (email gin_trgm_ops);