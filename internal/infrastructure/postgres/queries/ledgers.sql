-- name: CreateLedger :exec
INSERT INTO ledgers (id,name,created_at,created_by) VALUES ($1,$2,$3,$4);

-- name: GetLedgerById :one
SELECT * FROM ledgers WHERE id = $1 FOR UPDATE;

-- name: CreateLedgerMember :exec
INSERT INTO ledger_members (ledger_id,user_id,created_at,created_by,balance) 
VALUES ($1,$2,$3,$4,$5) 
ON CONFLICT(user_id) 
DO UPDATE
SET balance = EXCLUDED.balance
;

-- name: RemoveUserFromLedger :exec
DELETE FROM ledger_members WHERE user_id = $1;

-- name: GetLedgerMembers :many
SELECT * FROM ledger_members WHERE ledger_id = $1;

-- name: GetUserLedgers :many
SELECT ledgers.* FROM ledgers 
JOIN ledger_members ON 
    ledgers.id = ledger_members.ledger_id 
WHERE 
    ledgers.created_by = $1 OR
    ledger_members.user_id = $1 
ORDER BY 
    ledgers.created_at DESC;

-- name: LockLedgerForUpdate :exec
SELECT * FROM ledgers WHERE id = $1 FOR UPDATE;

-- name: CountLedgerUsers :one
SELECT COUNT(*) FROM ledgers WHERE id = $1;

-- name: LockUserForUpdate :exec
SELECT * FROM users WHERE id = $1 FOR UPDATE;

-- name: CountUserLedgers :one
SELECT COUNT(*) FROM ledgers WHERE created_by = $1;

-- name: UpdateLedger :many
UPDATE ledgers SET name = $1 WHERE id = $2 RETURNING id;

-- name: DeleteMembersNotIn :exec
DELETE FROM ledger_members WHERE user_id != ALL(sqlc.arg('ids')::uuid[]);