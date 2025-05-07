package domain

import (
	"fmt"
)

type (
	FieldErrorList []FieldError

	FormError FieldErrorList

	FieldErrorMetadata struct {
		Index int
	}

	FieldError struct {
		Cause    error
		Field    string
		Metadata FieldErrorMetadata
	}

	ValueLengthError struct {
		Max int
		Min int
	}

	StringError string
)

const (
	ErrConflict      = StringError("conflict")
	ErrForbidden     = StringError("forbidden")
	ErrInvalidValue  = StringError("invalid value")
	ErrNotFound      = StringError("not found")
	ErrRequiredValue = StringError("cannot be empty")
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

func NewFieldLengthError(field string, min, max int) FieldError {
	return FieldError{
		Field: field,
		Cause: &ValueLengthError{Min: min, Max: max},
	}
}

func (e FieldError) Error() string {
	return fmt.Sprintf("field '%s': %v", e.Field, e.Cause)
}

func (e FieldError) Unwrap() error {
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
	if len(el) == 1 {
		return el[0].Error()
	}
	return fmt.Sprintf("%v", []FieldError(el))
}

func (e *FormError) Validate() error {
	if len(*e) == 0 {
		return nil
	}
	return FieldErrorList(*e)
}

func (e *ValueLengthError) Error() string {
	return fmt.Sprintf("value length must be between %d and %d", e.Min, e.Max)
}

func (e StringError) Error() string {
	return string(e)
}
