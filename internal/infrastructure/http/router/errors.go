package router

import (
	"context"
	"errors"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/validate"
	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) NewError(ctx context.Context, err error) *server.ErrorResponseStatusCode {
	if target := new(ogenerrors.SecurityError); errors.As(err, &target) {
		return newErrorResponse(http.StatusUnauthorized, server.Error{
			Code:    server.ErrorCodeUnauthorized,
			Message: target.Err.Error(),
		})
	}

	if errors.Is(err, application.ErrForbidden) {
		return newErrorResponse(http.StatusForbidden, server.Error{
			Code:    server.ErrorCodeUnauthorized,
			Message: "not authorized to access resource",
		})
	}

	if target := new(validate.Error); errors.As(err, &target) {
		errs := make([]server.Error, 0, len(target.Fields))

		for _, fieldErr := range target.Fields {
			errs = append(errs, server.Error{
				Code:    server.ErrorCodeInvalidField,
				Message: fieldErr.Error.Error(),
				Metadata: server.NewOptErrorMetadata(server.ErrorMetadata{
					Field: server.NewOptString(fieldErr.Name),
				}),
			})
		}

		return newErrorResponse(http.StatusBadRequest, errs...)
	}

	if target := new(domain.FieldErrors); errors.As(err, target) {
		errs := make([]server.Error, 0, len(*target))

		for _, fieldErr := range *target {
			errs = append(errs, server.Error{
				Code:    server.ErrorCodeInvalidField,
				Message: fieldErr.Cause.Error(),
				Metadata: server.NewOptErrorMetadata(server.ErrorMetadata{
					Field: server.NewOptString(fieldErr.Field),
				}),
			})
		}

		return newErrorResponse(http.StatusBadRequest, errs...)
	}

	if errors.Is(err, application.ErrNotFound) {
		return newErrorResponse(http.StatusNotFound, server.Error{
			Code:    server.ErrorCodeNotFound,
			Message: err.Error(),
		})
	}

	if target := new(domain.FieldError); errors.As(err, target) {
		return newErrorResponse(http.StatusBadRequest, server.Error{
			Code:    server.ErrorCodeInvalidField,
			Message: target.Cause.Error(),
			Metadata: server.NewOptErrorMetadata(server.ErrorMetadata{
				Field: server.NewOptString(target.Field),
			}),
		})
	}

	if target := new(validate.InvalidContentTypeError); errors.As(err, &target) {
		return newErrorResponse(http.StatusUnsupportedMediaType, server.Error{
			Code:    server.ErrorCodeRequiredHeader,
			Message: target.Error(),
		})
	}

	if target := new(ogenerrors.DecodeParamError); errors.As(err, &target) {
		return newErrorResponse(http.StatusBadRequest, server.Error{
			Code:    server.ErrorCodeInvalidParameter,
			Message: target.Err.Error(),
			Metadata: server.NewOptErrorMetadata(server.ErrorMetadata{
				Field: server.NewOptString(target.Name),
			}),
		})
	}

	if target := new(ogenerrors.DecodeBodyError); errors.As(err, &target) {
		return newErrorResponse(http.StatusBadRequest, server.Error{
			Code:    server.ErrorCodeInvalidField,
			Message: target.Err.Error(),
			Metadata: server.NewOptErrorMetadata(server.ErrorMetadata{
				Field: server.NewOptString("body"),
			}),
		})
	}

	return newErrorResponse(http.StatusInternalServerError, server.Error{
		Code:    server.ErrorCodeInternalError,
		Message: "internal server error",
	})
}

func newErrorResponse(statusCode int, errs ...server.Error) *server.ErrorResponseStatusCode {
	return &server.ErrorResponseStatusCode{
		StatusCode: statusCode,
		Response: server.ErrorResponse{
			Errors: errs,
		},
	}
}
