package users

import (
	"github.com/sonalys/goshare/internal/application/controllers/identitycontroller"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
)

type Router struct {
	IdentityController *identitycontroller.Controller
	UserController     *usercontroller.Controller
}

func New(
	identityController *identitycontroller.Controller,
	userController *usercontroller.Controller,
) *Router {
	return &Router{
		IdentityController: identityController,
		UserController:     userController,
	}
}
