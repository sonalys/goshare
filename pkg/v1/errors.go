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
	ErrForbidden = StringError("forbidden")
	ErrNotFound  = StringError("not found")
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
