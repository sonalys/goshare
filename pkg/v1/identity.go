package v1

import "github.com/sonalys/goshare/internal/domain"

type Identity struct {
	Email  string
	UserID domain.ID
	Exp    int64
}

const (
	ErrAuthenticationExpired = domain.ErrorString("authentication expired")
)
