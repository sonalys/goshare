package usercontroller_test

import (
	"context"
	"testing"

	"github.com/sonalys/goshare/internal/application"
	applicationmock "github.com/sonalys/goshare/mocks/internal_/application"
)

type repositoryMock struct {
	user    applicationmock.UserRepository
	ledger  applicationmock.LedgerRepository
	expense applicationmock.ExpenseRepository
}

type databaseMock struct {
	repositories *repositoryMock
	tx           *repositoryMock
	db           applicationmock.Database
}

func setupDatabaseMock(_ *testing.T) *databaseMock {
	var repositories repositoryMock
	var tx repositoryMock

	return &databaseMock{
		repositories: &repositories,
		tx:           &tx,
		db: applicationmock.Database{
			ExpenseFunc: func() application.ExpenseQueries { return &repositories.expense },
			LedgerFunc:  func() application.LedgerQueries { return &repositories.ledger },
			UserFunc:    func() application.UserQueries { return &repositories.user },
			TransactionFunc: func(ctx context.Context, f func(tx application.Repositories) error) error {
				return f(&applicationmock.Repositories{
					ExpenseFunc: func() application.ExpenseRepository { return &tx.expense },
					LedgerFunc:  func() application.LedgerRepository { return &tx.ledger },
					UserFunc:    func() application.UserRepository { return &tx.user },
				})
			},
		},
	}
}
