-- name: UpsertEventParticipant :one
INSERT INTO event_participants (event_id, user_id, status)
VALUES ($1, $2, $3)
ON CONFLICT (event_id, user_id) DO UPDATE
    SET status     = EXCLUDED.status,
        updated_at = NOW()
RETURNING *;

-- name: GetEventParticipantByEventAndUser :one
SELECT * FROM event_participants
WHERE event_id = $1 AND user_id = $2;

-- name: ListEventParticipantsByEvent :many
SELECT * FROM event_participants
WHERE event_id = $1
ORDER BY created_at ASC;

-- name: ListConfirmedEventParticipantsByEvent :many
SELECT * FROM event_participants
WHERE event_id = $1 AND status = 'confirmed'
ORDER BY created_at ASC;
