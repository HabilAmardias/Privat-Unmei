CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    course_request_id UUID UNIQUE REFERENCES course_requests(id),
    subtotal NUMERIC NOT NULL,
    operational_cost NUMERIC NOT NULL,
    total_price NUMERIC NOT NULL,
    payment_method_name VARCHAR NOT NULL,
    account_number VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);