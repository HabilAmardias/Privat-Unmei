create table chatrooms (
	id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
	student_id UUID references students(id),
	mentor_id UUID references mentors(id),
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
	UNIQUE(student_id, mentor_id)
);