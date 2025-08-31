package v1

type (
	StringError string

	UserEmailConflictError struct {
		Email string
	}

	UserCredentialsMismatchError struct {
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

func (e *UserEmailConflictError) Error() string {
	return "email already registered"
}

func (e *UserCredentialsMismatchError) Error() string {
	return "email and/or password mismatch"
}

func (s StringError) Error() string {
	return string(s)
}
