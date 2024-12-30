package main

import "github.com/sonalys/goshare/internal/infrastructure/postgres"

type repositories struct {
	ParticipantRepository *postgres.ParticipantRepository
}

func loadRepositories(infrastructure *infrastructure) *repositories {
	return &repositories{
		ParticipantRepository: postgres.NewParticipantRepository(infrastructure.postgres),
	}
}
