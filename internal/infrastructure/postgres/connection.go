package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

type (
	connection interface {
		transaction(ctx context.Context, f func(q connection) error) error
		queries() *sqlcgen.Queries
		readWrite() *readWriteRepository
	}

	pgxConn interface {
		sqlcgen.DBTX
		Begin(ctx context.Context) (pgx.Tx, error)
	}

	conn[T pgxConn] struct {
		conn T
	}

	readWriteRepository struct{ connection }
)

func (c *readWriteRepository) Ledger() application.LedgerRepository {
	return &LedgerRepository{
		client: c.connection,
	}
}

func (c *readWriteRepository) User() application.UserRepository {
	return &UsersRepository{
		client: c.connection,
	}
}

func (c *readWriteRepository) Expense() application.ExpenseRepository {
	return &ExpenseRepository{
		client: c.connection,
	}
}

func (c *conn[T]) queries() *sqlcgen.Queries {
	return sqlcgen.New(c.conn)
}

func (c *conn[T]) transaction(ctx context.Context, f func(connection) error) error {
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

func (c *conn[T]) readWrite() *readWriteRepository {
	return &readWriteRepository{
		connection: c,
	}
}
