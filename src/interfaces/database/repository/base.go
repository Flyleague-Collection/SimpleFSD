// Package repository
package repository

import "github.com/half-nothing/simple-fsd/src/interfaces/database/entity"

type Base[T entity.Base] interface {
	GetById(id uint) (T, error)
	Save(entity T) error
	Delete(entity T) error
	Update(entity T, updates map[string]interface{}) error
}

type Enum[T comparable] struct {
	Value T      `json:"value"`
	Label string `json:"label"`
}

func NewEnum[T comparable](value T, label string) *Enum[T] {
	return &Enum[T]{Value: value, Label: label}
}

type EnumManagerInterface[T comparable] interface {
	IsValidEnum(value T) bool
	GetEnum(value T) *Enum[T]
	GetEnums() map[T]*Enum[T]
}

type EnumManager[T comparable] struct {
	enums map[T]*Enum[T]
}

func NewEnumManager[T comparable](enums ...*Enum[T]) *EnumManager[T] {
	m := &EnumManager[T]{}
	for _, e := range enums {
		m.enums[e.Value] = e
	}
	return m
}

func (manager *EnumManager[T]) IsValidEnum(value T) bool {
	return manager.enums[value] != nil
}

func (manager *EnumManager[T]) GetEnum(value T) *Enum[T] {
	return manager.enums[value]
}

func (manager *EnumManager[T]) GetEnums() map[T]*Enum[T] {
	return manager.enums
}

type Builder[T any] interface {
	Build() T
}
