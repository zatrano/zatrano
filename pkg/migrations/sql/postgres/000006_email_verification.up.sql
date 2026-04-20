-- Optional columns on users (when a users table exists).
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = 'users') THEN
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                      WHERE table_schema = current_schema() AND table_name = 'users' AND column_name = 'email_verified_at') THEN
            ALTER TABLE users ADD COLUMN email_verified_at TIMESTAMPTZ NULL;
        END IF;

        IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                      WHERE table_schema = current_schema() AND table_name = 'users' AND column_name = 'email_verification_token') THEN
            ALTER TABLE users ADD COLUMN email_verification_token VARCHAR(255) NULL;
        END IF;

        IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                      WHERE table_schema = current_schema() AND table_name = 'users' AND column_name = 'email_verification_expires_at') THEN
            ALTER TABLE users ADD COLUMN email_verification_expires_at TIMESTAMPTZ NULL;
        END IF;

        IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE schemaname = current_schema() AND tablename = 'users' AND indexname = 'idx_users_email_verification_token') THEN
            CREATE INDEX idx_users_email_verification_token ON users(email_verification_token);
        END IF;
    END IF;
END $$;
