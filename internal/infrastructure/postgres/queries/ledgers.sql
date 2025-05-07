-- name: CreateLedger :exec
INSERT INTO ledgers (id,name,created_at,created_by) VALUES ($1,$2,$3,$4);

-- name: FindLedgerById :one
SELECT * FROM ledgers WHERE id = $1;

-- name: AddUserToLedger :exec
INSERT INTO ledger_participants (id,ledger_id,user_id,created_at,created_by,balance) VALUES ($1,$2,$3,$4,$5,$6);

-- name: RemoveUserFromLedger :exec
DELETE FROM ledger_participants WHERE id = $1;

-- name: GetLedgerParticipants :many
SELECT * FROM ledger_participants WHERE ledger_id = $1;

-- name: GetUserLedgers :many
SELECT ledgers.* FROM ledgers JOIN ledger_participants ON ledgers.id = ledger_participants.ledger_id WHERE ledger_participants.user_id = $1 ORDER BY ledgers.created_at DESC;

-- name: LockLedgerForUpdate :exec
SELECT * FROM ledgers WHERE id = $1 FOR UPDATE;

-- name: CountLedgerUsers :one
SELECT COUNT(*) FROM ledgers WHERE id = $1;

-- name: LockUserForUpdate :exec
SELECT * FROM users WHERE id = $1 FOR UPDATE;

-- name: CountUserLedgers :one
SELECT COUNT(*) FROM ledgers WHERE created_by = $1;

-- name: UpdateLedger :exec
UPDATE ledgers SET name = $1 WHERE id = $2;