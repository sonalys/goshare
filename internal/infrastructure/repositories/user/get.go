package user

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
)

func (r *Repository) Get(ctx context.Context, id domain.ID) (*domain.User, error) {
	user, err := r.conn.Queries().GetUser(ctx, id)
	if err != nil {
		return nil, userError(err)
	}

	return mappers.NewUser(user), nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.conn.Queries().GetUserByEmail(ctx, email)
	if err != nil {
		return nil, userError(err)
	}

	return mappers.NewUser(user), nil
}
