package ledger_test

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

func Test_Ledger_ListByUser(t *testing.T) {
	client := repositories.New(testcontainers.Postgres(t))

	type testData struct {
		from   *domain.User
		to     *domain.User
		ledger *domain.Ledger
	}

	createTestData := func(t *testing.T) testData {
		from := testfixtures.User(t)
		to := testfixtures.User(t)

		ledger := testfixtures.Ledger(t, from)

		err := ledger.AddMember(from.ID, to.ID)
		require.NoError(t, err)

		return testData{
			from:   from,
			to:     to,
			ledger: ledger,
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

	t.Run("pass/created by", func(t *testing.T) {
		ctx := t.Context()
		td := createTestData(t)

		setup(t, td, func(r ports.LocalRepositories) {
			got, err := r.Ledger().ListByUser(ctx, td.ledger.CreatedBy)
			require.NoError(t, err)
			assert.Len(t, got, 1)
		})
	})

	t.Run("pass/member", func(t *testing.T) {
		ctx := t.Context()
		td := createTestData(t)

		var anotherMember *domain.ID

		for id := range td.ledger.Members {
			if id != td.ledger.CreatedBy {
				anotherMember = &id

				break
			}
		}
		require.NotNil(t, anotherMember)

		setup(t, td, func(r ports.LocalRepositories) {
			got, err := r.Ledger().ListByUser(ctx, *anotherMember)
			require.NoError(t, err)
			assert.Len(t, got, 1)
		})
	})

	t.Run("pass/empty", func(t *testing.T) {
		ctx := t.Context()
		td := createTestData(t)

		setup(t, td, func(r ports.LocalRepositories) {
			got, err := r.Ledger().ListByUser(ctx, domain.NewID())
			require.NoError(t, err)
			assert.Empty(t, got)
		})
	})
}
