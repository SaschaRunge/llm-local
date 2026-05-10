-- +goose Up
CREATE TABLE chat_subscriptions (
    id UUID PRIMARY KEY,
    chat_id UUID NOT NULL,
    character_id UUID NOT NULL,
    UNIQUE(chat_id, character_id),
    CONSTRAINT fk_chat_id FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    CONSTRAINT fk_character_id FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chat_subscriptions;