// Package entity
package entity

type Base interface {
	GetId() uint
}

type Comparable[T Base] interface {
	Equal(other T) bool
	Diff(other T) map[string]interface{}
}
