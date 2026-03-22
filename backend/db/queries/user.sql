-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, role)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByVerificationToken :one
SELECT * FROM users
WHERE email_verification_token = $1
  AND email_verification_expires_at > NOW();

-- name: GetUserByResetToken :one
SELECT * FROM users
WHERE password_reset_token = $1
  AND password_reset_expires_at > NOW();

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET name = $1, email = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: UpdateUserRole :one
UPDATE users
SET role = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: UpdateEmailVerified :one
UPDATE users
SET email_verified_at = NOW(),
    email_verification_token = NULL,
    email_verification_expires_at = NULL,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SetVerificationToken :one
UPDATE users
SET email_verification_token = $1,
    email_verification_expires_at = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: SetPasswordResetToken :one
UPDATE users
SET password_reset_token = $1,
    password_reset_expires_at = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: UpdatePassword :one
UPDATE users
SET password_hash = $1,
    password_reset_token = NULL,
    password_reset_expires_at = NULL,
    updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
