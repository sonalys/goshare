package main

import "github.com/sonalys/goshare/internal/application/participants"

type controllers struct {
	participantController *participants.ParticipantController
}

func loadControllers(
	repositories *repositories,
) *controllers {
	return &controllers{
		participantController: participants.NewParticipantController(repositories.ParticipantRepository),
	}
}
