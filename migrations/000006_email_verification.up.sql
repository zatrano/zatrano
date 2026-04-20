-- +migrate Up
-- Add email verification fields to users table (if exists)
-- This migration assumes you have a users table with id, email fields
-- Adjust according to your actual users table structure

-- Add email verification columns if users table exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users') THEN
        -- Add email_verified_at column if it doesn't exist
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                      WHERE table_name = 'users' AND column_name = 'email_verified_at') THEN
            ALTER TABLE users ADD COLUMN email_verified_at TIMESTAMP NULL;
        END IF;

        -- Add email_verification_token column if it doesn't exist
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                      WHERE table_name = 'users' AND column_name = 'email_verification_token') THEN
            ALTER TABLE users ADD COLUMN email_verification_token VARCHAR(255) NULL;
        END IF;

        -- Add email_verification_expires_at column if it doesn't exist
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                      WHERE table_name = 'users' AND column_name = 'email_verification_expires_at') THEN
            ALTER TABLE users ADD COLUMN email_verification_expires_at TIMESTAMP NULL;
        END IF;

        -- Create index for email verification token
        IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE tablename = 'users' AND indexname = 'idx_users_email_verification_token') THEN
            CREATE INDEX idx_users_email_verification_token ON users(email_verification_token);
        END IF;
    END IF;
END $$;

-- +migrate Down
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users') THEN
        -- Drop columns if they exist
        IF EXISTS (SELECT 1 FROM information_schema.columns
                  WHERE table_name = 'users' AND column_name = 'email_verification_expires_at') THEN
            ALTER TABLE users DROP COLUMN email_verification_expires_at;
        END IF;

        IF EXISTS (SELECT 1 FROM information_schema.columns
                  WHERE table_name = 'users' AND column_name = 'email_verification_token') THEN
            ALTER TABLE users DROP COLUMN email_verification_token;
        END IF;

        IF EXISTS (SELECT 1 FROM information_schema.columns
                  WHERE table_name = 'users' AND column_name = 'email_verified_at') THEN
            ALTER TABLE users DROP COLUMN email_verified_at;
        END IF;
    END IF;
END $$;