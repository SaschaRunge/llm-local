-- +goose Up
CREATE TABLE characters (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    name TEXT NOT NULL,
    system_prompt TEXT,
    is_user BOOLEAN DEFAULT FALSE
);

-- +goose Down
DROP TABLE characters;