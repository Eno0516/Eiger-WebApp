-- name: CreateUserProfile :one
INSERT INTO user_profiles (user_id, display_name, university, faculty, grade)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserProfileByUserID :one
SELECT * FROM user_profiles WHERE user_id = $1;

-- name: UpdateUserProfile :one
UPDATE user_profiles
SET display_name = $1,
    university   = $2,
    faculty      = $3,
    grade        = $4,
    updated_at   = NOW()
WHERE user_id = $5
RETURNING *;
