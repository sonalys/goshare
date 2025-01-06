package monoids

// Map applies a function to each element of a slice and returns a new slice with the results.
func Map[T1 any, T2 any](input []T1, f func(T1) T2) []T2 {
	result := make([]T2, len(input))
	for i, v := range input {
		result[i] = f(v)
	}
	return result
}

// Reduce applies a function to each element of a slice, accumulating the result.
func Reduce[T any, R any](input []T, initial R, f func(R, T) R) R {
	acc := initial
	for _, v := range input {
		acc = f(acc, v)
	}
	return acc
}
