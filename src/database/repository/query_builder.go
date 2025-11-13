// Package repository 数据库操作基础仓库
package repository

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// QueryBuilderInterface 查询构建器接口
//
// 该接口定义了查询构建器的基本功能，支持链式调用和条件构建。
type QueryBuilderInterface[T entity.Base] interface {
	// Where 添加查询条件
	//
	// 参数:
	//   - condition: 查询条件字符串
	//   - args: 查询条件参数
	//
	// 返回值:
	//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
	Where(condition string, args ...interface{}) QueryBuilderInterface[T]

	// Preload 添加预加载关联
	//
	// 参数:
	//   - column: 要预加载的关联字段名
	//   - args: 预加载的附加参数
	//
	// 返回值:
	//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
	Preload(column string, args ...interface{}) QueryBuilderInterface[T]

	// Order 添加排序条件
	//
	// 参数:
	//   - column: 排序字段名
	//   - desc: 是否降序排列
	//
	// 返回值:
	//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
	Order(column string, desc bool) QueryBuilderInterface[T]

	// Limit 限制返回记录数
	//
	// 参数:
	//   - limit: 限制返回的记录数量
	//
	// 返回值:
	//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
	Limit(limit int) QueryBuilderInterface[T]

	// Offset 设置查询偏移量
	//
	// 参数:
	//   - offset: 查询偏移量
	//
	// 返回值:
	//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
	Offset(offset int) QueryBuilderInterface[T]

	// Select 指定要查询的字段
	//
	// 参数:
	//   - fields: 要查询的字段列表
	//
	// 返回值:
	//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
	Select(fields string) QueryBuilderInterface[T]

	// Group 添加分组条件
	//
	// 参数:
	//   - group: 分组字段
	//
	// 返回值:
	//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
	Group(group string) QueryBuilderInterface[T]

	// Having 添加Having条件
	//
	// 参数:
	//   - having: Having条件字符串
	//   - args: Having条件参数
	//
	// 返回值:
	//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
	Having(having string, args ...interface{}) QueryBuilderInterface[T]

	// SetTransaction 设置是否启用事务
	//
	// 参数:
	//   - enable: 是否启用事务
	//
	// 返回值:
	//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
	SetTransaction(enable bool) QueryBuilderInterface[T]

	// GetTransaction 获取是否启用事务的设置
	//
	// 返回值:
	//   - bool: 是否启用事务，true表示启用，false表示不启用
	GetTransaction() bool

	// SetLock 设置锁表达式
	//
	// 参数:
	//   - lock: 锁表达式
	//
	// 返回值:
	//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
	SetLock(lock clause.Expression) QueryBuilderInterface[T]

	// Build 构建最终的查询语句
	//
	// 参数:
	//   - tx: GORM数据库连接对象
	//
	// 返回值:
	//   - *gorm.DB: 构建完成的查询对象
	Build(tx *gorm.DB) *gorm.DB
}

// QueryBuilder 查询构建器结构体
//
// 实现了QueryBuilderInterface接口，用于构建复杂的数据库查询条件。
type QueryBuilder[T entity.Base] struct {
	conditions  []func(*gorm.DB) *gorm.DB
	transaction bool
	lock        clause.Expression
}

// NewQueryBuilder 创建一个新的查询构建器实例
//
// 该函数用于初始化QueryBuilder结构体，创建一个空的查询条件列表。
//
// 返回值:
//   - *QueryBuilder[T]: 新创建的查询构建器实例指针
func NewQueryBuilder[T entity.Base]() *QueryBuilder[T] {
	return &QueryBuilder[T]{
		conditions:  make([]func(*gorm.DB) *gorm.DB, 0),
		transaction: false,
		lock:        nil,
	}
}

// Where 添加查询条件到构建器中
//
// 该函数将查询条件添加到条件列表中，支持链式调用。
//
// 参数:
//   - condition: 查询条件字符串
//   - args: 查询条件参数
//
// 返回值:
//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
func (builder *QueryBuilder[T]) Where(condition string, args ...interface{}) QueryBuilderInterface[T] {
	builder.conditions = append(builder.conditions, func(tx *gorm.DB) *gorm.DB {
		return tx.Where(condition, args...)
	})
	return builder
}

// Preload 添加预加载关联到构建器中
//
// 该函数将预加载关联添加到条件列表中，用于在查询时同时加载关联数据。
//
// 参数:
//   - column: 要预加载的关联字段名
//   - args: 预加载的附加参数
//
// 返回值:
//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
func (builder *QueryBuilder[T]) Preload(column string, args ...interface{}) QueryBuilderInterface[T] {
	builder.conditions = append(builder.conditions, func(tx *gorm.DB) *gorm.DB {
		return tx.Preload(column, args...)
	})
	return builder
}

