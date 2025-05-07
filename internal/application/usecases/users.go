package usecases

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers"
)

type (
	Users interface {
		Login(ctx context.Context, req controllers.LoginRequest) (*controllers.LoginResponse, error)
		Register(ctx context.Context, req controllers.RegisterRequest) (*controllers.RegisterResponse, error)
	}

	users struct {
		controller *controllers.Users
	}
)

func NewUsers(userController *controllers.Users) Users {
	return &users{
		controller: userController,
	}
}

func (u *users) Login(ctx context.Context, req controllers.LoginRequest) (*controllers.LoginResponse, error) {
	return u.controller.Login(ctx, req)
}

func (u *users) Register(ctx context.Context, req controllers.RegisterRequest) (*controllers.RegisterResponse, error) {
	return u.controller.Register(ctx, req)
}
