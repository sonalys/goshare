package expenses_test

import (
	"context"
	"testing"
	"time"

	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
	"github.com/sonalys/goshare/internal/infrastructure/http/testutils"
	"github.com/sonalys/goshare/internal/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LedgerExpenseList(t *testing.T) {
	t.Parallel()

	type testData struct {
		params  server.LedgerExpenseListParams
		expense *domain.Expense
	}

	getTestData := func(t *testing.T) testData {
		user := testfixtures.User(t)
		ledger := testfixtures.Ledger(t, user)
		expense := testfixtures.Expense(t, ledger, user.ID, domain.NewID())

		return testData{
			params: server.LedgerExpenseListParams{
				LedgerID: ledger.ID.UUID(),
				Cursor:   server.NewOptDateTime(time.Now()),
				Limit:    server.NewOptInt32(1),
			},
			expense: expense,
		}
	}

	assertController := func(t *testing.T, identity *application.Identity, td testData, req usercontroller.ListExpensesRequest) {
		assert.Equal(t, identity.UserID, req.ActorID)
		assert.Equal(t, td.params.LedgerID, req.LedgerID.UUID())
		assert.Equal(t, td.params.Cursor.Value, req.Cursor)
		assert.Equal(t, td.params.Limit.Value, req.Limit)
	}

	t.Run("pass", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		identity := testfixtures.Identity(t)
		router, mocks := testutils.Setup(t, testutils.WithIdentity(identity))

		td := getTestData(t)

		mocks.ExpenseController.ListFunc = func(ctx context.Context, req usercontroller.ListExpensesRequest) (*usercontroller.ListExpensesResponse, error) {
			assertController(t, identity, td, req)

			return &usercontroller.ListExpensesResponse{
				Expenses: []application.LedgerExpenseSummary{
					{},
				},
				Cursor: &time.Time{},
			}, nil
		}

		resp, err := router.LedgerExpenseList(ctx, td.params)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("fail/unauthenticated", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		router, _ := testutils.Setup(t, testutils.WithIdentity(nil))

		td := getTestData(t)

		resp, err := router.LedgerExpenseList(ctx, td.params)
		require.ErrorIs(t, err, assert.AnError)
		require.Nil(t, resp)
	})

	t.Run("fail/controller error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		identity := testfixtures.Identity(t)
		router, mocks := testutils.Setup(t, testutils.WithIdentity(identity))

		td := getTestData(t)

		mocks.ExpenseController.ListFunc = func(ctx context.Context, req usercontroller.ListExpensesRequest) (*usercontroller.ListExpensesResponse, error) {
			assertController(t, identity, td, req)

			return nil, assert.AnError
		}

		resp, err := router.LedgerExpenseList(ctx, td.params)
		require.ErrorIs(t, err, assert.AnError)
		require.Nil(t, resp)
	})
}
