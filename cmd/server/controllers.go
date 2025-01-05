package main

import (
	"github.com/sonalys/goshare/internal/application/ledgers"
	"github.com/sonalys/goshare/internal/application/users"
)

type controllers struct {
	userController   *users.Controller
	ledgerController *ledgers.Controller
}

func loadControllers(repositories *repositories) *controllers {
	return &controllers{
		userController:   users.NewController(repositories.UserRepository, repositories.JWTRepository),
		ledgerController: ledgers.NewController(repositories.LedgerRepository, repositories.ExpensesRepository, repositories.UserRepository),
	}
}
