-- +migrate Up
ALTER TABLE users
  ADD COLUMN password_hash                  VARCHAR(255) NOT NULL DEFAULT '',
  ADD COLUMN role                           VARCHAR(50)  NOT NULL DEFAULT 'pending',
  ADD COLUMN email_verified_at              TIMESTAMPTZ,
  ADD COLUMN email_verification_token       VARCHAR(255),
  ADD COLUMN email_verification_expires_at  TIMESTAMPTZ,
  ADD COLUMN password_reset_token           VARCHAR(255),
  ADD COLUMN password_reset_expires_at      TIMESTAMPTZ;

-- DEFAULT '' は既存レコードなし前提の workaround。
-- 実運用では password_hash NOT NULL のみで問題ない。

-- +migrate Down
ALTER TABLE users
  DROP COLUMN IF EXISTS password_hash,
  DROP COLUMN IF EXISTS role,
  DROP COLUMN IF EXISTS email_verified_at,
  DROP COLUMN IF EXISTS email_verification_token,
  DROP COLUMN IF EXISTS email_verification_expires_at,
  DROP COLUMN IF EXISTS password_reset_token,
  DROP COLUMN IF EXISTS password_reset_expires_at;
