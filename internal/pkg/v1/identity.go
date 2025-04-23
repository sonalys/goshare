package v1

type Identity struct {
	Email  string
	UserID ID
	Exp    int64
}

const (
	ErrAuthenticationExpired = StringError("authentication expired")
)
