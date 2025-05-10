package domain_test

import (
	"math"
	"testing"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpense_TotalDebt(t *testing.T) {
	t.Parallel()

	t.Run("pass/no records", func(t *testing.T) {
		t.Parallel()
		expense := domain.Expense{}

		debt := expense.TotalDebt()
		assert.Zero(t, debt)
	})

	t.Run("pass/with debts", func(t *testing.T) {
		t.Parallel()
		expense := domain.Expense{
			Records: map[domain.ID]*domain.Record{
				domain.NewID(): {Type: domain.RecordTypeDebt, Amount: 1},
			},
		}

		debt := expense.TotalDebt()
		assert.EqualValues(t, 1, debt)
	})
}

func TestExpense_TotalSettled(t *testing.T) {
	t.Parallel()

	t.Run("pass/no records", func(t *testing.T) {
		t.Parallel()
		expense := domain.Expense{}

		settled := expense.TotalSettled()
		assert.Zero(t, settled)
	})

	t.Run("pass/with settlement", func(t *testing.T) {
		t.Parallel()
		expense := domain.Expense{
			Records: map[domain.ID]*domain.Record{
				domain.NewID(): {Type: domain.RecordTypeSettlement, Amount: 1},
			},
		}

		settled := expense.TotalSettled()
		assert.EqualValues(t, 1, settled)
	})
}

func TestExpense_CreateRecords(t *testing.T) {
	t.Parallel()

	t.Run("pass/with records", func(t *testing.T) {
		t.Parallel()

		actorID := domain.NewID()
		memberID := domain.NewID()

		ledger := &domain.Ledger{
			ID: domain.NewID(),
			Members: map[domain.ID]*domain.LedgerMember{
				actorID:  {},
				memberID: {},
			},
		}

		expense := domain.Expense{
			LedgerID: ledger.ID,
			Records:  make(map[domain.ID]*domain.Record),
		}

		debt := domain.PendingRecord{
			Type:   domain.RecordTypeDebt,
			From:   actorID,
			To:     memberID,
			Amount: 32,
		}

		settlement := domain.PendingRecord{
			Type:   domain.RecordTypeSettlement,
			From:   actorID,
			To:     memberID,
			Amount: 10,
		}

		err := expense.CreateRecords(actorID, ledger, debt, settlement)
		require.NoError(t, err)

		assert.EqualValues(t, -debt.Amount+settlement.Amount, ledger.Members[actorID].Balance)
		assert.EqualValues(t, debt.Amount-settlement.Amount, ledger.Members[memberID].Balance)
		assert.EqualValues(t, debt.Amount, expense.Amount)

		assert.Len(t, expense.Records, 2)

		for _, record := range expense.Records {
			assert.NotZero(t, record.Type)
			assert.NotZero(t, record.CreatedAt)
			assert.Equal(t, actorID, record.CreatedBy)
			assert.NotZero(t, record.UpdatedAt)
			assert.Equal(t, actorID, record.UpdatedBy)
		}
	})

	t.Run("pass/no records", func(t *testing.T) {
		t.Parallel()

		actorID := domain.NewID()

		ledger := &domain.Ledger{
			ID: domain.NewID(),
			Members: map[domain.ID]*domain.LedgerMember{
				actorID: {},
			},
		}

		expense := domain.Expense{
			LedgerID: ledger.ID,
			Records:  make(map[domain.ID]*domain.Record),
		}

		err := expense.CreateRecords(actorID, ledger)
		require.NoError(t, err)
	})

	t.Run("fail/actor not a member", func(t *testing.T) {
		t.Parallel()

		actorID := domain.NewID()
		fromID := domain.NewID()
		toID := domain.NewID()

		ledger := &domain.Ledger{
			ID: domain.NewID(),
			Members: map[domain.ID]*domain.LedgerMember{
				fromID: {},
				toID:   {},
			},
		}

		expense := domain.Expense{
			LedgerID: ledger.ID,
			Records:  make(map[domain.ID]*domain.Record),
		}

		debt := domain.PendingRecord{
			Type:   domain.RecordTypeDebt,
			From:   fromID,
			To:     toID,
			Amount: 32,
		}

		err := expense.CreateRecords(actorID, ledger, debt)
		require.Error(t, err)

		var targetErr domain.FieldError
		require.ErrorAs(t, err, &targetErr)

		assert.Equal(t, "actor", targetErr.Field)
		assert.Equal(t, &domain.ErrLedgerUserNotMember{
			UserID:   actorID,
			LedgerID: ledger.ID,
		}, targetErr.Cause)
	})

	t.Run("fail/ledger mismatch", func(t *testing.T) {
		t.Parallel()

		actorID := domain.NewID()
		memberID := domain.NewID()

		ledger := &domain.Ledger{
			ID: domain.NewID(),
			Members: map[domain.ID]*domain.LedgerMember{
				actorID:  {},
				memberID: {},
			},
		}

		expense := domain.Expense{
			LedgerID: domain.NewID(),
			Records:  make(map[domain.ID]*domain.Record),
		}

		debt := domain.PendingRecord{
			Type:   domain.RecordTypeDebt,
			From:   actorID,
			To:     memberID,
			Amount: 32,
		}

		err := expense.CreateRecords(actorID, ledger, debt)
		require.Error(t, err)

		var targetErr domain.FieldError
		require.ErrorAs(t, err, &targetErr)

		assert.Equal(t, "ledger", targetErr.Field)
		assert.Equal(t, domain.ErrLedgerMismatch, targetErr.Cause)
	})

	t.Run("fail/max records", func(t *testing.T) {
		t.Parallel()

		actorID := domain.NewID()
		memberID := domain.NewID()

		ledger := &domain.Ledger{
			ID: domain.NewID(),
			Members: map[domain.ID]*domain.LedgerMember{
				actorID:  {},
				memberID: {},
			},
		}

		expense := domain.Expense{
			LedgerID: ledger.ID,
			Records:  make(map[domain.ID]*domain.Record, len(domain.ErrExpenseMaxRecords)),
		}

		for range domain.ExpenseMaxRecords {
			expense.Records[domain.NewID()] = &domain.Record{}
		}

		debt := domain.PendingRecord{
			Type:   domain.RecordTypeDebt,
			From:   actorID,
			To:     memberID,
			Amount: 32,
		}

		err := expense.CreateRecords(actorID, ledger, debt)
		require.ErrorIs(t, err, domain.FieldError{
			Field: "records",
			Cause: &domain.ValueLengthError{
				Min: 1,
				Max: 0,
			},
		})
	})

	t.Run("fail/invalid record type", func(t *testing.T) {
		t.Parallel()

		actorID := domain.NewID()
		memberID := domain.NewID()

		ledger := &domain.Ledger{
			ID: domain.NewID(),
			Members: map[domain.ID]*domain.LedgerMember{
				actorID:  {},
				memberID: {},
			},
		}

		expense := domain.Expense{
			LedgerID: ledger.ID,
			Records:  make(map[domain.ID]*domain.Record, 1),
		}

		debt := domain.PendingRecord{
			Type:   domain.RecordTypeUnknown,
			From:   actorID,
			To:     memberID,
			Amount: 32,
		}

		err := expense.CreateRecords(actorID, ledger, debt)
		require.ErrorIs(t, err, domain.FieldError{
			Field: "type",
			Cause: domain.ErrInvalid,
		})
	})

	t.Run("fail/invalid record amount", func(t *testing.T) {
		t.Parallel()

		actorID := domain.NewID()
		memberID := domain.NewID()

		ledger := &domain.Ledger{
			ID: domain.NewID(),
			Members: map[domain.ID]*domain.LedgerMember{
				actorID:  {},
				memberID: {},
			},
		}

		expense := domain.Expense{
			LedgerID: ledger.ID,
			Records:  make(map[domain.ID]*domain.Record, 1),
		}

		debt := domain.PendingRecord{
			Type:   domain.RecordTypeDebt,
			From:   actorID,
			To:     memberID,
			Amount: 0,
		}

		err := expense.CreateRecords(actorID, ledger, debt)
		require.ErrorIs(t, err, domain.FieldError{
			Field: "amount",
			Cause: &domain.ValueLengthError{
				Min: 1,
				Max: math.MaxInt32,
			},
		})
	})

	t.Run("fail/from and to are equal", func(t *testing.T) {
		t.Parallel()

		actorID := domain.NewID()

		ledger := &domain.Ledger{
			ID: domain.NewID(),
			Members: map[domain.ID]*domain.LedgerMember{
				actorID: {},
			},
		}

		expense := domain.Expense{
			LedgerID: ledger.ID,
			Records:  make(map[domain.ID]*domain.Record, 1),
		}

		debt := domain.PendingRecord{
			Type:   domain.RecordTypeDebt,
			From:   actorID,
			To:     actorID,
			Amount: 0,
		}

		err := expense.CreateRecords(actorID, ledger, debt)
		require.ErrorIs(t, err, domain.FieldError{
			Field: "to",
			Cause: domain.ErrLedgerFromToMatch,
		})
	})

	t.Run("fail/from is not a member", func(t *testing.T) {
		t.Parallel()

		actorID := domain.NewID()
		memberID := domain.NewID()

		ledger := &domain.Ledger{
			ID: domain.NewID(),
			Members: map[domain.ID]*domain.LedgerMember{
				actorID:  {},
				memberID: {},
			},
		}

		expense := domain.Expense{
			LedgerID: ledger.ID,
			Records:  make(map[domain.ID]*domain.Record, 1),
		}

		debt := domain.PendingRecord{
			Type:   domain.RecordTypeDebt,
			From:   domain.NewID(),
			To:     memberID,
			Amount: 0,
		}

		err := expense.CreateRecords(actorID, ledger, debt)
		require.ErrorIs(t, err, domain.FieldError{
			Field: "from",
			Cause: &domain.ErrLedgerUserNotMember{
				UserID:   debt.From,
				LedgerID: ledger.ID,
			},
		})
	})

	t.Run("fail/to is not a member", func(t *testing.T) {
		t.Parallel()

		actorID := domain.NewID()
		memberID := domain.NewID()

		ledger := &domain.Ledger{
			ID: domain.NewID(),
			Members: map[domain.ID]*domain.LedgerMember{
				actorID:  {},
				memberID: {},
			},
		}

		expense := domain.Expense{
			LedgerID: ledger.ID,
			Records:  make(map[domain.ID]*domain.Record, 1),
		}

		debt := domain.PendingRecord{
			Type:   domain.RecordTypeDebt,
			From:   memberID,
			To:     domain.NewID(),
			Amount: 0,
		}

		err := expense.CreateRecords(actorID, ledger, debt)
		require.ErrorIs(t, err, domain.FieldError{
			Field: "to",
			Cause: &domain.ErrLedgerUserNotMember{
				UserID:   debt.To,
				LedgerID: ledger.ID,
			},
		})
	})

	t.Run("fail/debt overflow", func(t *testing.T) {
		t.Parallel()

		actorID := domain.NewID()
		memberID := domain.NewID()

		ledger := &domain.Ledger{
			ID: domain.NewID(),
			Members: map[domain.ID]*domain.LedgerMember{
				actorID:  {},
				memberID: {},
			},
		}

		expense := domain.Expense{
			LedgerID: ledger.ID,
			Records: map[domain.ID]*domain.Record{
				domain.NewID(): {Type: domain.RecordTypeDebt, Amount: 1},
			},
		}

		debt := domain.PendingRecord{
			Type:   domain.RecordTypeDebt,
			From:   memberID,
			To:     domain.NewID(),
			Amount: math.MaxInt32,
		}

		err := expense.CreateRecords(actorID, ledger, debt)
		require.ErrorIs(t, err, domain.FieldError{
			Field: "amount",
			Cause: domain.ErrOverflow,
		})
	})

	t.Run("fail/settlement overflow", func(t *testing.T) {
		t.Parallel()

		actorID := domain.NewID()
		memberID := domain.NewID()

		ledger := &domain.Ledger{
			ID: domain.NewID(),
			Members: map[domain.ID]*domain.LedgerMember{
				actorID:  {},
				memberID: {},
			},
		}

		expense := domain.Expense{
			LedgerID: ledger.ID,
			Records: map[domain.ID]*domain.Record{
				domain.NewID(): {Type: domain.RecordTypeSettlement, Amount: 1},
			},
		}

		debt := domain.PendingRecord{
			Type:   domain.RecordTypeSettlement,
			From:   memberID,
			To:     domain.NewID(),
			Amount: math.MaxInt32,
		}

		err := expense.CreateRecords(actorID, ledger, debt)
		require.ErrorIs(t, err, domain.FieldError{
			Field: "amount",
			Cause: domain.ErrOverflow,
		})
	})

	t.Run("fail/settlement bigger than debt", func(t *testing.T) {
		t.Parallel()

		actorID := domain.NewID()
		memberID := domain.NewID()

		ledger := &domain.Ledger{
			ID: domain.NewID(),
			Members: map[domain.ID]*domain.LedgerMember{
				actorID:  {},
				memberID: {},
			},
		}

		expense := domain.Expense{
			LedgerID: ledger.ID,
			Records: map[domain.ID]*domain.Record{
				domain.NewID(): {Type: domain.RecordTypeDebt, Amount: 1},
			},
		}

		debt := domain.PendingRecord{
			Type:   domain.RecordTypeSettlement,
			From:   memberID,
			To:     domain.NewID(),
			Amount: 10,
		}

		err := expense.CreateRecords(actorID, ledger, debt)
		require.ErrorIs(t, err, domain.FieldError{
			Field: "pendingRecords",
			Cause: domain.ErrSettlementMismatch,
		})
	})
}
