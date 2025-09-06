package ledgers_test

import (
	"context"
	"testing"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	v1 "github.com/sonalys/goshare/internal/application/v1"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
	"github.com/sonalys/goshare/internal/infrastructure/http/testutils"
	"github.com/sonalys/goshare/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Create(t *testing.T) {
	t.Parallel()

	type testData struct {
		req *server.LedgerCreateReq
	}

	getTestData := func() testData {
		return testData{
			req: &server.LedgerCreateReq{
				Name: "my new ledger",
			},
		}
	}

	assertController := func(t *testing.T, identity *v1.Identity, td testData, got usercontroller.CreateLedgerRequest) {
		assert.Equal(t, td.req.Name, got.Name)
		assert.Equal(t, identity.UserID, got.ActorID)
	}

	t.Run("pass", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		identity := testfixtures.Identity(t)
		router, mocks := testutils.Setup(t, testutils.WithIdentity(identity))
		td := getTestData()

		ledgerID := domain.NewID()

		mocks.LedgerController.CreateFunc = func(ctx context.Context, got usercontroller.CreateLedgerRequest) (*usercontroller.CreateLedgerResponse, error) {
			assertController(t, identity, td, got)

			return &usercontroller.CreateLedgerResponse{
				ID: ledgerID,
			}, nil
		}

		resp, err := router.LedgerCreate(ctx, td.req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, ledgerID.UUID(), resp.ID)
	})

	t.Run("fail/controller error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		identity := testfixtures.Identity(t)
		router, mocks := testutils.Setup(t, testutils.WithIdentity(identity))
		td := getTestData()

		mocks.LedgerController.CreateFunc = func(ctx context.Context, got usercontroller.CreateLedgerRequest) (*usercontroller.CreateLedgerResponse, error) {
			assertController(t, identity, td, got)

			return nil, assert.AnError
		}

		_, err := router.LedgerCreate(ctx, td.req)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fail/unauthenticated", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		router, _ := testutils.Setup(t, testutils.WithIdentity(nil))
		_, err := router.LedgerCreate(ctx, &server.LedgerCreateReq{})
		require.ErrorIs(t, err, assert.AnError)
	})
}
