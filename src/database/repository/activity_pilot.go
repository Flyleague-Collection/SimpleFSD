// Package repository 活动飞行员仓库实现
package repository

import (
	"errors"
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"github.com/half-nothing/simple-fsd/src/interfaces/database/repository"
	"github.com/half-nothing/simple-fsd/src/interfaces/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ActivityPilotRepository 活动飞行员仓库结构体
// 继承自BaseRepository，提供活动飞行员相关的数据库操作功能
type ActivityPilotRepository struct {
	*BaseRepository[*entity.ActivityPilot]
}

// NewActivityPilotRepository 创建一个新的活动飞行员仓库实例
//
// 该函数用于初始化ActivityPilotRepository结构体，继承自BaseRepository，
// 为活动飞行员相关的数据库操作提供基础功能。
//
// 参数:
//   - lg: 日志接口，用于记录仓库操作日志
//   - db: GORM数据库连接实例
//   - queryTimeout: 查询超时时间
//
// 返回值:
//   - *ActivityPilotRepository: 新创建的活动飞行员仓库实例指针
func NewActivityPilotRepository(
	lg logger.Interface,
	db *gorm.DB,
	queryTimeout time.Duration,
) *ActivityPilotRepository {
	return &ActivityPilotRepository{
		BaseRepository: NewBaseRepository[*entity.ActivityPilot](lg, "ActivityPilotRepository", db, queryTimeout),
	}
}

// New 创建一个新的活动飞行员实体
//
// 该函数用于创建一个ActivityPilot实体，用于表示用户参与某个活动的飞行员身份。
// 如果传入的参数不符合要求（如ID为0或字符串为空），则返回nil。
//
// 参数:
//   - activityId: 活动ID，必须大于0
//   - userId: 用户ID，必须大于0
//   - callsign: 呼号，不能为空字符串
//   - aircraftType: 飞机型号，不能为空字符串
//
// 返回值:
//   - *entity.ActivityPilot: 创建成功的活动飞行员实体指针，如果参数无效则返回nil
func (repo *ActivityPilotRepository) New(
	activityId uint,
	userId uint,
	callsign string,
	aircraftType string,
) *entity.ActivityPilot {
	if activityId <= 0 || userId <= 0 || callsign == "" || aircraftType == "" {
		return nil
	}
	return &entity.ActivityPilot{
		ActivityId:   activityId,
		UserId:       userId,
		Callsign:     callsign,
		AircraftType: aircraftType,
		Status:       repository.ActivityPilotStatusSigned.Index,
	}
}

// GetById 根据ID获取活动飞行员信息
//
// 该函数通过主键ID查询数据库中的活动飞行员记录。
// 如果参数无效或查询不到匹配记录，则返回相应错误。
//
// 参数:
//   - id: 活动飞行员记录的主键ID，必须大于0
//
// 返回值:
//   - *entity.ActivityPilot: 查询到的活动飞行员实体指针
//   - error: 可能发生的错误，包括参数错误、记录不存在或数据库查询错误
func (repo *ActivityPilotRepository) GetById(id uint) (*entity.ActivityPilot, error) {
	if id <= 0 {
		return nil, repository.ErrArgument
	}
	activityPilot := &entity.ActivityPilot{ID: id}

	err := repo.query(func(tx *gorm.DB) error {
		return tx.First(activityPilot).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrActivityPilotNotFound
	}

	return activityPilot, err
}

// GetByActivityIdAndUserId 根据活动ID和用户ID获取活动飞行员信息
//
// 该函数通过活动ID和用户ID查询数据库中对应的活动飞行员记录。
// 如果参数无效或查询不到匹配记录，则返回相应错误。
//
// 参数:
//   - activityId: 活动ID，必须大于0
//   - userId: 用户ID，必须大于0
//
// 返回值:
//   - *entity.ActivityPilot: 查询到的活动飞行员实体指针
//   - error: 可能发生的错误，包括参数错误或数据库查询错误
func (repo *ActivityPilotRepository) GetByActivityIdAndUserId(
	activityId uint,
	userId uint,
) (*entity.ActivityPilot, error) {
	if activityId <= 0 || userId <= 0 {
		return nil, repository.ErrArgument
	}

	activityPilot := &entity.ActivityPilot{}
	err := repo.query(func(tx *gorm.DB) error {
		return tx.Where("activity_id = ? AND user_id = ?", activityId, userId).First(activityPilot).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrActivityPilotNotFound
	}

	return activityPilot, err
}

// GetByActivityIdAndCallsign 根据活动ID和呼号获取活动飞行员信息
// 该函数通过活动ID和呼号查询数据库中对应的活动飞行员记录
//
// 参数:
//   - activityId: 活动ID，必须大于0
//   - callsign: 呼号，不能为空字符串
//
// 返回值:
//   - *entity.ActivityPilot: 查询到的活动飞行员实体指针
//   - error: 可能发生的错误，包括参数错误、记录不存在或数据库查询错误
func (repo *ActivityPilotRepository) GetByActivityIdAndCallsign(
	activityId uint,
	callsign string,
) (*entity.ActivityPilot, error) {
	if activityId <= 0 || callsign == "" {
		return nil, repository.ErrArgument
	}

	activityPilot := &entity.ActivityPilot{}
	err := repo.query(func(tx *gorm.DB) error {
		return tx.Where("activity_id = ? AND callsign = ?", activityId, callsign).First(activityPilot).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrActivityPilotNotFound
	}

	return activityPilot, err
}

// VerifyUserIdAndCallsign 检查活动中的用户ID或呼号是否存在
// 该函数用于检查给定的活动ID中，指定的用户ID或呼号是否已经存在记录
//
// 参数:
//   - activityId: 活动ID，必须大于0
//   - userId: 用户ID，必须大于0
//   - callsign: 呼号，不能为空字符串
//
// 返回值:
//   - *entity.ActivityPilot: 查询到的活动飞行员实体指针
//   - error: 可能发生的错误，包括参数错误、记录不存在或数据库查询错误
func (repo *ActivityPilotRepository) VerifyUserIdAndCallsign(
	activityId uint,
	userId uint,
	callsign string,
) (*entity.ActivityPilot, error) {
	if activityId <= 0 || userId <= 0 || callsign == "" {
		return nil, repository.ErrArgument
	}

	activityPilot := &entity.ActivityPilot{}
	err := repo.query(func(tx *gorm.DB) error {
		return tx.Where("activity_id = ? AND (user_id = ? OR callsign = ?)", activityId, userId, callsign).First(activityPilot).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrActivityPilotNotFound
	}

	return activityPilot, err
}

// Save 保存活动飞行员信息
//
// 该函数用于保存活动飞行员实体到数据库。如果实体是新创建的（ID为0），则执行插入操作；
// 如果实体已存在（ID不为0），则执行更新操作。
//
// 参数:
//   - entity: 要保存的活动飞行员实体指针
//
// 返回值:
//   - error: 保存过程中可能发生的错误，如实体为空或数据库操作错误
func (repo *ActivityPilotRepository) Save(entity *entity.ActivityPilot) error {
	return repo.save(entity)
}

// Delete 删除活动飞行员信息
//
// 该函数用于从数据库中删除指定的活动飞行员记录。
// 删除前会验证实体是否有效以及是否包含有效的ID。
//
// 参数:
//   - entity: 要删除的活动飞行员实体指针
//
// 返回值:
//   - error: 删除过程中可能发生的错误，如实体为空、参数无效或数据库操作错误
func (repo *ActivityPilotRepository) Delete(entity *entity.ActivityPilot) error {
	return repo.delete(entity)
}

// Update 更新活动飞行员信息
//
// 该函数用于更新活动飞行员实体的部分字段。通过传入更新映射，
// 可以只更新指定的字段，而不是整个实体。
//
// 参数:
//   - entity: 要更新的活动飞行员实体指针
//   - updates: 包含要更新的字段和值的映射
//
// 返回值:
//   - error: 更新过程中可能发生的错误，如实体为空、参数无效或数据库操作错误
func (repo *ActivityPilotRepository) Update(entity *entity.ActivityPilot, updates map[string]interface{}) error {
	return repo.update(entity, updates)
}

// UpdateStatus 更新活动飞行员的状态
//
// 该函数用于更新活动飞行员实体的状态字段。
// 通过传入新的状态值，更新数据库中对应的记录。
//
// 参数:
//   - entity: 要更新状态的活动飞行员实体指针
//   - status: 新的状态值
//
// 返回值:
//   - error: 更新过程中可能发生的错误，如实体为空、参数无效或数据库操作错误
func (repo *ActivityPilotRepository) UpdateStatus(entity *entity.ActivityPilot, status repository.ActivityPilotStatus) error {
	return repo.update(entity, map[string]interface{}{"status": status.Index})
}

// JoinActivity 更新飞行员加入活动的信息
// 该函数检查飞行员是否已经注册、呼号是否已被使用，如果没有问题则保存新的飞行员活动信息
//
// 参数:
//   - activity: 活动实体指针
//   - user: 用户实体指针
//   - callsign: 飞行员呼号
//   - aircraftType: 飞机型号
//
// 返回值:
//   - error: 错误信息，可能包括参数错误、飞行员已注册、呼号已被使用等
func (repo *ActivityPilotRepository) JoinActivity(activity *entity.Activity, user *entity.User, callsign string, aircraftType string) error {
	if activity == nil || user == nil || callsign == "" || aircraftType == "" {
		return repository.ErrArgument
	}

	processFunc := func(tx *gorm.DB) error {
		activityPilot, err := repo.VerifyUserIdAndCallsign(activity.ID, user.ID, callsign)
		if err == nil {
			if activityPilot.UserId == user.ID {
				return repository.ErrPilotAlreadySigned
			} else if activityPilot.Callsign == callsign {
				return repository.ErrCallsignAlreadyUsed
			} else {
				// 这个分支理论上不可能被执行
				// 如果查询成功那么userId和callsign必定至少有一个匹配
				return repository.ErrDataConflicts
			}
		}
		if !errors.Is(err, repository.ErrActivityPilotNotFound) {
			return err
		}

		activityPilot = repo.New(activity.ID, user.ID, callsign, aircraftType)
		if activityPilot == nil {
			return repository.ErrArgument
		}
		return repo.save(activityPilot)
	}

	return repo.queryWithLock(processFunc, clause.Locking{Strength: "UPDATE"})
}

// LeaveActivity 处理飞行员退出活动的逻辑
// 该函数检查飞行员是否已报名参加活动，如果已报名则删除其报名记录
//
// 参数:
//   - activity: 活动实体指针
//   - user: 用户实体指针
//
// 返回值:
//   - error: 错误信息，可能包括参数错误、飞行员未报名等
func (repo *ActivityPilotRepository) LeaveActivity(activity *entity.Activity, user *entity.User) error {
	if activity == nil || user == nil {
		return repository.ErrArgument
	}

	processFunc := func(tx *gorm.DB) error {
		activityPilot, err := repo.GetByActivityIdAndUserId(activity.ID, user.ID)
		if err != nil {
			if errors.Is(err, repository.ErrActivityPilotNotFound) {
				return repository.ErrPilotUnsigned
			}
			return err
		}

		return repo.delete(activityPilot)
	}

	return repo.queryWithLock(processFunc, clause.Locking{Strength: "UPDATE"})
}
