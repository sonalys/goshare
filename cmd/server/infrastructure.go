package main

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/pkg/secrets"
)

type infrastructure struct {
	postgresConnection postgres.Connection
}

func loadInfrastructure(ctx context.Context, secrets secrets.Secrets) *infrastructure {
	postgresClient, err := postgres.New(ctx, secrets.PostgresConn)
	if err != nil {
		panic(fmt.Errorf("failed to load Postgres infrastructure: %w", err))
	}
	return &infrastructure{
		postgresConnection: postgresClient,
	}
}
