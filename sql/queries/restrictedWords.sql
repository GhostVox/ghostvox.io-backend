

-- name: GetAllRestrictedWords :one
Select ARRAY_AGG(word) from restrictedWords;

-- name: AddRestrictedWord :one
Insert into restrictedWords (word) values ($1) returning *;


