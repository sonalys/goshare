// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addUserToLedger = `-- name: AddUserToLedger :exec
INSERT INTO ledger_participants (id,ledger_id,user_id,created_at,created_by) VALUES ($1,$2,$3,$4,$5)
`

type AddUserToLedgerParams struct {
	ID        pgtype.UUID
	LedgerID  pgtype.UUID
	UserID    pgtype.UUID
	CreatedAt pgtype.Timestamp
	CreatedBy pgtype.UUID
}

func (q *Queries) AddUserToLedger(ctx context.Context, arg AddUserToLedgerParams) error {
	_, err := q.db.Exec(ctx, addUserToLedger,
		arg.ID,
		arg.LedgerID,
		arg.UserID,
		arg.CreatedAt,
		arg.CreatedBy,
	)
	return err
}

const appendLedgerRecord = `-- name: AppendLedgerRecord :exec
INSERT INTO ledger_records (id,ledger_id,expense_id,user_id,amount,created_at,created_by,description) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
`

type AppendLedgerRecordParams struct {
	ID          pgtype.UUID
	LedgerID    pgtype.UUID
	ExpenseID   pgtype.UUID
	UserID      pgtype.UUID
	Amount      int32
	CreatedAt   pgtype.Timestamp
	CreatedBy   pgtype.UUID
	Description string
}

func (q *Queries) AppendLedgerRecord(ctx context.Context, arg AppendLedgerRecordParams) error {
	_, err := q.db.Exec(ctx, appendLedgerRecord,
		arg.ID,
		arg.LedgerID,
		arg.ExpenseID,
		arg.UserID,
		arg.Amount,
		arg.CreatedAt,
		arg.CreatedBy,
		arg.Description,
	)
	return err
}

const createCategory = `-- name: CreateCategory :exec
INSERT INTO categories (id,ledger_id,name,parent_id,created_at,created_by) VALUES ($1,$2,$3,$4,$5,$6)
`

type CreateCategoryParams struct {
	ID        pgtype.UUID
	LedgerID  pgtype.UUID
	Name      string
	ParentID  pgtype.UUID
	CreatedAt pgtype.Timestamp
	CreatedBy pgtype.UUID
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) error {
	_, err := q.db.Exec(ctx, createCategory,
		arg.ID,
		arg.LedgerID,
		arg.Name,
		arg.ParentID,
		arg.CreatedAt,
		arg.CreatedBy,
	)
	return err
}

const createExpense = `-- name: CreateExpense :exec
INSERT INTO expenses (id,category_id,ledger_id,amount,name,expense_date,created_at,created_by,updated_at,updated_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
`

type CreateExpenseParams struct {
	ID          pgtype.UUID
	CategoryID  pgtype.UUID
	LedgerID    pgtype.UUID
	Amount      int32
	Name        string
	ExpenseDate pgtype.Timestamp
	CreatedAt   pgtype.Timestamp
	CreatedBy   pgtype.UUID
	UpdatedAt   pgtype.Timestamp
	UpdatedBy   pgtype.UUID
}

func (q *Queries) CreateExpense(ctx context.Context, arg CreateExpenseParams) error {
	_, err := q.db.Exec(ctx, createExpense,
		arg.ID,
		arg.CategoryID,
		arg.LedgerID,
		arg.Amount,
		arg.Name,
		arg.ExpenseDate,
		arg.CreatedAt,
		arg.CreatedBy,
		arg.UpdatedAt,
		arg.UpdatedBy,
	)
	return err
}

const createExpensePayment = `-- name: CreateExpensePayment :exec
INSERT INTO expense_payments (id,expense_id,user_id,amount,payment_date,created_at,created_by,updated_at,updated_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
`

type CreateExpensePaymentParams struct {
	ID          pgtype.UUID
	ExpenseID   pgtype.UUID
	UserID      pgtype.UUID
	Amount      int32
	PaymentDate pgtype.Timestamp
	CreatedAt   pgtype.Timestamp
	CreatedBy   pgtype.UUID
	UpdatedAt   pgtype.Timestamp
	UpdatedBy   pgtype.UUID
}

func (q *Queries) CreateExpensePayment(ctx context.Context, arg CreateExpensePaymentParams) error {
	_, err := q.db.Exec(ctx, createExpensePayment,
		arg.ID,
		arg.ExpenseID,
		arg.UserID,
		arg.Amount,
		arg.PaymentDate,
		arg.CreatedAt,
		arg.CreatedBy,
		arg.UpdatedAt,
		arg.UpdatedBy,
	)
	return err
}

const createLedger = `-- name: CreateLedger :exec
INSERT INTO ledgers (id,name,created_at,created_by) VALUES ($1,$2,$3,$4)
`

type CreateLedgerParams struct {
	ID        pgtype.UUID
	Name      string
	CreatedAt pgtype.Timestamp
	CreatedBy pgtype.UUID
}

func (q *Queries) CreateLedger(ctx context.Context, arg CreateLedgerParams) error {
	_, err := q.db.Exec(ctx, createLedger,
		arg.ID,
		arg.Name,
		arg.CreatedAt,
		arg.CreatedBy,
	)
	return err
}

const createLedgerParticipantBalance = `-- name: CreateLedgerParticipantBalance :exec
INSERT INTO ledger_participant_balances (id,ledger_id,user_id,last_timestamp,balance) VALUES ($1,$2,$3,$4,$5)
`

type CreateLedgerParticipantBalanceParams struct {
	ID            pgtype.UUID
	LedgerID      pgtype.UUID
	UserID        pgtype.UUID
	LastTimestamp pgtype.Timestamp
	Balance       int32
}

func (q *Queries) CreateLedgerParticipantBalance(ctx context.Context, arg CreateLedgerParticipantBalanceParams) error {
	_, err := q.db.Exec(ctx, createLedgerParticipantBalance,
		arg.ID,
		arg.LedgerID,
		arg.UserID,
		arg.LastTimestamp,
		arg.Balance,
	)
	return err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users (id,first_name,last_name,email,password_hash,created_at) VALUES ($1,$2,$3,$4,$5,$6)
`

type CreateUserParams struct {
	ID           pgtype.UUID
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	CreatedAt    pgtype.Timestamp
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.PasswordHash,
		arg.CreatedAt,
	)
	return err
}

const deleteExpense = `-- name: DeleteExpense :exec
DELETE FROM expenses WHERE id = $1
`

func (q *Queries) DeleteExpense(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteExpense, id)
	return err
}

