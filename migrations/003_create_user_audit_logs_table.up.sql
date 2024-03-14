CREATE TABLE user_audit_logs
(
    id         SERIAL PRIMARY KEY,
    user_id    INT          NOT NULL,
    action     VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);