-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE votes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    poll_id UUID NOT NULL,
    option_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now (),
    user_id UUID NOT NULL,
    CONSTRAINT votes_poll_id FOREIGN KEY (poll_id) REFERENCES polls (id) ON DELETE CASCADE,
    CONSTRAINT votes_option_id FOREIGN KEY (option_id) REFERENCES options (id) ON DELETE CASCADE,
    CONSTRAINT votes_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT votes_unique UNIQUE (poll_id, user_id)
);

-- +goose Down
DROP TABLE votes;
