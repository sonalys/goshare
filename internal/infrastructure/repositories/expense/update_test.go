package expense_test

import (
	"testing"
	"testing/synctest"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/repositories"
	"github.com/sonalys/goshare/internal/pkg/testcontainers"
	"github.com/sonalys/goshare/internal/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Expense_Update(t *testing.T) {
	client := repositories.New(testcontainers.Postgres(t))

	t.Run("pass/found", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			td := createTestData(t)

			setup(t, client, td, func(r ports.LocalRepositories) {
				td.expense.Name = "updated name"

				err := r.Expense().Update(ctx, td.expense)
				require.NoError(t, err)

				got, err := r.Expense().Get(ctx, td.expense.ID)
				require.NoError(t, err)
				assert.Equal(t, td.expense, got)
			})
		})
	})

	t.Run("fail/not found", func(t *testing.T) {
		ctx := t.Context()
		td := createTestData(t)

		err := client.Transaction(ctx, func(r ports.LocalRepositories) error {
			return r.Expense().Update(ctx, td.expense)
		})
		require.ErrorIs(t, err, domain.ErrExpenseNotFound)
	})
}
