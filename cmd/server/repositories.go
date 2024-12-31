package main

import "github.com/sonalys/goshare/internal/infrastructure/postgres"

type repositories struct {
	ParticipantRepository *postgres.UsersRepository
}

func loadRepositories(infrastructure *infrastructure) *repositories {
	return &repositories{
		ParticipantRepository: postgres.NewUsersRepository(infrastructure.postgres),
	}
}
