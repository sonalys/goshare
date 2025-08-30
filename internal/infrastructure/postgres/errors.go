package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
)

func DefaultErrorMapping(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, pgx.ErrNoRows):
		return v1.ErrNotFound
	case isConstraintError(err):
		return v1.ErrConflict
	default:
		return err
	}
}
