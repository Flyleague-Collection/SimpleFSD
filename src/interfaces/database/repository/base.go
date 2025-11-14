// Package repository
package repository

import "github.com/half-nothing/simple-fsd/src/interfaces/database/entity"

// Base 是一个泛型接口，定义了基本的数据库操作方法
// T 是实现了 entity.Base 接口的实体类型
type Base[T entity.Base] interface {
	GetById(id uint) (T, error)
	Save(entity T) error
	Delete(entity T) error
	Update(entity T, updates map[string]interface{}) error
}

// Enum 枚举类型结构体，用于表示具有值和标签的枚举项
// T 是可比较的类型
type Enum[T comparable] struct {
	Value T      `json:"value"`
	Label string `json:"label"`
}

// NewEnum 创建一个新的枚举实例
func NewEnum[T comparable](value T, label string) *Enum[T] {
	return &Enum[T]{Value: value, Label: label}
}

// EnumManagerInterface 枚举管理器接口，定义了枚举管理的基本操作
// T 是可比较的类型
type EnumManagerInterface[T comparable] interface {
	IsValidEnum(value T) bool
	GetEnum(value T) *Enum[T]
	GetEnums() map[T]*Enum[T]
}

// EnumManager 枚举管理器结构体，用于管理一组枚举项
// T 是可比较的类型
type EnumManager[T comparable] struct {
	enums map[T]*Enum[T]
}

func NewEnumManager[T comparable](enums ...*Enum[T]) *EnumManager[T] {
	manager := &EnumManager[T]{
		enums: make(map[T]*Enum[T]),
	}
	for _, e := range enums {
		manager.enums[e.Value] = e
	}
	return manager
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
