package pointers

func Coalesce[T any](from *T, fallback T) T {
	if from == nil {
		return fallback
	}

	return *from
}
