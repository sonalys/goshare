package users

import (
	"context"

	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type (
	Repository interface {
		Create(ctx context.Context, user *v1.User) error
		FindByEmail(ctx context.Context, email string) (*v1.User, error)
	}

	IdentityEncoder interface {
		Encode(identity *v1.Identity) (string, error)
	}

	Controller struct {
		identityEncoder IdentityEncoder
		repository      Repository
	}
)

func NewController(
	repository Repository,
	identityEncoder IdentityEncoder,
) *Controller {
	return &Controller{
		repository:      repository,
		identityEncoder: identityEncoder,
	}
}
