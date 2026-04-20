-- +migrate Up
CREATE TABLE IF NOT EXISTS user_line_accounts (
    user_id   UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    line_id   VARCHAR(255) NOT NULL,  -- 公式LINEアカウントでのメッセージ送信先
    linked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +migrate Down
DROP TABLE IF EXISTS user_line_accounts;
