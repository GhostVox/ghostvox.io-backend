-- +goose Up
CREATE TABLE refresh_tokens (
    user_id UUID Primary KEY NULL,
    token TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    expires_at TIMESTAMP NOT NULL,
    CONSTRAINT unique_token UNIQUE (token, user_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE refresh_tokens;
