-- Add email verification fields to users table
ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT false NOT NULL;
ALTER TABLE users ADD COLUMN email_verified_at TIMESTAMP WITH TIME ZONE;

-- For OAuth users, we'll trust the provider's email verification
-- Local users will need to verify their email separately
COMMENT ON COLUMN users.email_verified IS 'Whether the email address has been verified';
COMMENT ON COLUMN users.email_verified_at IS 'Timestamp when email was verified';

-- Create index for efficient lookups
CREATE INDEX idx_users_email_verified ON users(email_verified);
