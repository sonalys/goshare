package domain

type ErrorString string

const (
	ErrUnknown  = ErrorString("")
	ErrInvalid  = ErrorString("invalid")
	ErrRequired = ErrorString("required")
	ErrOverflow = ErrorString("overflow")
)

func (e ErrorString) Error() string {
	return string(e)
}

func (e ErrorString) Is(target error) bool {
	cast, ok := target.(ErrorString)
	return ok && e == cast
}
