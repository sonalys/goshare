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

		events []Event
	}

	LedgerParticipant struct {
		Balance   int32
		CreatedAt time.Time
		CreatedBy ID
		ID        ID
		Identity  ID
	}

	CreateExpenseRequest struct {
		Actor          ID
		Name           string
		ExpenseDate    time.Time
		PendingRecords []PendingRecord
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

func (l *Ledger) Events() []Event {
	return l.events
}

func (req *CreateExpenseRequest) validate() error {
	var errs FormError

	if req.Name == "" {
		errs = append(errs, NewRequiredFieldError("name"))
	}

	if req.ExpenseDate.IsZero() {
		errs = append(errs, NewRequiredFieldError("expenseDate"))
	}

	if recordsLen := len(req.PendingRecords); recordsLen < 1 || recordsLen > ExpenseMaxRecords {
		errs = append(errs, NewFieldLengthError("records", 1, ExpenseMaxRecords))
	}

	var totalAmount int32
	for i := range req.PendingRecords {
		var recordErrs FormError

		record := &req.PendingRecords[i]

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
	}

	return errs.Validate()
}

func (ledger *Ledger) CreateExpense(req CreateExpenseRequest) (*Expense, error) {
	if !ledger.IsParticipant(req.Actor) {
		return nil, ErrForbidden
	}

	if err := req.validate(); err != nil {
		return nil, err
	}

	var totalAmount int32
	now := time.Now()
	createdRecords := make([]Record, 0, len(req.PendingRecords))

	for i := range req.PendingRecords {
		record := &req.PendingRecords[i]

		if record.Type == RecordTypeDebt {
			totalAmount += record.Amount
		}

		createdRecords = append(createdRecords, Record{
			ID:        NewID(),
			Amount:    record.Amount,
			Type:      record.Type,
			From:      record.From,
			To:        record.To,
			CreatedAt: now,
			CreatedBy: req.Actor,
			UpdatedAt: now,
			UpdatedBy: req.Actor,
		})
	}

	expense := Expense{
		ID:          NewID(),
		LedgerID:    ledger.ID,
		Name:        req.Name,
		ExpenseDate: req.ExpenseDate,
		Records:     createdRecords,
		Amount:      totalAmount,
		CreatedAt:   now,
		CreatedBy:   req.Actor,
		UpdatedAt:   now,
		UpdatedBy:   req.Actor,
	}

	ledger.events = append(ledger.events, event[Expense]{
		topic: TopicLedgerExpenseCreated,
		data:  expense,
	})

	return &expense, nil
}

func (ledger *Ledger) IsParticipant(identity ID) bool {
	for _, p := range ledger.Participants {
		if p.Identity == identity {
			return true
		}
	}
	return false
}

func (ledger *Ledger) AddParticipants(actor ID, participants ...ID) error {
	if !ledger.IsParticipant(actor) {
		return ErrForbidden
	}

	participantsSet := kset.HashMapKeyValue(func(p LedgerParticipant) ID { return p.Identity }, ledger.Participants...)
	pendingParticipantsSet := kset.HashMapKeyValue(func(p LedgerParticipant) ID { return p.Identity })

	for _, id := range participants {
		pendingParticipantsSet.Append(LedgerParticipant{
			ID:        NewID(),
			Identity:  id,
			Balance:   0,
			CreatedAt: time.Now(),
			CreatedBy: actor,
		})
	}

	newParticipants := pendingParticipantsSet.Difference(participantsSet).ToSlice()
	ledger.Participants = append(ledger.Participants, newParticipants...)

	if len(ledger.Participants) >= LedgerMaxMembers {
		return ErrLedgerMaxUsers
	}

	events := make([]Event, 0, len(newParticipants))

	for i := range newParticipants {
		events = append(events, event[LedgerParticipant]{
			topic: TopicLedgerParticipantAdded,
			data:  newParticipants[i],
		})
	}

	ledger.events = append(ledger.events, events...)
	return nil
}
