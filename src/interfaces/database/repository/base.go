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
