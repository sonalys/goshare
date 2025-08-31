package main

import (
	"context"

	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/pkg/secrets"
	"github.com/sonalys/goshare/pkg/slog"
)

type infrastructure struct {
	postgresConnection postgres.Connection
}

func loadInfrastructure(ctx context.Context, secrets secrets.Secrets) *infrastructure {
	postgresClient, err := postgres.New(ctx, secrets.PostgresConn)
	if err != nil {
		slog.Panic(ctx, "initializing postgres", slog.WithError(err))
	}

	return &infrastructure{
		postgresConnection: postgresClient,
	}
}
