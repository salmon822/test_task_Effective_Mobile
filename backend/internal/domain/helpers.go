package domain

func MapSlice[T1, T2 any](x []T1, converter func(T1) T2) []T2 {
	res := make([]T2, len(x))
	for i := range x {
		res[i] = converter(x[i])
	}
	return res
}
