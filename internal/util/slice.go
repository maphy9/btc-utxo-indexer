package util

func Reverse[T any](data []T) []T {
	res := make([]T, len(data))
	for i, item := range data {
		res[len(data)-i-1] = item
	}
	return res
}
