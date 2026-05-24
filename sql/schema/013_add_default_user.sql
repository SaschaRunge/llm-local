-- +goose Up
INSERT INTO characters (id, created_at, updated_at, deleted_at, name, card, is_user, is_system) 
VALUES (
    '00000000-0000-0000-0000-000000000001',
    NOW(), 
    NOW(), 
    NULL, 
    'DefaultUser', 
    'DummyPrompt', 
    'true', 
    'false'
);


-- +goose Down
DELETE FROM characters 
WHERE id = '00000000-0000-0000-0000-000000000001';