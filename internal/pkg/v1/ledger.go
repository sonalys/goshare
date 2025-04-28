package v1

import (
	"fmt"
	"time"

	"github.com/sonalys/kset"
)

type (
	Ledger struct {
		ID           ID
		Name         string
		Participants []LedgerParticipant
		CreatedAt    time.Time
		CreatedBy    ID
	}

	LedgerParticipant struct {
		ID        ID
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

func (l Ledger) IsParticipant(userID ID) bool {
	for _, p := range l.Participants {
		if p.UserID == userID {
			return true
		}
	}
	return false
}

func (l *Ledger) AddParticipants(userID ID, ids ...ID) int {
	participantsSet := kset.NewKeyValue(func(p LedgerParticipant) ID { return p.UserID }, l.Participants...)
	addedCount := 0

	for _, id := range ids {
		addedCount += participantsSet.Append(LedgerParticipant{
			ID:        NewID(),
			UserID:    id,
			Balance:   0,
			CreatedAt: time.Now(),
			CreatedBy: userID,
		})
	}

	l.Participants = participantsSet.ToSlice()

	return addedCount
}
