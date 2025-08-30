package usercontroller_test

import (
	"context"
	"testing"
	"time"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/utils/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Expense_List(t *testing.T) {
	type testSetup struct {
		db *databaseMock
	}

	createTestData := func() usercontroller.ListExpensesRequest {
		return usercontroller.ListExpensesRequest{
			ActorID:  domain.NewID(),
			LedgerID: domain.NewID(),
			Cursor:   testfixtures.Now,
			Limit:    1,
		}
	}

	setup := func(t *testing.T, td usercontroller.ListExpensesRequest) (*usercontroller.Controller, testSetup) {
		mocks := testSetup{
			db: setupDatabaseMock(t),
		}

		user := testfixtures.User(t)
		user.ID = td.ActorID

		ledger := testfixtures.Ledger(t, user)
		ledger.ID = td.LedgerID

		mocks.db.repositories.ledger.GetFunc = func(ctx context.Context, id domain.ID) (*domain.Ledger, error) {
			assert.Equal(t, td.LedgerID, id)
			return ledger, nil
		}

		mocks.db.repositories.expense.ListByLedgerFunc = func(ctx context.Context, ledgerID domain.ID, cursor time.Time, limit int32) ([]v1.LedgerExpenseSummary, error) {
			return []v1.LedgerExpenseSummary{}, nil
		}

		controller := usercontroller.New(usercontroller.Dependencies{
			Database: &mocks.db.db,
		})

		return controller, mocks
	}

	t.Run("pass/empty results", func(t *testing.T) {
		ctx := t.Context()

		td := createTestData()
		controller, _ := setup(t, td)

		resp, err := controller.Expenses().List(ctx, td)
		require.NoError(t, err)
		require.Empty(t, resp)
	})

	t.Run("pass/pagination", func(t *testing.T) {
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.repositories.expense.ListByLedgerFunc = func(ctx context.Context, ledgerID domain.ID, cursor time.Time, limit int32) ([]v1.LedgerExpenseSummary, error) {
			return []v1.LedgerExpenseSummary{
				{},
				{},
			}, nil
		}

		resp, err := controller.Expenses().List(ctx, td)
		require.NoError(t, err)
		assert.NotNil(t, resp.Cursor)
		assert.Len(t, resp.Expenses, 1)
	})

	t.Run("fail/user is not authorized", func(t *testing.T) {
		ctx := t.Context()

		td := createTestData()
		controller, _ := setup(t, td)
		td.ActorID = domain.NewID()

		resp, err := controller.Expenses().List(ctx, td)
		require.ErrorIs(t, err, v1.ErrForbidden)
		require.Empty(t, resp)
	})

	t.Run("fail/ledger repository get errors", func(t *testing.T) {
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.repositories.ledger.GetFunc = func(ctx context.Context, id domain.ID) (*domain.Ledger, error) {
			return nil, assert.AnError
		}

		resp, err := controller.Expenses().List(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
		require.Empty(t, resp)
	})

	t.Run("fail/expense repository list errors", func(t *testing.T) {
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.repositories.expense.ListByLedgerFunc = func(ctx context.Context, ledgerID domain.ID, cursor time.Time, limit int32) ([]v1.LedgerExpenseSummary, error) {
			return nil, assert.AnError
		}

		resp, err := controller.Expenses().List(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
		require.Empty(t, resp)
	})
}
