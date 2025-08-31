package user

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/ports"
)

var constraintMapping = map[string]error{
	"unique_user_email": domain.ErrUserAlreadyRegistered,
}

type Repository struct {
	conn postgres.Connection
}

func New(client postgres.Connection) ports.UserRepository {
	return &Repository{
		conn: client,
	}
}

func userError(err error) error {
	if err == nil {
		return nil
	}

	if err := postgres.MapConstraintError(err, constraintMapping); err != nil {
		return err
	}

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return domain.ErrUserNotFound
	default:
		return postgres.DefaultErrorMapping(err)
	}
}
