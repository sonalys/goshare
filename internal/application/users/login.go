package users

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sonalys/goshare/internal/pkg/otel"
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
		slog.ErrorContext(ctx, "invalid login request", slog.Any("error", err))
		return nil, err
	}

	user, err := c.repository.FindByEmail(ctx, req.Email)
	if err != nil {
		slog.ErrorContext(ctx, "could not find user by email", slog.Any("error", err))
		return nil, v1.ErrEmailPasswordMismatch
	}
	span.AddEvent("user found")

	// Use bcrypt to compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		slog.ErrorContext(ctx, "password hash mismatch", slog.Any("error", err))
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
		slog.ErrorContext(ctx, "could not sign JWT token", slog.Any("error", err))
		return nil, fmt.Errorf("failed to sign token: %v", err)
	}

	slog.InfoContext(ctx, "user logged in", slog.String("user_id", user.ID.String()))

	return &LoginResponse{Token: token}, nil
}
