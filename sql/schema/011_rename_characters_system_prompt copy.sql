-- +goose Up
ALTER TABLE characters
RENAME COLUMN system_prompt TO card;

-- +goose Down
ALTER TABLE characters
RENAME COLUMN card TO system_prompt;