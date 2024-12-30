package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/oapi-codegen/runtime/types"
	"github.com/sonalys/goshare/cmd/server/handlers"
	"go.opentelemetry.io/otel/trace"
)

type errOpt interface {
	apply(*handlers.Error)
}

type metadata handlers.ErrorMetadata

func (m metadata) apply(err *handlers.Error) {
	err.Metadata = (*handlers.ErrorMetadata)(&m)
}

func newError(code handlers.ErrorCode, message string, opt ...errOpt) handlers.Error {
	err := handlers.Error{
		Code:     code,
		Message:  message,
		Metadata: &handlers.ErrorMetadata{},
	}

	for _, o := range opt {
		o.apply(&err)
	}

	return err
}

func newErrorResponse(r *http.Request, cause []handlers.Error) handlers.ErrorResponse {
	return handlers.ErrorResponse{
		TraceId: types.UUID(trace.SpanContextFromContext(r.Context()).TraceID()),
		Url:     r.URL.Path,
		Errors:  cause,
	}
}

// requestErrorHandler is a handler for the openapi request handling errors.
// It does not handle application errors, only errors related to the http request itself.
func requestErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	slog.ErrorContext(ctx, "request error", slog.Any("error", err))

	var cause handlers.Error

	switch err := err.(type) {
	case *handlers.RequiredHeaderError:
		cause = newError(handlers.RequiredHeader, fmt.Sprintf("missing required header: %s", err.ParamName), metadata{Field: &err.ParamName})
	case *handlers.InvalidParamFormatError:
		cause = newError(handlers.InvalidParameter, fmt.Sprintf("invalid format for parameter: %s", err.ParamName), metadata{Field: &err.ParamName})
	case *handlers.RequiredParamError:
		cause = newError(handlers.RequiredParameter, fmt.Sprintf("missing required parameter: %s", err.ParamName), metadata{Field: &err.ParamName})
	case *handlers.UnmarshalingParamError:
		cause = newError(handlers.InvalidParameter, fmt.Sprintf("failed to unmarshal parameter: %s", err.ParamName), metadata{Field: &err.ParamName})
	case *handlers.TooManyValuesForParamError:
		cause = newError(handlers.InvalidParameter, fmt.Sprintf("too many values for parameter: %s", err.ParamName), metadata{Field: &err.ParamName})
	default:
		switch {
		// Special case where the request body reader is nil.
		// The oapi-codegen handler will not wrap this error.
		case errors.Is(err, io.EOF):
			cause = newError(handlers.RequiredBody, "missing request body")
		default:
			cause = newError(handlers.InternalError, "internal server error")
		}
	}

	writeErrorResponse(ctx, w, http.StatusBadRequest, newErrorResponse(r, []handlers.Error{cause}))
}

// responseErrorHandler is a handler for the openapi response handling errors.
// It happens when the API returns an error, instead of a response.
// This does not include expected errors that are also returned as a response.
func responseErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	slog.ErrorContext(ctx, "response error", slog.Any("error", err))

	errList := []handlers.Error{
		{
			Code:    handlers.InternalError,
			Message: http.StatusText(http.StatusInternalServerError),
		},
	}

	writeErrorResponse(ctx, w, http.StatusInternalServerError, newErrorResponse(r, errList))
}
