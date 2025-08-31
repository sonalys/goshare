package user_test

import (
	"testing"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/repositories"
	"github.com/sonalys/goshare/internal/ports"
	"github.com/sonalys/goshare/pkg/testcontainers"
	"github.com/sonalys/goshare/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_User_ListByEmail(t *testing.T) {
	client := repositories.New(testcontainers.Postgres(t))

	t.Run("pass/results", func(t *testing.T) {
		ctx := t.Context()
		user := testfixtures.User(t)

		err := client.Transaction(ctx, func(r ports.LocalRepositories) error {
			return r.User().Create(ctx, user)
		})
		require.NoError(t, err)

		got, err := client.User().ListByEmail(ctx, []string{user.Email})
		require.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run("fail/one email doesn't exist", func(t *testing.T) {
		ctx := t.Context()

		got, err := client.User().ListByEmail(ctx, []string{"random"})
		require.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Empty(t, got)
	})
}
