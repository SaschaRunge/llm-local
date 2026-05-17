-- +goose Up
ALTER TABLE messages
ADD role TEXT NOT NULL;

-- +goose Down
ALTER TABLE messages
DROP COLUMN role;