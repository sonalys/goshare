package main

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/pkg/secrets"
)

type infrastructure struct {
	postgres *postgres.Client
}

func loadInfrastructure(ctx context.Context, secrets secrets.Secrets) *infrastructure {
	postgresClient, err := postgres.NewClient(ctx, secrets.PostgresConn)
	if err != nil {
		panic(fmt.Errorf("failed to load Postgres infrastructure: %w", err))
	}
	return &infrastructure{
		postgres: postgresClient,
	}
}
