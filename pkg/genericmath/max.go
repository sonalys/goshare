package genericmath

import "cmp"

func Max[T cmp.Ordered](values ...T) T {
	maxValue := values[0]

	for i := 1; i < len(values); i++ {
		if values[i] > maxValue {
			maxValue = values[i]
		}
	}

	return maxValue
}
