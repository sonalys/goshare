package postgres

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

func (r *UsersRepository) Create(ctx context.Context, user *domain.User) error {
	return mapUserErrors(r.client.queries().CreateUser(ctx, sqlcgen.CreateUserParams{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		LedgerCount:  user.LedgersCount,
		CreatedAt:    convertTime(user.CreatedAt),
	}))
}
