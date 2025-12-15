CREATE TABLE course_categories (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_categories_gin ON course_categories USING GIN (name gin_trgm_ops) WHERE deleted_at IS NULL;
CREATE INDEX idx_course_categories_active ON course_categories (id) WHERE deleted_at IS NULL;