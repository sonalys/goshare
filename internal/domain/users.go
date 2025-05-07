package domain

import (
	"fmt"
	"net/mail"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserCreated struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
	FirstName      string
	LastName       string
}

func NewUser(firstName string, lastName string, email string, password string) (*Event[UserCreated], error) {
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

	return &Event[UserCreated]{
		Topic: TopicUserCreated,
		Data: UserCreated{
			ID:             uuid.New(),
			FirstName:      firstName,
			LastName:       lastName,
			Email:          email,
			HashedPassword: string(hashedPassword),
		},
	}, nil
}
