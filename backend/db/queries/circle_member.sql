-- name: UpsertCircleMember :one
INSERT INTO circle_members (circle_id, user_id, role)
VALUES ($1, $2, $3)
ON CONFLICT (circle_id, user_id) DO UPDATE
    SET role = EXCLUDED.role
RETURNING *;

-- name: GetCircleMemberByCircleAndUser :one
SELECT * FROM circle_members
WHERE circle_id = $1 AND user_id = $2;

-- name: ListCircleMembersByCircle :many
SELECT * FROM circle_members
WHERE circle_id = $1
ORDER BY joined_at ASC;

-- name: ListCircleMembersByUser :many
SELECT * FROM circle_members
WHERE user_id = $1
ORDER BY joined_at ASC;

-- name: ListAdminsByCircle :many
SELECT * FROM circle_members
WHERE circle_id = $1 AND role = 'admin'
ORDER BY joined_at ASC;

-- name: DeleteCircleMember :exec
DELETE FROM circle_members
WHERE circle_id = $1 AND user_id = $2;
