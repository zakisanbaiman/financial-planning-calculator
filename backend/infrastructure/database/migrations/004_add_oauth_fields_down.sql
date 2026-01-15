-- Rollback OAuth fields from users table
DROP INDEX IF EXISTS idx_users_provider_user_id;

ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
ALTER TABLE users DROP COLUMN IF EXISTS name;
ALTER TABLE users DROP COLUMN IF EXISTS provider_user_id;
ALTER TABLE users DROP COLUMN IF EXISTS provider;

-- Restore password_hash NOT NULL constraint
-- Note: This will fail if there are OAuth users without passwords
-- You should migrate OAuth users before rolling back
ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;
