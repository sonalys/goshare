package domain

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
		Identity  ID
		Balance   int32
		CreatedAt time.Time
		CreatedBy ID
	}

	RecordType int

	Record struct {
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

	Expense struct {
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
)

const (
	RecordTypeUnknown RecordType = iota
	RecordTypeDebt
	RecordTypeSettlement
	recordTypeMaxBoundary

	UserMaxLedgers    = 5
	LedgerMaxMembers  = 100
	ExpenseMaxRecords = 100

	ErrUserAlreadyMember = StringError("user is already a member")
	ErrUserNotAMember    = StringError("user is not a member")
)

var (
	ErrLedgerMaxUsers = fmt.Errorf("ledger reached maximum number of members: %d", LedgerMaxMembers)
	ErrUserMaxLedgers = fmt.Errorf("user reached the maximum number of ledgers: %d", UserMaxLedgers)
)

func NewLedger(identity ID, name string, userLedgersCount int64) (*Event[Ledger], error) {
	var errs FormError

	if userLedgersCount+1 > UserMaxLedgers {
		return nil, ErrUserMaxLedgers
	}

	if nameLength := len(name); nameLength < 3 || nameLength > 255 {
		errs = append(errs, NewFieldLengthError("name", 3, 255))
	}

	if err := errs.Validate(); err != nil {
		return nil, err
	}

	return &Event[Ledger]{
		Topic: TopicLedgerCreated,
		Data: Ledger{
			ID:   NewID(),
			Name: name,
			Participants: []LedgerParticipant{
				{
					ID:        NewID(),
					Identity:  identity,
					Balance:   0,
					CreatedAt: time.Now(),
					CreatedBy: identity,
				},
			},
			CreatedAt: time.Now(),
			CreatedBy: identity,
		},
	}, nil
}

func NewRecordType(s string) RecordType {
	switch s {
	case "debt":
		return RecordTypeDebt
	case "settlement":
		return RecordTypeSettlement
	default:
		return RecordTypeUnknown
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
	return r <= RecordTypeUnknown || r >= recordTypeMaxBoundary
}

func NewLedgerExpense(
	identity ID,
	ledgerID ID,
	name string,
	expenseDate time.Time,
	records []Record,
) (*Event[Expense], error) {
	var errs FormError

	if name == "" {
		errs = append(errs, NewRequiredFieldError("name"))
	}

	if expenseDate.IsZero() {
		errs = append(errs, NewRequiredFieldError("expenseDate"))
	}

	if recordsLen := len(records); recordsLen < 1 || recordsLen > ExpenseMaxRecords {
		errs = append(errs, NewFieldLengthError("records", 1, ExpenseMaxRecords))
	}

	var totalAmount int32

	for i := range records {
		var recordErrs FormError

		record := &records[i]

		if !record.Type.IsValid() {
			recordErrs = append(recordErrs, NewInvalidFieldError("type"))
		}

		if err := recordErrs.Validate(); err != nil {
			errs = append(errs, FieldError{
				Field: fmt.Sprintf("records[%d]", i),
				Cause: err,
			})
			continue
		}

		if record.Type == RecordTypeDebt {
			totalAmount += record.Amount
		}
	}

	if err := errs.Validate(); err != nil {
		return nil, err
	}

	return &Event[Expense]{
		Topic: TopicLedgerExpenseCreated,
		Data: Expense{
			ID:          NewID(),
			LedgerID:    ledgerID,
			Name:        name,
			ExpenseDate: expenseDate,
			Records:     records,
			Amount:      totalAmount,
			CreatedAt:   time.Now(),
			CreatedBy:   identity,
			UpdatedAt:   time.Now(),
			UpdatedBy:   identity,
		},
	}, nil
}

func (l Ledger) IsParticipant(identity ID) bool {
	for _, p := range l.Participants {
		if p.Identity == identity {
			return true
		}
	}
	return false
}

func (l *Ledger) AddParticipants(identity ID, participants ...ID) ([]Event[LedgerParticipant], error) {
	participantsSet := kset.HashMapKeyValue(func(p LedgerParticipant) ID { return p.Identity }, l.Participants...)
	pendingParticipantsSet := kset.HashMapKeyValue(func(p LedgerParticipant) ID { return p.Identity })

	for _, id := range participants {
		pendingParticipantsSet.Append(LedgerParticipant{
			ID:        NewID(),
			Identity:  id,
			Balance:   0,
			CreatedAt: time.Now(),
			CreatedBy: identity,
		})
	}

	newParticipants := pendingParticipantsSet.Difference(participantsSet).ToSlice()
	l.Participants = append(l.Participants, newParticipants...)

	if len(l.Participants) >= LedgerMaxMembers {
		return nil, ErrLedgerMaxUsers
	}

	events := make([]Event[LedgerParticipant], 0, len(newParticipants))
	for i := range newParticipants {
		events = append(events, Event[LedgerParticipant]{
			Topic: TopicLedgerParticipantAdded,
			Data:  newParticipants[i],
		})
	}

	return events, nil
}
