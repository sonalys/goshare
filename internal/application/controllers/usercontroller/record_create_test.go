package usercontroller_test

import (
	"context"
	"testing"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/utils/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Record_Create(t *testing.T) {
	type testSetup struct {
		db *databaseMock
	}

	createTestData := func() usercontroller.CreateExpenseRecordRequest {
		return usercontroller.CreateExpenseRecordRequest{
			ActorID:   domain.NewID(),
			LedgerID:  domain.NewID(),
			ExpenseID: domain.NewID(),
			PendingRecords: []domain.PendingRecord{
				{
					From:   domain.NewID(),
					To:     domain.NewID(),
					Type:   domain.RecordTypeDebt,
					Amount: 42,
				},
			},
		}
	}

	setup := func(t *testing.T, td usercontroller.CreateExpenseRecordRequest) (*usercontroller.Controller, testSetup) {
		mocks := testSetup{
			db: setupDatabaseMock(t),
		}

		user := testfixtures.User(t)
		user.ID = td.ActorID

		ledger := testfixtures.Ledger(t, user)
		ledger.ID = td.LedgerID

		for _, record := range td.PendingRecords {
			ledger.Members[record.From] = &domain.LedgerMember{}
			ledger.Members[record.To] = &domain.LedgerMember{}
		}

		expense := testfixtures.Expense(t, ledger)

		mocks.db.tx.ledger.GetFunc = func(ctx context.Context, id domain.ID) (*domain.Ledger, error) {
			assert.Equal(t, td.LedgerID, id)
			return ledger, nil
		}

		mocks.db.tx.expense.GetFunc = func(ctx context.Context, id domain.ID) (*domain.Expense, error) {
			assert.Equal(t, td.ExpenseID, id)
			return expense, nil
		}

		mocks.db.tx.ledger.UpdateFunc = func(ctx context.Context, ledger *domain.Ledger) error {
			return nil
		}

		mocks.db.tx.expense.UpdateFunc = func(ctx context.Context, ledger *domain.Expense) error {
			return nil
		}

		controller := usercontroller.New(usercontroller.Dependencies{
			Database: &mocks.db.db,
		})

		return controller, mocks
	}

	t.Run("pass", func(t *testing.T) {
		ctx := t.Context()

		td := createTestData()
		controller, _ := setup(t, td)

		expense, err := controller.Records().Create(ctx, td)
		require.NoError(t, err)
		require.NotNil(t, expense)
	})

	t.Run("fail/user unauthorized to manage expenses", func(t *testing.T) {
		ctx := t.Context()

		td := createTestData()
		controller, _ := setup(t, td)
		td.ActorID = domain.NewID()

		expense, err := controller.Records().Create(ctx, td)
		require.ErrorIs(t, err, v1.ErrForbidden)
		require.Nil(t, expense)
	})

	t.Run("fail/ledger repository get error", func(t *testing.T) {
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.ledger.GetFunc = func(ctx context.Context, id domain.ID) (*domain.Ledger, error) {
			return nil, assert.AnError
		}

		expense, err := controller.Records().Create(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
		require.Nil(t, expense)
	})

	t.Run("fail/expense repository get error", func(t *testing.T) {
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.expense.GetFunc = func(ctx context.Context, id domain.ID) (*domain.Expense, error) {
			return nil, assert.AnError
		}

		expense, err := controller.Records().Create(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
		require.Nil(t, expense)
	})

	t.Run("fail/ledger repository update error", func(t *testing.T) {
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.ledger.UpdateFunc = func(ctx context.Context, ledger *domain.Ledger) error {
			return assert.AnError
		}

		expense, err := controller.Records().Create(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
		require.Nil(t, expense)
	})

	t.Run("fail/expense repository update error", func(t *testing.T) {
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.expense.UpdateFunc = func(ctx context.Context, ledger *domain.Expense) error {
			return assert.AnError
		}

		expense, err := controller.Records().Create(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
		require.Nil(t, expense)
	})
}
