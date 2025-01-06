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

-- name: UpsertLedgerParticipantBalance :exec
INSERT INTO ledger_participant_balances (id, ledger_id, user_id, last_timestamp, balance)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (ledger_id, user_id) 
DO UPDATE SET 
    last_timestamp = EXCLUDED.last_timestamp,
    balance = EXCLUDED.balance;

-- name: GetLedgerBalances :many
SELECT * FROM ledger_participant_balances WHERE ledger_id = $1;

-- name: GetUserLedgers :many
SELECT ledgers.* FROM ledgers JOIN ledger_participants ON ledgers.id = ledger_participants.ledger_id WHERE ledger_participants.user_id = $1 ORDER BY ledgers.created_at DESC;

-- name: GetLedgerParticipantsWithBalance :many
SELECT 
    lp.ledger_id,
    lp.user_id,
    lp.created_by,
    MAX(lr.created_at)::TIMESTAMP AS last_timestamp,
    COALESCE(lpb.balance, 0) + COALESCE(SUM(lr.amount), 0) AS balance
FROM 
    ledger_participants lp
LEFT JOIN 
    ledger_participant_balances lpb ON lp.ledger_id = lpb.ledger_id AND lp.user_id = lpb.user_id
LEFT JOIN 
    ledger_records lr ON lp.ledger_id = lr.ledger_id AND lp.user_id = lr.user_id AND lr.created_at > lpb.last_timestamp
WHERE 
    lp.ledger_id = $1
GROUP BY 
    lp.ledger_id, lp.user_id, lp.created_at, lp.created_by, lpb.balance
ORDER BY
    lp.user_id;

-- name: LockLedgerForUpdate :exec
SELECT * FROM ledgers WHERE id = $1 FOR UPDATE;

-- name: CountLedgerUsers :one
SELECT COUNT(*) FROM ledgers WHERE id = $1;

-- name: LockUserForUpdate :exec
SELECT * FROM users WHERE id = $1 FOR UPDATE;

-- name: CountUserLedgers :one
SELECT COUNT(*) FROM ledgers WHERE created_by = $1;

-- name: GetExpensesRecords :many
SELECT * FROM ledger_records WHERE expense_id IN (SELECT unnest($1::uuid[]));