package v1

import (
	"errors"
	"time"
)

type User struct {
	ID              ID
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
