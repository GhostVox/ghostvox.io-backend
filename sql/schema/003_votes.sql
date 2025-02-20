-- +goose Up
CREATE TABLE votes (
    id INT PRIMARY KEY auto_increment,
    poll_id INTEGER NOT NULL,
    option_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now (),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT votes_poll_id FOREIGN KEY (poll_id) REFERENCES polls (id) ON DELETE CASCADE,
    CONSTRAINT votes_option_id FOREIGN KEY (option_id) REFERENCES options (id) ON DELETE CASCADE,
    CONSTRAINT votes_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT votes_unique UNIQUE (poll_id, option_id, user_id)
);

CREATE TABLE options (
    id INT PRIMARY KEY auto_increment,
    poll_id INTEGER NOT NULL REFERENCES polls (id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now (),
    updated_at TIMESTAMP NOT NULL DEFAULT now ()
);
