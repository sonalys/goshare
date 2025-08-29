package identitycontroller

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/domain"
)

type (
	RegisterRequest struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
	}

	RegisterResponse struct {
		ID domain.ID
	}
)

func (c *Controller) Register(ctx context.Context, req RegisterRequest) (resp *RegisterResponse, err error) {
	ctx, span := c.tracer.Start(ctx, "register")
	defer span.End()

	err = c.db.Transaction(ctx, func(tx application.Repositories) error {
		user, err := domain.NewUser(domain.NewUserRequest{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  req.Password,
		})
		if err != nil {
			return fmt.Errorf("creating user: %w", err)
		}

		if err = tx.User().Save(ctx, user); err != nil {
			return fmt.Errorf("saving user: %w", err)
		}

		resp = &RegisterResponse{
			ID: user.ID,
		}

		return nil
	})
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "registering new user", err)
	}

	slog.Info(ctx, "user registered",
		slog.WithStringer("user_id", resp.ID),
	)

	return
}
