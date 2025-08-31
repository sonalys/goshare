package user_test

import (
	"testing"
	"testing/synctest"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/repositories"
	"github.com/sonalys/goshare/internal/ports"
	"github.com/sonalys/goshare/pkg/testcontainers"
	"github.com/sonalys/goshare/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_User_Get(t *testing.T) {
	client := repositories.New(testcontainers.Postgres(t))

	t.Run("pass/found", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			user := testfixtures.User(t)

			err := client.Transaction(ctx, func(r ports.LocalRepositories) error {
				return r.User().Create(ctx, user)
			})
			require.NoError(t, err)

			got, err := client.User().Get(ctx, user.ID)
			require.NoError(t, err)
			assert.Equal(t, user, got)
		})
	})

	t.Run("fail/not found", func(t *testing.T) {
		ctx := t.Context()

		_, err := client.User().Get(ctx, domain.NewID())
		require.ErrorIs(t, err, domain.ErrUserNotFound)
	})
}

func Test_User_GetByEmail(t *testing.T) {
	client := repositories.New(testcontainers.Postgres(t))

	t.Run("pass/found", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			user := testfixtures.User(t)

			err := client.Transaction(ctx, func(r ports.LocalRepositories) error {
				return r.User().Create(ctx, user)
			})
			require.NoError(t, err)

			got, err := client.User().GetByEmail(ctx, user.Email)
			require.NoError(t, err)
			assert.Equal(t, user, got)
		})
	})

	t.Run("fail/not found", func(t *testing.T) {
		ctx := t.Context()

		_, err := client.User().GetByEmail(ctx, domain.NewID().String())
		require.ErrorIs(t, err, domain.ErrUserNotFound)
	})
}
