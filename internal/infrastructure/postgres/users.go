package postgres

import (
	"context"
	"fmt"
	"slices"

	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type UsersRepository struct {
	client *Client
}

func NewUsersRepository(client *Client) *UsersRepository {
	return &UsersRepository{
		client: client,
	}
}

func (r *UsersRepository) Create(ctx context.Context, user *v1.User) error {
	return mapError(r.client.queries().CreateUser(ctx, sqlc.CreateUserParams{
		ID:           convertID(user.ID),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    convertTime(user.CreatedAt),
	}))
}

func (r *UsersRepository) FindByEmail(ctx context.Context, email string) (*v1.User, error) {
	user, err := r.client.queries().FindUserByEmail(ctx, email)
	if err != nil {
		return nil, mapError(err)
	}

	return mappers.NewUser(user), nil
}

func (r *UsersRepository) ListByEmail(ctx context.Context, emails []string) ([]v1.User, error) {
	emails = slices.Compact(emails)
	users, err := r.client.queries().ListByEmail(ctx, emails)
	if err != nil {
		return nil, mapError(err)
	}

	var errs v1.FormError
	for idx, email := range emails {
		if !slices.ContainsFunc(users, func(user sqlc.User) bool {
			return user.Email == email
		}) {
			errs = append(errs, v1.FieldError{
				Field: fmt.Sprintf("emails.%d", idx),
				Cause: v1.ErrNotFound,
			})
		}
	}

	if err := errs.Validate(); err != nil {
		return nil, fmt.Errorf("failed to get users by email: %w", err)
	}

	return mappers.NewUsers(users), nil
}
