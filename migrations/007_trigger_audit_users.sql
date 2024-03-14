-- 007_trigger_audit_users.sql

-- Create the function to log audit
CREATE OR REPLACE FUNCTION log_audit(
    _table_name TEXT,
    _action TEXT,
    _row_id INTEGER,
    _old_data JSONB,
    _new_data JSONB
) RETURNS VOID AS $$
BEGIN
    INSERT INTO user_audit_logs (table_name, action, row_id, old_data, new_data)
    VALUES (_table_name, _action, _row_id, _old_data, _new_data);
END;
$$ LANGUAGE plpgsql;

-- Create the trigger to audit insert on users
CREATE TRIGGER audit_insert_trigger
    AFTER INSERT ON users
    FOR EACH ROW
EXECUTE FUNCTION log_audit('users', 'INSERT', NEW.id, NULL, to_jsonb(NEW));

-- Create the trigger to audit update on users
CREATE TRIGGER audit_update_trigger
    AFTER UPDATE ON users
    FOR EACH ROW
EXECUTE FUNCTION log_audit('users', 'UPDATE', NEW.id, to_jsonb(OLD), to_jsonb(NEW));

-- Create the trigger to audit delete on users (example)
CREATE TRIGGER audit_delete_trigger
    AFTER DELETE ON users
    FOR EACH ROW
EXECUTE FUNCTION log_audit('users', 'DELETE', OLD.id, to_jsonb(OLD), NULL);
