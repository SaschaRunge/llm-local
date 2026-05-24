-- +goose Up
ALTER TABLE chats
ALTER COLUMN user_character_id SET NOT NULL;

-- +goose Down
ALTER TABLE chats
ALTER COLUMN user_character_id DROP NOT NULL;