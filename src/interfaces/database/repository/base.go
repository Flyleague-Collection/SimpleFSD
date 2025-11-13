// Package repository
package repository

type Base[T any] interface {
	GetById(id uint) (T, error)
	Save(entity T) error
	Delete(entity T) error
	Update(entity T, updates map[string]interface{}) error
}

type Enum struct {
	Index int    `json:"value"`
	Name  string `json:"name"`
}

func NewEnum(index int, name string) *Enum {
	return &Enum{Index: index, Name: name}
}

type Builder[T any] interface {
	Build() T
}
