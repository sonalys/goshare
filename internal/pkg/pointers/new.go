package pointers

func New[T any](value T) *T {
	return &value
}
