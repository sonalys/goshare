package expenses_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
	"github.com/sonalys/goshare/internal/infrastructure/http/testutils"
	"github.com/sonalys/goshare/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LedgerExpenseCreate(t *testing.T) {
	t.Parallel()

	type testData struct {
		expense *server.Expense
		params  server.LedgerExpenseCreateParams
	}

	getTestData := func(_ *testing.T) testData {
		expense := &server.Expense{
			Name:        "expense",
			ExpenseDate: domain.Now(),
			Records: []server.ExpenseRecord{
				{
					Type:       server.ExpenseRecordTypeDebt,
					FromUserID: uuid.New(),
					ToUserID:   uuid.New(),
					Amount:     10,
				},
			},
		}

		params := server.LedgerExpenseCreateParams{
			LedgerID: uuid.New(),
		}

		return testData{
			expense: expense,
			params:  params,
		}
	}

	assertCreateFunc := func(t *testing.T, identity *application.Identity, td testData, got usercontroller.CreateExpenseRequest) {
		expense := td.expense
		params := td.params

		assert.Equal(t, identity.UserID, got.ActorID)
		assert.Equal(t, params.LedgerID, got.LedgerID.UUID())
		assert.Equal(t, expense.Name, got.Name)
		assert.Equal(t, expense.ExpenseDate, got.ExpenseDate)
		require.Len(t, got.PendingRecords, 1)

		gotRecord := got.PendingRecords[0]
		expectedRecord := expense.Records[0]

		assert.Equal(t, expectedRecord.Amount, gotRecord.Amount)
		assert.Equal(t, expectedRecord.FromUserID, gotRecord.From.UUID())
		assert.Equal(t, expectedRecord.ToUserID, gotRecord.To.UUID())

		gotType, err := domain.NewRecordType(string(expectedRecord.Type))
		require.NoError(t, err)
		assert.Equal(t, gotType, gotRecord.Type)
	}

	t.Run("pass", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		identity := testfixtures.Identity(t)
		router, mocks := testutils.Setup(t, testutils.WithIdentity(identity))

		td := getTestData(t)

		expense := td.expense
		params := td.params

		mocks.ExpenseController.CreateFunc = func(ctx context.Context, got usercontroller.CreateExpenseRequest) (*usercontroller.CreateExpenseResponse, error) {
			assertCreateFunc(t, identity, td, got)

			return &usercontroller.CreateExpenseResponse{
				ID: domain.ConvertID(td.params.LedgerID),
			}, nil
		}

		resp, err := router.LedgerExpenseCreate(ctx, expense, params)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, td.params.LedgerID, resp.ID)
	})

	t.Run("fail/controller error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		identity := testfixtures.Identity(t)
		router, mocks := testutils.Setup(t, testutils.WithIdentity(identity))

		td := getTestData(t)

		mocks.ExpenseController.CreateFunc = func(ctx context.Context, got usercontroller.CreateExpenseRequest) (*usercontroller.CreateExpenseResponse, error) {
			assertCreateFunc(t, identity, td, got)

			return nil, assert.AnError
		}

		_, err := router.LedgerExpenseCreate(ctx, td.expense, td.params)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fail/unauthenticated", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		router, _ := testutils.Setup(t, testutils.WithIdentity(nil))
		td := getTestData(t)

		_, err := router.LedgerExpenseCreate(ctx, td.expense, td.params)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fail/invalid request", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		identity := testfixtures.Identity(t)
		router, _ := testutils.Setup(t, testutils.WithIdentity(identity))
		td := getTestData(t)
		td.expense.Records[0].Type = "invalid"

		_, err := router.LedgerExpenseCreate(ctx, td.expense, td.params)
		require.ErrorAs(t, err, &domain.FieldError{})
	})
}
