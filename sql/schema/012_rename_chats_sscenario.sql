-- +goose Up
ALTER TABLE chats
RENAME COLUMN scenario TO card;

-- +goose Down
ALTER TABLE chats
RENAME COLUMN card TO scenario;