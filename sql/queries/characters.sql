-- name: AddCharacter :one
INSERT INTO characters(id, created_at, updated_at, name, card, is_user)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetCharactersLikeName :many
SELECT id FROM characters 
WHERE name LIKE $1;

-- name: GetUserCharactersLikeName :many
SELECT id, name, card FROM characters 
WHERE 
    name LIKE $1 AND
    is_user = 1;