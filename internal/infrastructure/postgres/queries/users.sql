-- name: CreateUser :exec
INSERT INTO users (id,first_name,last_name,email,password_hash,created_at) VALUES ($1,$2,$3,$4,$5,$6);

-- name: FindUserByEmail :one
SELECT * FROM user_view WHERE email = $1;

-- name: ListByEmail :many
SELECT * FROM user_view WHERE email = ANY(@emails::text[]);

-- name: FindUser :one
SELECT * FROM user_view WHERE id = $1 FOR UPDATE;