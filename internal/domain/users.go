package domain

import (
	"fmt"
	"net/mail"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		CreatedAt       time.Time
		Email           string
		FirstName       string
		ID              ID
		IsEmailVerified bool
		LastName        string
		PasswordHash    string
		LedgersCount    int32
	}

	NewUserRequest struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
	}
)

const (
	UserMaxLedgers = 5

	ErrUserAlreadyRegistered = StringError("email already in use")
	ErrUserNotFound          = StringError("user not found")
)

func (req *NewUserRequest) validate() error {
	var form Form

	if req.FirstName == "" {
		form.Append(newRequiredFieldError("firstName"))
	}

	if req.LastName == "" {
		form.Append(newRequiredFieldError("lastName"))
	}

	if pwdLen := len(req.Password); pwdLen < 8 || pwdLen > 72 {
		form.Append(newFieldLengthError("password", 8, 72))
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		form.Append(newInvalidFieldError("email"))
	}

	return form.Close()
}

func NewUser(req NewUserRequest) (*User, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash user password: %w", err)
	}

	return &User{
		ID:              NewID(),
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		PasswordHash:    string(hashedPassword),
		IsEmailVerified: false,
		LedgersCount:    0,
		CreatedAt:       time.Now(),
	}, nil
}

func (user *User) CreateLedger(name string) (*Ledger, error) {
	var form Form

	if user.LedgersCount+1 > UserMaxLedgers {
		return nil, UserMaxLedgersError{
			UserID:     user.ID,
			MaxLedgers: UserMaxLedgers,
		}
	}

	if nameLength := len(name); nameLength < 3 || nameLength > 255 {
		form.Append(newFieldLengthError("name", 3, 255))
	}

	if err := form.Close(); err != nil {
		return nil, err
	}

	user.LedgersCount += 1

	return &Ledger{
		ID:   NewID(),
		Name: name,
		Members: map[ID]*LedgerMember{
			user.ID: {
				Balance:   0,
				CreatedAt: time.Now(),
				CreatedBy: user.ID,
			},
		},
		CreatedAt: time.Now(),
		CreatedBy: user.ID,
	}, nil
}
