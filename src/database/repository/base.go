// Package repository 数据库操作基础仓库
package repository

import (
	"context"
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"github.com/half-nothing/simple-fsd/src/interfaces/database/repository"
	"github.com/half-nothing/simple-fsd/src/interfaces/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BaseRepository 数据库操作基础仓库
// 该仓库提供了数据库操作的基础功能，如查询、保存、更新和删除。
// 使用泛型T来约束实体类型，确保类型安全。
type BaseRepository[T entity.Base] struct {
	logger       logger.Interface
	db           *gorm.DB
	queryTimeout time.Duration
}

// NewBaseRepository 创建一个新的基础仓库实例
//
// 该函数用于初始化BaseRepository结构体，为所有具体的仓库实现提供基础功能。
// 包括日志记录、数据库连接和查询超时设置等通用功能。
//
// 参数:
//   - lg: 日志接口，用于记录仓库操作日志
//   - loggerName: 日志记录器名称，用于区分不同仓库的日志输出
//   - db: GORM数据库连接实例
//   - queryTimeout: 查询超时时间，用于控制数据库查询的最大执行时间
//
// 返回值:
//   - *BaseRepository: 新创建的基础仓库实例指针
func NewBaseRepository[T entity.Base](
	lg logger.Interface,
	loggerName string,
	db *gorm.DB,
	queryTimeout time.Duration,
) *BaseRepository[T] {
	return &BaseRepository[T]{
		logger:       logger.NewLoggerAdapter(lg, loggerName),
		db:           db,
		queryTimeout: queryTimeout,
	}
}

// query 执行数据库查询操作
//
// 该函数在指定的超时时间内执行传入的查询函数，为查询操作提供统一的上下文管理。
//
// 参数:
//   - queryFunc: 要执行的查询函数，接收*gorm.DB作为参数并返回error
//
// 返回值:
//   - error: 查询过程中可能发生的错误
func (repo *BaseRepository[T]) query(queryFunc func(tx *gorm.DB) error) error {
	if queryFunc == nil {
		return repository.ErrArgument
	}

	ctx, cancel := context.WithTimeout(context.Background(), repo.queryTimeout)
	defer cancel()

	return queryFunc(repo.db.WithContext(ctx))
}

// queryEntityWithBuilder 使用查询构建器查询单个实体
//
// 该函数通过传入的查询构建器构造查询条件，并将查询结果填充到目标实体对象中。
// 适用于需要根据复杂条件查询单个实体的场景。
//
// 参数:
//   - queryBuilder: 查询构建器接口，用于构造特定的查询条件
//   - dest: 目标实体对象指针，用于存储查询结果
//
// 返回值:
//   - error: 查询过程中可能发生的错误
func (repo *BaseRepository[T]) queryEntityWithBuilder(queryBuilder QueryBuilderInterface[T], dest T) error {
	queryFunc := func(tx *gorm.DB) error { return queryBuilder.Build(tx).First(dest).Error }
	if queryBuilder.GetTransaction() {
		return repo.queryWithTransaction(queryFunc)
	}
	return repo.query(queryFunc)
}

// queryEntitiesWithBuilder 使用查询构建器查询多个实体
//
// 该函数通过传入的查询构建器构造查询条件，并将查询结果填充到目标实体切片中。
// 适用于需要根据复杂条件查询多个实体的场景。
//
// 参数:
//   - queryBuilder: 查询构建器接口，用于构造特定的查询条件
//   - dest: 目标实体对象切片指针，用于存储查询结果
//
// 返回值:
//   - error: 查询过程中可能发生的错误
func (repo *BaseRepository[T]) queryEntitiesWithBuilder(queryBuilder QueryBuilderInterface[T], dest []T) error {
	queryFunc := func(tx *gorm.DB) error { return queryBuilder.Build(tx).Find(dest).Error }
	if queryBuilder.GetTransaction() {
		return repo.queryWithTransaction(queryFunc)
	}
	return repo.query(queryFunc)
}

// queryWithTransaction 在事务中执行数据库查询操作
//
// 该函数在指定的超时时间内启动一个数据库事务，并在该事务中执行传入的查询函数。
// 如果查询函数执行成功，事务将被提交；否则事务将被回滚。
//
// 参数:
//   - queryFunc: 要在事务中执行的查询函数，接收*gorm.DB作为参数并返回error
//
// 返回值:
//   - error: 查询过程中可能发生的错误
func (repo *BaseRepository[T]) queryWithTransaction(queryFunc func(tx *gorm.DB) error) error {
	if queryFunc == nil {
		return repository.ErrArgument
	}

	return repo.query(func(tx *gorm.DB) error {
		return tx.Transaction(queryFunc)
	})
}

// queryWithPagination 执行分页查询操作
//
// 该函数使用分页器接口执行分页查询，计算总记录数并获取当前页数据。
//
// 参数:
//   - paginator: 分页器接口，负责执行分页逻辑
//   - page: 页面接口，包含分页参数和查询配置
//
// 返回值:
//   - int64: 总记录数
//   - error: 查询过程中可能发生的错误
func (repo *BaseRepository[T]) queryWithPagination(paginator PageableInterface[T], page PageInterface[T]) (int64, error) {
	if paginator == nil {
		return 0, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), repo.queryTimeout)
	defer cancel()

	return paginator.Paginate(ctx, page)
}

