-- +migrate Up
CREATE TABLE IF NOT EXISTS circles (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                 VARCHAR(255) NOT NULL,
    description          TEXT NOT NULL DEFAULT '',
    application_form_url TEXT NOT NULL DEFAULT '',
    spreadsheet_id       VARCHAR(255) NOT NULL DEFAULT '',
    uuid_column          VARCHAR(10) NOT NULL DEFAULT 'A',
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +migrate Down
DROP TABLE IF EXISTS circles;
