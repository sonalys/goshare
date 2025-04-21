package postgres

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/jackc/pgx/v5/pgconn"
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

func (r *UsersRepository) Create(ctx context.Context, participant *v1.User) error {
	return mapError(r.client.queries().CreateUser(ctx, sqlc.CreateUserParams{
		ID:           convertUUID(participant.ID),
		FirstName:    participant.FirstName,
		LastName:     participant.LastName,
		Email:        participant.Email,
		PasswordHash: participant.PasswordHash,
		CreatedAt:    convertTime(participant.CreatedAt),
	}))
}

func (r *UsersRepository) FindByEmail(ctx context.Context, email string) (*v1.User, error) {
	user, err := r.client.queries().FindUserByEmail(ctx, email)
	if err != nil {
		return nil, mapError(err)
	}

	return convertUser(user), nil
}

func convertUser(user sqlc.User) *v1.User {
	return &v1.User{
		ID:              newUUID(user.ID),
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		IsEmailVerified: false,
		PasswordHash:    user.PasswordHash,
		CreatedAt:       user.CreatedAt.Time,
	}
}

func convertUsers(from []sqlc.User) []v1.User {
	to := make([]v1.User, 0, len(from))

	for i := range from {
		to = append(to, *convertUser(from[i]))
	}

	return to
}

func (r *UsersRepository) GetByEmail(ctx context.Context, emails []string) ([]v1.User, error) {
	emails = slices.Compact(emails)
	users, err := r.client.queries().GetByEmail(ctx, emails)
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

	return convertUsers(users), nil
}

func isViolatingConstraint(err error, constraintName string) bool {
	if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
		return pgErr.ConstraintName == constraintName
	}
	return false
}
