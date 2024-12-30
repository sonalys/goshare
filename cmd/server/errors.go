package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/oapi-codegen/runtime/types"
	"github.com/sonalys/goshare/cmd/server/handlers"
	"go.opentelemetry.io/otel/trace"
)

// requestErrorHandler is a handler for the openapi request handling errors.
// It does not handle application errors, only errors related to the http request itself.
func requestErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	slog.ErrorContext(ctx, "request error", slog.Any("error", err))

	var resp handlers.ErrorResponse
	var cause handlers.Error

	switch err := err.(type) {
	case *handlers.RequiredHeaderError:
		cause = handlers.Error{
			Code:    handlers.RequiredHeader,
			Message: fmt.Sprintf("missing required header: %s", err.ParamName),
			Url:     r.URL.Path,
			Metadata: &handlers.ErrorMetadata{
				Field: &err.ParamName,
			},
		}
	case *handlers.InvalidParamFormatError:
		cause = handlers.Error{
			Code:    handlers.InvalidParameter,
			Message: fmt.Sprintf("invalid format for parameter: %s", err.ParamName),
			Url:     r.URL.Path,
			Metadata: &handlers.ErrorMetadata{
				Field: &err.ParamName,
			},
		}
	case *handlers.RequiredParamError:
		cause = handlers.Error{
			Code:    handlers.RequiredParameter,
			Message: fmt.Sprintf("missing required parameter: %s", err.ParamName),
			Url:     r.URL.Path,
			Metadata: &handlers.ErrorMetadata{
				Field: &err.ParamName,
			},
		}
	case *handlers.UnmarshalingParamError:
		cause = handlers.Error{
			Code:    handlers.InvalidParameter,
			Message: fmt.Sprintf("failed to unmarshal parameter: %s", err.ParamName),
			Url:     r.URL.Path,
			Metadata: &handlers.ErrorMetadata{
				Field: &err.ParamName,
			},
		}
	case *handlers.TooManyValuesForParamError:
		cause = handlers.Error{
			Code:    handlers.InvalidParameter,
			Message: fmt.Sprintf("too many values for parameter: %s", err.ParamName),
			Url:     r.URL.Path,
			Metadata: &handlers.ErrorMetadata{
				Field: &err.ParamName,
			},
		}
	default:
		switch {
		// Special case where the request body reader is nil.
		// The oapi-codegen handler will not wrap this error.
		case errors.Is(err, io.EOF):
			cause = handlers.Error{
				Code:     handlers.RequiredBody,
				Message:  "missing request body",
				Url:      r.URL.Path,
				Metadata: &handlers.ErrorMetadata{},
			}
		default:
			cause = handlers.Error{
				Code:     handlers.InternalError,
				Message:  "internal server error",
				Url:      r.URL.Path,
				Metadata: &handlers.ErrorMetadata{},
			}
		}
	}

	w.WriteHeader(http.StatusBadRequest)

	cause.TraceId = types.UUID(trace.SpanContextFromContext(ctx).TraceID())
	resp.Errors = append(resp.Errors, cause)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.ErrorContext(ctx, "failed to write response", slog.Any("error", err))
	}
}

// responseErrorHandler is a handler for the openapi response handling errors.
// It happens when the API returns an error, instead of a response.
// This does not include expected errors that are also returned as a response.
func responseErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	slog.ErrorContext(ctx, "response error", slog.Any("error", err))

	var resp handlers.ErrorResponse

	w.WriteHeader(http.StatusInternalServerError)

	cause := handlers.Error{
		TraceId: types.UUID(trace.SpanContextFromContext(ctx).TraceID()),
		Code:    handlers.InternalError,
		Message: http.StatusText(http.StatusInternalServerError),
		Url:     r.URL.Path,
	}
	resp.Errors = append(resp.Errors, cause)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.ErrorContext(ctx, "failed to write response", slog.Any("error", err))
	}
}
