package v1

import (
	"errors"
	"fmt"
)

type (
	FieldErrorList []FieldError

	FormError struct {
		Fields FieldErrorList
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

func (e *FieldError) Error() string {
	return fmt.Sprintf("field %s: %v", e.Field, e.Cause)
}

func (e *FieldError) Unwrap() error {
	return e.Cause
}

func (el FieldErrorList) Unwrap() []error {
	errs := make([]error, 0, len(el))
	for _, f := range el {
		errs = append(errs, f.Cause)
	}
	return errs
}

func (el FieldErrorList) Error() string {
	return fmt.Sprintf("%v", []FieldError(el))
}

func (e *FormError) Validate() error {
	if len(e.Fields) == 0 {
		return nil
	}
	return e.Fields
}

func (e *ValueRangeError) Error() string {
	return fmt.Sprintf("value must be between %d and %d", e.Min, e.Max)
}

func (e *ValueLengthError) Error() string {
	return fmt.Sprintf("value length must be between %d and %d", e.Min, e.Max)
}
