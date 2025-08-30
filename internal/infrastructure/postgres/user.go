package postgres

import (
	"github.com/sonalys/goshare/internal/domain"
)

var userConstraintMapping = map[string]error{
	"unique_user_email": domain.ErrUserAlreadyRegistered,
}

type UsersRepository struct {
	client connection
}

func NewUsersRepository(client connection) *UsersRepository {
	return &UsersRepository{
		client: client,
	}
}

func mapUserErrors(err error) error {
	if err := constraintErrorMap(err, userConstraintMapping); err != nil {
		return err
	}

	return defaultErrorMapping(err)
}
