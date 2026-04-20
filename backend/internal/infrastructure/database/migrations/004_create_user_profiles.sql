-- +migrate Up
CREATE TABLE IF NOT EXISTS user_profiles (
    user_id      UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    display_name VARCHAR(255) NOT NULL,
    university   VARCHAR(255),
    faculty      VARCHAR(255),
    grade        SMALLINT,           -- 学年（1〜4）。留年時は本人が更新
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +migrate Down
DROP TABLE IF EXISTS user_profiles;
