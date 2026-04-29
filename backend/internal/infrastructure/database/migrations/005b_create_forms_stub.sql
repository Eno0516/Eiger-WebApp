-- +migrate Up
-- forms テーブルは将来的なフォームビルダー機能のためのスタブ。
-- events.form_id の外部キー参照を解決するために先に作成する。
CREATE TABLE IF NOT EXISTS forms (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +migrate Down
DROP TABLE IF EXISTS forms;
