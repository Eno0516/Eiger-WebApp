-- name: CreateEvent :one
INSERT INTO events (circle_id, title, date, description, form_id, liff_url, line_group_url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetEventByID :one
SELECT * FROM events WHERE id = $1;

-- name: ListEventsByCircle :many
SELECT * FROM events
WHERE circle_id = $1
ORDER BY date DESC;

-- name: UpdateEvent :one
UPDATE events
SET title          = $1,
    date           = $2,
    description    = $3,
    form_id        = $4,
    liff_url       = $5,
    line_group_url = $6,
    confirmed_at   = $7,
    updated_at     = NOW()
WHERE id = $8
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events WHERE id = $1;
