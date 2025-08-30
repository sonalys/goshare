package postgres

import (
	"context"
	"fmt"
	"slices"

	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

type UsersRepository struct {
	client connection
}

func NewUsersRepository(client connection) *UsersRepository {
	return &UsersRepository{
		client: client,
	}
}

func (r *UsersRepository) Save(ctx context.Context, user *domain.User) error {
	return mapError(r.client.queries().SaveUser(ctx, sqlcgen.SaveUserParams{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		LedgerCount:  user.LedgersCount,
		CreatedAt:    convertTime(user.CreatedAt),
	}))
}

func (r *UsersRepository) Get(ctx context.Context, id domain.ID) (*domain.User, error) {
	user, err := r.client.queries().GetUser(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	return mappers.NewUser(user), nil
}

func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.client.queries().GetUserByEmail(ctx, email)
	if err != nil {
		return nil, mapError(err)
	}

	return mappers.NewUser(user), nil
}

func (r *UsersRepository) ListByEmail(ctx context.Context, emails []string) ([]domain.User, error) {
	emails = slices.Compact(emails)
	users, err := r.client.queries().ListByEmail(ctx, emails)
	if err != nil {
		return nil, mapError(err)
	}

	var errs domain.Form
	for idx, email := range emails {
		if !slices.ContainsFunc(users, func(user sqlcgen.User) bool {
			return user.Email == email
		}) {
			errs = append(errs, domain.FieldError{
				Field: fmt.Sprintf("emails.%d", idx),
				Cause: v1.ErrNotFound,
			})
		}
	}

	if err := errs.Close(); err != nil {
		return nil, fmt.Errorf("failed to get users by email: %w", err)
	}

	return mappers.NewUsers(users), nil
}
