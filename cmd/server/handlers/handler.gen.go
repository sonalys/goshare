//go:build go1.22

// Package handlers provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	CookieAuthScopes = "cookieAuth.Scopes"
)

// Defines values for ErrorCode.
const (
	EmailPasswordMismatch ErrorCode = "email_password_mismatch"
	Forbidden             ErrorCode = "forbidden"
	InternalError         ErrorCode = "internal_error"
	InvalidField          ErrorCode = "invalid_field"
	InvalidParameter      ErrorCode = "invalid_parameter"
	NotFound              ErrorCode = "not_found"
	RequiredBody          ErrorCode = "required_body"
	RequiredField         ErrorCode = "required_field"
	RequiredHeader        ErrorCode = "required_header"
	RequiredParameter     ErrorCode = "required_parameter"
	Unauthorized          ErrorCode = "unauthorized"
)

// Error defines model for Error.
type Error struct {
	Code ErrorCode `json:"code"`

	// Message Human readable error message
	Message  string         `json:"message"`
	Metadata *ErrorMetadata `json:"metadata,omitempty"`
}

// ErrorCode defines model for Error.Code.
type ErrorCode string

// ErrorMetadata defines model for ErrorMetadata.
type ErrorMetadata struct {
	Field *string `json:"field,omitempty"`
}

// Ledger defines model for Ledger.
type Ledger struct {
	// CreatedAt Date and time the ledger was created
	CreatedAt time.Time `json:"created_at"`

	// CreatedBy User ID of the creator
	CreatedBy openapi_types.UUID `json:"created_by"`
	Id        openapi_types.UUID `json:"id"`

	// Name Ledger's name
	Name string `json:"name"`
}

