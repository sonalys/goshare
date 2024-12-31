-- name: CreateUser :exec
INSERT INTO users (id,first_name,last_name,email,password_hash,created_at) VALUES ($1,$2,$3,$4,$5,$6);

-- name: FindUserByEmail :one
SELECT * FROM users WHERE email = $1;