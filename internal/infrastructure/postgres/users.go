package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/sonalys/goshare/internal/domain"
)

type UsersRepository struct {
	client connection
}

func NewUsersRepository(client connection) *UsersRepository {
	return &UsersRepository{
		client: client,
	}
}

func mapUserErrors(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, pgx.ErrNoRows):
		return domain.ErrUserNotFound
	case isViolatingConstraint(err, constraintUserUniqueEmail):
		return domain.ErrUserAlreadyRegistered
	default:
		return mapError(err)
	}
}