// LedgerParticipantBalance defines model for LedgerParticipantBalance.
type LedgerParticipantBalance struct {
	Balance float32            `json:"balance"`
	UserId  openapi_types.UUID `json:"user_id"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Errors []Error `json:"errors"`

	// TraceId Unique identifier for the error instance
	TraceId openapi_types.UUID `json:"trace_id"`

	// Url URL of the failed request
	Url string `json:"url"`
}

// LoginJSONBody defines parameters for Login.
type LoginJSONBody struct {
	// Email User's email address
	Email openapi_types.Email `json:"email"`

	// Password User's password
	Password string `json:"password"`
}

// CreateLedgerJSONBody defines parameters for CreateLedger.
type CreateLedgerJSONBody struct {
	// Name Ledger's name
	Name string `json:"name"`
}

// RegisterUserJSONBody defines parameters for RegisterUser.
type RegisterUserJSONBody struct {
	// Email User's email address
	Email openapi_types.Email `json:"email"`

	// FirstName User's first name
	FirstName string `json:"first_name"`

	// LastName User's last name
	LastName string `json:"last_name"`

	// Password User's password
	Password string `json:"password"`
}

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody LoginJSONBody

// CreateLedgerJSONRequestBody defines body for CreateLedger for application/json ContentType.
type CreateLedgerJSONRequestBody CreateLedgerJSONBody

// RegisterUserJSONRequestBody defines body for RegisterUser for application/json ContentType.
type RegisterUserJSONRequestBody RegisterUserJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /authentication/login)
	Login(w http.ResponseWriter, r *http.Request)

	// (GET /authentication/whoami)
	GetIdentity(w http.ResponseWriter, r *http.Request)
	// Healthcheck
	// (GET /healthcheck)
	GetHealthcheck(w http.ResponseWriter, r *http.Request)

	// (GET /ledgers)
	ListLedgers(w http.ResponseWriter, r *http.Request)

	// (POST /ledgers)
	CreateLedger(w http.ResponseWriter, r *http.Request)

	// (GET /ledgers/{ledgerID}/balances)
	ListLedgerBalances(w http.ResponseWriter, r *http.Request, ledgerID openapi_types.UUID)

	// (POST /users)
	RegisterUser(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// Login operation middleware
func (siw *ServerInterfaceWrapper) Login(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Login(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetIdentity operation middleware
func (siw *ServerInterfaceWrapper) GetIdentity(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, CookieAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetIdentity(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetHealthcheck operation middleware
func (siw *ServerInterfaceWrapper) GetHealthcheck(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHealthcheck(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// ListLedgers operation middleware
func (siw *ServerInterfaceWrapper) ListLedgers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, CookieAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListLedgers(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CreateLedger operation middleware
func (siw *ServerInterfaceWrapper) CreateLedger(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, CookieAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateLedger(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// ListLedgerBalances operation middleware
func (siw *ServerInterfaceWrapper) ListLedgerBalances(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "ledgerID" -------------
	var ledgerID openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "ledgerID", r.PathValue("ledgerID"), &ledgerID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "ledgerID", Err: err})
		return
	}

	ctx := r.Context()

	ctx = context.WithValue(ctx, CookieAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListLedgerBalances(w, r, ledgerID)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// RegisterUser operation middleware
func (siw *ServerInterfaceWrapper) RegisterUser(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RegisterUser(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

// ServeMux is an abstraction of http.ServeMux.
type ServeMux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("POST "+options.BaseURL+"/authentication/login", wrapper.Login)
	m.HandleFunc("GET "+options.BaseURL+"/authentication/whoami", wrapper.GetIdentity)
	m.HandleFunc("GET "+options.BaseURL+"/healthcheck", wrapper.GetHealthcheck)
	m.HandleFunc("GET "+options.BaseURL+"/ledgers", wrapper.ListLedgers)
	m.HandleFunc("POST "+options.BaseURL+"/ledgers", wrapper.CreateLedger)
	m.HandleFunc("GET "+options.BaseURL+"/ledgers/{ledgerID}/balances", wrapper.ListLedgerBalances)
	m.HandleFunc("POST "+options.BaseURL+"/users", wrapper.RegisterUser)

	return m
}

type ErrorResponseJSONResponse struct {
	Errors []Error `json:"errors"`

	// TraceId Unique identifier for the error instance
	TraceId openapi_types.UUID `json:"trace_id"`

	// Url URL of the failed request
	Url string `json:"url"`
}

type LoginRequestObject struct {
	Body *LoginJSONRequestBody
}

type LoginResponseObject interface {
	VisitLoginResponse(w http.ResponseWriter) error
}

type Login200ResponseHeaders struct {
	SetCookie string
}

type Login200Response struct {
	Headers Login200ResponseHeaders
}

func (response Login200Response) VisitLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Set-Cookie", fmt.Sprint(response.Headers.SetCookie))
	w.WriteHeader(200)
	return nil
}

type LogindefaultJSONResponse struct {
	Body struct {
		Errors []Error `json:"errors"`

		// TraceId Unique identifier for the error instance
		TraceId openapi_types.UUID `json:"trace_id"`

		// Url URL of the failed request
		Url string `json:"url"`
	}
	StatusCode int
}

func (response LogindefaultJSONResponse) VisitLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetIdentityRequestObject struct {
}

type GetIdentityResponseObject interface {
	VisitGetIdentityResponse(w http.ResponseWriter) error
}

type GetIdentity200JSONResponse struct {
	// Email User's email address
	Email  openapi_types.Email `json:"email"`
	UserId openapi_types.UUID  `json:"user_id"`
}

func (response GetIdentity200JSONResponse) VisitGetIdentityResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetIdentitydefaultJSONResponse struct {
	Body struct {
		Errors []Error `json:"errors"`

		// TraceId Unique identifier for the error instance
		TraceId openapi_types.UUID `json:"trace_id"`

		// Url URL of the failed request
		Url string `json:"url"`
	}
	StatusCode int
}

func (response GetIdentitydefaultJSONResponse) VisitGetIdentityResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetHealthcheckRequestObject struct {
}

type GetHealthcheckResponseObject interface {
	VisitGetHealthcheckResponse(w http.ResponseWriter) error
}

type GetHealthcheck200Response struct {
}

func (response GetHealthcheck200Response) VisitGetHealthcheckResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type GetHealthcheckdefaultJSONResponse struct {
	Body struct {
		Errors []Error `json:"errors"`

		// TraceId Unique identifier for the error instance
		TraceId openapi_types.UUID `json:"trace_id"`

		// Url URL of the failed request
		Url string `json:"url"`
	}
	StatusCode int
}

func (response GetHealthcheckdefaultJSONResponse) VisitGetHealthcheckResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type ListLedgersRequestObject struct {
}

type ListLedgersResponseObject interface {
	VisitListLedgersResponse(w http.ResponseWriter) error
}

type ListLedgers200JSONResponse struct {
	Ledgers []Ledger `json:"ledgers"`
}

func (response ListLedgers200JSONResponse) VisitListLedgersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type ListLedgersdefaultJSONResponse struct {
	Body struct {
		Errors []Error `json:"errors"`

		// TraceId Unique identifier for the error instance
		TraceId openapi_types.UUID `json:"trace_id"`

		// Url URL of the failed request
		Url string `json:"url"`
	}
	StatusCode int
}

func (response ListLedgersdefaultJSONResponse) VisitListLedgersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type CreateLedgerRequestObject struct {
	Body *CreateLedgerJSONRequestBody
}

type CreateLedgerResponseObject interface {
	VisitCreateLedgerResponse(w http.ResponseWriter) error
}

type CreateLedger200JSONResponse struct {
	Id openapi_types.UUID `json:"id"`
}

func (response CreateLedger200JSONResponse) VisitCreateLedgerResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type CreateLedgerdefaultJSONResponse struct {
	Body struct {
		Errors []Error `json:"errors"`

		// TraceId Unique identifier for the error instance
		TraceId openapi_types.UUID `json:"trace_id"`

		// Url URL of the failed request
		Url string `json:"url"`
	}
	StatusCode int
}

func (response CreateLedgerdefaultJSONResponse) VisitCreateLedgerResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type ListLedgerBalancesRequestObject struct {
	LedgerID openapi_types.UUID `json:"ledgerID"`
}

type ListLedgerBalancesResponseObject interface {
	VisitListLedgerBalancesResponse(w http.ResponseWriter) error
}

type ListLedgerBalances200JSONResponse struct {
	Balances []LedgerParticipantBalance `json:"balances"`
}

func (response ListLedgerBalances200JSONResponse) VisitListLedgerBalancesResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type ListLedgerBalancesdefaultJSONResponse struct {
	Body struct {
		Errors []Error `json:"errors"`

		// TraceId Unique identifier for the error instance
		TraceId openapi_types.UUID `json:"trace_id"`

		// Url URL of the failed request
		Url string `json:"url"`
	}
	StatusCode int
}

func (response ListLedgerBalancesdefaultJSONResponse) VisitListLedgerBalancesResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type RegisterUserRequestObject struct {
	Body *RegisterUserJSONRequestBody
}

type RegisterUserResponseObject interface {
	VisitRegisterUserResponse(w http.ResponseWriter) error
}

type RegisterUser200JSONResponse struct {
	Id openapi_types.UUID `json:"id"`
}

func (response RegisterUser200JSONResponse) VisitRegisterUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type RegisterUserdefaultJSONResponse struct {
	Body struct {
		Errors []Error `json:"errors"`

		// TraceId Unique identifier for the error instance
		TraceId openapi_types.UUID `json:"trace_id"`

		// Url URL of the failed request
		Url string `json:"url"`
	}
	StatusCode int
}

func (response RegisterUserdefaultJSONResponse) VisitRegisterUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (POST /authentication/login)
	Login(ctx context.Context, request LoginRequestObject) (LoginResponseObject, error)

	// (GET /authentication/whoami)
	GetIdentity(ctx context.Context, request GetIdentityRequestObject) (GetIdentityResponseObject, error)
	// Healthcheck
	// (GET /healthcheck)
	GetHealthcheck(ctx context.Context, request GetHealthcheckRequestObject) (GetHealthcheckResponseObject, error)

	// (GET /ledgers)
	ListLedgers(ctx context.Context, request ListLedgersRequestObject) (ListLedgersResponseObject, error)

	// (POST /ledgers)
	CreateLedger(ctx context.Context, request CreateLedgerRequestObject) (CreateLedgerResponseObject, error)

	// (GET /ledgers/{ledgerID}/balances)
	ListLedgerBalances(ctx context.Context, request ListLedgerBalancesRequestObject) (ListLedgerBalancesResponseObject, error)

	// (POST /users)
	RegisterUser(ctx context.Context, request RegisterUserRequestObject) (RegisterUserResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// Login operation middleware
func (sh *strictHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequestObject

	var body LoginJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Login(ctx, request.(LoginRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Login")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(LoginResponseObject); ok {
		if err := validResponse.VisitLoginResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetIdentity operation middleware
func (sh *strictHandler) GetIdentity(w http.ResponseWriter, r *http.Request) {
	var request GetIdentityRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetIdentity(ctx, request.(GetIdentityRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetIdentity")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetIdentityResponseObject); ok {
		if err := validResponse.VisitGetIdentityResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetHealthcheck operation middleware
func (sh *strictHandler) GetHealthcheck(w http.ResponseWriter, r *http.Request) {
	var request GetHealthcheckRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetHealthcheck(ctx, request.(GetHealthcheckRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetHealthcheck")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetHealthcheckResponseObject); ok {
		if err := validResponse.VisitGetHealthcheckResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// ListLedgers operation middleware
func (sh *strictHandler) ListLedgers(w http.ResponseWriter, r *http.Request) {
	var request ListLedgersRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.ListLedgers(ctx, request.(ListLedgersRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ListLedgers")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ListLedgersResponseObject); ok {
		if err := validResponse.VisitListLedgersResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// CreateLedger operation middleware
func (sh *strictHandler) CreateLedger(w http.ResponseWriter, r *http.Request) {
	var request CreateLedgerRequestObject

	var body CreateLedgerJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.CreateLedger(ctx, request.(CreateLedgerRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateLedger")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreateLedgerResponseObject); ok {
		if err := validResponse.VisitCreateLedgerResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// ListLedgerBalances operation middleware
func (sh *strictHandler) ListLedgerBalances(w http.ResponseWriter, r *http.Request, ledgerID openapi_types.UUID) {
	var request ListLedgerBalancesRequestObject

	request.LedgerID = ledgerID

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.ListLedgerBalances(ctx, request.(ListLedgerBalancesRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ListLedgerBalances")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ListLedgerBalancesResponseObject); ok {
		if err := validResponse.VisitListLedgerBalancesResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// RegisterUser operation middleware
func (sh *strictHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var request RegisterUserRequestObject

	var body RegisterUserJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.RegisterUser(ctx, request.(RegisterUserRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "RegisterUser")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(RegisterUserResponseObject); ok {
		if err := validResponse.VisitRegisterUserResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
