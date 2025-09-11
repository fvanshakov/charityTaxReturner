CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL,
    email_hmac TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    oauth_token BYTEA
);

CREATE INDEX IF NOT EXISTS idx_users_email_hmac ON users(email_hmac);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);