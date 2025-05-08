package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/sonalys/goshare/internal/application/pkg/otel"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	identityEncoder IdentityEncoder
	subscriber      *Subscriber
	db              Database
}

type (
	LoginRequest struct {
		Email    string
		Password string
	}

	LoginResponse struct {
		Token string
	}

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

func (c *Users) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "user.login")
	defer span.End()

	user, err := c.db.User().FindByEmail(ctx, req.Email)
	if err != nil {
		slog.Error(ctx, "could not find user by email", err)
		return nil, domain.ErrEmailPasswordMismatch
	}
	span.AddEvent("user found")

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		slog.Error(ctx, "password hash mismatch", err)
		return nil, domain.ErrEmailPasswordMismatch
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

func (c *Users) Register(ctx context.Context, req RegisterRequest) (resp *RegisterResponse, err error) {
	ctx, span := otel.Tracer.Start(ctx, "user.register")
	defer span.End()

	err = c.db.Transaction(ctx, func(db Database) error {
		user, err := domain.NewUser(domain.NewUserRequest{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  req.Password,
		})
		if err != nil {
			return fmt.Errorf("creating user: %w", err)
		}

		if err = db.User().Create(ctx, user); err != nil {
			return fmt.Errorf("storing user: %w", err)
		}

		resp = &RegisterResponse{
			ID: user.ID,
		}

		return c.subscriber.Handle(ctx, c.db, user.Events()...)
	})
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "registering new user", err)
	}
	return
}
