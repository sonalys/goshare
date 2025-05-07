package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/sonalys/goshare/internal/application/controllers"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

type pgxConn interface {
	sqlc.DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}

type conn[T pgxConn] struct {
	conn T
}

func (c *conn[T]) Ledger() controllers.LedgerRepository {
	return &LedgerRepository{
		client: c,
	}
}

func (c *conn[T]) User() controllers.UserRepository {
	return &UsersRepository{
		client: c,
	}
}

func (c *conn[T]) Expense() controllers.ExpenseRepository {
	return &ExpenseRepository{
		client: c,
	}
}

func (c *conn[T]) queries() *sqlc.Queries {
	return sqlc.New(c.conn)
}

func (c *conn[T]) transaction(ctx context.Context, f func(*sqlc.Queries) error) error {
	tx, err := c.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error(ctx, "failed to rollback transaction", err)
		}
	}()

	if err := f(sqlc.New(tx)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (c *conn[T]) Transaction(ctx context.Context, f func(controllers.Database) error) error {
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
