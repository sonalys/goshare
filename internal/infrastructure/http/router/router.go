package router

import (
	"github.com/sonalys/goshare/internal/application/controllers/identitycontroller"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/infrastructure/http/router/ledgers"
	"github.com/sonalys/goshare/internal/infrastructure/http/router/ledgers/expenses"
	"github.com/sonalys/goshare/internal/infrastructure/http/router/ledgers/expenses/records"
	"github.com/sonalys/goshare/internal/infrastructure/http/router/users"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
	"github.com/sonalys/goshare/internal/ports"
)

type (
	Router struct {
		server.LedgersHandler
		server.UsersHandler
		server.RecordsHandler
		server.ExpensesHandler
	}
)

func New(
	securityHandler ports.SecurityHandler,
	identityController identitycontroller.Controller,
	userController usercontroller.Controller,
) server.Handler {
	return &Router{
		LedgersHandler:  ledgers.New(securityHandler, userController),
		UsersHandler:    users.New(securityHandler, identityController, userController),
		ExpensesHandler: expenses.New(securityHandler, userController),
		RecordsHandler:  records.New(securityHandler, userController),
	}
}
