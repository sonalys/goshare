package users

import (
	"context"
	"fmt"
	"time"

	"github.com/sonalys/goshare/internal/pkg/otel"
	"github.com/sonalys/goshare/internal/pkg/slog"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
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

func (r LoginRequest) Validate() error {
	var errs v1.FormError

	if r.Email == "" {
		errs = append(errs, v1.NewRequiredFieldError("email"))
	}

	if r.Password == "" {
		errs = append(errs, v1.NewRequiredFieldError("password"))
	}

	return errs.Validate()
}

func (c *Controller) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "user.login")
	defer span.End()

	if err := req.Validate(); err != nil {
		return nil, slog.ErrorReturn(ctx, "invalid login request", err)
	}

	user, err := c.repository.FindByEmail(ctx, req.Email)
	if err != nil {
		slog.Error(ctx, "could not find user by email", err)
		return nil, v1.ErrEmailPasswordMismatch
	}
	span.AddEvent("user found")

	// Use bcrypt to compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		slog.Error(ctx, "password hash mismatch", err)
		return nil, v1.ErrEmailPasswordMismatch
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
