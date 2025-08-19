CREATE TABLE mentor_availability (
    id BIGSERIAL PRIMARY KEY,
    mentor_id UUID REFERENCES mentors(id),
    day_of_week INT CHECK (day_of_week IN (1, 2, 3, 4, 5, 6, 7)) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    UNIQUE(mentor_id, day_of_week, start_time, end_time)
);