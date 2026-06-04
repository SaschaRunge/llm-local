-- +goose Up
ALTER TABLE messages
ALTER COLUMN idx TYPE BIGINT,
ALTER COLUMN idx SET DEFAULT 0, 
ALTER COLUMN idx SET NOT NULL;

-- +goose Down
ALTER TABLE messages
ALTER COLUMN idx SET DEFAULT nextval('messages_idx_seq'::regclass), 
ALTER COLUMN idx DROP NOT NULL;