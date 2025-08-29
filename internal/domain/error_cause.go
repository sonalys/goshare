package domain

type Cause string

const (
	CauseUnknown  = Cause("")
	CauseInvalid  = Cause("invalid")
	CauseRequired = Cause("required")
	CauseNotFound = Cause("not found")
	CauseOverflow = Cause("overflow")
)

func (e Cause) Error() string {
	return string(e)
}

func (e Cause) Is(target error) bool {
	cast, ok := target.(Cause)
	return ok && e == cast
}
