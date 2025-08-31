package repositories

import (
	"context"

	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/infrastructure/repositories/expense"
	"github.com/sonalys/goshare/internal/infrastructure/repositories/ledger"
	"github.com/sonalys/goshare/internal/infrastructure/repositories/user"
	"github.com/sonalys/goshare/internal/ports"
)

type (
	LocalRepository struct {
		conn postgres.Connection
	}

	localTransaction struct {
		conn postgres.Connection
	}
)

func New(conn postgres.Connection) LocalRepository {
	return LocalRepository{
		conn: conn,
	}
}

func (r LocalRepository) User() ports.UserQueries {
	return user.New(r.conn)
}

func (r localTransaction) User() ports.UserRepository {
	return user.New(r.conn)
}

func (r LocalRepository) Expense() ports.ExpenseQueries {
	return expense.New(r.conn)
}

func (r localTransaction) Expense() ports.ExpenseRepository {
	return expense.New(r.conn)
}

func (r LocalRepository) Ledger() ports.LedgerQueries {
	return ledger.New(r.conn)
}

func (r localTransaction) Ledger() ports.LedgerRepository {
	return ledger.New(r.conn)
}

func (r LocalRepository) Transaction(ctx context.Context, handler func(r ports.LocalRepositories) error) error {
	return r.conn.Transaction(ctx, func(conn postgres.Connection) error {
		return handler(localTransaction{conn: conn})
	})
}
