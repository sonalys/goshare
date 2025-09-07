package expense_test

import (
	"testing"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/repositories"
	"github.com/sonalys/goshare/internal/pkg/testcontainers"
	"github.com/sonalys/goshare/internal/pkg/testfixtures"
	"github.com/sonalys/goshare/internal/ports"
	"github.com/stretchr/testify/require"
)

func Test_Expense_Create(t *testing.T) {
	client := repositories.New(testcontainers.Postgres(t))

	type testData struct {
		from    *domain.User
		to      *domain.User
		ledger  *domain.Ledger
		expense *domain.Expense
	}

	createTestData := func(t *testing.T) testData {
		from := testfixtures.User(t)
		to := testfixtures.User(t)
		ledger := testfixtures.Ledger(t, from)
		expense := testfixtures.Expense(t, ledger, from.ID, to.ID)

		return testData{
			from:    from,
			to:      to,
			ledger:  ledger,
			expense: expense,
		}
	}

	setup := func(t *testing.T, td testData, handler func(r ports.LocalRepositories)) {
		ctx := t.Context()
		var run bool
		err := client.Transaction(ctx, func(r ports.LocalRepositories) error {
			err := r.User().Create(ctx, td.from)
			require.NoError(t, err)

			err = r.User().Create(ctx, td.to)
			require.NoError(t, err)

			err = r.Ledger().Create(ctx, td.ledger)
			require.NoError(t, err)

			handler(r)

			run = true

			return nil
		})
		require.NoError(t, err)
		require.True(t, run)
	}

	t.Run("pass", func(t *testing.T) {
		ctx := t.Context()
		td := createTestData(t)

		setup(t, td, func(r ports.LocalRepositories) {
			err := r.Expense().Create(ctx, td.expense)
			require.NoError(t, err)
		})
	})
}
