package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/validate"
	"github.com/sonalys/goshare/cmd/server/handlers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
	"go.opentelemetry.io/otel/trace"
)

type (
	API struct {
		handlers.UnimplementedHandler
		dependencies Dependencies
	}
)

func New(dependencies Dependencies) *API {
	return &API{
		dependencies: dependencies,
	}
}

func newErrorResponse(ctx context.Context, statusCode int, errs ...handlers.Error) *handlers.ErrorResponseStatusCode {
	var traceID trace.TraceID
	if span := trace.SpanFromContext(ctx); span != nil {
		traceID = span.SpanContext().TraceID()
	}

	return &handlers.ErrorResponseStatusCode{
		StatusCode: statusCode,
		Response: handlers.ErrorResponse{
			TraceID: uuid.UUID(traceID),
			Errors:  errs,
		},
	}
}

func (a *API) NewError(ctx context.Context, err error) *handlers.ErrorResponseStatusCode {
	if target := new(validate.Error); errors.As(err, &target) {
		errs := make([]handlers.Error, 0, len(target.Fields))

		for _, fieldErr := range target.Fields {
			errs = append(errs, handlers.Error{
				Code:    handlers.ErrorCodeInvalidField,
				Message: fieldErr.Error.Error(),
				Metadata: handlers.NewOptErrorMetadata(handlers.ErrorMetadata{
					Field: handlers.NewOptString(fieldErr.Name),
				}),
			})
		}

		return newErrorResponse(ctx, http.StatusBadRequest, errs...)
	}

	if target := new(v1.FieldErrorList); errors.As(err, target) {
		errs := make([]handlers.Error, 0, len(*target))

		for _, fieldErr := range *target {
			errs = append(errs, handlers.Error{
				Code:    handlers.ErrorCodeInvalidField,
				Message: fieldErr.Cause.Error(),
				Metadata: handlers.NewOptErrorMetadata(handlers.ErrorMetadata{
					Field: handlers.NewOptString(fieldErr.Field),
				}),
			})
		}

		return newErrorResponse(ctx, http.StatusBadRequest, errs...)
	}

	if errors.Is(err, v1.ErrNotFound) {
		return newErrorResponse(ctx, http.StatusNotFound, handlers.Error{
			Code:    handlers.ErrorCodeNotFound,
			Message: err.Error(),
		})
	}

	if target := new(v1.FieldError); errors.As(err, target) {
		return newErrorResponse(ctx, http.StatusBadRequest, handlers.Error{
			Code:    handlers.ErrorCodeInvalidField,
			Message: target.Cause.Error(),
			Metadata: handlers.NewOptErrorMetadata(handlers.ErrorMetadata{
				Field: handlers.NewOptString(target.Field),
			}),
		})
	}

	if target := new(validate.InvalidContentTypeError); errors.As(err, &target) {
		return newErrorResponse(ctx, http.StatusUnsupportedMediaType, handlers.Error{
			Code:    handlers.ErrorCodeRequiredHeader,
			Message: target.Error(),
		})
	}

	if target := new(ogenerrors.DecodeParamError); errors.As(err, &target) {
		return newErrorResponse(ctx, http.StatusBadRequest, handlers.Error{
			Code:    handlers.ErrorCodeInvalidParameter,
			Message: target.Err.Error(),
			Metadata: handlers.NewOptErrorMetadata(handlers.ErrorMetadata{
				Field: handlers.NewOptString(target.Name),
			}),
		})
	}

	if target := new(ogenerrors.DecodeBodyError); errors.As(err, &target) {
		return newErrorResponse(ctx, http.StatusBadRequest, handlers.Error{
			Code:    handlers.ErrorCodeInvalidField,
			Message: target.Err.Error(),
			Metadata: handlers.NewOptErrorMetadata(handlers.ErrorMetadata{
				Field: handlers.NewOptString("body"),
			}),
		})
	}

	return newErrorResponseWithStatusCode(ctx, http.StatusInternalServerError, handlers.Error{
		Code:    handlers.ErrorCodeInternalError,
		Message: "internal server error",
	})
}

func newErrorResponseWithStatusCode(ctx context.Context, statusCode int, errs ...handlers.Error) *handlers.ErrorResponseStatusCode {
	var traceID trace.TraceID
	if span := trace.SpanFromContext(ctx); span != nil {
		traceID = span.SpanContext().TraceID()
	}

	return &handlers.ErrorResponseStatusCode{
		StatusCode: statusCode,
		Response: handlers.ErrorResponse{
			TraceID: uuid.UUID(traceID),
			Errors:  errs,
		},
	}
}
