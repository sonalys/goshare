-- name: CreateUser :exec
INSERT INTO users (id,first_name,last_name,email,password_hash,created_at) VALUES ($1,$2,$3,$4,$5,$6);

-- name: FindUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateLedger :exec
INSERT INTO ledgers (id,name,created_at,created_by) VALUES ($1,$2,$3,$4);

-- name: FindLedgerById :one
SELECT * FROM ledgers WHERE id = $1;

-- name: AddUserToLedger :exec
INSERT INTO ledger_participants (id,ledger_id,user_id,created_at,created_by) VALUES ($1,$2,$3,$4,$5);

-- name: GetLedgerParticipants :many
SELECT * FROM ledger_participants WHERE ledger_id = $1;

-- name: CreateCategory :exec
INSERT INTO categories (id,ledger_id,name,parent_id,created_at,created_by) VALUES ($1,$2,$3,$4,$5,$6);

-- name: GetLedgerCategories :many
SELECT * FROM categories WHERE ledger_id = $1;

-- name: CreateExpense :exec
INSERT INTO expenses (id,category_id,ledger_id,amount,name,expense_date,created_at,created_by,updated_at,updated_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);

-- name: GetLedgerExpenses :many
SELECT * FROM expenses WHERE ledger_id = $1 ORDER BY expense_date DESC;

-- name: FindExpenseById :one
SELECT * FROM expenses WHERE id = $1;

-- name: UpdateExpense :exec
UPDATE expenses SET category_id = $1, amount = $2, name = $3, expense_date = $4, updated_at = $5, updated_by = $6 WHERE id = $7;

-- name: DeleteExpense :exec
DELETE FROM expenses WHERE id = $1;

-- name: CreateExpensePayment :exec
INSERT INTO expense_payments (id,expense_id,user_id,amount,payment_date,created_at,created_by,updated_at,updated_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);

-- name: GetExpensePayments :many
SELECT * FROM expense_payments WHERE expense_id = $1;

-- name: FindExpensePaymentById :one
SELECT * FROM expense_payments WHERE id = $1;

-- name: UpdateExpensePayment :exec
UPDATE expense_payments SET user_id = $1, amount = $2, payment_date = $3, updated_at = $4, updated_by = $5 WHERE id = $6;

-- name: DeleteExpensePayment :exec
DELETE FROM expense_payments WHERE id = $1;

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

-- name: GetLedgerParticipantsBalances :many
SELECT * FROM ledger_participant_balances WHERE ledger_id = $1;

-- name: GetUserLedgers :many
SELECT ledgers.* FROM ledgers JOIN ledger_participants ON ledgers.id = ledger_participants.ledger_id WHERE ledger_participants.user_id = $1;