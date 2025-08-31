package domain

import "fmt"

type (
	LedgerUserAlreadyMemberError struct {
		UserID   ID
		LedgerID ID
	}

	LedgerUserNotMemberError struct {
		UserID   ID
		LedgerID ID
	}

	LedgerMaxMembersError struct {
		LedgerID   ID
		MaxMembers int
	}
)

func (e LedgerUserAlreadyMemberError) Error() string {
	return fmt.Sprintf("user '%s' is already a member of the ledger '%s'", e.UserID, e.LedgerID)
}

func (e LedgerUserNotMemberError) Error() string {
	return fmt.Sprintf("user '%s' is not a member of the ledger '%s'", e.UserID, e.LedgerID)
}

func (e LedgerMaxMembersError) Error() string {
	return fmt.Sprintf("ledger '%s' has reached maximum number of members: %d", e.LedgerID, e.MaxMembers)
}
