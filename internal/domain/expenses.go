package domain

import (
	"fmt"
	"math"
	"time"
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

	ErrExpenseMaxRecords  = ErrCause("expense reached maximum number of records")
	ErrSettlementMismatch = ErrCause("settlement cannot be greater than debt")
	ErrLedgerFromToMatch  = ErrCause("from and to cannot be the equal")
	ErrLedgerMismatch     = ErrCause("ledger mismatch")
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

func (e *Expense) CreateRecords(actor ID, ledger *Ledger, records ...PendingRecord) error {
	if !ledger.IsMember(actor) {
		return FieldError{
			Field: "actor",
			Cause: &ErrLedgerUserNotMember{
				UserID:   actor,
				LedgerID: ledger.ID,
			},
		}
	}

	if len(records) == 0 {
		return nil
	}

	if ledger.ID != e.LedgerID {
		return FieldError{
			Field: "ledger",
			Cause: ErrLedgerMismatch,
		}
	}

	var errs FormError

	if len(e.Records)+len(records) > ExpenseMaxRecords {
		errs.Append(newFieldLengthError("records", 1, ExpenseMaxRecords-len(e.Records)))
	}

	totalDebt := e.TotalDebt()
	totalSettled := e.TotalSettled()

	for i := range records {
		var recordErrs FormError
		record := &records[i]

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

		if !ledger.IsMember(record.From) {
			recordErrs.Append(FieldError{
				Field: "from",
				Cause: &ErrLedgerUserNotMember{
					UserID:   record.From,
					LedgerID: ledger.ID,
				},
			})
		}

		if !ledger.IsMember(record.To) {
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

	if err := errs.Close(); err != nil {
		return err
	}

	e.Amount = totalDebt
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

	return nil
}

func (e *Expense) DeleteRecord(actor ID, ledger *Ledger, recordID ID) error {
	if !ledger.IsMember(actor) {
		return &ErrLedgerUserNotMember{
			UserID:   actor,
			LedgerID: ledger.ID,
		}
	}

	oldRecord, ok := e.Records[recordID]
	if !ok {
		return FieldError{
			Field: "recordID",
			Cause: ErrNotFound,
		}
	}

	switch oldRecord.Type {
	case RecordTypeDebt:
		ledger.Members[oldRecord.From].Balance += oldRecord.Amount
		ledger.Members[oldRecord.To].Balance -= oldRecord.Amount
	case RecordTypeSettlement:
		ledger.Members[oldRecord.From].Balance -= oldRecord.Amount
		ledger.Members[oldRecord.To].Balance += oldRecord.Amount
	}

	delete(e.Records, recordID)

	return nil
}
