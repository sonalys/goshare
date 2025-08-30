-- name: SaveUser :exec
INSERT INTO users (id,first_name,last_name,email,password_hash,ledger_count,created_at) 
VALUES ($1,$2,$3,$4,$5,$6,$7)
ON CONFLICT (id)
DO UPDATE
SET
first_name = EXCLUDED.first_name,
last_name = EXCLUDED.last_name,
email = EXCLUDED.email,
password_hash = EXCLUDED.password_hash,
ledger_count = EXCLUDED.ledger_count
;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListByEmail :many
SELECT * FROM users WHERE email = ANY(@emails::text[]);

-- name: GetUser :one
SELECT * FROM users WHERE id = $1 FOR UPDATE;