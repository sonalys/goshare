package api

import (
	"github.com/sonalys/goshare/cmd/server/handlers"
)

type (
	API struct {
		handlers.UnimplementedHandler
		dependencies Dependencies
	}
)

func New(dependencies Dependencies) *API {
	return &API{
		dependencies: dependencies,
	}
}
