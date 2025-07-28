CREATE TABLE cancel_requests (
    id BIGSERIAL PRIMARY KEY,
    course_request_id BIGINT REFERENCES course_requests(id),
    reason TEXT,
    attachment TEXT,
    status TEXT CHECK (status IN ('pending', 'approved', 'rejected')) DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);