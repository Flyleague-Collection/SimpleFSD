// Package repository
package repository

import (
	"context"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"gorm.io/gorm"
)

// PageInterface 分页接口定义
type PageInterface[T entity.Base] interface {
	GetOffset() int
	GetTotalQuery(tx *gorm.DB) *gorm.DB
	GetPageQuery(tx *gorm.DB) *gorm.DB
}

// Page 分页结构体
// 使用泛型T来约束实体类型，确保类型安全
type Page[T entity.Base] struct {
	pageNumber int
	pageSize   int
	dest       []T
	model      T
	queryFunc  func(tx *gorm.DB) *gorm.DB
}

// NewPage 创建一个新的分页对象
//
// 该函数用于初始化Page结构体，配置分页查询所需的各种参数。
//
// 参数:
//   - pageNumber: 页码（从1开始）
//   - pageSize: 每页记录数
//   - dest: 查询结果存储目标
//   - model: 查询的模型类型
//   - queryFunc: 自定义查询函数
//
// 返回值:
//   - *Page: 新创建的分页对象指针
func NewPage[T entity.Base](pageNumber int,
	pageSize int,
	dest []T,
	model T,
	queryFunc func(tx *gorm.DB) *gorm.DB,
) *Page[T] {
	return &Page[T]{
		pageNumber: pageNumber,
		pageSize:   pageSize,
		dest:       dest,
		model:      model,
		queryFunc:  queryFunc,
	}
}

// GetOffset 计算分页偏移量
//
// 返回值:
//   - int: 分页查询所需的偏移量
func (p *Page[T]) GetOffset() int {
	return (p.pageNumber - 1) * p.pageSize
}

// GetTotalQuery 构造计算总记录数的查询
//
// 参数:
//   - tx: 数据库连接对象
//
// 返回值:
//   - *gorm.DB: 构造好的查询对象
func (p *Page[T]) GetTotalQuery(tx *gorm.DB) *gorm.DB {
	// 应用自定义查询条件
	if p.queryFunc != nil {
		tx = p.queryFunc(tx)
	}
	return tx.Model(p.model).Select("id")
}

// GetPageQuery 构造分页查询
//
// 参数:
//   - tx: 数据库连接对象
//
// 返回值:
//   - *gorm.DB: 构造好的查询对象
func (p *Page[T]) GetPageQuery(tx *gorm.DB) *gorm.DB {
	// 应用自定义查询条件
	if p.queryFunc != nil {
		tx = p.queryFunc(tx)
	}
	return tx.Offset(p.GetOffset()).Limit(p.pageSize).Find(p.dest)
}

// PageableInterface 可分页接口定义
type PageableInterface[T entity.Base] interface {
	Paginate(ctx context.Context, page PageInterface[T]) (total int64, err error)
}

// PageRequest 分页器结构体
type PageRequest[T entity.Base] struct {
	database *gorm.DB
}

// NewPageRequest 创建一个新的分页器
//
// 参数:
//   - db: 数据库连接对象
//
// 返回值:
//   - *PageRequest: 新创建的分页器对象指针
func NewPageRequest[T entity.Base](db *gorm.DB) *PageRequest[T] {
	return &PageRequest[T]{
		database: db,
	}
}

// Paginate 执行分页查询
//
// 该函数执行完整的分页查询流程，包括计算总记录数和获取当前页数据。
//
// 参数:
//   - ctx: 上下文对象
//   - page: 页面接口，包含分页参数和查询配置
//
// 返回值:
//   - total: 总记录数
//   - err: 查询过程中可能发生的错误
func (p *PageRequest[T]) Paginate(ctx context.Context, page PageInterface[T]) (total int64, err error) {
	// 获取数据库连接并应用上下文
	query := p.database.WithContext(ctx)

	// 查询总记录数
	if err = page.GetTotalQuery(query).Count(&total).Error; err != nil {
		return
	}

	// 查询当前页数据
	err = page.GetPageQuery(query).Error
	return
}
