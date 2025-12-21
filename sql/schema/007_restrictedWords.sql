-- +goose Up
Create Table RESTRICTEDWORDS (
    ID UUID Primary Key Default uuid_generate_v4(),
    WORD TEXT Not Null Unique,
    CREATED_AT TIMESTAMP Not Null Default now(),
    UPDATED_AT TIMESTAMP Not Null Default now()
);
