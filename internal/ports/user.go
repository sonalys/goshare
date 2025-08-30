package ports

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
)

type (
	UserQueries interface {
		Get(ctx context.Context, id domain.ID) (*domain.User, error)
		GetByEmail(ctx context.Context, email string) (*domain.User, error)
		ListByEmail(ctx context.Context, emails []string) ([]domain.User, error)
	}

	UserCommands interface {
		Create(ctx context.Context, user *domain.User) error
	}

	UserRepository interface {
		UserQueries
		UserCommands
	}
)
