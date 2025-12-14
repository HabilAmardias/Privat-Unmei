-- request by user (student)
CREATE TABLE course_requests (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    student_id UUID REFERENCES students(id),
    course_id BIGINT REFERENCES courses(id),
    status VARCHAR DEFAULT 'reserved' CHECK(status IN ('reserved','pending payment', 'scheduled', 'completed', 'cancelled')),
    number_of_participant INT DEFAULT 1,
    number_of_sessions INT NOT NULL DEFAULT 7,
    expired_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_course_request_student ON course_requests (student_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_course_request_course ON course_requests (course_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_course_request_status ON course_requests (status) WHERE deleted_at IS NULL;