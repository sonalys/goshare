package pointers

func From[T any](value T) *T {
	return &value
}
