CREATE TABLE course_category_assignments (
    course_id BIGINT REFERENCES courses(id),
    category_id BIGINT REFERENCES course_categories(id),
    PRIMARY KEY (course_id, category_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);