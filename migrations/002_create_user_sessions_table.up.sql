CREATE TABLE user_sessions
(
    id            SERIAL PRIMARY KEY,
    user_id       INT          NOT NULL,
    session_token VARCHAR(255) NOT NULL,
    ip_address    VARCHAR(255) NOT NULL,
    location    VARCHAR(255) NOT NULL,
    device_connected VARCHAR(255),
    browser_used    VARCHAR(255),
    is_active     BOOLEAN   DEFAULT true,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
