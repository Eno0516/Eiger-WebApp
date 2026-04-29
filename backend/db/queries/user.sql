-- name: CreateUser :one
INSERT INTO users (line_user_id, is_system_admin)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByLineUserID :one
SELECT * FROM users WHERE line_user_id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET is_registered = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
