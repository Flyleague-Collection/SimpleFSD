// Package repository
package repository

type Base[T any] interface {
	GetById(id uint) (T, error)
	Save(entity T) error
	Delete(entity T) error
	Update(entity T, updates map[string]interface{}) error
}

type Enum[T any] struct {
	Value T      `json:"value"`
	Label string `json:"label"`
}

func NewEnum[T any](value T, label string) *Enum[T] {
	return &Enum[T]{Value: value, Label: label}
}

type Builder[T any] interface {
	Build() T
}
