CREATE TABLE course_category_assignments (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL REFERENCES courses(id),
    category_id BIGINT NOT NULL REFERENCES course_categories(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    UNIQUE(course_id, category_id)
);

CREATE INDEX idx_cca_course_category_active 
ON course_category_assignments (course_id, category_id) 
WHERE deleted_at IS NULL;
