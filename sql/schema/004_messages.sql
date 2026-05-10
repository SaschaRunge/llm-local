-- +goose Up
CREATE TABLE messages (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    content_thoughts TEXT,
    content_answer TEXT NOT NULL,

    chat_id UUID NOT NULL,
    author_id UUID NOT NULL,

    CONSTRAINT fk_chat_id 
    FOREIGN KEY (chat_id) 
    REFERENCES chats(id) ON DELETE CASCADE,

    CONSTRAINT fk_author_id 
    FOREIGN KEY (author_id) 
    REFERENCES characters(id) ON DELETE CASCADE
);

CREATE INDEX idx_messages_composite ON messages(chat_id, created_at);

-- +goose Down
DROP TABLE messages;