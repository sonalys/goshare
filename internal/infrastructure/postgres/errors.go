package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/sonalys/goshare/internal/domain"
)

func mapError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, pgx.ErrNoRows):
		return domain.ErrNotFound
	case isConstraintError(err):
		return domain.ErrConflict
	default:
		return err
	}
}
