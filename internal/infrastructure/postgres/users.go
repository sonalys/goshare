package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/queries"
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
	return mapError(r.client.queries().CreateUser(ctx, queries.CreateUserParams{
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

	return &v1.User{
		ID:              newUUID(user.ID),
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		IsEmailVerified: false,
		PasswordHash:    user.PasswordHash,
		CreatedAt:       user.CreatedAt.Time,
	}, nil
}

func isViolatingConstraint(err error, constraintName string) bool {
	if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
		return pgErr.ConstraintName == constraintName
	}
	return false
}

func mapUserError(err error) error {
	switch {
	case err == nil:
		return nil
	case isViolatingConstraint(err, constraintParticipantUniqueEmail):
		return v1.ErrEmailAlreadyRegistered
	default:
		return mapError(err)
	}
}
