package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/oapi-codegen/runtime/types"
	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/pkg/pointers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
	"go.opentelemetry.io/otel/trace"
)

func WriteErrorResponse(ctx context.Context, w http.ResponseWriter, code int, resp handlers.ErrorResponse) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/problem+json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.ErrorContext(ctx, "failed to write response", slog.Any("error", err))
	}
}

func newErrorResponse(ctx context.Context, cause []handlers.Error) handlers.ErrorResponse {
	return handlers.ErrorResponse{
		TraceId: types.UUID(trace.SpanContextFromContext(ctx).TraceID()),
		Url:     getURL(ctx),
		Errors:  cause,
	}
}

func getFieldErrorCode(from v1.FieldError) handlers.ErrorCode {
	switch {
	case errors.Is(from.Cause, v1.ErrInvalidValue):
		return handlers.InvalidField
	case errors.Is(from.Cause, v1.ErrRequiredValue):
		return handlers.RequiredField
	default:
		return handlers.ErrorCode("")
	}
}

func newCauseFromFieldError(from v1.FieldError) handlers.Error {
	return handlers.Error{
		Message: from.Error(),
		Code:    getFieldErrorCode(from),
		Metadata: &handlers.ErrorMetadata{
			Field: pointers.From(from.Field),
		},
	}
}

func getCausesFromFieldErrors(from v1.FieldErrorList) []handlers.Error {
	resp := make([]handlers.Error, 0, len(from))

	for i := range from {
		resp = append(resp, newCauseFromFieldError(from[i]))
	}

	return resp
}
