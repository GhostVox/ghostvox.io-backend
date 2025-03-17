-- +goose Up
CREATE TABLE options (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name TEXT NOT NULL,
    poll_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now (),
    updated_at TIMESTAMP NOT NULL DEFAULT now (),
    CONSTRAINT options_poll_id FOREIGN KEY (poll_id) REFERENCES polls (id) ON DELETE CASCADE,
    CONSTRAINT unique_option_per_poll UNIQUE (poll_id, name)
);

CREATE INDEX idx_options_poll_id ON options (poll_id);

-- +goose Down
DROP TABLE options;
