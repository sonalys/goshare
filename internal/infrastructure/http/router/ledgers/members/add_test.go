package members_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
	"github.com/sonalys/goshare/internal/infrastructure/http/testutils"
	"github.com/sonalys/goshare/internal/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LedgerMemberAdd(t *testing.T) {
	t.Parallel()

	type testData struct {
		req    *server.LedgerMemberAddReq
		params server.LedgerMemberAddParams
	}

	getTestData := func(_ *testing.T) testData {
		req := &server.LedgerMemberAddReq{
			Emails: []string{"member@example.com"},
		}

		params := server.LedgerMemberAddParams{
			LedgerID: uuid.New(),
		}

		return testData{
			req:    req,
			params: params,
		}
	}

	assertCreateFunc := func(t *testing.T, identity *application.Identity, td testData, got usercontroller.AddMembersRequest) {
		req := td.req
		params := td.params

		assert.Equal(t, req.Emails, got.Emails)
		assert.Equal(t, identity.UserID, got.ActorID)
		assert.Equal(t, params.LedgerID, got.LedgerID.UUID())
	}

	t.Run("pass", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		identity := testfixtures.Identity(t)
		router, mocks := testutils.Setup(t, testutils.WithIdentity(identity))

		td := getTestData(t)

		expense := td.req
		params := td.params

		mocks.LedgerController.MembersAddFunc = func(ctx context.Context, got usercontroller.AddMembersRequest) error {
			assertCreateFunc(t, identity, td, got)

			return nil
		}

		err := router.LedgerMemberAdd(ctx, expense, params)
		require.NoError(t, err)
	})

	t.Run("fail/controller error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		identity := testfixtures.Identity(t)
		router, mocks := testutils.Setup(t, testutils.WithIdentity(identity))

		td := getTestData(t)

		mocks.LedgerController.MembersAddFunc = func(ctx context.Context, got usercontroller.AddMembersRequest) error {
			assertCreateFunc(t, identity, td, got)

			return assert.AnError
		}

		err := router.LedgerMemberAdd(ctx, td.req, td.params)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fail/unauthenticated", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		router, _ := testutils.Setup(t, testutils.WithIdentity(nil))
		td := getTestData(t)

		err := router.LedgerMemberAdd(ctx, td.req, td.params)
		require.ErrorIs(t, err, assert.AnError)
	})
}
