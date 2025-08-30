package repositories_test

import (
	"testing"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/ports"
	"github.com/sonalys/goshare/pkg/testfixtures"
	"github.com/stretchr/testify/require"
)

func Test_Users_Create(t *testing.T) {
	client := initializePostgres(t)

	t.Run("pass", func(t *testing.T) {
		ctx := t.Context()

		err := client.Transaction(ctx, func(r ports.Repositories) error {
			user := testfixtures.User(t)
			return r.User().Create(ctx, user)
		})
		require.NoError(t, err)
	})

	t.Run("fail/email conflict", func(t *testing.T) {
		ctx := t.Context()

		err := client.Transaction(ctx, func(r ports.Repositories) error {
			user := testfixtures.User(t)

			err := r.User().Create(ctx, user)
			require.NoError(t, err)

			user.ID = domain.NewID()

			err = r.User().Create(ctx, user)
			require.ErrorIs(t, err, domain.ErrUserAlreadyRegistered)

			return nil
		})
		require.NoError(t, err)
	})
}
