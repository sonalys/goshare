package main

import (
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/pkg/jwt"
	"github.com/sonalys/goshare/internal/pkg/secrets"
)

type repositories struct {
	UserRepository     *postgres.UsersRepository
	LedgerRepository   *postgres.LedgerRepository
	ExpensesRepository *postgres.ExpensesRepository
	JWTRepository      *jwt.Client
}

func loadRepositories(
	secrets secrets.Secrets,
	infrastructure *infrastructure,
) *repositories {
	return &repositories{
		UserRepository:     postgres.NewUsersRepository(infrastructure.postgres),
		LedgerRepository:   postgres.NewLedgerRepository(infrastructure.postgres),
		ExpensesRepository: postgres.NewExpensesRepository(infrastructure.postgres),
		JWTRepository:      jwt.NewClient(secrets.JWTSignKey),
	}
}
