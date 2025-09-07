package main

import (
	"github.com/sonalys/goshare/internal/infrastructure/repositories"
	"github.com/sonalys/goshare/internal/pkg/jwt"
	"github.com/sonalys/goshare/internal/pkg/secrets"
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
