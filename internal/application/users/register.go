package users

import (
	"context"
	"fmt"
	"log/slog"
	"net/mail"
	"time"

	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (r RegisterRequest) Validate() error {
	var errs v1.FormError

	if r.FirstName == "" {
		errs = append(errs, v1.NewRequiredFieldError("first_name"))
	}

	if r.LastName == "" {
		errs = append(errs, v1.NewRequiredFieldError("last_name"))
	}

	if r.Email == "" {
		errs = append(errs, v1.NewRequiredFieldError("email"))
	} else if _, err := mail.ParseAddress(r.Email); err == nil {
		errs = append(errs, v1.NewInvalidFieldError("email"))
	}

	if r.Password == "" {
		errs = append(errs, v1.NewRequiredFieldError("password"))
	} else if passwordLength := len(r.Password); passwordLength < 8 || passwordLength > 64 {
		errs = append(errs, v1.NewFieldLengthError("password", 8, 64))
	}

	return errs.Validate()
}

type RegisterResponse struct {
	ID v1.ID
}

func (c *Controller) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "user.register")
	defer span.End()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed to hash password", slog.Any("error", err))
		return nil, err
	}

	user := &v1.User{
		ID:              v1.NewID(),
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		PasswordHash:    string(hashedPassword),
		IsEmailVerified: false,
		CreatedAt:       time.Now(),
	}

	if err := c.repository.Create(ctx, user); err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, err.Error())
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	span.SetStatus(codes.Ok, "")
	slog.InfoContext(ctx, "user registered", slog.String("user_id", user.ID.String()))

	return &RegisterResponse{
		ID: user.ID,
	}, nil
}
