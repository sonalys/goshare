package v1

import (
	"errors"
	"fmt"
)

type (
	FormError struct {
		Fields []FieldError
	}

	FieldError struct {
		Field string
		Cause error
	}

	ValueRangeError struct {
		Min int
		Max int
	}

	ValueLengthError struct {
		Min int
		Max int
	}
)

var (
	ErrRequiredValue = errors.New("cannot be empty")
	ErrInvalidValue  = errors.New("invalid value")
)

func (e *FieldError) Error() string {
	return fmt.Sprintf("field %s: %v", e.Field, e.Cause)
}

func (e *FieldError) Unwrap() error {
	return e.Cause
}

func (e *FormError) Error() string {
	return fmt.Sprintf("form error: %v", e.Fields)
}

func (e *FormError) Unwrap() []error {
	errs := make([]error, 0, len(e.Fields))
	for _, f := range e.Fields {
		errs = append(errs, f.Cause)
	}
	return errs
}

func (e *FormError) Validate() error {
	if len(e.Fields) == 0 {
		return nil
	}
	return e
}

func NewRequiredFieldError(field string) FieldError {
	return FieldError{
		Field: field,
		Cause: ErrRequiredValue,
	}
}

func NewInvalidFieldError(field string) FieldError {
	return FieldError{
		Field: field,
		Cause: ErrInvalidValue,
	}
}

func NewFieldRangeError(field string, min, max int) FieldError {
	return FieldError{
		Field: field,
		Cause: &ValueRangeError{Min: min, Max: max},
	}
}

func NewFieldLengthError(field string, min, max int) FieldError {
	return FieldError{
		Field: field,
		Cause: &ValueLengthError{Min: min, Max: max},
	}
}

func (e *ValueRangeError) Error() string {
	return fmt.Sprintf("value must be between %d and %d", e.Min, e.Max)
}

func (e *ValueLengthError) Error() string {
	return fmt.Sprintf("value length must be between %d and %d", e.Min, e.Max)
}
