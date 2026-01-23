-- 006_add_two_factor_auth_down.sql
-- 2段階認証（2FA）機能のロールバック

DROP INDEX IF EXISTS idx_users_two_factor_enabled;

ALTER TABLE users DROP COLUMN IF EXISTS two_factor_backup_codes;
ALTER TABLE users DROP COLUMN IF EXISTS two_factor_secret;
ALTER TABLE users DROP COLUMN IF EXISTS two_factor_enabled;
