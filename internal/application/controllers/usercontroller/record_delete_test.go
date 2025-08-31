package usercontroller_test

import (
	"context"
	"testing"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/pkg/testfixtures"
	v1 "github.com/sonalys/goshare/pkg/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Record_Delete(t *testing.T) {
	t.Parallel()

	type testSetup struct {
		db *databaseMock
	}

	createTestData := func() usercontroller.DeleteExpenseRecordRequest {
		return usercontroller.DeleteExpenseRecordRequest{
			ActorID:   domain.NewID(),
			LedgerID:  domain.NewID(),
			ExpenseID: domain.NewID(),
			RecordID:  domain.NewID(),
		}
	}

	setup := func(t *testing.T, td usercontroller.DeleteExpenseRecordRequest) (usercontroller.Controller, testSetup) {
		mocks := testSetup{
			db: setupDatabaseMock(t),
		}

		user := testfixtures.User(t)
		user.ID = td.ActorID

		ledger := testfixtures.Ledger(t, user)
		ledger.ID = td.LedgerID

		expense := testfixtures.Expense(t, ledger, domain.NewID(), domain.NewID())
		expense.Records[td.RecordID] = &domain.Record{}

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
			LocalDatabase: &mocks.db.db,
		})

		return controller, mocks
	}

	t.Run("pass", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		td := createTestData()
		controller, _ := setup(t, td)

		err := controller.Records().Delete(ctx, td)
		require.NoError(t, err)
	})

	t.Run("fail/user unauthorized to manage expenses", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		td := createTestData()
		controller, _ := setup(t, td)
		td.ActorID = domain.NewID()

		err := controller.Records().Delete(ctx, td)
		require.ErrorIs(t, err, v1.ErrForbidden)
	})

	t.Run("fail/ledger repository get error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.ledger.GetFunc = func(ctx context.Context, id domain.ID) (*domain.Ledger, error) {
			return nil, assert.AnError
		}

		err := controller.Records().Delete(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fail/expense repository get error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.expense.GetFunc = func(ctx context.Context, id domain.ID) (*domain.Expense, error) {
			return nil, assert.AnError
		}

		err := controller.Records().Delete(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fail/ledger repository update error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.ledger.UpdateFunc = func(ctx context.Context, ledger *domain.Ledger) error {
			return assert.AnError
		}

		err := controller.Records().Delete(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fail/expense repository update error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.expense.UpdateFunc = func(ctx context.Context, ledger *domain.Expense) error {
			return assert.AnError
		}

		err := controller.Records().Delete(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
	})
}
