package users

import (
	"github.com/sonalys/goshare/internal/application/controllers/identitycontroller"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/ports"
)

type Router struct {
	IdentityController *identitycontroller.Controller
	UserController     usercontroller.Controller
	SecurityHandler    ports.SecurityHandler
}

func New(
	securityHandler ports.SecurityHandler,
	identityController *identitycontroller.Controller,
	userController usercontroller.Controller,
) *Router {
	return &Router{
		SecurityHandler:    securityHandler,
		IdentityController: identityController,
		UserController:     userController,
	}
}
