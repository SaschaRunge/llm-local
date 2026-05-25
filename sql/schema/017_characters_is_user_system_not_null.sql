-- +goose Up
ALTER TABLE characters
ALTER COLUMN is_user SET NOT NULL,
ALTER COLUMN is_system SET NOT NULL;

-- +goose Down
ALTER TABLE characters
ALTER COLUMN is_user DROP NOT NULL,
ALTER COLUMN is_system DROP NOT NULL;