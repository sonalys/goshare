package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonalys/goshare/internal/application/controllers"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

type connection interface {
	Transaction(ctx context.Context, f func(uow controllers.Repositories) error) error
	transaction(ctx context.Context, f func(query *sqlc.Queries) error) error
	queries() *sqlc.Queries
}

type Postgres struct {
	*conn[*pgxpool.Pool]
}

type pgxConn interface {
	sqlc.DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}

func New(ctx context.Context, connStr string) (*Postgres, error) {
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connStr: %w", err)
	}

	cfg.ConnConfig.Tracer = tracer{}

	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &Postgres{
		conn: &conn[*pgxpool.Pool]{
			Pool: dbpool,
		},
	}, nil
}

type conn[T pgxConn] struct {
	Pool T
}

func (r *conn[T]) Ledger() controllers.LedgerRepository {
	return &LedgerRepository{
		client: r,
	}
}

func (r *conn[T]) User() controllers.UserRepository {
	return &UsersRepository{
		client: r,
	}
}

func (r *conn[T]) Expense() controllers.ExpenseRepository {
	return &ExpenseRepository{
		client: r,
	}
}

func (c *conn[T]) queries() *sqlc.Queries {
	return sqlc.New(c.Pool)
}

func (c *conn[T]) transaction(ctx context.Context, f func(*sqlc.Queries) error) error {
	tx, err := c.Pool.Begin(ctx)
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

func (c *conn[T]) Transaction(ctx context.Context, f func(controllers.Repositories) error) error {
	tx, err := c.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error(ctx, "failed to rollback transaction", err)
		}
	}()

	if err := f(&conn[pgx.Tx]{Pool: tx}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
