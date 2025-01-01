package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	constraintParticipantUniqueEmail = "participant_unique_email"
)

func isConstraintError(err error) bool {
	if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
