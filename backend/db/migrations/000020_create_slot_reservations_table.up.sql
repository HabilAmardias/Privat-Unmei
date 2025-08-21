CREATE TABLE slot_reservations (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT REFERENCES courses(id) NOT NULL,
    student_id UUID REFERENCES students(id) NOT NULL,
    reserved_slots JSONB NOT NULL, -- Array of {date, start_time, end_time}
    reserved_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    status VARCHAR DEFAULT 'booked' CHECK(status IN ('active', 'expired', 'booked')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);