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
		LedgersCount    int64

		events []Event
	}
)

const (
	ErrEmailAlreadyRegistered = StringError("email already registered")
	ErrEmailPasswordMismatch  = StringError("email and/or password mismatch")
)

type NewUserRequest struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (req *NewUserRequest) validate() error {
	var errs FormError

	if _, err := mail.ParseAddress(req.Email); err != nil {
		errs = append(errs, NewInvalidFieldError("email"))
	}

	if req.FirstName == "" {
		errs = append(errs, NewRequiredFieldError("firstName"))
	}

	if req.LastName == "" {
		errs = append(errs, NewRequiredFieldError("lastName"))
	}

	if pwdLen := len(req.Password); pwdLen < 8 || pwdLen > 72 {
		errs = append(errs, NewFieldLengthError("password", 8, 72))
	}

	return errs.Validate()
}

func NewUser(req NewUserRequest) (*User, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash user password: %w", err)
	}

	user := User{
		ID:              NewID(),
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		PasswordHash:    string(hashedPassword),
		IsEmailVerified: false,
		LedgersCount:    0,
		CreatedAt:       time.Now(),
	}

	user.events = append(user.events, event[User]{
		topic: TopicUserCreated,
		data:  user,
	})

	return &user, nil
}

func (user *User) Events() []Event {
	return user.events
}

func (user *User) CreateLedger(name string) (*Ledger, error) {
	var errs FormError

	if user.LedgersCount+1 > UserMaxLedgers {
		return nil, ErrUserMaxLedgers
	}

	if nameLength := len(name); nameLength < 3 || nameLength > 255 {
		errs = append(errs, NewFieldLengthError("name", 3, 255))
	}

	if err := errs.Validate(); err != nil {
		return nil, err
	}

	ledger := Ledger{
		ID:   NewID(),
		Name: name,
		Participants: []LedgerParticipant{
			{
				ID:        NewID(),
				Identity:  user.ID,
				Balance:   0,
				CreatedAt: time.Now(),
				CreatedBy: user.ID,
			},
		},
		CreatedAt: time.Now(),
		CreatedBy: user.ID,
	}

	user.LedgersCount += 1

	ledger.events = append(ledger.events, event[Ledger]{
		topic: TopicLedgerCreated,
		data:  ledger,
	})

	return &ledger, nil
}
