CREATE TABLE course_requests (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES students(id),
    course_id BIGINT REFERENCES courses(id),
    status INT CHECK (status IN (1, 2, 3, 4, 5)) NOT NULL,
    price NUMERIC NOT NULL,
    accepted_at TIMESTAMPTZ,
    payment_due TIMESTAMPTZ,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);