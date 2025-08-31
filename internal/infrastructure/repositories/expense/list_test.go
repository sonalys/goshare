package expense_test

import (
	"testing"
	"time"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/repositories"
	"github.com/sonalys/goshare/internal/ports"
	"github.com/sonalys/goshare/pkg/testcontainers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Expense_ListByLedger(t *testing.T) {
	client := repositories.New(testcontainers.Postgres(t))

	t.Run("pass/found", func(t *testing.T) {
		ctx := t.Context()
		td := createTestData(t)

		setup(t, client, td, func(r ports.LocalRepositories) {
			got, err := r.Expense().ListByLedger(ctx, td.ledger.ID, domain.Now().Add(time.Second), 1)
			require.NoError(t, err)
			assert.Len(t, got, 1)
		})
	})

	t.Run("pass/empty", func(t *testing.T) {
		ctx := t.Context()

		got, err := client.Expense().ListByLedger(ctx, domain.NewID(), domain.Now().Add(time.Second), 1)
		require.NoError(t, err)
		assert.Empty(t, got)
	})
}
