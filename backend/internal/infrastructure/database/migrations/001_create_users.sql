-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    line_user_id    VARCHAR(255) NOT NULL UNIQUE,
    is_system_admin BOOLEAN NOT NULL DEFAULT false,
    is_registered   BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +migrate Down
DROP TABLE IF EXISTS users;
