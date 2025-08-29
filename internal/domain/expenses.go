package domain

import (
	"fmt"
	"math"
	"time"

	"github.com/sonalys/goshare/internal/utils/genericmath"
)

type (
	Expense struct {
		Amount      int32
		CreatedAt   time.Time
		CreatedBy   ID
		ExpenseDate time.Time
		ID          ID
		LedgerID    ID
		Name        string
		Records     map[ID]*Record
		UpdatedAt   time.Time
		UpdatedBy   ID
	}
)

const (
	ExpenseMaxRecords = 100

	ErrExpenseMaxRecords  = Cause("expense reached maximum number of records")
	ErrSettlementMismatch = Cause("settlement cannot be greater than debt")
	ErrLedgerFromToMatch  = Cause("from and to cannot be the equal")
	ErrLedgerMismatch     = Cause("ledger mismatch")
)

func (e *Expense) sumRecords(t RecordType) int32 {
	var sum int32

	for i := range e.Records {
		if e.Records[i].Type == t {
			sum += e.Records[i].Amount
		}
	}

	return sum
}

func (e *Expense) TotalDebt() int32 {
	return e.sumRecords(RecordTypeDebt)
}

func (e *Expense) TotalSettled() int32 {
	return e.sumRecords(RecordTypeSettlement)
}

func (e *Expense) validateCreateRecords(actor ID, ledger *Ledger, records ...PendingRecord) error {
	if !ledger.IsMember(actor) {
		return FieldError{
			Field: "actor",
			Cause: ErrLedgerUserNotMember{
				UserID:   actor,
				LedgerID: ledger.ID,
			},
		}
	}

	if ledger.ID != e.LedgerID {
		return FieldError{
			Field: "ledger",
			Cause: ErrLedgerMismatch,
		}
	}

	var errs Form

	if recordsLen := len(records); len(e.Records)+recordsLen > ExpenseMaxRecords {
		errs.Append(
			FieldError{
				Cause: fmt.Errorf("%w: %w", ErrExpenseMaxRecords, RangeError{
					Min: 0,
					Max: genericmath.Max(ExpenseMaxRecords-(len(e.Records)+recordsLen), 0),
				}),
				Field: "records",
			},
		)
	}

	totalDebt := e.TotalDebt()
	totalSettled := e.TotalSettled()

	for i := range records {
		var recordsForm Form
		record := &records[i]

		if !record.Type.IsValid() {
			recordsForm.Append(newInvalidFieldError("type"))
		}

		if record.Amount <= 0 {
			recordsForm.Append(newFieldLengthError("amount", 1, math.MaxInt32))
		}

		if record.From == record.To {
			recordsForm.Append(FieldError{
				Field: "to",
				Cause: ErrLedgerFromToMatch,
			})
		}

		if !ledger.IsMember(record.From) {
			recordsForm.Append(FieldError{
				Field: "from",
				Cause: ErrLedgerUserNotMember{
					UserID:   record.From,
					LedgerID: ledger.ID,
				},
			})
		}

		if !ledger.IsMember(record.To) {
			recordsForm.Append(FieldError{
				Field: "to",
				Cause: ErrLedgerUserNotMember{
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
					recordsForm.Append(FieldError{
						Cause: CauseOverflow,
						Field: "amount",
					})
				}
				totalDebt = newAmount
			case RecordTypeSettlement:
				newAmount := totalSettled + record.Amount
				if newAmount < totalSettled {
					recordsForm.Append(FieldError{
						Cause: CauseOverflow,
						Field: "amount",
					})
				}
				totalSettled = newAmount
			}
		}

		if err := recordsForm.Close(); err != nil {
			errs.Append(FieldError{
				Field: fmt.Sprintf("records[%d]", i),
				Cause: err,
			})
			continue
		}
	}

	return errs.Close()
}

func (e *Expense) CreateRecords(actor ID, ledger *Ledger, records ...PendingRecord) error {
	if err := e.validateCreateRecords(actor, ledger, records...); err != nil {
		return err
	}

	now := time.Now()

	for i := range records {
		record := &records[i]
		e.Records[NewID()] = &Record{
			Type:      record.Type,
			Amount:    record.Amount,
			From:      record.From,
			To:        record.To,
			CreatedAt: now,
			CreatedBy: actor,
			UpdatedAt: now,
			UpdatedBy: actor,
		}

		switch record.Type {
		case RecordTypeDebt:
			ledger.Members[record.From].Balance -= record.Amount
			ledger.Members[record.To].Balance += record.Amount
		case RecordTypeSettlement:
			ledger.Members[record.From].Balance += record.Amount
			ledger.Members[record.To].Balance -= record.Amount
		}
	}

	e.Amount = e.TotalDebt()

	return nil
}

func (e *Expense) validateDeleteRecord(actor ID, ledger *Ledger, recordID ID) error {
	if !ledger.IsMember(actor) {
		return FieldError{
			Field: "actor",
			Cause: ErrLedgerUserNotMember{
				UserID:   actor,
				LedgerID: ledger.ID,
			},
		}
	}

	if ledger.ID != e.LedgerID {
		return FieldError{
			Field: "ledger",
			Cause: ErrLedgerMismatch,
		}
	}

	if _, ok := e.Records[recordID]; !ok {
		return FieldError{
			Field: "recordID",
			Cause: CauseNotFound,
		}
	}

	return nil
}

func (e *Expense) DeleteRecord(actor ID, ledger *Ledger, recordID ID) error {
	if err := e.validateDeleteRecord(actor, ledger, recordID); err != nil {
		return err
	}

	switch record := e.Records[recordID]; record.Type {
	case RecordTypeDebt:
		ledger.Members[record.From].Balance += record.Amount
		ledger.Members[record.To].Balance -= record.Amount
	case RecordTypeSettlement:
		ledger.Members[record.From].Balance -= record.Amount
		ledger.Members[record.To].Balance += record.Amount
	}

	delete(e.Records, recordID)

	return nil
}
