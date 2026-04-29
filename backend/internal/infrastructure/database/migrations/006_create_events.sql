-- +migrate Up
CREATE TABLE IF NOT EXISTS events (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    circle_id      UUID NOT NULL REFERENCES circles(id) ON DELETE CASCADE,
    title          VARCHAR(255) NOT NULL,
    date           DATE NOT NULL,
    description    TEXT NOT NULL DEFAULT '',
    form_id        UUID REFERENCES forms(id) ON DELETE SET NULL,
    liff_url       TEXT NOT NULL DEFAULT '',
    line_group_url TEXT NOT NULL DEFAULT '',
    confirmed_at   TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +migrate Down
DROP TABLE IF EXISTS events;
