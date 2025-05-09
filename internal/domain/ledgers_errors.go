package domain

import "fmt"

type (
	ErrLedgerUserAlreadyMember struct {
		UserID   ID
		LedgerID ID
	}

	ErrLedgerUserNotMember struct {
		UserID   ID
		LedgerID ID
	}

	ErrLedgerMaxMembers struct {
		LedgerID   ID
		MaxMembers int
	}
)

func (e *ErrLedgerUserAlreadyMember) Error() string {
	return fmt.Sprintf("user '%s' is already a member of the ledger '%s'", e.UserID, e.LedgerID)
}

func (e *ErrLedgerUserNotMember) Error() string {
	return fmt.Sprintf("user '%s' is not a member of the ledger '%s'", e.UserID, e.LedgerID)
}

func (e *ErrLedgerMaxMembers) Error() string {
	return fmt.Sprintf("ledger '%s' has reached maximum number of members: %d", e.LedgerID, e.MaxMembers)
}
