CREATE TABLE course_ratings (
    id BIGSERIAL PRIMARY KEY,
    course_request_id BIGINT REFERENCES course_requests(id),
    user_id UUID REFERENCES students(id),
    mentor_id UUID REFERENCES mentors(id),
    rating INTEGER CHECK (rating BETWEEN 1 AND 5),
    feedback TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);