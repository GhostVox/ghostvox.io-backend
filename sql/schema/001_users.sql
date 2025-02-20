-- +goose Up
Create table users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    created_at TIMESTAMP NOT NULL DEFAULT now (),
    updated_at TIMESTAMP NOT NULL DEFAULT now (),
    email TEXT NOT NULL UNIQUE,
    user_token TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL DEFAULT 'user'
);

-- +goose Down
DROP TABLE users;
