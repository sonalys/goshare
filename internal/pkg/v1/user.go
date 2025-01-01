package v1

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID
	FirstName       string
	LastName        string
	Email           string
	IsEmailVerified bool
	PasswordHash    string
	CreatedAt       time.Time
}

var (
	ErrEmailPasswordMismatch  = errors.New("email and/or password mismatch")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
)
