package domain

import (
	"fmt"
	"net/mail"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	CreatedAt       time.Time
	Email           string
	FirstName       string
	ID              ID
	IsEmailVerified bool
	LastName        string
	PasswordHash    string
}

const (
	ErrEmailAlreadyRegistered = StringError("email already registered")
	ErrEmailPasswordMismatch  = StringError("email and/or password mismatch")
)

func NewUser(firstName string, lastName string, email string, password string) (*Event[User], error) {
	var errs FormError

	if _, err := mail.ParseAddress(email); err != nil {
		errs = append(errs, NewInvalidFieldError("email"))
	}

	if firstName == "" {
		errs = append(errs, NewRequiredFieldError("firstName"))
	}

	if lastName == "" {
		errs = append(errs, NewRequiredFieldError("lastName"))
	}

	if pwdLen := len(password); pwdLen < 8 || pwdLen > 72 {
		errs = append(errs, NewFieldLengthError("password", 8, 72))
	}

	if err := errs.Validate(); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash user password: %w", err)
	}

	return &Event[User]{
		Topic: TopicUserCreated,
		Data: User{
			ID:              NewID(),
			FirstName:       firstName,
			LastName:        lastName,
			Email:           email,
			PasswordHash:    string(hashedPassword),
			IsEmailVerified: false,
			CreatedAt:       time.Now(),
		},
	}, nil
}
