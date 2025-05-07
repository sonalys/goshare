package main

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/application/pkg/secrets"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
)

type infrastructure struct {
	postgres *postgres.Postgres
}

func loadInfrastructure(ctx context.Context, secrets secrets.Secrets) *infrastructure {
	postgresClient, err := postgres.New(ctx, secrets.PostgresConn)
	if err != nil {
		panic(fmt.Errorf("failed to load Postgres infrastructure: %w", err))
	}
	return &infrastructure{
		postgres: postgresClient,
	}
}
