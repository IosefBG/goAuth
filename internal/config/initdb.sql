-- Grant select privilege to the user
CREATE USER limited_user WITH PASSWORD 'password';
GRANT SELECT ON ALL TABLES IN SCHEMA public TO limited_user;