CREATE TABLE course_requests (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES students(id),
    course_id BIGINT REFERENCES courses(id),
    status TEXT CHECK (status IN ('pending_mentor', 'pending_payment', 'in_progress', 'completed', 'cancelled')) NOT NULL,
    price NUMERIC NOT NULL,
    accepted_at TIMESTAMPTZ,
    payment_due TIMESTAMPTZ,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);