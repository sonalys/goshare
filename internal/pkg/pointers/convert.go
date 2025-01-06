package pointers

func Convert[T1, T2 any](from *T1, cast func(T1) T2) *T2 {
	if from == nil {
		return nil
	}

	to := cast(*from)

	return &to
}
