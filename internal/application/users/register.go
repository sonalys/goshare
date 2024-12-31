package users

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
	"go.opentelemetry.io/otel/codes"
)

type (
	Repository interface {
		Create(ctx context.Context, participant *v1.User) error
	}

	UserController struct {
		repository Repository
	}
)

func NewParticipantController(repository Repository) *UserController {
	return &UserController{
		repository: repository,
	}
}

type RegisterRequest struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (r RegisterRequest) Validate() error {
	var errs v1.FormError

	if r.FirstName == "" {
		errs.Fields = append(errs.Fields, v1.NewRequiredFieldError("first_name"))
	}

	if r.LastName == "" {
		errs.Fields = append(errs.Fields, v1.NewRequiredFieldError("last_name"))
	}

	if r.Email == "" {
		errs.Fields = append(errs.Fields, v1.NewRequiredFieldError("email"))
	} else if _, err := mail.ParseAddress(r.Email); err == nil {
		errs.Fields = append(errs.Fields, v1.NewInvalidFieldError("email"))
	}

	if r.Password == "" {
		errs.Fields = append(errs.Fields, v1.NewRequiredFieldError("password"))
	} else if passwordLength := len(r.Password); passwordLength < 8 || passwordLength > 64 {
		errs.Fields = append(errs.Fields, v1.NewFieldLengthError("password", 8, 64))
	}

	return errs.Validate()
}

type RegisterResponse struct {
	ID uuid.UUID
}

func hashPassword(password, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(salt + password))
	return hex.EncodeToString(hash.Sum(nil))
}

func (c *UserController) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "user.register")
	defer span.End()

	user := &v1.User{
		ID:              uuid.New(),
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		PasswordHash:    hashPassword(req.Password, req.Email),
		IsEmailVerified: false,
		CreatedAt:       time.Now(),
	}

	if err := c.repository.Create(ctx, user); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	span.SetStatus(codes.Ok, "")
	slog.InfoContext(ctx, "user registered", slog.String("id", user.ID.String()))

	return &RegisterResponse{
		ID: user.ID,
	}, nil
}
