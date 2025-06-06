package domain

import (
	"fmt"
	"time"
)

const (
	LedgerMaxMembers = 100
)

type (
	Ledger struct {
		CreatedAt time.Time
		CreatedBy ID
		ID        ID
		Name      string
		Members   map[ID]*LedgerMember
	}

	LedgerMember struct {
		Balance   int32
		CreatedAt time.Time
		CreatedBy ID
	}

	CreateExpenseRequest struct {
		Actor          ID
		Name           string
		ExpenseDate    time.Time
		PendingRecords []PendingRecord
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

	return errs.Close()
}

func (ledger *Ledger) CreateExpense(req CreateExpenseRequest) (*Expense, error) {
	if !ledger.IsMember(req.Actor) {
		return nil, &ErrLedgerUserNotMember{
			UserID:   req.Actor,
			LedgerID: ledger.ID,
		}
	}

	if err := req.validate(); err != nil {
		return nil, err
	}

	now := time.Now()

	expense := &Expense{
		ID:          NewID(),
		LedgerID:    ledger.ID,
		Name:        req.Name,
		ExpenseDate: req.ExpenseDate,
		Records:     make(map[ID]*Record, len(req.PendingRecords)),
		Amount:      0,
		CreatedAt:   now,
		CreatedBy:   req.Actor,
		UpdatedAt:   now,
		UpdatedBy:   req.Actor,
	}

	if err := expense.CreateRecords(req.Actor, ledger, req.PendingRecords...); err != nil {
		return nil, fmt.Errorf("creating expense records: %w", err)
	}

	return expense, nil
}

func (ledger *Ledger) IsMember(identity ID) bool {
	_, ok := ledger.Members[identity]
	return ok
}

func (ledger *Ledger) AddMember(actor ID, newMembers ...ID) error {
	if !ledger.IsMember(actor) {
		return &ErrLedgerUserNotMember{
			UserID:   actor,
			LedgerID: ledger.ID,
		}
	}

	var errs FormError
	for _, id := range newMembers {
		if !ledger.IsMember(id) {
			continue
		}
		errs.Append(FieldError{
			Field: "members",
			Cause: &ErrLedgerUserAlreadyMember{
				UserID:   id,
				LedgerID: ledger.ID,
			},
		})
	}

	if err := errs.Close(); err != nil {
		return err
	}

	if len(ledger.Members)+len(newMembers) >= LedgerMaxMembers {
		return &ErrLedgerMaxMembers{
			LedgerID:   ledger.ID,
			MaxMembers: LedgerMaxMembers,
		}
	}

	for _, id := range newMembers {
		ledger.Members[id] = &LedgerMember{
			Balance:   0,
			CreatedAt: time.Now(),
			CreatedBy: actor,
		}
	}

	return nil
}
