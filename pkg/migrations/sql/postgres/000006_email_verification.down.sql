DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = 'users') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns
                  WHERE table_schema = current_schema() AND table_name = 'users' AND column_name = 'email_verification_expires_at') THEN
            ALTER TABLE users DROP COLUMN email_verification_expires_at;
        END IF;

        IF EXISTS (SELECT 1 FROM information_schema.columns
                  WHERE table_schema = current_schema() AND table_name = 'users' AND column_name = 'email_verification_token') THEN
            ALTER TABLE users DROP COLUMN email_verification_token;
        END IF;

        IF EXISTS (SELECT 1 FROM information_schema.columns
                  WHERE table_schema = current_schema() AND table_name = 'users' AND column_name = 'email_verified_at') THEN
            ALTER TABLE users DROP COLUMN email_verified_at;
        END IF;
    END IF;
END $$;
