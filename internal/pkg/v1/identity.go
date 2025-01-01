package v1

import (
	"errors"

	"github.com/google/uuid"
)

type Identity struct {
	Email  string
	UserID uuid.UUID
	Exp    int64
}

var (
	ErrAuthorizationExpired = errors.New("authorization expired")
)
