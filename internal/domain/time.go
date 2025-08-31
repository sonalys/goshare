package domain

import "time"

// Now returns a non-monotonic time.Time.
// Useful for avoiding mistakes when comparing distinct instants,
// Also helpful for the Postgres time serialization.
func Now() time.Time {
	return time.Now().Round(0)
}
