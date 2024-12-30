package v1

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID          uuid.UUID
	CategoryID  uuid.UUID
	LedgerID    uuid.UUID
	Amount      int
	Name        string
	ExpenseDate time.Time
	CreatedAt   string
	CreatedBy   uuid.UUID
	UpdatedAt   string
	UpdatedBy   uuid.UUID
}

type ExpensePayment struct {
	ID        uuid.UUID
	ExpenseID uuid.UUID
	LedgerID  uuid.UUID
	PaidByID  uuid.UUID
	Amount    int
}
