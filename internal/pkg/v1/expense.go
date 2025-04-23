package v1

import (
	"time"
)

type RecordType int

const (
	RecordTypeUnknown RecordType = iota
	RecordTypeExpense
	RecordTypeSettlement
	recordTypeMaxBoundary
)

func (r RecordType) String() string {
	switch r {
	case RecordTypeExpense:
		return "expense"
	case RecordTypeSettlement:
		return "settlement"
	default:
		return "unknown"
	}
}

func (r RecordType) IsValid() bool {
	return r <= RecordTypeUnknown || r >= recordTypeMaxBoundary
}

func NewRecordType(s string) RecordType {
	switch s {
	case "expense":
		return RecordTypeExpense
	case "settlement":
		return RecordTypeSettlement
	default:
		return RecordTypeUnknown
	}
}

type Record struct {
	ID        ID
	Type      RecordType
	Amount    int32
	From      ID
	To        ID
	CreatedAt time.Time
	CreatedBy ID
	UpdatedAt time.Time
	UpdatedBy ID
}

type Expense struct {
	ID          ID
	LedgerID    ID
	Amount      int32
	Name        string
	ExpenseDate time.Time
	Records     []Record

	CreatedAt time.Time
	CreatedBy ID
	UpdatedAt time.Time
	UpdatedBy ID
}
