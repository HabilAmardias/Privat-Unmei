CREATE TABLE mentor_reports (
    id BIGSERIAL PRIMARY KEY,
    course_schedule_id BIGINT REFERENCES course_schedule(id),
    mentor_id UUID REFERENCES mentors(id),
    topic_id BIGINT REFERENCES topics(id),
    materials_given TEXT,
    date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);