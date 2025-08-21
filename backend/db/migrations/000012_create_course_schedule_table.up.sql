
-- scheduled session
CREATE TABLE course_schedule (
    id BIGSERIAL PRIMARY KEY,
    course_request_id BIGINT REFERENCES course_requests(id),
    session_number INT NOT NULL,
    scheduled_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_course_schedule_unique_active_session 
ON course_schedule (course_request_id, session_number) 
WHERE deleted_at IS NULL;