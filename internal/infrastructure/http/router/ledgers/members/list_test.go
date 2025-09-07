package members_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
	"github.com/sonalys/goshare/internal/infrastructure/http/testutils"
	"github.com/sonalys/goshare/internal/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LedgerMemberList(t *testing.T) {
	t.Parallel()

	type testData struct {
		params server.LedgerMemberListParams
	}

	getTestData := func(_ *testing.T) testData {
		params := server.LedgerMemberListParams{
			LedgerID: uuid.New(),
		}

		return testData{
			params: params,
		}
	}

	assertController := func(t *testing.T, identity *application.Identity, td testData, got usercontroller.GetLedgerRequest) {
		params := td.params

		assert.Equal(t, identity.UserID, got.ActorID)
		assert.Equal(t, params.LedgerID, got.LedgerID.UUID())
	}

	t.Run("pass", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		identity := testfixtures.Identity(t)
		router, mocks := testutils.Setup(t, testutils.WithIdentity(identity))

		td := getTestData(t)

		params := td.params

		mocks.LedgerController.GetFunc = func(ctx context.Context, got usercontroller.GetLedgerRequest) (*domain.Ledger, error) {
			assertController(t, identity, td, got)

			return &domain.Ledger{}, nil
		}

		resp, err := router.LedgerMemberList(ctx, params)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("fail/controller error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		identity := testfixtures.Identity(t)
		router, mocks := testutils.Setup(t, testutils.WithIdentity(identity))

		td := getTestData(t)

		mocks.LedgerController.GetFunc = func(ctx context.Context, got usercontroller.GetLedgerRequest) (*domain.Ledger, error) {
			assertController(t, identity, td, got)

			return nil, assert.AnError
		}

		_, err := router.LedgerMemberList(ctx, td.params)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fail/unauthenticated", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		router, _ := testutils.Setup(t, testutils.WithIdentity(nil))
		td := getTestData(t)

		_, err := router.LedgerMemberList(ctx, td.params)
		require.ErrorIs(t, err, assert.AnError)
	})
}
