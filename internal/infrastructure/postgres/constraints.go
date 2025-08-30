package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func isConstraintError(err error) bool {
	if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}

func MapConstraintError(err error, mapper map[string]error) error {
	pgErr := new(pgconn.PgError)
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return nil
	}
	return mapper[pgErr.ConstraintName]
}
