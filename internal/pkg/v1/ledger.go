package v1

import (
	"fmt"
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
		Balance   int32
		CreatedAt time.Time
		CreatedBy ID
	}

	LedgerExpenseSummary struct {
		ID          ID
		Amount      int32
		Name        string
		ExpenseDate time.Time
		CreatedAt   time.Time
		CreatedBy   ID
		UpdatedAt   time.Time
		UpdatedBy   ID
	}
)

const (
	ErrUserAlreadyMember = StringError("user is already a member")
	ErrUserNotAMember    = StringError("user is not a member")

	LedgerMaxUsers = 100
	UserMaxLedgers = 5
)

var (
	ErrLedgerMaxUsers = fmt.Errorf("ledger reached maximum number of users: %d", LedgerMaxUsers)
	ErrUserMaxLedgers = fmt.Errorf("user reached the maximum number of ledgers: %d", UserMaxLedgers)
)
