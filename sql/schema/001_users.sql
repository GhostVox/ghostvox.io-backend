-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

Create table users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    created_at TIMESTAMP NOT NULL DEFAULT now (),
    updated_at TIMESTAMP NOT NULL DEFAULT now (),
    email TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT default NULL,
    user_token TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL DEFAULT 'user'
);

-- +goose Down
DROP TABLE users;

Drop extension "uuid-ossp";
