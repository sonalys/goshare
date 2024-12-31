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

	Controller struct {
		repository Repository
	}
)

func NewController(repository Repository) *Controller {
	return &Controller{
		repository: repository,
	}
}
