// Package enum
package enum

// Enum 枚举类型结构体，用于表示具有值和标签的枚举项
// T 是可比较的类型
type Enum[T comparable] struct {
	Value T      `json:"value"`
	Label string `json:"label"`
}

// New 创建一个新的枚举实例
func New[T comparable](value T, label string) *Enum[T] {
	return &Enum[T]{Value: value, Label: label}
}

// ManagerInterface 枚举管理器接口，定义了枚举管理的基本操作
// T 是可比较的类型
type ManagerInterface[T comparable] interface {
	IsValidEnum(value T) bool
	GetEnum(value T) *Enum[T]
	GetEnums() map[T]*Enum[T]
}

// Manager 枚举管理器结构体，用于管理一组枚举项
// T 是可比较的类型
type Manager[T comparable] struct {
	enums map[T]*Enum[T]
}

func NewManager[T comparable](enums ...*Enum[T]) *Manager[T] {
	manager := &Manager[T]{
		enums: make(map[T]*Enum[T]),
	}
	for _, e := range enums {
		manager.enums[e.Value] = e
	}
	return manager
}

func (manager *Manager[T]) IsValidEnum(value T) bool {
	return manager.enums[value] != nil
}

func (manager *Manager[T]) GetEnum(value T) *Enum[T] {
	return manager.enums[value]
}

func (manager *Manager[T]) GetEnums() map[T]*Enum[T] {
	return manager.enums
}
