package ports

import (
	"context"

	"github.com/sonalys/goshare/internal/application"
)

type (
	SecurityHandler interface {
		GetIdentity(ctx context.Context) (*application.Identity, error)
	}
)
