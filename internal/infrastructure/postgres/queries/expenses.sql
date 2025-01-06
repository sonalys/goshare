-- name: CreateExpense :exec
INSERT INTO expenses (id,category_id,ledger_id,amount,name,expense_date,created_at,created_by,updated_at,updated_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);

-- name: GetLedgerExpenses :many
SELECT * FROM expenses WHERE ledger_id = $1 AND created_at < $2 ORDER BY created_at DESC LIMIT $3;

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