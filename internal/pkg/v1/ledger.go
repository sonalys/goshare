package v1

import (
	"errors"
	"time"
)

type (
	Ledger struct {
		ID        ID
		Name      string
		CreatedAt time.Time
		CreatedBy ID
	}

	LedgerParticipant struct {
		ID        ID
		LedgerID  ID
		UserID    ID
		CreatedAt time.Time
		CreatedBy ID
	}

	LedgerRecord struct {
		ID          ID
		LedgerID    ID
		ExpenseID   ID
		UserID      ID
		Amount      int32
		Description string
		CreatedAt   time.Time
		CreatedBy   ID
	}

	LedgerParticipantBalance struct {
		LedgerID ID
		UserID   ID
		Balance  int32
	}
)

var (
	ErrUserAlreadyMember = errors.New("user is already a member")
)