const deleteExpensePayment = `-- name: DeleteExpensePayment :exec
DELETE FROM expense_payments WHERE id = $1
`

func (q *Queries) DeleteExpensePayment(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteExpensePayment, id)
	return err
}

const findExpenseById = `-- name: FindExpenseById :one
SELECT id, category_id, ledger_id, amount, name, expense_date, created_at, created_by, updated_at, updated_by FROM expenses WHERE id = $1
`

func (q *Queries) FindExpenseById(ctx context.Context, id pgtype.UUID) (Expense, error) {
	row := q.db.QueryRow(ctx, findExpenseById, id)
	var i Expense
	err := row.Scan(
		&i.ID,
		&i.CategoryID,
		&i.LedgerID,
		&i.Amount,
		&i.Name,
		&i.ExpenseDate,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
	)
	return i, err
}

const findExpensePaymentById = `-- name: FindExpensePaymentById :one
SELECT id, expense_id, user_id, ledger_id, amount, payment_date, created_at, created_by, updated_at, updated_by FROM expense_payments WHERE id = $1
`

func (q *Queries) FindExpensePaymentById(ctx context.Context, id pgtype.UUID) (ExpensePayment, error) {
	row := q.db.QueryRow(ctx, findExpensePaymentById, id)
	var i ExpensePayment
	err := row.Scan(
		&i.ID,
		&i.ExpenseID,
		&i.UserID,
		&i.LedgerID,
		&i.Amount,
		&i.PaymentDate,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
	)
	return i, err
}

