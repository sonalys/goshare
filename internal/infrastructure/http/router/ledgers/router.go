package ledgers

import (
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
)

type Router struct {
	UserController *usercontroller.Controller
}

func New(
	userController *usercontroller.Controller,
) *Router {
	return &Router{
		UserController: userController,
	}
}
