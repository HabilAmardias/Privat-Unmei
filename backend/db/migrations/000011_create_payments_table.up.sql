CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    course_request_id BIGINT REFERENCES course_requests(id),
    subtotal NUMERIC NOT NULL,
    payment_method_id BIGINT REFERENCES payment_methods(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);