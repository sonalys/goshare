package identitycontroller

import (
	"context"
	"errors"
	"time"

	v1 "github.com/sonalys/goshare/internal/application/v1"
	"github.com/sonalys/goshare/pkg/slog"
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

func (c *controller) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	ctx, span := c.tracer.Start(ctx, "login")
	defer span.End()

	user, err := c.db.User().GetByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, v1.ErrNotFound) {
			return nil, slog.ErrorReturn(ctx, "getting user by email", err)
		}
		slog.Warn(ctx, "could not find user by email", err)

		return nil, &v1.UserCredentialsMismatchError{
			Email: req.Email,
		}
	}
	span.AddEvent("user found")

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		slog.Error(ctx, "password hash mismatch", err)

		return nil, &v1.UserCredentialsMismatchError{
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
		return nil, slog.ErrorReturn(ctx, "signing jwt token", err)
	}

	slog.Info(ctx, "user logged in", slog.WithStringer("user_id", user.ID))

	return &LoginResponse{Token: token}, nil
}
