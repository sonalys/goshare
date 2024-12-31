package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/participants"
)

type (
	ParticipantRegister interface {
		Register(ctx context.Context, req participants.RegisterRequest) (*participants.RegisterResponse, error)
	}

	Dependencies struct {
		ParticipantRegister
	}

	API struct {
		dependencies Dependencies
	}
)

var (
	_ handlers.StrictServerInterface = (*API)(nil)
)

func New(dependencies Dependencies) *API {
	return &API{
		dependencies: dependencies,
	}
}
