CREATE TABLE students (
    id UUID PRIMARY KEY NOT NULL REFERENCES users(id),
    verify_token VARCHAR UNIQUE,
    reset_token VARCHAR UNIQUE,
    login_token VARCHAR UNIQUE,
    otp INTEGER,
    otp_last_updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);