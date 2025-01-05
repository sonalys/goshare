package v1

import (
	"errors"
)

type Identity struct {
	Email  string
	UserID ID
	Exp    int64
}

var (
	ErrAuthenticationExpired = errors.New("authentication expired")
)
