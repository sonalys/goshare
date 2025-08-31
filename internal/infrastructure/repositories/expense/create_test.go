package expense_test

import (
	"testing"

	"github.com/sonalys/goshare/internal/infrastructure/repositories"
	"github.com/sonalys/goshare/internal/ports"
	"github.com/sonalys/goshare/pkg/testcontainers"
	"github.com/sonalys/goshare/pkg/testfixtures"
	"github.com/stretchr/testify/require"
)

func Test_Expense_Create(t *testing.T) {
	client := repositories.New(testcontainers.Postgres(t))

	t.Run("pass", func(t *testing.T) {
		ctx := t.Context()

		from := testfixtures.User(t)
		to := testfixtures.User(t)
		ledger := testfixtures.Ledger(t, from)
		expense := testfixtures.Expense(t, ledger, from.ID, to.ID)

		err := client.Transaction(ctx, func(r ports.LocalRepositories) error {
			err := r.User().Create(ctx, from)
			require.NoError(t, err)

			err = r.User().Create(ctx, to)
			require.NoError(t, err)

			err = r.Ledger().Create(ctx, ledger)
			require.NoError(t, err)

			return r.Expense().Create(ctx, ledger.ID, expense)
		})
		require.NoError(t, err)
	})
}
