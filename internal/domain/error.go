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

	ErrCause string
)

const (
	ErrCauseInvalid  = ErrCause("invalid")
	ErrCauseRequired = ErrCause("required")
)

func newRequiredFieldError(field string) FieldError {
	return FieldError{
		Field: field,
		Cause: ErrCauseRequired,
	}
}

func newInvalidFieldError(field string) FieldError {
	return FieldError{
		Field: field,
		Cause: ErrCauseInvalid,
	}
}

func newFieldLengthError(field string, min, max int) FieldError {
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
		errs = append(errs, f)
	}
	return errs
}

func (el FieldErrorList) Error() string {
	if len(el) == 1 {
		return el[0].Error()
	}
	return fmt.Sprintf("%v", []FieldError(el))
}

func (e *ValueLengthError) Error() string {
	return fmt.Sprintf("value length must be between %d and %d", e.Min, e.Max)
}

func (e ErrCause) Error() string {
	return string(e)
}

func (e *FormError) Close() error {
	if len(*e) == 0 {
		return nil
	}
	return FieldErrorList(*e)
}

func (e *FormError) Append(err FieldError) {
	*e = append(*e, err)
}
