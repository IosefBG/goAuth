-- Create the users table
CREATE TABLE users
(
    id             SERIAL PRIMARY KEY,
    username       VARCHAR(255) UNIQUE NOT NULL,
    password       VARCHAR(255) NOT NULL,
    email          VARCHAR(255) UNIQUE NOT NULL,
    is_blocked     BOOLEAN   DEFAULT false,
    login_attempts INT       DEFAULT 0,
    last_login     TIMESTAMP,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active     BOOLEAN   DEFAULT true
);

