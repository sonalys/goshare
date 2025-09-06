package expenses_test

import (
	"context"
	"testing"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	v1 "github.com/sonalys/goshare/internal/application/v1"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/router/testutils"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
	"github.com/sonalys/goshare/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LedgerExpenseGet(t *testing.T) {
	t.Parallel()

	type testData struct {
		params  server.LedgerExpenseGetParams
		expense *domain.Expense
	}

	getTestData := func(t *testing.T) testData {
		user := testfixtures.User(t)
		ledger := testfixtures.Ledger(t, user)
		expense := testfixtures.Expense(t, ledger, user.ID, domain.NewID())

		return testData{
			params: server.LedgerExpenseGetParams{
				LedgerID:  ledger.ID.UUID(),
				ExpenseID: expense.ID.UUID(),
			},
			expense: expense,
		}
	}

	assertCreateFunc := func(t *testing.T, identity *v1.Identity, td testData) {

	}

	t.Run("pass", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		identity := testfixtures.Identity(t)
		router, mocks := testutils.Setup(t, testutils.WithIdentity(identity))

		td := getTestData(t)

		mocks.ExpenseController.GetFunc = func(ctx context.Context, req usercontroller.GetExpenseRequest) (*domain.Expense, error) {
			assertCreateFunc(t, identity, td)

			return td.expense, nil
		}

		resp, err := router.LedgerExpenseGet(ctx, td.params)
		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, td.expense.ID.UUID(), resp.ID.Value)
		assert.Equal(t, td.expense.Name, resp.Name)
		assert.Equal(t, td.expense.ExpenseDate, resp.ExpenseDate)
		assert.Len(t, resp.Records, len(td.expense.Records))

		for _, gotRecord := range resp.Records {
			expectedRecord, ok := td.expense.Records[domain.ConvertID(gotRecord.ID.Value)]
			require.True(t, ok)

			assert.EqualValues(t, expectedRecord.Type.String(), gotRecord.Type)
			assert.Equal(t, expectedRecord.From.UUID(), gotRecord.FromUserID)
			assert.Equal(t, expectedRecord.To.UUID(), gotRecord.ToUserID)
			assert.Equal(t, expectedRecord.Amount, gotRecord.Amount)
		}
	})

	t.Run("fail/unauthenticated", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		router, _ := testutils.Setup(t, testutils.WithIdentity(nil))

		td := getTestData(t)

		resp, err := router.LedgerExpenseGet(ctx, td.params)
		require.ErrorIs(t, err, assert.AnError)
		require.Nil(t, resp)
	})

	t.Run("fail/controller error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		identity := testfixtures.Identity(t)
		router, mocks := testutils.Setup(t, testutils.WithIdentity(identity))

		td := getTestData(t)

		mocks.ExpenseController.GetFunc = func(ctx context.Context, req usercontroller.GetExpenseRequest) (*domain.Expense, error) {
			assertCreateFunc(t, identity, td)

			return nil, assert.AnError
		}

		resp, err := router.LedgerExpenseGet(ctx, td.params)
		require.ErrorIs(t, err, assert.AnError)
		require.Nil(t, resp)
	})
}