// queryWithLock 在带锁的事务中执行数据库查询操作
//
// 该函数在指定的超时时间内启动一个带锁的数据库事务，并在该事务中执行传入的查询函数。
// 适用于需要加锁的并发安全操作。
//
// 参数:
//   - lock: 锁表达式，指定要使用的锁类型
//   - queryFunc: 要在事务中执行的查询函数，接收*gorm.DB作为参数并返回error
//
// 返回值:
//   - error: 查询过程中可能发生的错误
func (repo *BaseRepository[T]) queryWithLock(queryFunc func(tx *gorm.DB) error, lock clause.Expression) error {
	if queryFunc == nil {
		return repository.ErrArgument
	}

	return repo.query(func(tx *gorm.DB) error {
		return tx.Clauses(lock).Transaction(queryFunc)
	})
}

func (repo *BaseRepository[T]) saveWithOptionalLock(model T, lock clause.Expression) error {
	if model == nil {
		return repository.ErrArgument
	}

	processFunc := func(tx *gorm.DB) error {
		if model.GetId() <= 0 {
			return tx.Create(model).Error
		}
		return tx.Save(model).Error
	}

	if lock == nil {
		return repo.queryWithTransaction(processFunc)
	}
	return repo.queryWithLock(processFunc, lock)
}

// save 保存实体到数据库
//
// 该函数根据实体的ID字段判断是创建新记录还是更新现有记录。
// 如果ID为0，则执行创建操作；否则执行更新操作。
//
// 参数:
//   - model: 要保存的实体对象
//
// 返回值:
//   - error: 保存过程中可能发生的错误，包括实体为空或数据库操作错误
func (repo *BaseRepository[T]) save(model T) error {
	return repo.saveWithOptionalLock(model, nil)
}

func (repo *BaseRepository[T]) saveWithLock(model T, lock clause.Expression) error {
	if lock == nil {
		lock = clause.Locking{Strength: "UPDATE"}
	}

	return repo.saveWithOptionalLock(model, lock)
}

func (repo *BaseRepository[T]) updateWithOptionalLock(model T, updates map[string]interface{}, lock clause.Expression) error {
	if updates == nil || len(updates) == 0 {
		return nil
	}

	if model == nil || model.GetId() <= 0 {
		return repository.ErrArgument
	}

	processFunc := func(tx *gorm.DB) error {
		return tx.Model(model).Updates(updates).Error
	}

	if lock != nil {
		return repo.queryWithTransaction(processFunc)
	}
	return repo.queryWithLock(processFunc, lock)
}

// update 更新实体的部分字段
//
// 该函数用于更新实体的部分字段，通过updates映射指定要更新的字段和值。
// 更新前会验证实体是否有效以及是否包含有效的ID。
//
// 参数:
//   - model: 要更新的实体对象
//   - updates: 包含要更新的字段名和值的映射
//
// 返回值:
//   - error: 更新过程中可能发生的错误
func (repo *BaseRepository[T]) update(model T, updates map[string]interface{}) error {
	return repo.updateWithOptionalLock(model, updates, nil)
}

func (repo *BaseRepository[T]) updateWithLock(model T, updates map[string]interface{}, lock clause.Expression) error {
	if lock == nil {
		lock = clause.Locking{Strength: "UPDATE"}
	}

	return repo.updateWithOptionalLock(model, updates, lock)
}

// deleteWithOptionalLock 删除指定的模型实体
// 参数:
//   - model: 要删除的模型实体，必须实现GetId()方法且ID大于0
//   - lock: 锁定条件表达式，可为nil
//
// 返回值:
//   - error: 删除操作可能产生的错误，如果参数无效则返回ErrArgument
func (repo *BaseRepository[T]) deleteWithOptionalLock(model T, lock clause.Expression) error {
	if model == nil || model.GetId() <= 0 {
		return repository.ErrArgument
	}

	processFunc := func(tx *gorm.DB) error {
		return tx.Delete(model).Error
	}

	if lock == nil {
		return repo.queryWithTransaction(processFunc)
	}
	return repo.queryWithLock(processFunc, lock)
}

// delete 从数据库中删除实体
//
// 该函数用于删除指定的实体记录。
//
// 参数:
//   - model: 要删除的实体对象
//
// 返回值:
//   - error: 删除过程中可能发生的错误
func (repo *BaseRepository[T]) delete(model T) error {
	return repo.deleteWithOptionalLock(model, nil)
}

// deleteWithLock 带锁删除指定的模型实体
// 当未提供锁时，默认使用UPDATE锁
// 参数:
//   - model: 要删除的模型实体，必须实现GetId()方法且ID大于0
//   - lock: 锁定条件表达式，若为nil则默认使用UPDATE锁
//
// 返回值:
//   - error: 删除操作可能产生的错误，如果参数无效则返回ErrArgument
func (repo *BaseRepository[T]) deleteWithLock(model T, lock clause.Expression) error {
	if lock == nil {
		lock = clause.Locking{Strength: "UPDATE"}
	}

	return repo.deleteWithOptionalLock(model, lock)
}
