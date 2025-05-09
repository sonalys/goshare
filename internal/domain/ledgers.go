package domain

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/sonalys/kset"
)

const (
	ExpenseMaxRecords = 100
	LedgerMaxMembers  = 100
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

func (req *CreateExpenseRequest) validate() error {
	var errs FormError

	if req.Name == "" {
		errs.Append(newRequiredFieldError("name"))
	}

	if req.ExpenseDate.IsZero() {
		errs.Append(newRequiredFieldError("expenseDate"))
	}

	if recordsLen := len(req.PendingRecords); recordsLen < 1 || recordsLen > ExpenseMaxRecords {
		errs.Append(newFieldLengthError("records", 1, ExpenseMaxRecords))
	}

	var totalAmount int32
	for i := range req.PendingRecords {
		var recordErrs FormError

		record := &req.PendingRecords[i]

		if !record.Type.IsValid() {
			recordErrs = append(recordErrs, newInvalidFieldError("type"))
		}

		if record.Amount <= 0 {
			recordErrs = append(recordErrs, newFieldLengthError("amount", 1, math.MaxInt32))
		}

		if record.From == record.To {
			recordErrs = append(recordErrs, FieldError{
				Field: "to",
				Cause: errors.New("from and to cannot be the same user"),
			})
		}

		if err := recordErrs.Close(); err != nil {
			errs.Append(FieldError{
				Field: fmt.Sprintf("records[%d]", i),
				Cause: err,
			})
			continue
		}

		if record.Type == RecordTypeDebt {
			newAmount := totalAmount + record.Amount
			if newAmount < totalAmount {
				errs.Append(FieldError{
					Cause: errors.New("expense amount overflowed"),
					Field: "amount",
				})
				break
			}
			totalAmount = newAmount
		}
	}

	return errs.Close()
}

func (ledger *Ledger) CreateExpense(req CreateExpenseRequest) (*Expense, error) {
	if !ledger.IsParticipant(req.Actor) {
		return nil, &ErrLedgerUserNotMember{
			UserID:   req.Actor,
			LedgerID: ledger.ID,
		}
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

	return &Expense{
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
	}, nil
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
		return &ErrLedgerUserNotMember{
			UserID:   actor,
			LedgerID: ledger.ID,
		}
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
		return &ErrLedgerMaxMembers{
			LedgerID:   ledger.ID,
			MaxMembers: LedgerMaxMembers,
		}
	}

	return nil
}
