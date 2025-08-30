package identitycontroller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sonalys/goshare/pkg/slog"
	v1 "github.com/sonalys/goshare/pkg/v1"
	"golang.org/x/crypto/bcrypt"
)

type (
	LoginRequest struct {
		Email    string
		Password string
	}

	LoginResponse struct {
		Token string
	}
)

func (c *Controller) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	ctx, span := c.tracer.Start(ctx, "login")
	defer span.End()

	user, err := c.db.User().GetByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, v1.ErrNotFound) {
			return nil, err
		}
		slog.Error(ctx, "could not find user by email", err)
		return nil, &v1.ErrUserCredentialsMismatch{
			Email: req.Email,
		}
	}
	span.AddEvent("user found")

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		slog.Error(ctx, "password hash mismatch", err)
		return nil, &v1.ErrUserCredentialsMismatch{
			Email: req.Email,
		}
	}
	span.AddEvent("hash compared")

	identity := &v1.Identity{
		Email:  user.Email,
		UserID: user.ID,
		Exp:    time.Now().Add(72 * time.Hour).Unix(),
	}
	token, err := c.identityEncoder.Encode(identity)
	if err != nil {
		slog.Error(ctx, "could not sign JWT token", err)
		return nil, fmt.Errorf("failed to sign token: %v", err)
	}

	slog.Info(ctx, "user logged in", slog.WithStringer("user_id", user.ID))

	return &LoginResponse{Token: token}, nil
}
