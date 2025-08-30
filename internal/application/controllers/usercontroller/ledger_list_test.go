package usercontroller_test

import (
	"context"
	"testing"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Ledger_ListByUser(t *testing.T) {
	type testSetup struct {
		db *databaseMock
	}

	createTestData := func() domain.ID {
		return domain.NewID()
	}

	setup := func(t *testing.T, td domain.ID) (*usercontroller.Controller, testSetup) {
		mocks := testSetup{
			db: setupDatabaseMock(t),
		}

		mocks.db.repositories.ledger.ListByUserFunc = func(ctx context.Context, identity domain.ID) ([]domain.Ledger, error) {
			assert.Equal(t, td, identity)
			return []domain.Ledger{}, nil
		}

		controller := usercontroller.New(usercontroller.Dependencies{
			LocalDatabase: &mocks.db.db,
		})

		return controller, mocks
	}

	t.Run("pass", func(t *testing.T) {
		ctx := t.Context()
		td := createTestData()

		controller, _ := setup(t, td)

		resp, err := controller.Ledgers().ListByUser(ctx, td)
		require.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("fail/ledger repository error", func(t *testing.T) {
		ctx := t.Context()
		td := createTestData()

		controller, mocks := setup(t, td)

		mocks.db.repositories.ledger.ListByUserFunc = func(ctx context.Context, identity domain.ID) ([]domain.Ledger, error) {
			return nil, assert.AnError
		}

		resp, err := controller.Ledgers().ListByUser(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, resp)
	})
}
