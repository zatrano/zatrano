-- Requires a users table (create it in an earlier app migration if needed).
ALTER TABLE users ADD COLUMN email_verified_at TEXT;
ALTER TABLE users ADD COLUMN email_verification_token TEXT;
ALTER TABLE users ADD COLUMN email_verification_expires_at TEXT;
CREATE INDEX IF NOT EXISTS idx_users_email_verification_token ON users (email_verification_token);
