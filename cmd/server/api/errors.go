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

func extractErrorCauses(err error) ([]v1.FieldError, bool) {
	var resp []v1.FieldError

	if fieldErr := new(v1.FieldError); errors.As(err, fieldErr) {
		resp = []v1.FieldError{*fieldErr}
	} else if errList := new(v1.FieldErrorList); errors.As(err, errList) {
		resp = *errList
	}

	return resp, len(resp) > 0
}

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
	switch cause := from.Cause; {
	case errors.Is(cause, v1.ErrInvalidValue), errors.Is(cause, v1.ErrNotFound):
		return handlers.InvalidField
	case errors.Is(cause, v1.ErrRequiredValue):
		return handlers.RequiredField
	case errors.Is(cause, v1.ErrUserAlreadyMember):
		return handlers.UserAlreadyMember
	case errors.Is(cause, v1.ErrUserNotAMember):
		return handlers.UserNotMember
	default:
		return handlers.ErrorCode("")
	}
}

func newCauseFromFieldError(from v1.FieldError) handlers.Error {
	return handlers.Error{
		Message: from.Error(),
		Code:    getFieldErrorCode(from),
		Metadata: &handlers.ErrorMetadata{
			Field: pointers.New(from.Field),
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
