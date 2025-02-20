-- +goose Up
CREATE TABLE polls (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now (),
    updated_at TIMESTAMP NOT NULL DEFAULT now (),
    expires_at TIMESTAMP NOT NULL DEFAULT now () + INTERVAL '1 day',
    status TEXT NOT NULL DEFAULT 'active',
);

-- +goose Down
DROP TABLE polls;