const findLedgerById = `-- name: FindLedgerById :one
SELECT id, name, created_at, created_by FROM ledgers WHERE id = $1
`

func (q *Queries) FindLedgerById(ctx context.Context, id pgtype.UUID) (Ledger, error) {
	row := q.db.QueryRow(ctx, findLedgerById, id)
	var i Ledger
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.CreatedBy,
	)
	return i, err
}

const findUserByEmail = `-- name: FindUserByEmail :one
SELECT id, first_name, last_name, email, password_hash, created_at FROM users WHERE email = $1
`

func (q *Queries) FindUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, findUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
	)
	return i, err
}

const getExpensePayments = `-- name: GetExpensePayments :many
SELECT id, expense_id, user_id, ledger_id, amount, payment_date, created_at, created_by, updated_at, updated_by FROM expense_payments WHERE expense_id = $1
`

func (q *Queries) GetExpensePayments(ctx context.Context, expenseID pgtype.UUID) ([]ExpensePayment, error) {
	rows, err := q.db.Query(ctx, getExpensePayments, expenseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExpensePayment
	for rows.Next() {
		var i ExpensePayment
		if err := rows.Scan(
			&i.ID,
			&i.ExpenseID,
			&i.UserID,
			&i.LedgerID,
			&i.Amount,
			&i.PaymentDate,
			&i.CreatedAt,
			&i.CreatedBy,
			&i.UpdatedAt,
			&i.UpdatedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLedgerCategories = `-- name: GetLedgerCategories :many
SELECT id, ledger_id, name, parent_id, created_at, created_by FROM categories WHERE ledger_id = $1
`

func (q *Queries) GetLedgerCategories(ctx context.Context, ledgerID pgtype.UUID) ([]Category, error) {
	rows, err := q.db.Query(ctx, getLedgerCategories, ledgerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.LedgerID,
			&i.Name,
			&i.ParentID,
			&i.CreatedAt,
			&i.CreatedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLedgerExpenses = `-- name: GetLedgerExpenses :many
SELECT id, category_id, ledger_id, amount, name, expense_date, created_at, created_by, updated_at, updated_by FROM expenses WHERE ledger_id = $1 ORDER BY expense_date DESC
`

func (q *Queries) GetLedgerExpenses(ctx context.Context, ledgerID pgtype.UUID) ([]Expense, error) {
	rows, err := q.db.Query(ctx, getLedgerExpenses, ledgerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Expense
	for rows.Next() {
		var i Expense
		if err := rows.Scan(
			&i.ID,
			&i.CategoryID,
			&i.LedgerID,
			&i.Amount,
			&i.Name,
			&i.ExpenseDate,
			&i.CreatedAt,
			&i.CreatedBy,
			&i.UpdatedAt,
			&i.UpdatedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLedgerParticipants = `-- name: GetLedgerParticipants :many
SELECT id, ledger_id, user_id, created_at, created_by FROM ledger_participants WHERE ledger_id = $1
`

func (q *Queries) GetLedgerParticipants(ctx context.Context, ledgerID pgtype.UUID) ([]LedgerParticipant, error) {
	rows, err := q.db.Query(ctx, getLedgerParticipants, ledgerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LedgerParticipant
	for rows.Next() {
		var i LedgerParticipant
		if err := rows.Scan(
			&i.ID,
			&i.LedgerID,
			&i.UserID,
			&i.CreatedAt,
			&i.CreatedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLedgerRecords = `-- name: GetLedgerRecords :many
SELECT id, ledger_id, expense_id, user_id, amount, created_at, created_by, description FROM ledger_records WHERE ledger_id = $1 ORDER BY created_at DESC
`

func (q *Queries) GetLedgerRecords(ctx context.Context, ledgerID pgtype.UUID) ([]LedgerRecord, error) {
	rows, err := q.db.Query(ctx, getLedgerRecords, ledgerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LedgerRecord
	for rows.Next() {
		var i LedgerRecord
		if err := rows.Scan(
			&i.ID,
			&i.LedgerID,
			&i.ExpenseID,
			&i.UserID,
			&i.Amount,
			&i.CreatedAt,
			&i.CreatedBy,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLedgerUserRecords = `-- name: GetLedgerUserRecords :many
SELECT id, ledger_id, expense_id, user_id, amount, created_at, created_by, description FROM ledger_records WHERE ledger_id = $1 AND user_id = $2 ORDER BY created_at DESC
`

type GetLedgerUserRecordsParams struct {
	LedgerID pgtype.UUID
	UserID   pgtype.UUID
}

func (q *Queries) GetLedgerUserRecords(ctx context.Context, arg GetLedgerUserRecordsParams) ([]LedgerRecord, error) {
	rows, err := q.db.Query(ctx, getLedgerUserRecords, arg.LedgerID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LedgerRecord
	for rows.Next() {
		var i LedgerRecord
		if err := rows.Scan(
			&i.ID,
			&i.LedgerID,
			&i.ExpenseID,
			&i.UserID,
			&i.Amount,
			&i.CreatedAt,
			&i.CreatedBy,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserLedgers = `-- name: GetUserLedgers :many
SELECT ledgers.id, ledgers.name, ledgers.created_at, ledgers.created_by FROM ledgers JOIN ledger_participants ON ledgers.id = ledger_participants.ledger_id WHERE ledger_participants.user_id = $1
`

func (q *Queries) GetUserLedgers(ctx context.Context, userID pgtype.UUID) ([]Ledger, error) {
	rows, err := q.db.Query(ctx, getUserLedgers, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Ledger
	for rows.Next() {
		var i Ledger
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.CreatedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateExpense = `-- name: UpdateExpense :exec
UPDATE expenses SET category_id = $1, amount = $2, name = $3, expense_date = $4, updated_at = $5, updated_by = $6 WHERE id = $7
`

type UpdateExpenseParams struct {
	CategoryID  pgtype.UUID
	Amount      int32
	Name        string
	ExpenseDate pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
	UpdatedBy   pgtype.UUID
	ID          pgtype.UUID
}

func (q *Queries) UpdateExpense(ctx context.Context, arg UpdateExpenseParams) error {
	_, err := q.db.Exec(ctx, updateExpense,
		arg.CategoryID,
		arg.Amount,
		arg.Name,
		arg.ExpenseDate,
		arg.UpdatedAt,
		arg.UpdatedBy,
		arg.ID,
	)
	return err
}

const updateExpensePayment = `-- name: UpdateExpensePayment :exec
UPDATE expense_payments SET user_id = $1, amount = $2, payment_date = $3, updated_at = $4, updated_by = $5 WHERE id = $6
`

type UpdateExpensePaymentParams struct {
	UserID      pgtype.UUID
	Amount      int32
	PaymentDate pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
	UpdatedBy   pgtype.UUID
	ID          pgtype.UUID
}

func (q *Queries) UpdateExpensePayment(ctx context.Context, arg UpdateExpensePaymentParams) error {
	_, err := q.db.Exec(ctx, updateExpensePayment,
		arg.UserID,
		arg.Amount,
		arg.PaymentDate,
		arg.UpdatedAt,
		arg.UpdatedBy,
		arg.ID,
	)
	return err
}

const updateLedgerParticipantBalance = `-- name: UpdateLedgerParticipantBalance :exec
UPDATE ledger_participant_balances SET last_timestamp = $1, balance = $2 WHERE ledger_id = $3 AND user_id = $4
`

type UpdateLedgerParticipantBalanceParams struct {
	LastTimestamp pgtype.Timestamp
	Balance       int32
	LedgerID      pgtype.UUID
	UserID        pgtype.UUID
}

func (q *Queries) UpdateLedgerParticipantBalance(ctx context.Context, arg UpdateLedgerParticipantBalanceParams) error {
	_, err := q.db.Exec(ctx, updateLedgerParticipantBalance,
		arg.LastTimestamp,
		arg.Balance,
		arg.LedgerID,
		arg.UserID,
	)
	return err
}
