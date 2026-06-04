-- +goose Up
ALTER TABLE messages
    ADD CONSTRAINT fk_parent_msg_id
    FOREIGN KEY (parent_msg_id) 
    REFERENCES messages(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE messages
DROP CONSTRAINT fk_parent_msg_id;