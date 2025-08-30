package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

type (
	Connection interface {
		Transaction(ctx context.Context, f func(q Connection) error) error
		Queries() *sqlcgen.Queries
	}

	pgxConn interface {
		sqlcgen.DBTX
		Begin(ctx context.Context) (pgx.Tx, error)
	}

	conn[T pgxConn] struct {
		conn T
	}
)

func (c *conn[T]) Queries() *sqlcgen.Queries {
	return sqlcgen.New(c.conn)
}

func (c *conn[T]) Transaction(ctx context.Context, f func(Connection) error) error {
	tx, err := c.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error(ctx, "failed to rollback transaction", err)
		}
	}()

	if err := f(&conn[pgx.Tx]{conn: tx}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
