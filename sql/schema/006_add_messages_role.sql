-- +goose Up
ALTER TABLE messages
ADD role TEXT;

-- +goose Down
ALTER TABLE
DROP COLUMN role;