
-- scheduled session
CREATE TABLE course_schedule (
    id BIGSERIAL PRIMARY KEY,
    course_request_id UUID REFERENCES course_requests(id),
    scheduled_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    status VARCHAR DEFAULT 'reserved' CHECK(status IN ('reserved', 'scheduled', 'completed', 'cancelled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_course_schedule_unique_active_session 
ON course_schedule (course_request_id, scheduled_date) 
WHERE deleted_at IS NULL;