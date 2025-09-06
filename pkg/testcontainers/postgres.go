package testcontainers

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
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
	conn postgres.Connection
	lock sync.Mutex
)

// watcher will kill the container once the pid is stopped.
func watcher(cid string) {
	pid := os.Getpid()
	//nolint:gosec,noctx // safe.
	watcher := exec.Command("sh", "-c",
		fmt.Sprintf(`(while kill -0 %d 2>/dev/null; do sleep 2; done; docker rm -f %s) & disown`, pid, cid),
	)

	// Detach stdio so watcher doesn't block
	watcher.Stdout = nil
	watcher.Stderr = nil
	watcher.Stdin = nil

	if err := watcher.Start(); err != nil {
		log.Fatalf("failed to start watcher: %v", err)
	}
}

func Postgres(t *testing.T) postgres.Connection {
	lock.Lock()
	defer lock.Unlock()

	if conn != nil {
		return conn
	}

	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	var err error
	container, err := module.Run(ctx,
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

	conn, err = postgres.New(ctx, connStr)
	require.NoError(t, err)

	err = migrations.MigrateUp(ctx, connStr)
	require.NoError(t, err)

	watcher(container.GetContainerID())

	return conn
}