// Order 添加排序条件到构建器中
//
// 该函数将排序条件添加到条件列表中，支持链式调用。
//
// 参数:
//   - column: 排序字段名
//   - desc: 是否降序排列
//
// 返回值:
//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
func (builder *QueryBuilder[T]) Order(column string, desc bool) QueryBuilderInterface[T] {
	builder.conditions = append(builder.conditions, func(tx *gorm.DB) *gorm.DB {
		order := "ASC"
		if desc {
			order = "DESC"
		}
		return tx.Order(column + " " + order)
	})
	return builder
}

// Limit 添加限制条件到构建器中
//
// 该函数将限制条件添加到条件列表中，用于限制查询结果的数量。
//
// 参数:
//   - limit: 限制返回的记录数量
//
// 返回值:
//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
func (builder *QueryBuilder[T]) Limit(limit int) QueryBuilderInterface[T] {
	builder.conditions = append(builder.conditions, func(tx *gorm.DB) *gorm.DB {
		return tx.Limit(limit)
	})
	return builder
}

// Offset 添加偏移条件到构建器中
//
// 该函数将偏移条件添加到条件列表中，用于设置查询的起始位置。
//
// 参数:
//   - offset: 查询偏移量
//
// 返回值:
//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
func (builder *QueryBuilder[T]) Offset(offset int) QueryBuilderInterface[T] {
	builder.conditions = append(builder.conditions, func(tx *gorm.DB) *gorm.DB {
		return tx.Offset(offset)
	})
	return builder
}

// Select 添加字段选择条件到构建器中
//
// 该函数将字段选择条件添加到条件列表中，用于指定查询的字段。
//
// 参数:
//   - fields: 要查询的字段列表
//
// 返回值:
//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
func (builder *QueryBuilder[T]) Select(fields string) QueryBuilderInterface[T] {
	builder.conditions = append(builder.conditions, func(tx *gorm.DB) *gorm.DB {
		return tx.Select(fields)
	})
	return builder
}

// Group 添加分组条件到构建器中
//
// 该函数将分组条件添加到条件列表中，用于对查询结果进行分组。
//
// 参数:
//   - group: 分组字段
//
// 返回值:
//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
func (builder *QueryBuilder[T]) Group(group string) QueryBuilderInterface[T] {
	builder.conditions = append(builder.conditions, func(tx *gorm.DB) *gorm.DB {
		return tx.Group(group)
	})
	return builder
}

// Having 添加Having条件到构建器中
//
// 该函数将Having条件添加到条件列表中，用于对分组后的结果进行筛选。
//
// 参数:
//   - having: Having条件字符串
//   - args: Having条件参数
//
// 返回值:
//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
func (builder *QueryBuilder[T]) Having(having string, args ...interface{}) QueryBuilderInterface[T] {
	builder.conditions = append(builder.conditions, func(tx *gorm.DB) *gorm.DB {
		return tx.Having(having, args...)
	})
	return builder
}

// SetTransaction 设置是否启用事务
//
// 该函数用于设置查询是否在事务中执行。
//
// 参数:
//   - transaction: 是否启用事务，true表示启用，false表示不启用
//
// 返回值:
//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
func (builder *QueryBuilder[T]) SetTransaction(transaction bool) QueryBuilderInterface[T] {
	builder.transaction = transaction
	return builder
}

// GetTransaction 获取是否启用事务的设置
//
// 返回值:
//   - bool: 是否启用事务，true表示启用，false表示不启用
func (builder *QueryBuilder[T]) GetTransaction() bool {
	return builder.transaction
}

// SetLock 设置锁表达式
//
// 该函数用于设置查询时使用的锁表达式，如共享锁或排他锁。
//
// 参数:
//   - lock: 锁表达式，指定查询时要使用的锁类型
//
// 返回值:
//   - QueryBuilderInterface[T]: 返回自身以支持链式调用
func (builder *QueryBuilder[T]) SetLock(lock clause.Expression) QueryBuilderInterface[T] {
	builder.lock = lock
	return builder
}

// Build 构建最终的查询语句
//
// 该函数将所有已添加的查询条件应用到数据库连接对象上，
// 构建出完整的查询语句。
//
// 参数:
//   - tx: GORM数据库连接对象
//
// 返回值:
//   - *gorm.DB: 构建完成的查询对象
func (builder *QueryBuilder[T]) Build(tx *gorm.DB) *gorm.DB {
	if builder.lock != nil {
		tx = tx.Clauses(builder.lock)
	}

	for _, condition := range builder.conditions {
		tx = condition(tx)
	}
	return tx
}
