package postgres

import (
	"context"
	"fmt"
	"slices"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

type UsersRepository struct {
	client connection
}

func NewUsersRepository(client connection) *UsersRepository {
	return &UsersRepository{
		client: client,
	}
}

func (r *UsersRepository) Create(ctx context.Context, user *domain.User) error {
	return mapError(r.client.queries().CreateUser(ctx, sqlc.CreateUserParams{
		ID:           convertID(user.ID),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    convertTime(user.CreatedAt),
	}))
}

func (r *UsersRepository) Find(ctx context.Context, id domain.ID) (*domain.User, error) {
	user, err := r.client.queries().FindUser(ctx, convertID(id))
	if err != nil {
		return nil, mapError(err)
	}

	return mappers.NewUser(user), nil
}

func (r *UsersRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.client.queries().FindUserByEmail(ctx, email)
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

	var errs domain.FormError
	for idx, email := range emails {
		if !slices.ContainsFunc(users, func(user sqlc.UserView) bool {
			return user.Email == email
		}) {
			errs = append(errs, domain.FieldError{
				Field: fmt.Sprintf("emails.%d", idx),
				Cause: domain.ErrNotFound,
			})
		}
	}

	if err := errs.Validate(); err != nil {
		return nil, fmt.Errorf("failed to get users by email: %w", err)
	}

	return mappers.NewUsers(users), nil
}
