package v1

import (
	"time"

	"github.com/google/uuid"
)

type Ledger struct {
	ID   uuid.UUID
	Name string
}

type LedgerParticipant struct {
	ID        uuid.UUID
	LedgerID  uuid.UUID
	UserID    uuid.UUID
	CreatedAt string
	CreatedBy uuid.UUID
}

type LedgerRecord struct {
	ID        uuid.UUID
	LedgerID  uuid.UUID
	ExpenseID uuid.UUID
	UserID    uuid.UUID
	Amount    int
	CreatedAt string
	CreatedBy uuid.UUID
}

type LedgerParticipantBalance struct {
	ID            uuid.UUID
	LedgerID      uuid.UUID
	UserID        uuid.UUID
	LastTimestamp time.Time
	Balance       int
}
