-- name: CreateExpense :exec
INSERT INTO expenses (id,ledger_id,amount,name,expense_date,created_at,created_by,updated_at,updated_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);

-- name: GetLedgerExpenses :many
SELECT * FROM expenses WHERE ledger_id = $1 AND created_at < $2 ORDER BY created_at DESC LIMIT $3;

-- name: FindExpenseById :one
SELECT * FROM expenses WHERE id = $1;

-- name: UpdateExpense :exec
UPDATE expenses SET amount = $$, name = $$, expense_date = $$, updated_at = $$, updated_by = $$ WHERE id = $$;

-- name: DeleteExpense :exec
DELETE FROM expenses WHERE id = $1;

-- name: CreateExpenseRecord :exec
INSERT INTO expense_records (id,expense_id,record_type,amount,from_user_id,to_user_id,created_at,created_by,updated_at,updated_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);

-- name: GetExpenseRecords :many
SELECT * FROM expense_records WHERE expense_id = $1 ORDER BY created_at DESC LIMIT $2;

