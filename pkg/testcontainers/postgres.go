package testcontainers

import (
	"context"
	"io"
	"log"
	"sync"
	"testing"

	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/migrations"
	"github.com/sonalys/goshare/pkg/slog"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	module "github.com/testcontainers/testcontainers-go/modules/postgres"
)

//nolint:gochecknoinits // test utility
func init() {
	slog.Init(slog.LevelDebug)
}

var (
	conn      postgres.Connection
	container *module.PostgresContainer

	lock       sync.Mutex
	references int
)

func Postgres(t *testing.T) postgres.Connection {
	lock.Lock()
	defer lock.Unlock()

	references++
	t.Cleanup(func() {
		lock.Lock()
		defer lock.Unlock()

		references--

		if references > 0 || container == nil {
			return
		}

		err := container.Terminate(context.Background())
		require.NoError(t, err)

		container = nil
	})

	if conn != nil {
		return conn
	}

	ctx := t.Context()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	var err error
	container, err = module.Run(ctx,
		"postgres:17-alpine",
		module.WithDatabase(dbName),
		module.WithUsername(dbUser),
		module.WithPassword(dbPassword),
		module.BasicWaitStrategies(),
		testcontainers.WithReuseByName("goshare-test-postgres"),
		testcontainers.WithLogger(log.New(io.Discard, "", 0)),
	)
	require.NoError(t, err)

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	slog.Debug(ctx, "postgres started", slog.WithString("connStr", connStr))

	conn, err := postgres.New(ctx, connStr)
	require.NoError(t, err)

	err = migrations.MigrateUp(ctx, connStr)
	require.NoError(t, err)

	return conn
}
