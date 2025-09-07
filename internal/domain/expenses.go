package domain

import (
	"fmt"
	"math"
	"time"

	"github.com/sonalys/goshare/internal/pkg/genericmath"
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

	ErrExpenseMaxRecords  = StringError("expense reached maximum number of records")
	ErrSettlementMismatch = StringError("settlement cannot be greater than debt")
	ErrLedgerFromToMatch  = StringError("from and to cannot be the equal")
	ErrLedgerMismatch     = StringError("ledger mismatch")
	ErrExpenseNotFound    = StringError("expense not found")
)

func (e *Expense) TotalDebt() int32 {
	return e.sumRecords(RecordTypeDebt)
}

func (e *Expense) TotalSettled() int32 {
	return e.sumRecords(RecordTypeSettlement)
}

func (e *Expense) CreateRecords(creator ID, ledger *Ledger, records ...PendingRecord) error {
	if err := e.validateCreateRecords(creator, ledger, records...); err != nil {
		return err
	}

	now := Now()

	for i := range records {
		record := &records[i]
		e.Records[NewID()] = &Record{
			Type:      record.Type,
			Amount:    record.Amount,
			From:      record.From,
			To:        record.To,
			CreatedAt: now,
			CreatedBy: creator,
			UpdatedAt: now,
			UpdatedBy: creator,
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

func (e *Expense) DeleteRecord(ledger *Ledger, recordID ID) error {
	if err := e.validateDeleteRecord(ledger, recordID); err != nil {
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

func (e *Expense) sumRecords(t RecordType) int32 {
	var sum int32

	for i := range e.Records {
		if e.Records[i].Type == t {
			sum += e.Records[i].Amount
		}
	}

	return sum
}

func (e *Expense) validateDeleteRecord(ledger *Ledger, recordID ID) error {
	if ledger.ID != e.LedgerID {
		return FieldError{
			Field: "ledger",
			Cause: ErrLedgerMismatch,
		}
	}

	if _, ok := e.Records[recordID]; !ok {
		return FieldError{
			Field: "recordID",
			Cause: ErrRecordNotFound,
		}
	}

	return nil
}

func (e *Expense) validateCreateRecords(creator ID, ledger *Ledger, records ...PendingRecord) error {
	if !ledger.HasMember(creator) {
		return FieldError{
			Field: "creator",
			Cause: LedgerUserNotMemberError{
				UserID:   creator,
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

	var form Form

	if recordsLen := len(records); len(e.Records)+recordsLen > ExpenseMaxRecords {
		form.Append(
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

		if !ledger.HasMember(record.From) {
			recordsForm.Append(FieldError{
				Field: "from",
				Cause: LedgerUserNotMemberError{
					UserID:   record.From,
					LedgerID: ledger.ID,
				},
			})
		}

		if !ledger.HasMember(record.To) {
			recordsForm.Append(FieldError{
				Field: "to",
				Cause: LedgerUserNotMemberError{
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
						Cause: ErrOverflow,
						Field: "amount",
					})
				}
				totalDebt = newAmount
			case RecordTypeSettlement:
				newAmount := totalSettled + record.Amount
				if newAmount < totalSettled {
					recordsForm.Append(FieldError{
						Cause: ErrOverflow,
						Field: "amount",
					})
				}
				totalSettled = newAmount
			}
		}

		if err := recordsForm.Close(); err != nil {
			form.Append(FieldError{
				Field: fmt.Sprintf("records[%d]", i),
				Cause: err,
			})
		}
	}

	return form.Close()
}
