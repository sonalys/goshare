-- name: CreateLedger :exec
INSERT INTO ledgers (id,name,created_at,created_by) VALUES ($1,$2,$3,$4);

-- name: FindLedgerById :one
SELECT * FROM ledgers WHERE id = $1;

-- name: AddUserToLedger :exec
INSERT INTO ledger_participants (id,ledger_id,user_id,created_at,created_by) VALUES ($1,$2,$3,$4,$5);

-- name: GetLedgerParticipants :many
SELECT * FROM ledger_participants WHERE ledger_id = $1;

-- name: AppendLedgerRecord :exec
INSERT INTO ledger_records (id,ledger_id,expense_id,user_id,amount,created_at,created_by,description) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);

-- name: GetLedgerRecords :many
SELECT * FROM ledger_records WHERE ledger_id = $1 ORDER BY created_at DESC;

-- name: GetLedgerRecordsFromTimestamp :many
SELECT * FROM ledger_records WHERE ledger_id = $1 AND created_at > $2 ORDER BY created_at ASC;

-- name: GetLedgerUserRecords :many
SELECT * FROM ledger_records WHERE ledger_id = $1 AND user_id = $2 ORDER BY created_at DESC;

-- name: CreateLedgerParticipantBalance :exec
INSERT INTO ledger_participant_balances (id,ledger_id,user_id,last_timestamp,balance) VALUES ($1,$2,$3,$4,$5);

-- name: UpdateLedgerParticipantBalance :exec
UPDATE ledger_participant_balances SET last_timestamp = $1, balance = $2 WHERE ledger_id = $3 AND user_id = $4;

-- name: GetLedgerBalances :many
SELECT * FROM ledger_participant_balances WHERE ledger_id = $1;

-- name: GetUserLedgers :many
SELECT ledgers.* FROM ledgers JOIN ledger_participants ON ledgers.id = ledger_participants.ledger_id WHERE ledger_participants.user_id = $1;