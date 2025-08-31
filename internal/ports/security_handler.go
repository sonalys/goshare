package ports

import (
	"context"

	v1 "github.com/sonalys/goshare/pkg/v1"
)

type (
	SecurityHandler interface {
		GetIdentity(ctx context.Context) (*v1.Identity, error)
	}
)
