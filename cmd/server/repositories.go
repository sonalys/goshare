package main

import (
	"github.com/sonalys/goshare/internal/application/pkg/jwt"
	"github.com/sonalys/goshare/internal/application/pkg/secrets"
	"github.com/sonalys/goshare/internal/infrastructure/repositories"
)

type repos struct {
	Database      repositories.LocalRepository
	JWTRepository *jwt.Client
}

func loadRepositories(secrets secrets.Secrets, infrastructure *infrastructure) *repos {
	return &repos{
		Database:      repositories.New(infrastructure.postgresConnection),
		JWTRepository: jwt.NewClient(secrets.JWTSignKey),
	}
}
