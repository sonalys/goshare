package genericmath

import "golang.org/x/exp/constraints"

func Max[T constraints.Integer | constraints.Float](values ...T) T {
	max := values[0]

	for i := 1; i < len(values); i++ {
		if values[i] > max {
			max = values[i]
		}
	}

	return max
}
