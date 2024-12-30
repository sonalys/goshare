package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/queries"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type ParticipantRepository struct {
	client *Client
}

func NewParticipantRepository(client *Client) *ParticipantRepository {
	return &ParticipantRepository{
		client: client,
	}
}

func convertUUID(from uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: from,
		Valid: true,
	}
}

func convertTime(from time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  from,
		Valid: true,
	}
}

func (r *ParticipantRepository) Create(ctx context.Context, participant *v1.Participant) error {
	return mapError(r.client.queries().CreateParticipant(ctx, queries.CreateParticipantParams{
		ID:           convertUUID(participant.ID),
		FirstName:    participant.FirstName,
		LastName:     participant.LastName,
		Email:        participant.Email,
		PasswordHash: participant.PasswordHash,
		CreatedAt:    convertTime(participant.CreatedAt),
	}))
}

func isConstraintError(err error, constraintName string) bool {
	if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
		return pgErr.ConstraintName == constraintName
	}
	return false
}

func mapError(err error) error {
	switch {
	case err == nil:
		return nil
	case isConstraintError(err, constraintParticipantUniqueEmail):
		return v1.ErrParticipantEmailAlreadyExists
	default:
		return err
	}
}
