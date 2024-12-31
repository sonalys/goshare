package users

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
	"go.opentelemetry.io/otel/codes"
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
		errs.Fields = append(errs.Fields, v1.NewRequiredFieldError("email"))
	}

	if r.Password == "" {
		errs.Fields = append(errs.Fields, v1.NewRequiredFieldError("password"))
	}

	return errs.Validate()
}

func (c *Controller) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "user.login")
	defer span.End()

	if err := req.Validate(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "invalid login request", slog.Any("error", err))
		return nil, err
	}
	span.AddEvent("request validated")

	user, err := c.repository.FindByEmail(ctx, req.Email)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "could not find user by email", slog.Any("error", err))
		return nil, v1.ErrEmailPasswordMismatch
	}
	span.AddEvent("user found")

	// Use bcrypt to compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "password hash mismatch", slog.Any("error", err))
		return nil, v1.ErrEmailPasswordMismatch
	}
	span.AddEvent("hash compared")

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  req.Email,
		"userID": user.ID.String(),
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})
	span.AddEvent("jwt generated")

	// Get the secret key from environment variables
	secretKey := "my-secret-key"

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "could not sign JWT token", slog.Any("error", err))
		return nil, fmt.Errorf("failed to sign token: %v", err)
	}
	span.AddEvent("jwt signed")

	span.SetStatus(codes.Ok, "")
	slog.InfoContext(ctx, "user logged in", slog.String("id", user.ID.String()))

	return &LoginResponse{Token: tokenString}, nil
}
