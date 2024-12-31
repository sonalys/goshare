package main

import (
	"github.com/sonalys/goshare/internal/application/users"
	"github.com/sonalys/goshare/internal/pkg/secrets"
)

type controllers struct {
	userController *users.Controller
}

func loadControllers(
	secrets secrets.Secrets,
	repositories *repositories,
) *controllers {
	return &controllers{
		userController: users.NewController(repositories.ParticipantRepository, secrets.JWTSignKey),
	}
}
