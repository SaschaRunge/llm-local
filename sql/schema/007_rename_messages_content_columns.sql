-- +goose Up
ALTER TABLE messages
RENAME COLUMN content_answer TO content;
ALTER TABLE messages
RENAME COLUMN content_thoughts TO reasoning;

-- +goose Down
ALTER TABLE messages
RENAME COLUMN content TO content_answer;
ALTER TABLE messages
RENAME COLUMN reasoning TO content_thoughts;