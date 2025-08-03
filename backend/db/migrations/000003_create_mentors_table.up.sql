CREATE TABLE mentors (
    id UUID PRIMARY KEY NOT NULL REFERENCES users(id),
    total_rating NUMERIC DEFAULT 0,
    rating_count INTEGER DEFAULT 0,
    resume_url VARCHAR NOT NULL,
    years_of_experience INT not null,
    whatsapp_number VARCHAR not null unique,
    degree VARCHAR NOT NULL CHECK(degree in ('bachelor','diploma','high school','master','professor')),
    major VARCHAR NOT NULL,
    campus VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);