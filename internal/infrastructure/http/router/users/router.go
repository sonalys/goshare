package users

import (
	"github.com/sonalys/goshare/internal/application/controllers/identitycontroller"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/ports"
)

type Router struct {
	identityController identitycontroller.Controller
	controller         usercontroller.Controller
	securityHandler    ports.SecurityHandler
}

func New(
	securityHandler ports.SecurityHandler,
	identityController identitycontroller.Controller,
	userController usercontroller.Controller,
) *Router {
	return &Router{
		securityHandler:    securityHandler,
		identityController: identityController,
		controller:         userController,
	}
}
