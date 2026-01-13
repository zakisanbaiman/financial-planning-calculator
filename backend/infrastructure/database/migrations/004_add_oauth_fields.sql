-- Add OAuth provider fields to users table
ALTER TABLE users ADD COLUMN provider VARCHAR(50) DEFAULT 'local' NOT NULL;
ALTER TABLE users ADD COLUMN provider_user_id VARCHAR(255);
ALTER TABLE users ADD COLUMN name VARCHAR(255);
ALTER TABLE users ADD COLUMN avatar_url TEXT;

-- Make password_hash nullable for OAuth users (they don't have passwords)
ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;

-- Create index for efficient OAuth user lookups
CREATE INDEX idx_users_provider_user_id ON users(provider, provider_user_id);

-- Add comment for documentation
COMMENT ON COLUMN users.provider IS 'Authentication provider: local, github, google, etc.';
COMMENT ON COLUMN users.provider_user_id IS 'User ID from OAuth provider';
COMMENT ON COLUMN users.name IS 'Display name from OAuth provider';
COMMENT ON COLUMN users.avatar_url IS 'Profile avatar URL from OAuth provider';
