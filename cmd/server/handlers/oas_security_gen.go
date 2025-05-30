// Code generated by ogen, DO NOT EDIT.

package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/ogenerrors"
)

// SecurityHandler is handler for security parameters.
type SecurityHandler interface {
	// HandleCookieAuth handles cookieAuth security.
	HandleCookieAuth(ctx context.Context, operationName OperationName, t CookieAuth) (context.Context, error)
}

func findAuthorization(h http.Header, prefix string) (string, bool) {
	v, ok := h["Authorization"]
	if !ok {
		return "", false
	}
	for _, vv := range v {
		scheme, value, ok := strings.Cut(vv, " ")
		if !ok || !strings.EqualFold(scheme, prefix) {
			continue
		}
		return value, true
	}
	return "", false
}

func (s *Server) securityCookieAuth(ctx context.Context, operationName OperationName, req *http.Request) (context.Context, bool, error) {
	var t CookieAuth
	const parameterName = "SESSIONID"
	var value string
	switch cookie, err := req.Cookie(parameterName); {
	case err == nil: // if NO error
		value = cookie.Value
	case errors.Is(err, http.ErrNoCookie):
		return ctx, false, nil
	default:
		return nil, false, errors.Wrap(err, "get cookie value")
	}
	t.APIKey = value
	rctx, err := s.sec.HandleCookieAuth(ctx, operationName, t)
	if errors.Is(err, ogenerrors.ErrSkipServerSecurity) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return rctx, true, err
}

// SecuritySource is provider of security values (tokens, passwords, etc.).
type SecuritySource interface {
	// CookieAuth provides cookieAuth security value.
	CookieAuth(ctx context.Context, operationName OperationName) (CookieAuth, error)
}

func (s *Client) securityCookieAuth(ctx context.Context, operationName OperationName, req *http.Request) error {
	t, err := s.sec.CookieAuth(ctx, operationName)
	if err != nil {
		return errors.Wrap(err, "security source \"CookieAuth\"")
	}
	req.AddCookie(&http.Cookie{
		Name:  "SESSIONID",
		Value: t.APIKey,
	})
	return nil
}
