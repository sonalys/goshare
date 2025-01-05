-- name: CreateCategory :exec
INSERT INTO categories (id,ledger_id,name,parent_id,created_at,created_by) VALUES ($1,$2,$3,$4,$5,$6);

-- name: GetLedgerCategories :many
SELECT * FROM categories WHERE ledger_id = $1;