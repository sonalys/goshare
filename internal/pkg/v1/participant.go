package v1

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	ID              uuid.UUID
	FirstName       string
	LastName        string
	Email           string
	IsEmailVerified bool
	PasswordHash    string
	CreatedAt       time.Time
}

var (
	ErrParticipantEmailAlreadyExists = errors.New("participant with email already exists")
)
