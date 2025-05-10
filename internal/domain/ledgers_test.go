package domain_test

import (
	"math"
	"testing"
	"time"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLedger_CreateExpense(t *testing.T) {
	t.Parallel()

	type testData struct {
		ledger *domain.Ledger
		req    domain.CreateExpenseRequest
	}

	factory := func(hooks ...func(*testData)) testData {
		actorID := domain.NewID()
		memberID := domain.NewID()
		td := testData{
			ledger: &domain.Ledger{
				ID: domain.NewID(),
				Members: map[domain.ID]*domain.LedgerMember{
					actorID:  {},
					memberID: {},
				},
			},
			req: domain.CreateExpenseRequest{
				Actor:       actorID,
				Name:        "expense",
				ExpenseDate: time.Now(),
				PendingRecords: []domain.PendingRecord{
					{
						From:   actorID,
						To:     memberID,
						Type:   domain.RecordTypeDebt,
						Amount: 32,
					},
				},
			},
		}

		for _, hook := range hooks {
			hook(&td)
		}

		return td
	}

	t.Run("pass", func(t *testing.T) {
		t.Parallel()
		data := factory()

		expense, err := data.ledger.CreateExpense(data.req)
		require.NoError(t, err)
		require.NotNil(t, expense)

		record := data.req.PendingRecords[0]

		assert.Equal(t, -record.Amount, data.ledger.Members[record.From].Balance)
		assert.Equal(t, record.Amount, data.ledger.Members[record.To].Balance)
	})

	t.Run("fail/actor is not a member", func(t *testing.T) {
		t.Parallel()
		data := factory(func(td *testData) {
			td.req.Actor = domain.NewID()
		})

		expense, err := data.ledger.CreateExpense(data.req)
		require.Error(t, err)
		assert.Nil(t, expense)

		var targetErr *domain.ErrLedgerUserNotMember
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, data.req.Actor, targetErr.UserID)
		assert.Equal(t, data.ledger.ID, targetErr.LedgerID)
	})

	t.Run("fail/from is not a member", func(t *testing.T) {
		t.Parallel()
		data := factory(func(td *testData) {
			td.req.PendingRecords[0].From = domain.NewID()
		})

		expense, err := data.ledger.CreateExpense(data.req)
		require.Error(t, err)
		assert.Nil(t, expense)

		var targetErr *domain.ErrLedgerUserNotMember
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, data.req.PendingRecords[0].From, targetErr.UserID)
		assert.Equal(t, data.ledger.ID, targetErr.LedgerID)
	})

	t.Run("fail/to is not a member", func(t *testing.T) {
		t.Parallel()
		data := factory(func(td *testData) {
			td.req.PendingRecords[0].To = domain.NewID()
		})

		expense, err := data.ledger.CreateExpense(data.req)
		require.Error(t, err)
		assert.Nil(t, expense)

		var targetErr *domain.ErrLedgerUserNotMember
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, data.req.PendingRecords[0].To, targetErr.UserID)
		assert.Equal(t, data.ledger.ID, targetErr.LedgerID)
	})

	t.Run("fail/from and to must be different", func(t *testing.T) {
		t.Parallel()
		data := factory(func(td *testData) {
			td.req.PendingRecords[0].To = td.req.PendingRecords[0].From
		})

		expense, err := data.ledger.CreateExpense(data.req)
		require.ErrorIs(t, err, domain.ErrLedgerFromToMatch)
		assert.Nil(t, expense)
	})

	t.Run("fail/expense type is invalid", func(t *testing.T) {
		t.Parallel()
		data := factory(func(td *testData) {
			td.req.PendingRecords[0].Type = domain.RecordTypeUnknown
		})

		expense, err := data.ledger.CreateExpense(data.req)
		require.ErrorIs(t, err, domain.ErrInvalid)
		assert.Nil(t, expense)
	})

	t.Run("fail/expense amount is negative", func(t *testing.T) {
		t.Parallel()
		data := factory(func(td *testData) {
			td.req.PendingRecords[0].Amount = -1
		})

		expense, err := data.ledger.CreateExpense(data.req)
		require.Error(t, err)
		assert.Nil(t, expense)

		var targetErr *domain.ValueLengthError
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, targetErr.Min, 1)
	})

	t.Run("fail/expense amount is zero", func(t *testing.T) {
		t.Parallel()
		data := factory(func(td *testData) {
			td.req.PendingRecords[0].Amount = 0
		})

		expense, err := data.ledger.CreateExpense(data.req)
		require.Error(t, err)
		assert.Nil(t, expense)

		var targetErr *domain.ValueLengthError
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, targetErr.Min, 1)
	})

	t.Run("fail/debt overflow", func(t *testing.T) {
		t.Parallel()
		data := factory(func(td *testData) {
			td.req.PendingRecords[0].Amount = math.MaxInt32
			td.req.PendingRecords = append(td.req.PendingRecords, domain.PendingRecord{
				From:   td.req.PendingRecords[0].From,
				To:     td.req.PendingRecords[0].To,
				Type:   domain.RecordTypeDebt,
				Amount: 1,
			})
		})

		expense, err := data.ledger.CreateExpense(data.req)
		require.ErrorIs(t, err, domain.ErrOverflow)
		assert.Nil(t, expense)
	})

	t.Run("fail/settlement overflow", func(t *testing.T) {
		t.Parallel()
		data := factory(func(td *testData) {
			td.req.PendingRecords[0].Amount = math.MaxInt32
			td.req.PendingRecords[0].Type = domain.RecordTypeSettlement
			td.req.PendingRecords = append(td.req.PendingRecords, domain.PendingRecord{
				From:   td.req.PendingRecords[0].From,
				To:     td.req.PendingRecords[0].To,
				Type:   domain.RecordTypeSettlement,
				Amount: 1,
			})
		})

		expense, err := data.ledger.CreateExpense(data.req)
		require.ErrorIs(t, err, domain.ErrOverflow)
		assert.Nil(t, expense)
	})

	t.Run("fail/name required", func(t *testing.T) {
		t.Parallel()
		data := factory(func(td *testData) {
			td.req.Name = ""
		})

		expense, err := data.ledger.CreateExpense(data.req)
		assert.Nil(t, expense)
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrRequired)
	})

	t.Run("fail/expenseDate required", func(t *testing.T) {
		t.Parallel()
		data := factory(func(td *testData) {
			td.req.ExpenseDate = time.Time{}
		})

		expense, err := data.ledger.CreateExpense(data.req)
		assert.Nil(t, expense)
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrRequired)
	})

	t.Run("fail/no records", func(t *testing.T) {
		t.Parallel()
		data := factory(func(td *testData) {
			td.req.PendingRecords = nil
		})

		expense, err := data.ledger.CreateExpense(data.req)
		assert.Nil(t, expense)
		require.Error(t, err)

		var targetErr *domain.ValueLengthError
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, targetErr.Min, 1)
		assert.Equal(t, targetErr.Max, domain.ExpenseMaxRecords)
	})

	t.Run("fail/too many records", func(t *testing.T) {
		t.Parallel()
		data := factory(func(td *testData) {
			td.req.PendingRecords = nil

			for range domain.ExpenseMaxRecords + 1 {
				td.req.PendingRecords = append(td.req.PendingRecords, domain.PendingRecord{})
			}
		})

		expense, err := data.ledger.CreateExpense(data.req)
		assert.Nil(t, expense)
		require.Error(t, err)

		var targetErr *domain.ValueLengthError
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, targetErr.Min, 1)
		assert.Equal(t, targetErr.Max, domain.ExpenseMaxRecords)
	})
}

