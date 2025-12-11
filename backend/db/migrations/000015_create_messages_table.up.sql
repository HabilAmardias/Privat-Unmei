CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    sender_id UUID REFERENCES users(id),
    chatroom_id UUID REFERENCES chatrooms(id),
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);