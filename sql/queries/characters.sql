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
WHERE 
    name LIKE $1 AND
    deleted_at IS NULL;

-- name: GetUserCharactersLikeName :many
SELECT id, name, card FROM characters 
WHERE 
    name LIKE $1 AND
    is_user = 'true' AND
    deleted_at IS NULL;

-- name: GetAvailableCharacters :many
SELECT id, name, card FROM characters
WHERE
    is_user = 'false' AND
    is_system = 'false' AND
    deleted_at IS NULL;

-- name: UpdateCharacterCard :one
UPDATE characters
SET 
    updated_at = NOW(),
    card = $2
WHERE id = $1
RETURNING id;
