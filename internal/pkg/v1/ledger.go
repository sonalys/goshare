package v1

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type (
	Ledger struct {
		ID        uuid.UUID
		Name      string
		CreatedAt time.Time
		CreatedBy uuid.UUID
	}

	LedgerParticipant struct {
		ID        uuid.UUID
		LedgerID  uuid.UUID
		UserID    uuid.UUID
		CreatedAt time.Time
		CreatedBy uuid.UUID
	}

	LedgerRecord struct {
		ID          uuid.UUID
		LedgerID    uuid.UUID
		ExpenseID   uuid.UUID
		UserID      uuid.UUID
		Amount      int32
		Description string
		CreatedAt   time.Time
		CreatedBy   uuid.UUID
	}

	LedgerParticipantBalance struct {
		LedgerID uuid.UUID
		UserID   uuid.UUID
		Balance  int32
	}
)

var (
	ErrUserAlreadyMember = errors.New("user is already a member")
)
