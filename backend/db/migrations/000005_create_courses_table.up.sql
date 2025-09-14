
-- course detail
CREATE TABLE courses (
    id BIGSERIAL PRIMARY KEY,
    mentor_id UUID REFERENCES mentors(id),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    domicile VARCHAR NOT NULL,
    method VARCHAR NOT NULL CHECK(method in ('offline', 'online', 'hybrid')),
    price NUMERIC NOT NULL,
    session_duration_minutes INT NOT NULL,
    max_total_session INT NOT NULL,
    transaction_count INT DEFAULT 0,
    is_active BOOL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_courses_mentor ON courses (mentor_id);
CREATE INDEX idx_courses_gin ON courses USING GIN (title gin_trgm_ops);