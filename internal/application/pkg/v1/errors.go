package v1

type (
	StringError string

	ErrUserEmailConflict struct {
		Email string
	}

	ErrUserCredentialsMismatch struct {
		Email string
	}
)

const (
	ErrConflict      = StringError("conflict")
	ErrForbidden     = StringError("forbidden")
	ErrInvalidValue  = StringError("invalid value")
	ErrNotFound      = StringError("not found")
	ErrRequiredValue = StringError("cannot be empty")
)

func (e *ErrUserEmailConflict) Error() string {
	return "email already registered"
}

func (e *ErrUserCredentialsMismatch) Error() string {
	return "email and/or password mismatch"
}

func (s StringError) Error() string {
	return string(s)
}
