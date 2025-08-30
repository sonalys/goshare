package postgres

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
)

func (r *UsersRepository) Get(ctx context.Context, id domain.ID) (*domain.User, error) {
	user, err := r.client.queries().GetUser(ctx, id)
	if err != nil {
		return nil, mapUserErrors(err)
	}

	return mappers.NewUser(user), nil
}

func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.client.queries().GetUserByEmail(ctx, email)
	if err != nil {
		return nil, mapUserErrors(err)
	}

	return mappers.NewUser(user), nil
}
