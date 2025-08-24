-- request by user (student)
CREATE TABLE course_requests (
    id BIGSERIAL PRIMARY KEY,
    student_id UUID REFERENCES students(id),
    course_id BIGINT REFERENCES courses(id),
    status VARCHAR DEFAULT 'reserved' CHECK(status IN ('reserved','pending payment', 'scheduled', 'completed', 'cancelled')),
    total_price NUMERIC NOT NULL,
    number_of_sessions INT NOT NULL DEFAULT 7,
    accepted_at TIMESTAMPTZ,
    payment_due TIMESTAMPTZ,
    expired_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);