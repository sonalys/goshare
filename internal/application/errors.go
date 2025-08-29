package application

type (
	ErrorString string
)

const (
	ErrUnauthorized = ErrorString("unauthorized")
)

func (c ErrorString) Error() string {
	return string(c)
}
