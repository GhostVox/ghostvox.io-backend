-- +goose Up
Create Table restrictedWords (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    word TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
);
