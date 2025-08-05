
-- course detail
CREATE TABLE courses (
    id BIGSERIAL PRIMARY KEY,
    mentor_id UUID REFERENCES mentors(id),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    domicile VARCHAR NOT NULL,
    min_price NUMERIC NOT NULL,
    max_price NUMERIC NOT NULL,
    min_duration_days INT NOT NULL,
    max_duration_days INT NOT NULL,
    transaction_count INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);