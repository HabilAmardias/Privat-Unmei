CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    sender_id UUID REFERENCES users(id),
    course_id BIGINT REFERENCES courses(id),
    chatroom_id BIGINT REFERENCES chatrooms(id),
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);