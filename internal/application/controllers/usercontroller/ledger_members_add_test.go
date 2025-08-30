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

func Test_Ledger_MembersAdd(t *testing.T) {
	type testSetup struct {
		db *databaseMock
	}

	createTestData := func() usercontroller.AddMembersRequest {
		return usercontroller.AddMembersRequest{
			ActorID:  domain.NewID(),
			LedgerID: domain.NewID(),
			Emails: []string{
				"email1",
				"email2",
			},
		}
	}

	setup := func(t *testing.T, td usercontroller.AddMembersRequest) (*usercontroller.Controller, testSetup) {
		mocks := testSetup{
			db: setupDatabaseMock(t),
		}

		user := testfixtures.User(t)
		user.ID = td.ActorID

		ledger := testfixtures.Ledger(t, user)
		ledger.ID = td.LedgerID

		mocks.db.tx.ledger.GetFunc = func(ctx context.Context, id domain.ID) (*domain.Ledger, error) {
			assert.Equal(t, td.LedgerID, id)
			return ledger, nil
		}

		members := []domain.User{
			*testfixtures.User(t),
			*testfixtures.User(t),
		}

		mocks.db.tx.user.ListByEmailFunc = func(ctx context.Context, emails []string) ([]domain.User, error) {
			assert.Equal(t, td.Emails, emails)
			return members, nil
		}

		previousLen := len(ledger.Members)

		mocks.db.tx.ledger.UpdateFunc = func(ctx context.Context, ledger *domain.Ledger) error {
			assert.Len(t, ledger.Members, previousLen+2)
			return nil
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

		err := controller.Ledgers().MembersAdd(ctx, td)
		require.NoError(t, err)
	})

	t.Run("fail/user is not authorized to manage members", func(t *testing.T) {
		ctx := t.Context()

		td := createTestData()
		controller, _ := setup(t, td)
		td.ActorID = domain.NewID()

		err := controller.Ledgers().MembersAdd(ctx, td)
		require.ErrorIs(t, err, v1.ErrForbidden)
	})

	t.Run("fail/ledger repository get error", func(t *testing.T) {
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.ledger.GetFunc = func(ctx context.Context, id domain.ID) (*domain.Ledger, error) {
			return nil, assert.AnError
		}

		err := controller.Ledgers().MembersAdd(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fail/user repository list error", func(t *testing.T) {
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.user.ListByEmailFunc = func(ctx context.Context, emails []string) ([]domain.User, error) {
			return nil, assert.AnError
		}

		err := controller.Ledgers().MembersAdd(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fail/ledger repository update error", func(t *testing.T) {
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.ledger.UpdateFunc = func(ctx context.Context, ledger *domain.Ledger) error {
			return assert.AnError
		}

		err := controller.Ledgers().MembersAdd(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
	})
}
