package v1

import (
	"time"

	"github.com/google/uuid"
)

type ExpenseUserBalance struct {
	UserID  uuid.UUID
	Balance int32
}

type Expense struct {
	ID           uuid.UUID
	CategoryID   *uuid.UUID
	LedgerID     uuid.UUID
	Amount       int32
	Name         string
	ExpenseDate  time.Time
	UserBalances []ExpenseUserBalance

	CreatedAt time.Time
	CreatedBy uuid.UUID
	UpdatedAt time.Time
	UpdatedBy uuid.UUID
}

type ExpensePayment struct {
	ID          uuid.UUID
	ExpenseID   uuid.UUID
	LedgerID    uuid.UUID
	PaidByID    uuid.UUID
	Amount      int32
	PaymentDate time.Time

	CreatedAt time.Time
	CreatedBy uuid.UUID
	UpdatedAt time.Time
	UpdatedBy uuid.UUID
}
