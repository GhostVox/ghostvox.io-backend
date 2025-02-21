-- +goose Up
CREATE TABLE options (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name TEXT NOT NULL,
    poll_id UUID NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now (),
    updated_at TIMESTAMP NOT NULL DEFAULT now (),
    CONSTRAINT options_poll_id FOREIGN KEY (poll_id) REFERENCES polls (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE options;
