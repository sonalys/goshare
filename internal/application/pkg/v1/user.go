package v1

import (
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

const (
	ErrEmailPasswordMismatch  = StringError("email and/or password mismatch")
	ErrEmailAlreadyRegistered = StringError("email already registered")
)
