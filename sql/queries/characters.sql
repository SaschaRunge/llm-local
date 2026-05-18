-- name: AddCharacter :one
INSERT INTO characters(id, created_at, updated_at, name, system_prompt, is_user)
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

-- name: GetPersonasLikeName :many
SELECT id, name FROM characters 
WHERE 
    name LIKE $1 AND
    is_user = 1;