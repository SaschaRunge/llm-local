-- +goose Up
ALTER TABLE chats 
    ALTER COLUMN card TYPE jsonb USING jsonb_build_object('description', card);
ALTER TABLE characters 
    ALTER COLUMN card TYPE jsonb USING jsonb_build_object('description', card);

-- +goose Down
ALTER TABLE chats 
    ALTER COLUMN card TYPE text USING card ->> 'description';
ALTER TABLE characters 
    ALTER COLUMN card TYPE text USING card ->> 'description';