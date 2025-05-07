package main

import (
	"github.com/sonalys/goshare/internal/application/pkg/jwt"
	"github.com/sonalys/goshare/internal/application/pkg/secrets"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
)

type repositories struct {
	Database      *postgres.Postgres
	JWTRepository *jwt.Client
}

func loadRepositories(
	secrets secrets.Secrets,
	infrastructure *infrastructure,
) *repositories {
	return &repositories{
		Database:      infrastructure.postgres,
		JWTRepository: jwt.NewClient(secrets.JWTSignKey),
	}
}
