package expenses

import (
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/ports"
)

type Router struct {
	ports.SecurityHandler
	usercontroller.Controller
}

func New(
	securityHandler ports.SecurityHandler,
	userController usercontroller.Controller,
) *Router {
	return &Router{
		SecurityHandler: securityHandler,
		Controller:      userController,
	}
}
