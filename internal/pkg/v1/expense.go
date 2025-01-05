package v1

import (
	"time"
)

type ExpenseUserBalance struct {
	UserID  ID
	Balance int32
}

type Expense struct {
	ID           ID
	CategoryID   *ID
	LedgerID     ID
	Amount       int32
	Name         string
	ExpenseDate  time.Time
	UserBalances []ExpenseUserBalance

	CreatedAt time.Time
	CreatedBy ID
	UpdatedAt time.Time
	UpdatedBy ID
}

type ExpensePayment struct {
	ID          ID
	ExpenseID   ID
	LedgerID    ID
	PaidByID    ID
	Amount      int32
	PaymentDate time.Time

	CreatedAt time.Time
	CreatedBy ID
	UpdatedAt time.Time
	UpdatedBy ID
}
