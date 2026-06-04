-- +goose Up
INSERT INTO characters(id, created_at, updated_at, name, system_prompt, is_user, is_system)
VALUES (
    '00000000-0000-0000-0000-000000000000',
    NOW(),
    NOW(),
    'System',
    'You are the system.',
    false,
    true
);

-- +goose Down
DELETE FROM characters WHERE id = '00000000-0000-0000-0000-000000000000';