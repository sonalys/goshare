package repositories

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return mapUserErrors(r.conn.Queries().CreateUser(ctx, sqlcgen.CreateUserParams{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		LedgerCount:  user.LedgersCount,
		CreatedAt:    postgres.ConvertTime(user.CreatedAt),
	}))
}
