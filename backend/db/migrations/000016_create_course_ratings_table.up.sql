CREATE TABLE course_ratings (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT REFERENCES courses(id),
    student_id UUID REFERENCES students(id),
    rating INTEGER CHECK (rating BETWEEN 1 AND 5),
    feedback TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);