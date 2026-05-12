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