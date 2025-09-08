CREATE TABLE mentor_payments (
    payment_method_id BIGINT REFERENCES payment_methods(id),
    mentor_id UUID REFERENCES mentors(id),
    account_number VARCHAR NOT NULL,
    PRIMARY KEY (payment_method_id, mentor_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);