func TestLedger_AddMember(t *testing.T) {
	t.Parallel()

	t.Run("pass", func(t *testing.T) {
		t.Parallel()
		actorID := domain.NewID()
		ledger := &domain.Ledger{
			Members: map[domain.ID]*domain.LedgerMember{
				actorID: {},
			},
		}

		err := ledger.AddMember(actorID, domain.NewID())
		require.NoError(t, err)
	})

	t.Run("fail/actor not a member", func(t *testing.T) {
		t.Parallel()
		actor := domain.NewID()
		ledger := &domain.Ledger{
			Members: nil,
		}

		err := ledger.AddMember(actor, domain.NewID())
		require.Error(t, err)

		var targetErr *domain.ErrLedgerUserNotMember
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, actor, targetErr.UserID)
		assert.Equal(t, ledger.ID, targetErr.LedgerID)
	})

	t.Run("fail/already a member", func(t *testing.T) {
		t.Parallel()
		actorID := domain.NewID()
		ledger := &domain.Ledger{
			ID: domain.NewID(),
			Members: map[domain.ID]*domain.LedgerMember{
				actorID: {},
			},
		}

		err := ledger.AddMember(actorID, actorID)
		require.Error(t, err)

		var targetErr *domain.ErrLedgerUserAlreadyMember
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, actorID, targetErr.UserID)
		assert.Equal(t, ledger.ID, targetErr.LedgerID)
	})

	t.Run("fail/maximum member capacity", func(t *testing.T) {
		t.Parallel()
		actorID := domain.NewID()
		ledger := &domain.Ledger{
			ID: domain.NewID(),
			Members: map[domain.ID]*domain.LedgerMember{
				actorID: {},
			},
		}

		for range domain.LedgerMaxMembers - 1 {
			ledger.Members[domain.NewID()] = &domain.LedgerMember{}
		}

		err := ledger.AddMember(actorID, domain.NewID())
		require.Error(t, err)

		var targetErr *domain.ErrLedgerMaxMembers
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, ledger.ID, targetErr.LedgerID)
		assert.Equal(t, domain.LedgerMaxMembers, targetErr.MaxMembers)
	})
}
