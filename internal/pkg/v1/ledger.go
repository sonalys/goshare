package v1

import (
	"errors"
	"fmt"
	"time"
)

const (
	LedgerMaxUsers = 100
	UserMaxLedgers = 5
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
		Balance   int32
	}
)

var (
	ErrUserAlreadyMember = errors.New("user is already a member")
	ErrLedgerMaxUsers    = fmt.Errorf("ledger reached maximum number of users: %d", LedgerMaxUsers)
	ErrUserMaxLedgers    = fmt.Errorf("user reached the maximum number of ledgers: %d", UserMaxLedgers)
	ErrUserNotAMember    = errors.New("user is not a member")
)
