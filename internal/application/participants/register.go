package participants

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type (
	Repository interface {
		Create(ctx context.Context, participant *v1.Participant) error
	}

	ParticipantController struct {
		repository Repository
	}
)

func NewParticipantController(repository Repository) *ParticipantController {
	return &ParticipantController{
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

func (c *ParticipantController) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	participant := &v1.Participant{
		ID:              uuid.New(),
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		PasswordHash:    hashPassword(req.Password, req.Email),
		IsEmailVerified: false,
		CreatedAt:       time.Now(),
	}

	if err := c.repository.Create(ctx, participant); err != nil {
		return nil, fmt.Errorf("failed to create participant: %w", err)
	}

	return &RegisterResponse{
		ID: participant.ID,
	}, nil
}
