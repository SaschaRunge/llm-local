-- +goose Up
CREATE TABLE chats (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    name TEXT NOT NULL
);

-- +goose Down
DROP TABLE chats;