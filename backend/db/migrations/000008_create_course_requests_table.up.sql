-- request by user (student)
CREATE TABLE course_requests (
    id BIGSERIAL PRIMARY KEY,
    student_id UUID REFERENCES students(id),
    course_id BIGINT REFERENCES courses(id),
    status INT CHECK (status IN (1, 2, 3)) NOT NULL,
    total_price NUMERIC NOT NULL,
    number_of_sessions INT NOT NULL DEFAULT 7,
    accepted_at TIMESTAMPTZ,
    payment_due TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);