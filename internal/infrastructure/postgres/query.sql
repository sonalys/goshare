-- name: CreateParticipant :exec
INSERT INTO participant (id,first_name,last_name,email,password_hash,created_at) VALUES ($1,$2,$3,$4,$5,$6);