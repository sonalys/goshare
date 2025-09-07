package usercontroller_test

import (
	"context"
	"testing"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Ledger_Create(t *testing.T) {
	t.Parallel()
	type testSetup struct {
		db *databaseMock
	}

	createTestData := func() usercontroller.CreateLedgerRequest {
		return usercontroller.CreateLedgerRequest{
			ActorID: domain.NewID(),
			Name:    domain.NewID().String(),
		}
	}

	setup := func(t *testing.T, td usercontroller.CreateLedgerRequest) (usercontroller.Controller, testSetup) {
		mocks := testSetup{
			db: setupDatabaseMock(t),
		}

		user := testfixtures.User(t)
		user.ID = td.ActorID

		mocks.db.tx.user.GetFunc = func(ctx context.Context, id domain.ID) (*domain.User, error) {
			assert.Equal(t, td.ActorID, id)

			return user, nil
		}

		mocks.db.tx.ledger.CreateFunc = func(ctx context.Context, ledger *domain.Ledger) error {
			assert.Equal(t, td.Name, ledger.Name)

			return nil
		}

		previousLedgerCount := user.LedgersCount

		mocks.db.tx.user.CreateFunc = func(ctx context.Context, user *domain.User) error {
			assert.Equal(t, previousLedgerCount+1, user.LedgersCount)

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

		resp, err := controller.Ledgers().Create(ctx, td)
		require.NoError(t, err)
		assert.NotNil(t, resp.ID)
	})

	t.Run("fail/user repository find error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.user.GetFunc = func(ctx context.Context, id domain.ID) (*domain.User, error) {
			assert.Equal(t, td.ActorID, id)

			return nil, assert.AnError
		}

		resp, err := controller.Ledgers().Create(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, resp)
	})

	t.Run("fail/ledger repository create error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.ledger.CreateFunc = func(ctx context.Context, ledger *domain.Ledger) error {
			assert.Equal(t, td.Name, ledger.Name)

			return assert.AnError
		}

		resp, err := controller.Ledgers().Create(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, resp)
	})
}
