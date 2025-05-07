package domain

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/sonalys/kset"
)

type (
	Ledger struct {
		CreatedAt    time.Time
		CreatedBy    ID
		ID           ID
		Name         string
		Participants []LedgerParticipant
	}

	LedgerParticipant struct {
		Balance   int32
		CreatedAt time.Time
		CreatedBy ID
		ID        ID
		Identity  ID
	}

	RecordType int

	Record struct {
		Amount    int32
		CreatedAt time.Time
		CreatedBy ID
		From      ID
		ID        ID
		To        ID
		Type      RecordType
		UpdatedAt time.Time
		UpdatedBy ID
	}

	Expense struct {
		Amount      int32
		CreatedAt   time.Time
		CreatedBy   ID
		ExpenseDate time.Time
		ID          ID
		LedgerID    ID
		Name        string
		Records     []Record
		UpdatedAt   time.Time
		UpdatedBy   ID
	}
)

const (
	RecordTypeUnknown RecordType = iota
	RecordTypeDebt
	RecordTypeSettlement
	recordTypeMaxBoundary

	ExpenseMaxRecords = 100
	LedgerMaxMembers  = 100
	UserMaxLedgers    = 5

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
	return r > RecordTypeUnknown && r < recordTypeMaxBoundary
}

type NewRecord struct {
	From   ID
	To     ID
	Type   RecordType
	Amount int32
}

func NewLedgerExpense(
	identity ID,
	ledgerID ID,
	name string,
	expenseDate time.Time,
	pendingRecords []NewRecord,
) (*Event[Expense], error) {
	var errs FormError

	if name == "" {
		errs = append(errs, NewRequiredFieldError("name"))
	}

	if expenseDate.IsZero() {
		errs = append(errs, NewRequiredFieldError("expenseDate"))
	}

	if recordsLen := len(pendingRecords); recordsLen < 1 || recordsLen > ExpenseMaxRecords {
		errs = append(errs, NewFieldLengthError("records", 1, ExpenseMaxRecords))
	}

	var totalAmount int32
	createdRecords := make([]Record, 0, len(pendingRecords))

	now := time.Now()

	for i := range pendingRecords {
		var recordErrs FormError

		record := &pendingRecords[i]

		if !record.Type.IsValid() {
			recordErrs = append(recordErrs, NewInvalidFieldError("type"))
		}

		if record.Amount <= 0 {
			recordErrs = append(recordErrs, NewFieldLengthError("amount", 1, math.MaxInt32))
		}

		if record.From == record.To {
			recordErrs = append(recordErrs, FieldError{
				Field: "to",
				Cause: errors.New("from and to cannot be the same user"),
			})
		}

		if err := recordErrs.Validate(); err != nil {
			errs = append(errs, FieldError{
				Field: fmt.Sprintf("records[%d]", i),
				Cause: err,
			})
			continue
		}

		if record.Type == RecordTypeDebt {
			newAmount := totalAmount + record.Amount
			if newAmount < totalAmount {
				errs = append(errs, FieldError{
					Cause: errors.New("expense amount overflowed"),
					Field: "amount",
				})
				break
			}
			totalAmount = newAmount
		}

		createdRecords = append(createdRecords, Record{
			ID:        NewID(),
			Amount:    record.Amount,
			Type:      record.Type,
			From:      record.From,
			To:        record.To,
			CreatedAt: now,
			CreatedBy: identity,
			UpdatedAt: now,
			UpdatedBy: identity,
		})
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
			Records:     createdRecords,
			Amount:      totalAmount,
			CreatedAt:   now,
			CreatedBy:   identity,
			UpdatedAt:   now,
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
