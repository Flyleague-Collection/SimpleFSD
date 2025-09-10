// Package utils
package utils

func ReverseForEach[T any](slice []T, f func(index int, value T)) {
	for i := len(slice) - 1; i >= 0; i-- {
		f(i, slice[i])
	}
}
