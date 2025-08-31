package domain

type StringError string

const (
	ErrUnknown  = StringError("")
	ErrInvalid  = StringError("invalid")
	ErrRequired = StringError("required")
	ErrOverflow = StringError("overflow")
	ErrConflict = StringError("conflict")
)

func (e StringError) Error() string {
	return string(e)
}

func (e StringError) Is(target error) bool {
	cast, ok := target.(StringError)

	return ok && e == cast
}
