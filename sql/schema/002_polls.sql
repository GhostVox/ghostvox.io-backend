-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE poll_status AS ENUM ('Active', 'Inactive', 'Archived');

CREATE TABLE polls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    category TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now (),
    updated_at TIMESTAMP NOT NULL DEFAULT now (),
    expires_at TIMESTAMP NOT NULL DEFAULT now () + INTERVAL '1 day',
    status poll_status NOT NULL DEFAULT 'Active',
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_poll_user ON polls (user_id);

CREATE INDEX idx_poll_status ON polls (status);

-- +goose Down
DROP TABLE polls;

DROP TYPE poll_status;
