package usercontroller_test

import (
	"context"
	"testing"

	"github.com/sonalys/goshare/internal/mocks/portsmock"
	"github.com/sonalys/goshare/internal/ports"
)

type repositoryMock struct {
	user    portsmock.UserRepository
	ledger  portsmock.LedgerRepository
	expense portsmock.ExpenseRepository
}

type databaseMock struct {
	repositories *repositoryMock
	tx           *repositoryMock
	db           portsmock.LocalDatabase
}

func setupDatabaseMock(_ *testing.T) *databaseMock {
	var repositories repositoryMock
	var tx repositoryMock

	return &databaseMock{
		repositories: &repositories,
		tx:           &tx,
		db: portsmock.LocalDatabase{
			ExpenseFunc: func() ports.ExpenseQueries { return &repositories.expense },
			LedgerFunc:  func() ports.LedgerQueries { return &repositories.ledger },
			UserFunc:    func() ports.UserQueries { return &repositories.user },
			TransactionFunc: func(ctx context.Context, f func(tx ports.LocalRepositories) error) error {
				return f(&portsmock.LocalRepositories{
					ExpenseFunc: func() ports.ExpenseRepository { return &tx.expense },
					LedgerFunc:  func() ports.LedgerRepository { return &tx.ledger },
					UserFunc:    func() ports.UserRepository { return &tx.user },
				})
			},
		},
	}
}
