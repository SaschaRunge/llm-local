-- +goose Up
ALTER TABLE messages
ADD idx BIGSERIAL;

-- +goose Down
ALTER TABLE
DROP COLUMN idx;