-- 004_add_foreign_key_to_user_sessions.up.sql

-- Add foreign key constraint to user_sessions table
ALTER TABLE user_sessions
    ADD CONSTRAINT fk_user_sessions_user_id FOREIGN KEY (user_id) REFERENCES users(id);
