package domain

import "fmt"

type (
	ErrUserMaxLedgers struct {
		UserID     ID
		MaxLedgers int
	}
)

func (e ErrUserMaxLedgers) Error() string {
	return fmt.Sprintf("user '%s' has reached maximum number of ledgers: %d", e.UserID, e.MaxLedgers)
}
