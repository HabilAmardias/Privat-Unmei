CREATE TABLE courses (
    id BIGSERIAL PRIMARY KEY,
    mentor_id UUID REFERENCES mentors(id),
    title TEXT NOT NULL,
    description TEXT,
    price INTEGER NOT NULL,
    duration_days INTEGER NOT NULL,
    transaction_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);