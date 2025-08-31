package domain

import "fmt"

type (
	UserMaxLedgersError struct {
		UserID     ID
		MaxLedgers int
	}
)

func (e UserMaxLedgersError) Error() string {
	return fmt.Sprintf("user '%s' has reached maximum number of ledgers: %d", e.UserID, e.MaxLedgers)
}
