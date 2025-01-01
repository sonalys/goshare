package v1

import (
	"time"

	"github.com/google/uuid"
)

type Ledger struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	CreatedBy uuid.UUID
}

type LedgerParticipant struct {
	ID        uuid.UUID
	LedgerID  uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	CreatedBy uuid.UUID
}

type LedgerRecord struct {
	ID          uuid.UUID
	LedgerID    uuid.UUID
	ExpenseID   uuid.UUID
	UserID      uuid.UUID
	Amount      int32
	Description string
	CreatedAt   time.Time
	CreatedBy   uuid.UUID
}

type LedgerParticipantBalance struct {
	ID            uuid.UUID
	LedgerID      uuid.UUID
	UserID        uuid.UUID
	LastTimestamp time.Time
	Balance       int32
}
