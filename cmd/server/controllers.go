package main

import "github.com/sonalys/goshare/internal/application/users"

type controllers struct {
	userController *users.Controller
}

func loadControllers(
	repositories *repositories,
) *controllers {
	return &controllers{
		userController: users.NewController(repositories.ParticipantRepository),
	}
}
