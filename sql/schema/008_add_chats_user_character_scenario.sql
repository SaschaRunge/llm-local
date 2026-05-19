-- +goose Up
ALTER TABLE chats
ADD COLUMN user_character_id UUID,
ADD COLUMN scenario TEXT,
ADD CONSTRAINT fk_user_character_id 
FOREIGN KEY (user_character_id) 
REFERENCES characters(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE chats
DROP COLUMN user_character_id, scenario;