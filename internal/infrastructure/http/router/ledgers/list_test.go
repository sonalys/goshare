package ledgers_test

import (
	"context"
	"testing"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/testutils"
	"github.com/sonalys/goshare/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LedgerList(t *testing.T) {
	t.Parallel()

	t.Run("pass", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		identity := testfixtures.Identity(t)
		router, mocks := testutils.Setup(t, testutils.WithIdentity(identity))

		mocks.LedgerController.ListByUserFunc = func(ctx context.Context, actorID domain.ID) ([]domain.Ledger, error) {
			assert.Equal(t, identity.UserID, actorID)

			return []domain.Ledger{}, nil
		}

		resp, err := router.LedgerList(ctx)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("fail/controller error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		identity := testfixtures.Identity(t)
		router, mocks := testutils.Setup(t, testutils.WithIdentity(identity))

		mocks.LedgerController.ListByUserFunc = func(ctx context.Context, actorID domain.ID) ([]domain.Ledger, error) {
			assert.Equal(t, identity.UserID, actorID)

			return nil, assert.AnError
		}

		_, err := router.LedgerList(ctx)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fail/unauthenticated", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		router, _ := testutils.Setup(t, testutils.WithIdentity(nil))
		_, err := router.LedgerList(ctx)
		require.ErrorIs(t, err, assert.AnError)
	})
}
