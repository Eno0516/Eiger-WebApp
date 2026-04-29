-- name: CreateCircle :one
INSERT INTO circles (name, description, application_form_url, spreadsheet_id, uuid_column)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCircleByID :one
SELECT * FROM circles WHERE id = $1;

-- name: ListCircles :many
SELECT * FROM circles ORDER BY created_at DESC;

-- name: UpdateCircle :one
UPDATE circles
SET name                 = $1,
    description          = $2,
    application_form_url = $3,
    spreadsheet_id       = $4,
    uuid_column          = $5,
    updated_at           = NOW()
WHERE id = $6
RETURNING *;

-- name: DeleteCircle :exec
DELETE FROM circles WHERE id = $1;
