CREATE TABLE mentor_availability (
    id BIGSERIAL PRIMARY KEY,
    mentor_id UUID REFERENCES mentors(id),
    day_of_week INT CHECK (day_of_week IN (0, 1, 2, 3, 4, 5, 6)) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    UNIQUE(mentor_id, day_of_week, start_time, end_time)
);

CREATE INDEX idx_mentor_availability ON mentor_availability (mentor_id);