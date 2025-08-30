package router

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/validate"
	"github.com/sonalys/goshare/internal/application/controllers/identitycontroller"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/router/ledgers"
	"github.com/sonalys/goshare/internal/infrastructure/http/router/users"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
	v1 "github.com/sonalys/goshare/pkg/v1"
	"go.opentelemetry.io/otel/trace"
)

type (
	Router struct {
		server.LedgersHandler
		server.UsersHandler
	}
)

func New(
	identityController *identitycontroller.Controller,
	userController *usercontroller.Controller,
) server.Handler {
	return &Router{
		LedgersHandler: ledgers.New(userController),
		UsersHandler:   users.New(identityController, userController),
	}
}

func newErrorResponse(ctx context.Context, statusCode int, errs ...server.Error) *server.ErrorResponseStatusCode {
	var traceID trace.TraceID
	if span := trace.SpanFromContext(ctx); span != nil {
		traceID = span.SpanContext().TraceID()
	}

	return &server.ErrorResponseStatusCode{
		StatusCode: statusCode,
		Response: server.ErrorResponse{
			TraceID: uuid.UUID(traceID),
			Errors:  errs,
		},
	}
}

func (a *Router) NewError(ctx context.Context, err error) *server.ErrorResponseStatusCode {
	if resp, ok := err.(*server.ErrorResponseStatusCode); ok {
		var traceID trace.TraceID
		if span := trace.SpanFromContext(ctx); span != nil {
			traceID = span.SpanContext().TraceID()
		}
		resp.Response.TraceID = uuid.UUID(traceID)
		return resp
	}

	if target := new(ogenerrors.SecurityError); errors.As(err, &target) {
		return newErrorResponse(ctx, http.StatusUnauthorized, server.Error{
			Code:    server.ErrorCodeUnauthorized,
			Message: target.Err.Error(),
		})
	}

	if errors.Is(err, v1.ErrForbidden) {
		return newErrorResponse(ctx, http.StatusForbidden, server.Error{
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

		return newErrorResponse(ctx, http.StatusBadRequest, errs...)
	}

	if target := new(domain.FieldErrorList); errors.As(err, target) {
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

		return newErrorResponse(ctx, http.StatusBadRequest, errs...)
	}

	if errors.Is(err, v1.ErrNotFound) {
		return newErrorResponse(ctx, http.StatusNotFound, server.Error{
			Code:    server.ErrorCodeNotFound,
			Message: err.Error(),
		})
	}

	if target := new(domain.FieldError); errors.As(err, target) {
		return newErrorResponse(ctx, http.StatusBadRequest, server.Error{
			Code:    server.ErrorCodeInvalidField,
			Message: target.Cause.Error(),
			Metadata: server.NewOptErrorMetadata(server.ErrorMetadata{
				Field: server.NewOptString(target.Field),
			}),
		})
	}

	if target := new(validate.InvalidContentTypeError); errors.As(err, &target) {
		return newErrorResponse(ctx, http.StatusUnsupportedMediaType, server.Error{
			Code:    server.ErrorCodeRequiredHeader,
			Message: target.Error(),
		})
	}

	if target := new(ogenerrors.DecodeParamError); errors.As(err, &target) {
		return newErrorResponse(ctx, http.StatusBadRequest, server.Error{
			Code:    server.ErrorCodeInvalidParameter,
			Message: target.Err.Error(),
			Metadata: server.NewOptErrorMetadata(server.ErrorMetadata{
				Field: server.NewOptString(target.Name),
			}),
		})
	}

	if target := new(ogenerrors.DecodeBodyError); errors.As(err, &target) {
		return newErrorResponse(ctx, http.StatusBadRequest, server.Error{
			Code:    server.ErrorCodeInvalidField,
			Message: target.Err.Error(),
			Metadata: server.NewOptErrorMetadata(server.ErrorMetadata{
				Field: server.NewOptString("body"),
			}),
		})
	}

	return newErrorResponseWithStatusCode(ctx, http.StatusInternalServerError, server.Error{
		Code:    server.ErrorCodeInternalError,
		Message: "internal server error",
	})
}

func newErrorResponseWithStatusCode(ctx context.Context, statusCode int, errs ...server.Error) *server.ErrorResponseStatusCode {
	var traceID trace.TraceID
	if span := trace.SpanFromContext(ctx); span != nil {
		traceID = span.SpanContext().TraceID()
	}

	return &server.ErrorResponseStatusCode{
		StatusCode: statusCode,
		Response: server.ErrorResponse{
			TraceID: uuid.UUID(traceID),
			Errors:  errs,
		},
	}
}
