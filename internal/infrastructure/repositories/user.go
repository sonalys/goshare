package repositories

import (
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/ports"
)

var userConstraintMapping = map[string]error{
	"unique_user_email": domain.ErrUserAlreadyRegistered,
}

type UserRepository struct {
	conn postgres.Connection
}

func newUserRepository(client postgres.Connection) ports.UserRepository {
	return &UserRepository{
		conn: client,
	}
}

func mapUserErrors(err error) error {
	if err := postgres.MapConstraintError(err, userConstraintMapping); err != nil {
		return err
	}

	return postgres.DefaultErrorMapping(err)
}
