package domain

import "fmt"

type RangeError struct {
	Max int
	Min int
}

func (e RangeError) Error() string {
	return fmt.Sprintf("value length must be between %d and %d", e.Min, e.Max)
}

func (e RangeError) Is(target error) bool {
	cast, ok := target.(RangeError)
	return ok && cast == e
}
