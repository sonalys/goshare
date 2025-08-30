package repositories

import (
	"context"

	"github.com/sonalys/goshare/internal/infrastructure/postgres"
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
	return newUserRepository(r.conn)
}

func (r localTransaction) User() ports.UserRepository {
	return newUserRepository(r.conn)
}

func (r LocalRepository) Expense() ports.ExpenseQueries {
	return newExpenseRepository(r.conn)
}

func (r localTransaction) Expense() ports.ExpenseRepository {
	return newExpenseRepository(r.conn)
}

func (r LocalRepository) Ledger() ports.LedgerQueries {
	return newLedgerRepository(r.conn)
}

func (r localTransaction) Ledger() ports.LedgerRepository {
	return newLedgerRepository(r.conn)
}

func (r LocalRepository) Transaction(ctx context.Context, handler func(ports.Repositories) error) error {
	return r.conn.Transaction(ctx, func(conn postgres.Connection) error {
		return handler(localTransaction{conn: conn})
	})
}
