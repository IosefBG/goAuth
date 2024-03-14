-- 005_add_foreign_key_to_user_audit_logs.up.sql

-- Add foreign key constraint to user_audit_logs table
ALTER TABLE user_audit_logs
    ADD CONSTRAINT fk_user_audit_logs_user_id FOREIGN KEY (user_id) REFERENCES users(id);
