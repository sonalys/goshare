// Code generated by ogen, DO NOT EDIT.

package handlers

import (
	"fmt"
	"time"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

func (s *ErrorResponseStatusCode) Error() string {
	return fmt.Sprintf("code %d: %+v", s.StatusCode, s.Response)
}

// AuthenticationLoginOK is response for AuthenticationLogin operation.
type AuthenticationLoginOK struct {
	SetCookie OptString
}

// GetSetCookie returns the value of SetCookie.
func (s *AuthenticationLoginOK) GetSetCookie() OptString {
	return s.SetCookie
}

// SetSetCookie sets the value of SetCookie.
func (s *AuthenticationLoginOK) SetSetCookie(val OptString) {
	s.SetCookie = val
}

type AuthenticationLoginReq struct {
	// User's email address.
	Email string `json:"email"`
	// User's password.
	Password string `json:"password"`
}

// GetEmail returns the value of Email.
func (s *AuthenticationLoginReq) GetEmail() string {
	return s.Email
}

// GetPassword returns the value of Password.
func (s *AuthenticationLoginReq) GetPassword() string {
	return s.Password
}

// SetEmail sets the value of Email.
func (s *AuthenticationLoginReq) SetEmail(val string) {
	s.Email = val
}

// SetPassword sets the value of Password.
func (s *AuthenticationLoginReq) SetPassword(val string) {
	s.Password = val
}

type AuthenticationWhoAmIOK struct {
	// User's email address.
	Email  string    `json:"email"`
	UserID uuid.UUID `json:"user_id"`
}

// GetEmail returns the value of Email.
func (s *AuthenticationWhoAmIOK) GetEmail() string {
	return s.Email
}

// GetUserID returns the value of UserID.
func (s *AuthenticationWhoAmIOK) GetUserID() uuid.UUID {
	return s.UserID
}

// SetEmail sets the value of Email.
func (s *AuthenticationWhoAmIOK) SetEmail(val string) {
	s.Email = val
}

// SetUserID sets the value of UserID.
func (s *AuthenticationWhoAmIOK) SetUserID(val uuid.UUID) {
	s.UserID = val
}

type CookieAuth struct {
	APIKey string
}

// GetAPIKey returns the value of APIKey.
func (s *CookieAuth) GetAPIKey() string {
	return s.APIKey
}

// SetAPIKey sets the value of APIKey.
func (s *CookieAuth) SetAPIKey(val string) {
	s.APIKey = val
}

// Ref: #/components/schemas/Error
type Error struct {
	Code ErrorCode `json:"code"`
	// Human readable error message.
	Message  string           `json:"message"`
	Metadata OptErrorMetadata `json:"metadata"`
}

// GetCode returns the value of Code.
func (s *Error) GetCode() ErrorCode {
	return s.Code
}

// GetMessage returns the value of Message.
func (s *Error) GetMessage() string {
	return s.Message
}

// GetMetadata returns the value of Metadata.
func (s *Error) GetMetadata() OptErrorMetadata {
	return s.Metadata
}

// SetCode sets the value of Code.
func (s *Error) SetCode(val ErrorCode) {
	s.Code = val
}

// SetMessage sets the value of Message.
func (s *Error) SetMessage(val string) {
	s.Message = val
}

// SetMetadata sets the value of Metadata.
func (s *Error) SetMetadata(val OptErrorMetadata) {
	s.Metadata = val
}

type ErrorCode string

const (
	ErrorCodeNotFound              ErrorCode = "not_found"
	ErrorCodeRequiredHeader        ErrorCode = "required_header"
	ErrorCodeRequiredParameter     ErrorCode = "required_parameter"
	ErrorCodeInvalidParameter      ErrorCode = "invalid_parameter"
	ErrorCodeInternalError         ErrorCode = "internal_error"
	ErrorCodeRequiredBody          ErrorCode = "required_body"
	ErrorCodeInvalidField          ErrorCode = "invalid_field"
	ErrorCodeRequiredField         ErrorCode = "required_field"
	ErrorCodeEmailPasswordMismatch ErrorCode = "email_password_mismatch"
	ErrorCodeForbidden             ErrorCode = "forbidden"
	ErrorCodeUnauthorized          ErrorCode = "unauthorized"
	ErrorCodeLedgerMaxUsers        ErrorCode = "ledger_max_users"
	ErrorCodeUserMaxLedgers        ErrorCode = "user_max_ledgers"
	ErrorCodeAuthenticationExpired ErrorCode = "authentication_expired"
	ErrorCodeUserAlreadyMember     ErrorCode = "user_already_member"
	ErrorCodeUserNotMember         ErrorCode = "user_not_member"
)

// AllValues returns all ErrorCode values.
func (ErrorCode) AllValues() []ErrorCode {
	return []ErrorCode{
		ErrorCodeNotFound,
		ErrorCodeRequiredHeader,
		ErrorCodeRequiredParameter,
		ErrorCodeInvalidParameter,
		ErrorCodeInternalError,
		ErrorCodeRequiredBody,
		ErrorCodeInvalidField,
		ErrorCodeRequiredField,
		ErrorCodeEmailPasswordMismatch,
		ErrorCodeForbidden,
		ErrorCodeUnauthorized,
		ErrorCodeLedgerMaxUsers,
		ErrorCodeUserMaxLedgers,
		ErrorCodeAuthenticationExpired,
		ErrorCodeUserAlreadyMember,
		ErrorCodeUserNotMember,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s ErrorCode) MarshalText() ([]byte, error) {
	switch s {
	case ErrorCodeNotFound:
		return []byte(s), nil
	case ErrorCodeRequiredHeader:
		return []byte(s), nil
	case ErrorCodeRequiredParameter:
		return []byte(s), nil
	case ErrorCodeInvalidParameter:
		return []byte(s), nil
	case ErrorCodeInternalError:
		return []byte(s), nil
	case ErrorCodeRequiredBody:
		return []byte(s), nil
	case ErrorCodeInvalidField:
		return []byte(s), nil
	case ErrorCodeRequiredField:
		return []byte(s), nil
	case ErrorCodeEmailPasswordMismatch:
		return []byte(s), nil
	case ErrorCodeForbidden:
		return []byte(s), nil
	case ErrorCodeUnauthorized:
		return []byte(s), nil
	case ErrorCodeLedgerMaxUsers:
		return []byte(s), nil
	case ErrorCodeUserMaxLedgers:
		return []byte(s), nil
	case ErrorCodeAuthenticationExpired:
		return []byte(s), nil
	case ErrorCodeUserAlreadyMember:
		return []byte(s), nil
	case ErrorCodeUserNotMember:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *ErrorCode) UnmarshalText(data []byte) error {
	switch ErrorCode(data) {
	case ErrorCodeNotFound:
		*s = ErrorCodeNotFound
		return nil
	case ErrorCodeRequiredHeader:
		*s = ErrorCodeRequiredHeader
		return nil
	case ErrorCodeRequiredParameter:
		*s = ErrorCodeRequiredParameter
		return nil
	case ErrorCodeInvalidParameter:
		*s = ErrorCodeInvalidParameter
		return nil
	case ErrorCodeInternalError:
		*s = ErrorCodeInternalError
		return nil
	case ErrorCodeRequiredBody:
		*s = ErrorCodeRequiredBody
		return nil
	case ErrorCodeInvalidField:
		*s = ErrorCodeInvalidField
		return nil
	case ErrorCodeRequiredField:
		*s = ErrorCodeRequiredField
		return nil
	case ErrorCodeEmailPasswordMismatch:
		*s = ErrorCodeEmailPasswordMismatch
		return nil
	case ErrorCodeForbidden:
		*s = ErrorCodeForbidden
		return nil
	case ErrorCodeUnauthorized:
		*s = ErrorCodeUnauthorized
		return nil
	case ErrorCodeLedgerMaxUsers:
		*s = ErrorCodeLedgerMaxUsers
		return nil
	case ErrorCodeUserMaxLedgers:
		*s = ErrorCodeUserMaxLedgers
		return nil
	case ErrorCodeAuthenticationExpired:
		*s = ErrorCodeAuthenticationExpired
		return nil
	case ErrorCodeUserAlreadyMember:
		*s = ErrorCodeUserAlreadyMember
		return nil
	case ErrorCodeUserNotMember:
		*s = ErrorCodeUserNotMember
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Ref: #/components/schemas/ErrorMetadata
type ErrorMetadata struct {
	Field OptString `json:"field"`
}

// GetField returns the value of Field.
func (s *ErrorMetadata) GetField() OptString {
	return s.Field
}

// SetField sets the value of Field.
func (s *ErrorMetadata) SetField(val OptString) {
	s.Field = val
}

type ErrorResponse struct {
	// Unique identifier for the error instance.
	TraceID uuid.UUID `json:"trace_id"`
	Errors  []Error   `json:"errors"`
}

// GetTraceID returns the value of TraceID.
func (s *ErrorResponse) GetTraceID() uuid.UUID {
	return s.TraceID
}

// GetErrors returns the value of Errors.
func (s *ErrorResponse) GetErrors() []Error {
	return s.Errors
}

// SetTraceID sets the value of TraceID.
func (s *ErrorResponse) SetTraceID(val uuid.UUID) {
	s.TraceID = val
}

// SetErrors sets the value of Errors.
func (s *ErrorResponse) SetErrors(val []Error) {
	s.Errors = val
}

// ErrorResponseStatusCode wraps ErrorResponse with StatusCode.
type ErrorResponseStatusCode struct {
	StatusCode int
	Response   ErrorResponse
}

// GetStatusCode returns the value of StatusCode.
func (s *ErrorResponseStatusCode) GetStatusCode() int {
	return s.StatusCode
}

// GetResponse returns the value of Response.
func (s *ErrorResponseStatusCode) GetResponse() ErrorResponse {
	return s.Response
}

// SetStatusCode sets the value of StatusCode.
func (s *ErrorResponseStatusCode) SetStatusCode(val int) {
	s.StatusCode = val
}

// SetResponse sets the value of Response.
func (s *ErrorResponseStatusCode) SetResponse(val ErrorResponse) {
	s.Response = val
}

// Ref: #/components/schemas/Expense
type Expense struct {
	// Unique identifier for the expense.
	ID OptUUID `json:"id"`
	// Name of the expense.
	Name string `json:"name"`
	// Date and time of the expense.
	ExpenseDate time.Time       `json:"expense_date"`
	Records     []ExpenseRecord `json:"records"`
}

// GetID returns the value of ID.
func (s *Expense) GetID() OptUUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *Expense) GetName() string {
	return s.Name
}

// GetExpenseDate returns the value of ExpenseDate.
func (s *Expense) GetExpenseDate() time.Time {
	return s.ExpenseDate
}

// GetRecords returns the value of Records.
func (s *Expense) GetRecords() []ExpenseRecord {
	return s.Records
}

// SetID sets the value of ID.
func (s *Expense) SetID(val OptUUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *Expense) SetName(val string) {
	s.Name = val
}

// SetExpenseDate sets the value of ExpenseDate.
func (s *Expense) SetExpenseDate(val time.Time) {
	s.ExpenseDate = val
}

// SetRecords sets the value of Records.
func (s *Expense) SetRecords(val []ExpenseRecord) {
	s.Records = val
}

// Ref: #/components/schemas/ExpenseRecord
type ExpenseRecord struct {
	// Unique identifier for the record.
	ID OptUUID `json:"id"`
	// Type of the record.
	Type ExpenseRecordType `json:"type"`
	// User ID of the person who owes money.
	FromUserID uuid.UUID `json:"from_user_id"`
	// User ID of the person who is owed money.
	ToUserID uuid.UUID `json:"to_user_id"`
	// Amount of money involved in the transaction.
	Amount int32 `json:"amount"`
}

// GetID returns the value of ID.
func (s *ExpenseRecord) GetID() OptUUID {
	return s.ID
}

// GetType returns the value of Type.
func (s *ExpenseRecord) GetType() ExpenseRecordType {
	return s.Type
}

// GetFromUserID returns the value of FromUserID.
func (s *ExpenseRecord) GetFromUserID() uuid.UUID {
	return s.FromUserID
}

// GetToUserID returns the value of ToUserID.
func (s *ExpenseRecord) GetToUserID() uuid.UUID {
	return s.ToUserID
}

// GetAmount returns the value of Amount.
func (s *ExpenseRecord) GetAmount() int32 {
	return s.Amount
}

// SetID sets the value of ID.
func (s *ExpenseRecord) SetID(val OptUUID) {
	s.ID = val
}

// SetType sets the value of Type.
func (s *ExpenseRecord) SetType(val ExpenseRecordType) {
	s.Type = val
}

// SetFromUserID sets the value of FromUserID.
func (s *ExpenseRecord) SetFromUserID(val uuid.UUID) {
	s.FromUserID = val
}

// SetToUserID sets the value of ToUserID.
func (s *ExpenseRecord) SetToUserID(val uuid.UUID) {
	s.ToUserID = val
}

// SetAmount sets the value of Amount.
func (s *ExpenseRecord) SetAmount(val int32) {
	s.Amount = val
}

// Type of the record.
type ExpenseRecordType string

const (
	ExpenseRecordTypeDebt       ExpenseRecordType = "debt"
	ExpenseRecordTypeSettlement ExpenseRecordType = "settlement"
)

// AllValues returns all ExpenseRecordType values.
func (ExpenseRecordType) AllValues() []ExpenseRecordType {
	return []ExpenseRecordType{
		ExpenseRecordTypeDebt,
		ExpenseRecordTypeSettlement,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s ExpenseRecordType) MarshalText() ([]byte, error) {
	switch s {
	case ExpenseRecordTypeDebt:
		return []byte(s), nil
	case ExpenseRecordTypeSettlement:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *ExpenseRecordType) UnmarshalText(data []byte) error {
	switch ExpenseRecordType(data) {
	case ExpenseRecordTypeDebt:
		*s = ExpenseRecordTypeDebt
		return nil
	case ExpenseRecordTypeSettlement:
		*s = ExpenseRecordTypeSettlement
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Ref: #/components/schemas/ExpenseSummary
type ExpenseSummary struct {
	// Unique identifier for the expense.
	ID uuid.UUID `json:"id"`
	// Name of the expense.
	Name string `json:"name"`
	// Date and time of the expense.
	ExpenseDate time.Time `json:"expense_date"`
	// Total amount of the expense.
	Amount int32 `json:"amount"`
	// Date and time the expense was created.
	CreatedAt time.Time `json:"created_at"`
	// User ID of the creator.
	CreatedBy uuid.UUID `json:"created_by"`
	// Date and time the expense was last updated.
	UpdatedAt time.Time `json:"updated_at"`
	// User ID of the last updater.
	UpdatedBy uuid.UUID `json:"updated_by"`
}

// GetID returns the value of ID.
func (s *ExpenseSummary) GetID() uuid.UUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *ExpenseSummary) GetName() string {
	return s.Name
}

// GetExpenseDate returns the value of ExpenseDate.
func (s *ExpenseSummary) GetExpenseDate() time.Time {
	return s.ExpenseDate
}

// GetAmount returns the value of Amount.
func (s *ExpenseSummary) GetAmount() int32 {
	return s.Amount
}

// GetCreatedAt returns the value of CreatedAt.
func (s *ExpenseSummary) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetCreatedBy returns the value of CreatedBy.
func (s *ExpenseSummary) GetCreatedBy() uuid.UUID {
	return s.CreatedBy
}

// GetUpdatedAt returns the value of UpdatedAt.
func (s *ExpenseSummary) GetUpdatedAt() time.Time {
	return s.UpdatedAt
}

// GetUpdatedBy returns the value of UpdatedBy.
func (s *ExpenseSummary) GetUpdatedBy() uuid.UUID {
	return s.UpdatedBy
}

// SetID sets the value of ID.
func (s *ExpenseSummary) SetID(val uuid.UUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *ExpenseSummary) SetName(val string) {
	s.Name = val
}

// SetExpenseDate sets the value of ExpenseDate.
func (s *ExpenseSummary) SetExpenseDate(val time.Time) {
	s.ExpenseDate = val
}

// SetAmount sets the value of Amount.
func (s *ExpenseSummary) SetAmount(val int32) {
	s.Amount = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *ExpenseSummary) SetCreatedAt(val time.Time) {
	s.CreatedAt = val
}

// SetCreatedBy sets the value of CreatedBy.
func (s *ExpenseSummary) SetCreatedBy(val uuid.UUID) {
	s.CreatedBy = val
}

// SetUpdatedAt sets the value of UpdatedAt.
func (s *ExpenseSummary) SetUpdatedAt(val time.Time) {
	s.UpdatedAt = val
}

// SetUpdatedBy sets the value of UpdatedBy.
func (s *ExpenseSummary) SetUpdatedBy(val uuid.UUID) {
	s.UpdatedBy = val
}

// HealthcheckOK is response for Healthcheck operation.
type HealthcheckOK struct{}

// Ref: #/components/schemas/Ledger
type Ledger struct {
	ID uuid.UUID `json:"id"`
	// Ledger's name.
	Name string `json:"name"`
	// Date and time the ledger was created.
	CreatedAt time.Time `json:"created_at"`
	// User ID of the creator.
	CreatedBy uuid.UUID `json:"created_by"`
}

// GetID returns the value of ID.
func (s *Ledger) GetID() uuid.UUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *Ledger) GetName() string {
	return s.Name
}

// GetCreatedAt returns the value of CreatedAt.
func (s *Ledger) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetCreatedBy returns the value of CreatedBy.
func (s *Ledger) GetCreatedBy() uuid.UUID {
	return s.CreatedBy
}

// SetID sets the value of ID.
func (s *Ledger) SetID(val uuid.UUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *Ledger) SetName(val string) {
	s.Name = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *Ledger) SetCreatedAt(val time.Time) {
	s.CreatedAt = val
}

// SetCreatedBy sets the value of CreatedBy.
func (s *Ledger) SetCreatedBy(val uuid.UUID) {
	s.CreatedBy = val
}

type LedgerCreateOK struct {
	ID uuid.UUID `json:"id"`
}

// GetID returns the value of ID.
func (s *LedgerCreateOK) GetID() uuid.UUID {
	return s.ID
}

// SetID sets the value of ID.
func (s *LedgerCreateOK) SetID(val uuid.UUID) {
	s.ID = val
}

type LedgerCreateReq struct {
	// Ledger's name.
	Name string `json:"name"`
}

// GetName returns the value of Name.
func (s *LedgerCreateReq) GetName() string {
	return s.Name
}

// SetName sets the value of Name.
func (s *LedgerCreateReq) SetName(val string) {
	s.Name = val
}

type LedgerExpenseCreateOK struct {
	ID uuid.UUID `json:"id"`
}

// GetID returns the value of ID.
func (s *LedgerExpenseCreateOK) GetID() uuid.UUID {
	return s.ID
}

// SetID sets the value of ID.
func (s *LedgerExpenseCreateOK) SetID(val uuid.UUID) {
	s.ID = val
}

type LedgerExpenseListOK struct {
	// Cursor for pagination.
	Cursor   OptDateTime      `json:"cursor"`
	Expenses []ExpenseSummary `json:"expenses"`
}

// GetCursor returns the value of Cursor.
func (s *LedgerExpenseListOK) GetCursor() OptDateTime {
	return s.Cursor
}

// GetExpenses returns the value of Expenses.
func (s *LedgerExpenseListOK) GetExpenses() []ExpenseSummary {
	return s.Expenses
}

// SetCursor sets the value of Cursor.
func (s *LedgerExpenseListOK) SetCursor(val OptDateTime) {
	s.Cursor = val
}

// SetExpenses sets the value of Expenses.
func (s *LedgerExpenseListOK) SetExpenses(val []ExpenseSummary) {
	s.Expenses = val
}

type LedgerExpenseRecordCreateReq struct {
	Records []ExpenseRecord `json:"records"`
}

// GetRecords returns the value of Records.
func (s *LedgerExpenseRecordCreateReq) GetRecords() []ExpenseRecord {
	return s.Records
}

// SetRecords sets the value of Records.
func (s *LedgerExpenseRecordCreateReq) SetRecords(val []ExpenseRecord) {
	s.Records = val
}

// LedgerExpenseRecordDeleteOK is response for LedgerExpenseRecordDelete operation.
type LedgerExpenseRecordDeleteOK struct{}

type LedgerListOK struct {
	Ledgers []Ledger `json:"ledgers"`
}

// GetLedgers returns the value of Ledgers.
func (s *LedgerListOK) GetLedgers() []Ledger {
	return s.Ledgers
}

// SetLedgers sets the value of Ledgers.
func (s *LedgerListOK) SetLedgers(val []Ledger) {
	s.Ledgers = val
}

// Ref: #/components/schemas/LedgerMember
type LedgerMember struct {
	UserID uuid.UUID `json:"user_id"`
	// Date and time the member was created.
	CreatedAt time.Time `json:"created_at"`
	// User ID of the creator.
	CreatedBy uuid.UUID `json:"created_by"`
	// User's balance in the ledger.
	Balance int32 `json:"balance"`
}

// GetUserID returns the value of UserID.
func (s *LedgerMember) GetUserID() uuid.UUID {
	return s.UserID
}

// GetCreatedAt returns the value of CreatedAt.
func (s *LedgerMember) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetCreatedBy returns the value of CreatedBy.
func (s *LedgerMember) GetCreatedBy() uuid.UUID {
	return s.CreatedBy
}

// GetBalance returns the value of Balance.
func (s *LedgerMember) GetBalance() int32 {
	return s.Balance
}

// SetUserID sets the value of UserID.
func (s *LedgerMember) SetUserID(val uuid.UUID) {
	s.UserID = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *LedgerMember) SetCreatedAt(val time.Time) {
	s.CreatedAt = val
}

// SetCreatedBy sets the value of CreatedBy.
func (s *LedgerMember) SetCreatedBy(val uuid.UUID) {
	s.CreatedBy = val
}

// SetBalance sets the value of Balance.
func (s *LedgerMember) SetBalance(val int32) {
	s.Balance = val
}

// LedgerMemberAddAccepted is response for LedgerMemberAdd operation.
type LedgerMemberAddAccepted struct{}

type LedgerMemberAddReq struct {
	// Invite up to 99 other users. The limit is by Ledger.
	Emails []string `json:"emails"`
}

// GetEmails returns the value of Emails.
func (s *LedgerMemberAddReq) GetEmails() []string {
	return s.Emails
}

// SetEmails sets the value of Emails.
func (s *LedgerMemberAddReq) SetEmails(val []string) {
	s.Emails = val
}

type LedgerMemberListOK struct {
	Members []LedgerMember `json:"members"`
}

// GetMembers returns the value of Members.
func (s *LedgerMemberListOK) GetMembers() []LedgerMember {
	return s.Members
}

// SetMembers sets the value of Members.
func (s *LedgerMemberListOK) SetMembers(val []LedgerMember) {
	s.Members = val
}

// NewOptDateTime returns new OptDateTime with value set to v.
func NewOptDateTime(v time.Time) OptDateTime {
	return OptDateTime{
		Value: v,
		Set:   true,
	}
}

// OptDateTime is optional time.Time.
type OptDateTime struct {
	Value time.Time
	Set   bool
}

// IsSet returns true if OptDateTime was set.
func (o OptDateTime) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptDateTime) Reset() {
	var v time.Time
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptDateTime) SetTo(v time.Time) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptDateTime) Get() (v time.Time, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptDateTime) Or(d time.Time) time.Time {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptErrorMetadata returns new OptErrorMetadata with value set to v.
func NewOptErrorMetadata(v ErrorMetadata) OptErrorMetadata {
	return OptErrorMetadata{
		Value: v,
		Set:   true,
	}
}

// OptErrorMetadata is optional ErrorMetadata.
type OptErrorMetadata struct {
	Value ErrorMetadata
	Set   bool
}

// IsSet returns true if OptErrorMetadata was set.
func (o OptErrorMetadata) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptErrorMetadata) Reset() {
	var v ErrorMetadata
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptErrorMetadata) SetTo(v ErrorMetadata) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptErrorMetadata) Get() (v ErrorMetadata, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptErrorMetadata) Or(d ErrorMetadata) ErrorMetadata {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptInt32 returns new OptInt32 with value set to v.
func NewOptInt32(v int32) OptInt32 {
	return OptInt32{
		Value: v,
		Set:   true,
	}
}

// OptInt32 is optional int32.
type OptInt32 struct {
	Value int32
	Set   bool
}

// IsSet returns true if OptInt32 was set.
func (o OptInt32) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptInt32) Reset() {
	var v int32
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptInt32) SetTo(v int32) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptInt32) Get() (v int32, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptInt32) Or(d int32) int32 {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptString returns new OptString with value set to v.
func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

// OptString is optional string.
type OptString struct {
	Value string
	Set   bool
}

// IsSet returns true if OptString was set.
func (o OptString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptUUID returns new OptUUID with value set to v.
func NewOptUUID(v uuid.UUID) OptUUID {
	return OptUUID{
		Value: v,
		Set:   true,
	}
}

// OptUUID is optional uuid.UUID.
type OptUUID struct {
	Value uuid.UUID
	Set   bool
}

// IsSet returns true if OptUUID was set.
func (o OptUUID) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptUUID) Reset() {
	var v uuid.UUID
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptUUID) SetTo(v uuid.UUID) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptUUID) Get() (v uuid.UUID, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptUUID) Or(d uuid.UUID) uuid.UUID {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

type UserRegisterOK struct {
	ID uuid.UUID `json:"id"`
}

// GetID returns the value of ID.
func (s *UserRegisterOK) GetID() uuid.UUID {
	return s.ID
}

// SetID sets the value of ID.
func (s *UserRegisterOK) SetID(val uuid.UUID) {
	s.ID = val
}

type UserRegisterReq struct {
	// User's first name.
	FirstName string `json:"first_name"`
	// User's last name.
	LastName string `json:"last_name"`
	// User's email address.
	Email string `json:"email"`
	// User's password.
	Password string `json:"password"`
}

// GetFirstName returns the value of FirstName.
func (s *UserRegisterReq) GetFirstName() string {
	return s.FirstName
}

// GetLastName returns the value of LastName.
func (s *UserRegisterReq) GetLastName() string {
	return s.LastName
}

// GetEmail returns the value of Email.
func (s *UserRegisterReq) GetEmail() string {
	return s.Email
}

// GetPassword returns the value of Password.
func (s *UserRegisterReq) GetPassword() string {
	return s.Password
}

// SetFirstName sets the value of FirstName.
func (s *UserRegisterReq) SetFirstName(val string) {
	s.FirstName = val
}

// SetLastName sets the value of LastName.
func (s *UserRegisterReq) SetLastName(val string) {
	s.LastName = val
}

// SetEmail sets the value of Email.
func (s *UserRegisterReq) SetEmail(val string) {
	s.Email = val
}

// SetPassword sets the value of Password.
func (s *UserRegisterReq) SetPassword(val string) {
	s.Password = val
}
