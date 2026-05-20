-- +goose Up
ALTER TABLE characters
ADD COLUMN is_system BOOLEAN DEFAULT false;

-- +goose Down
ALTER TABLE characters
DROP COLUMN is_system;