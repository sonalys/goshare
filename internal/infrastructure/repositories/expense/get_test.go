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

func Test_Expense_Get(t *testing.T) {
	client := repositories.New(testcontainers.Postgres(t))

	t.Run("pass/found", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			td := createTestData(t)

			setup(t, client, td, func(r ports.LocalRepositories) {
				got, err := r.Expense().Get(ctx, td.expense.ID)
				require.NoError(t, err)
				assert.Equal(t, td.expense, got)
			})
		})
	})

	t.Run("fail/not found", func(t *testing.T) {
		ctx := t.Context()

		_, err := client.Expense().Get(ctx, domain.NewID())
		require.ErrorIs(t, err, domain.ErrExpenseNotFound)
	})
}
