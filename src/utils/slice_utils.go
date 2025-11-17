// Package utils
package utils

type SliceOperation[T any] struct {
	value []T
}

func NewSliceOperation[T any](src []T) *SliceOperation[T] {
	return &SliceOperation[T]{src}
}

func (operation *SliceOperation[T]) Append(element T) *SliceOperation[T] {
	operation.value = append(operation.value, element)
	return operation
}

func (operation *SliceOperation[T]) Clone() *SliceOperation[T] {
	return NewSliceOperation(operation.value)
}

func (operation *SliceOperation[T]) Value() []T { return operation.value }

func (operation *SliceOperation[T]) Find(comparator func(element T) bool) T {
	return Find(operation.value, comparator)
}

func (operation *SliceOperation[T]) Filter(filter func(element T) bool) *SliceOperation[T] {
	operation.value = Filter(operation.value, filter)
	return operation
}

func (operation *SliceOperation[T]) Map(mapper func(element T)) *SliceOperation[T] {
	Map(operation.value, mapper)
	return operation
}

func (operation *SliceOperation[T]) ForEach(callback func(index int, element T)) *SliceOperation[T] {
	ForEach(operation.value, callback)
	return operation
}

func FilterNotNull[T any](element T) bool {
	return element != nil
}

func ReverseForEach[T any](slice []T, f func(index int, value T)) {
	for i := len(slice) - 1; i >= 0; i-- {
		f(i, slice[i])
	}
}

func Any[T any](src []T, comparator func(element T) bool) bool {
	for _, v := range src {
		if comparator(v) {
			return true
		}
	}
	return false
}

func Find[T any](src []T, comparator func(element T) bool) T {
	for _, v := range src {
		if comparator(v) {
			return v
		}
	}
	var zero T
	return zero
}

func Filter[T any](src []T, filter func(element T) bool) (result []T) {
	result = make([]T, 0, len(src))
	for _, v := range src {
		if filter(v) {
			result = append(result, v)
		}
	}
	return
}

func Map[T any](src []T, mapper func(element T)) {
	for _, v := range src {
		mapper(v)
	}
}

func ForEach[T any](src []T, callback func(index int, element T)) {
	for i, v := range src {
		callback(i, v)
	}
}
