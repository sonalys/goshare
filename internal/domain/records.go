package domain

import "time"

type (
	RecordType int

	Record struct {
		Amount    int32
		CreatedAt time.Time
		CreatedBy ID
		From      ID
		To        ID
		Type      RecordType
		UpdatedAt time.Time
		UpdatedBy ID
	}

	PendingRecord struct {
		From   ID
		To     ID
		Type   RecordType
		Amount int32
	}
)

const (
	RecordTypeUnknown RecordType = iota
	RecordTypeDebt
	RecordTypeSettlement
	recordTypeMaxBoundary
)

func NewRecordType(s string) (RecordType, error) {
	switch s {
	case "debt":
		return RecordTypeDebt, nil
	case "settlement":
		return RecordTypeSettlement, nil
	default:
		return RecordTypeUnknown, ErrInvalid
	}
}

func (r RecordType) String() string {
	switch r {
	case RecordTypeDebt:
		return "debt"
	case RecordTypeSettlement:
		return "settlement"
	default:
		return "unknown"
	}
}

func (r RecordType) IsValid() bool {
	return r > RecordTypeUnknown && r < recordTypeMaxBoundary
}
