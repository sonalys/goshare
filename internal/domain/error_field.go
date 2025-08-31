package domain

import (
	"errors"
	"fmt"
)

type (
	FieldErrorMetadata struct {
		Index int
	}

	FieldError struct {
		Cause    error
		Field    string
		Metadata *FieldErrorMetadata
	}

	FieldErrors []FieldError
)

func newRequiredFieldError(field string) FieldError {
	return FieldError{
		Field: field,
		Cause: ErrRequired,
	}
}

func newInvalidFieldError(field string) FieldError {
	return FieldError{
		Field: field,
		Cause: ErrInvalid,
	}
}

func newFieldLengthError(field string, minLength, maxLength int) FieldError {
	return FieldError{
		Field: field,
		Cause: RangeError{Min: minLength, Max: maxLength},
	}
}

func (e FieldError) Error() string {
	return fmt.Sprintf("field '%s': %v", e.Field, e.Cause)
}

func (e FieldError) Is(target error) bool {
	cast, ok := target.(FieldError)

	return ok && e.Field == cast.Field && errors.Is(e.Cause, cast.Cause)
}

func (e FieldError) Unwrap() error {
	return e.Cause
}

func (el FieldErrors) Unwrap() []error {
	errs := make([]error, 0, len(el))
	for _, f := range el {
		errs = append(errs, f)
	}

	return errs
}

func (el FieldErrors) Error() string {
	if len(el) == 1 {
		return el[0].Error()
	}

	return fmt.Sprintf("%v", []FieldError(el))
}
