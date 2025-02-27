-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP NOT NULL DEFAULT now (),
    updated_at TIMESTAMP NOT NULL DEFAULT now (),
    email TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT DEFAULT NULL,
    hashed_password TEXT DEFAULT NULL, -- Nullable for OAuth users
    provider TEXT DEFAULT NULL, -- e.g., 'google', 'github'
    provider_id TEXT DEFAULT NULL, -- Stores OAuth provider’s unique user ID
    refresh_token TEXT DEFAULT NULL, -- For JWT refresh flow (if needed)
    role TEXT NOT NULL DEFAULT 'user'
);

CREATE UNIQUE INDEX unique_provider_id ON users (provider, provider_id);

-- +goose Down
DROP TABLE users;

Drop extension "uuid-ossp";
