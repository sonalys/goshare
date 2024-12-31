package main

import "github.com/sonalys/goshare/internal/application/users"

type controllers struct {
	participantController *users.UserController
}

func loadControllers(
	repositories *repositories,
) *controllers {
	return &controllers{
		participantController: users.NewParticipantController(repositories.ParticipantRepository),
	}
}
