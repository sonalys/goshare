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
	type testData struct {
		ledger *domain.Ledger
		req    domain.CreateExpenseRequest
	}

	factory := func(hooks ...func(*testData)) testData {
		actorID := domain.NewID()
		memberID := domain.NewID()
		td := testData{
			ledger: &domain.Ledger{
				Participants: []domain.LedgerParticipant{
					{Identity: actorID},
					{Identity: memberID},
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
		data := factory()

		expense, err := data.ledger.CreateExpense(data.req)
		require.NoError(t, err)
		require.NotNil(t, expense)
	})

	t.Run("error/actor is not a participant", func(t *testing.T) {
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

	t.Run("error/from is not a participant", func(t *testing.T) {
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

	t.Run("error/to is not a participant", func(t *testing.T) {
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

	t.Run("error/from and to must be different", func(t *testing.T) {
		data := factory(func(td *testData) {
			td.req.PendingRecords[0].To = td.req.PendingRecords[0].From
		})

		expense, err := data.ledger.CreateExpense(data.req)
		require.ErrorIs(t, err, domain.ErrLedgerFromToMatch)
		assert.Nil(t, expense)
	})

	t.Run("error/expense type is invalid", func(t *testing.T) {
		data := factory(func(td *testData) {
			td.req.PendingRecords[0].Type = domain.RecordTypeUnknown
		})

		expense, err := data.ledger.CreateExpense(data.req)
		require.ErrorIs(t, err, domain.ErrCauseInvalid)
		assert.Nil(t, expense)
	})

	t.Run("error/expense amount is negative", func(t *testing.T) {
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

	t.Run("error/expense amount is zero", func(t *testing.T) {
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

	t.Run("error/debt overflow", func(t *testing.T) {
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

	t.Run("error/settlement greater than debt", func(t *testing.T) {
		data := factory(func(td *testData) {
			td.req.PendingRecords[0].Type = domain.RecordTypeSettlement
			td.req.PendingRecords[0].Amount = math.MaxInt32
		})

		expense, err := data.ledger.CreateExpense(data.req)
		require.ErrorIs(t, err, domain.ErrSettlementMismatch)
		assert.Nil(t, expense)
	})
}

func TestLedger_AddParticipants(t *testing.T) {
	t.Run("pass", func(t *testing.T) {
		actor := domain.NewID()
		ledger := &domain.Ledger{
			Participants: []domain.LedgerParticipant{
				{
					Identity: actor,
				},
			},
		}

		err := ledger.AddParticipants(actor, domain.NewID())
		require.NoError(t, err)
	})
}
