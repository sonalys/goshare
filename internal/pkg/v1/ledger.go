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
	ID            uuid.UUID
	LedgerID      uuid.UUID
	ParticipantID uuid.UUID
	CreatedAt     string
	CreatedBy     uuid.UUID
}

type LedgerRecord struct {
	ID            uuid.UUID
	LedgerID      uuid.UUID
	ExpenseID     uuid.UUID
	ParticipantID uuid.UUID
	Amount        int
	CreatedAt     string
	CreatedBy     uuid.UUID
}

type LedgerSnapshot struct {
	ID            uuid.UUID
	LedgerID      uuid.UUID
	LastTimestamp time.Time
}

type LedgerSnapshotParticipantBalance struct {
	ID               uuid.UUID
	LedgerID         uuid.UUID
	LedgerSnapshotID uuid.UUID
	ParticipantID    uuid.UUID
	Balance          int
}
