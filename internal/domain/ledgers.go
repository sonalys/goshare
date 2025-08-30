package domain

import (
	"fmt"
	"slices"
	"time"
)

const (
	LedgerMaxMembers = 100

	ErrLedgerNotFound = ErrorString("ledger not found")
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
		Creator        ID
		Name           string
		ExpenseDate    time.Time
		PendingRecords []PendingRecord
	}
)

func (req *CreateExpenseRequest) validate() error {
	var form Form

	if req.Name == "" {
		form.Append(newRequiredFieldError("name"))
	}

	if req.ExpenseDate.IsZero() {
		form.Append(newRequiredFieldError("expenseDate"))
	}

	if recordsLen := len(req.PendingRecords); recordsLen < 1 || recordsLen > ExpenseMaxRecords {
		form.Append(newFieldLengthError("records", 1, ExpenseMaxRecords))
	}

	return form.Close()
}

func (ledger *Ledger) CreateExpense(req CreateExpenseRequest) (*Expense, error) {
	if !ledger.HasMember(req.Creator) {
		return nil, ErrLedgerUserNotMember{
			UserID:   req.Creator,
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
		CreatedBy:   req.Creator,
		UpdatedAt:   now,
		UpdatedBy:   req.Creator,
	}

	if err := expense.CreateRecords(req.Creator, ledger, req.PendingRecords...); err != nil {
		return nil, fmt.Errorf("creating expense records: %w", err)
	}

	return expense, nil
}

func (ledger *Ledger) HasMember(identity ID) bool {
	_, ok := ledger.Members[identity]
	return ok
}

func (ledger *Ledger) AddMember(inviter ID, newMembers ...ID) error {
	if !ledger.HasMember(inviter) {
		return ErrLedgerUserNotMember{
			UserID:   inviter,
			LedgerID: ledger.ID,
		}
	}

	newMembers = slices.Compact(newMembers)

	var form Form
	for _, id := range newMembers {
		if !ledger.HasMember(id) {
			continue
		}
		form.Append(FieldError{
			Field: "members",
			Cause: ErrLedgerUserAlreadyMember{
				UserID:   id,
				LedgerID: ledger.ID,
			},
		})
	}

	if err := form.Close(); err != nil {
		return err
	}

	if len(ledger.Members)+len(newMembers) >= LedgerMaxMembers {
		return ErrLedgerMaxMembers{
			LedgerID:   ledger.ID,
			MaxMembers: LedgerMaxMembers,
		}
	}

	for _, id := range newMembers {
		ledger.Members[id] = &LedgerMember{
			Balance:   0,
			CreatedAt: time.Now(),
			CreatedBy: inviter,
		}
	}

	return nil
}

func (ledger *Ledger) CanView(actor ID) bool {
	return ledger.CreatedBy == actor || ledger.HasMember(actor)
}

func (ledger *Ledger) CanManageMembers(actor ID) bool {
	return ledger.CreatedBy == actor || ledger.HasMember(actor)
}

func (ledger *Ledger) CanManageExpenses(actor ID) bool {
	return ledger.CreatedBy == actor || ledger.HasMember(actor)
}
