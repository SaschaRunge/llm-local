-- +goose Up
ALTER TABLE messages
AdD parent_msg_id UUID,
ADD sequence_idx BIGINT DEFAULT 0 NOT NULL;


-- +goose Down
ALTER TABLE messages
DROP COLUMN parent_msg_id,
DROP COLUMN sequence_idx;