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
		Field    string
		Cause    error
		Metadata FieldErrorMetadata
	}

	ValueLengthError struct {
		Min int
		Max int
	}

	StringError string
)

const (
	ErrRequiredValue = StringError("cannot be empty")
	ErrInvalidValue  = StringError("invalid value")
	ErrConflict      = StringError("conflict")
	ErrNotFound      = StringError("not found")
	ErrForbidden     = StringError("forbidden")
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
	return fmt.Sprintf("field %s: %v", e.Field, e.Cause)
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
