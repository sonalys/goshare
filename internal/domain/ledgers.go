package domain

import (
	"fmt"
	"math"
	"time"

	"github.com/sonalys/kset"
)

const (
	ExpenseMaxRecords = 100
	LedgerMaxMembers  = 100

	ErrLedgerFromToMatch  = ErrCause("from and to cannot be the same user")
	ErrOverflow           = ErrCause("overflow")
	ErrSettlementMismatch = ErrCause("settlement cannot be greater than debt")
)

type (
	Ledger struct {
		CreatedAt time.Time
		CreatedBy ID
		ID        ID
		Name      string
		Members   []LedgerMember
	}

	LedgerMember struct {
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

func (req *CreateExpenseRequest) validate(ledger *Ledger) error {
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

	var totalDebt int32
	var totalSettled int32

	for i := range req.PendingRecords {
		var recordErrs FormError

		record := &req.PendingRecords[i]

		if !record.Type.IsValid() {
			recordErrs.Append(newInvalidFieldError("type"))
		}

		if record.Amount <= 0 {
			recordErrs.Append(newFieldLengthError("amount", 1, math.MaxInt32))
		}

		if record.From == record.To {
			recordErrs.Append(FieldError{
				Field: "to",
				Cause: ErrLedgerFromToMatch,
			})
		}

		if !ledger.IsParticipant(record.From) {
			recordErrs.Append(FieldError{
				Field: "from",
				Cause: &ErrLedgerUserNotMember{
					UserID:   record.From,
					LedgerID: ledger.ID,
				},
			})
		}

		if !ledger.IsParticipant(record.To) {
			recordErrs.Append(FieldError{
				Field: "to",
				Cause: &ErrLedgerUserNotMember{
					UserID:   record.To,
					LedgerID: ledger.ID,
				},
			})
		}

		if record.Amount > 0 {
			switch record.Type {
			case RecordTypeDebt:
				newAmount := totalDebt + record.Amount
				if newAmount < totalDebt {
					recordErrs.Append(FieldError{
						Cause: ErrOverflow,
						Field: "amount",
					})
				}
				totalDebt = newAmount
			case RecordTypeSettlement:
				newAmount := totalSettled + record.Amount
				if newAmount < totalSettled {
					recordErrs.Append(FieldError{
						Cause: ErrOverflow,
						Field: "amount",
					})
				}
				totalSettled = newAmount
			}
		}

		if err := recordErrs.Close(); err != nil {
			errs.Append(FieldError{
				Field: fmt.Sprintf("records[%d]", i),
				Cause: err,
			})
			continue
		}
	}

	if totalSettled > totalDebt {
		errs.Append(FieldError{
			Field: "pendingRecords",
			Cause: ErrSettlementMismatch,
		})
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

	if err := req.validate(ledger); err != nil {
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
	for _, p := range ledger.Members {
		if p.Identity == identity {
			return true
		}
	}
	return false
}

func (ledger *Ledger) AddMember(actor ID, members ...ID) error {
	if !ledger.IsParticipant(actor) {
		return &ErrLedgerUserNotMember{
			UserID:   actor,
			LedgerID: ledger.ID,
		}
	}

	currentMembersSet := kset.HashMapKeyValue(func(p LedgerMember) ID { return p.Identity }, ledger.Members...)
	pendingMembersSet := kset.HashMapKey(members...)

	var errs FormError
	for conflictID := range pendingMembersSet.Intersect(currentMembersSet).Iter() {
		errs.Append(FieldError{
			Field: "participants",
			Cause: &ErrLedgerUserAlreadyMember{
				UserID:   conflictID,
				LedgerID: ledger.ID,
			},
		})
	}

	if err := errs.Close(); err != nil {
		return err
	}

	if len(ledger.Members)+pendingMembersSet.Len() >= LedgerMaxMembers {
		return &ErrLedgerMaxMembers{
			LedgerID:   ledger.ID,
			MaxMembers: LedgerMaxMembers,
		}
	}

	for id := range pendingMembersSet.Iter() {
		ledger.Members = append(ledger.Members, LedgerMember{
			ID:        NewID(),
			Identity:  id,
			Balance:   0,
			CreatedAt: time.Now(),
			CreatedBy: actor,
		})
	}

	return nil
}
