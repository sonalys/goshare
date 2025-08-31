package testcontainers

import (
	"context"
	"testing"

	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/migrations"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	module "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func Postgres(t *testing.T) postgres.Connection {
	ctx := t.Context()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	container, err := module.Run(ctx,
		"postgres:17-alpine",
		module.WithDatabase(dbName),
		module.WithUsername(dbUser),
		module.WithPassword(dbPassword),
		module.BasicWaitStrategies(),
		testcontainers.WithReuseByName("goshare-test-postgres"),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := container.Terminate(context.Background())
		require.NoError(t, err)
	})

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	conn, err := postgres.New(ctx, connStr)
	require.NoError(t, err)

	err = migrations.MigrateUp(ctx, connStr)
	require.NoError(t, err)

	return conn
}
