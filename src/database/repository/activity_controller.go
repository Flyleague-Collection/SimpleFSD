// Package repository 实现了活动管制员相关的数据库操作
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

// ActivityControllerRepository 活动管制员仓库结构体
//
// 继承自BaseRepository，额外包含了席位仓库和飞行员仓库的引用，
// 用于处理活动管制员相关的复合业务逻辑。
type ActivityControllerRepository struct {
	*BaseRepository[*entity.ActivityController]
	facilityRepo repository.ActivityFacilityInterface
	pilotRepo    repository.ActivityPilotInterface
}

// NewActivityControllerRepository 创建一个新的活动管制员仓库实例
//
// 初始化ActivityControllerRepository结构体，设置基础仓库和依赖的其他仓库。
//
// 参数:
//   - lg: 日志接口
//   - db: GORM数据库连接
//   - queryTimeout: 查询超时时间
//   - facilityRepo: 活动席位仓库接口
//   - pilotRepo: 活动飞行员仓库接口
//
// 返回值:
//   - *ActivityControllerRepository: 新创建的活动管制员仓库实例
func NewActivityControllerRepository(
	lg logger.Interface,
	db *gorm.DB,
	queryTimeout time.Duration,
	facilityRepo repository.ActivityFacilityInterface,
	pilotRepo repository.ActivityPilotInterface,
) *ActivityControllerRepository {
	return &ActivityControllerRepository{
		BaseRepository: NewBaseRepository[*entity.ActivityController](lg, "ActivityControllerRepository", db, queryTimeout),
		facilityRepo:   facilityRepo,
		pilotRepo:      pilotRepo,
	}
}

// New 创建一个新的活动管制员实体
//
// 该函数用于创建一个ActivityController实体，用于表示用户在某个活动中控制某个席位的关系。
// 如果传入的参数不符合要求（如ID为0），则返回nil。
//
// 参数:
//   - activityId: 活动ID，必须大于0
//   - facilityId: 席位ID，必须大于0
//   - userId: 用户ID，必须大于0
//
// 返回值:
//   - *entity.ActivityController: 创建成功的活动管制员实体指针，如果参数无效则返回nil
func (repo *ActivityControllerRepository) New(activityId uint, facilityId uint, userId uint) *entity.ActivityController {
	if activityId <= 0 || facilityId <= 0 || userId <= 0 {
		return nil
	}

	return &entity.ActivityController{
		ActivityId: activityId,
		FacilityId: facilityId,
		UserId:     userId,
	}
}

// GetById 根据ID获取活动管制员信息
//
// 该函数通过主键ID查询数据库中的活动管制员记录。
// 如果参数无效或查询不到匹配记录，则返回相应错误。
//
// 参数:
//   - id: 活动管制员记录的主键ID，必须大于0
//
// 返回值:
//   - *entity.ActivityController: 查询到的活动管制员实体指针
//   - error: 可能发生的错误，包括参数错误、记录不存在或数据库查询错误
func (repo *ActivityControllerRepository) GetById(id uint) (*entity.ActivityController, error) {
	if id <= 0 {
		return nil, repository.ErrArgument
	}

	activityController := &entity.ActivityController{ID: id}
	err := repo.query(func(tx *gorm.DB) error {
		return tx.First(activityController).Error
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrActivityControllerNotFound
	}

	return activityController, err
}

// GetByActivityIdAndFacilityIdAndUserId 根据活动ID、席位ID和用户ID获取活动管制员信息
//
// 该函数通过活动ID、席位ID和用户ID联合条件查询数据库中的活动管制员记录。
// 如果参数无效或查询不到匹配记录，则返回相应错误。
//
// 参数:
//   - activityId: 活动ID，必须大于0
//   - facilityId: 席位ID，必须大于0
//   - userId: 用户ID，必须大于0
//
// 返回值:
//   - *entity.ActivityController: 查询到的活动管制员实体指针
//   - error: 可能发生的错误，包括参数错误、记录不存在或数据库查询错误
func (repo *ActivityControllerRepository) GetByActivityIdAndFacilityIdAndUserId(
	activityId uint,
	facilityId uint,
	userId uint,
) (*entity.ActivityController, error) {
	if activityId <= 0 || facilityId <= 0 || userId <= 0 {
		return nil, repository.ErrArgument
	}

	activityController := &entity.ActivityController{}
	err := repo.query(func(tx *gorm.DB) error {
		return tx.Where("activity_id = ? AND facility_id = ? AND user_id = ?", activityId, facilityId, userId).First(activityController).Error
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrActivityControllerNotFound
	}

	return activityController, err
}

func (repo *ActivityControllerRepository) GetByActivityIdAndUserId(
	activityId uint,
	userId uint,
) (*entity.ActivityController, error) {
	if activityId <= 0 || userId <= 0 {
		return nil, repository.ErrArgument
	}

	activityController := &entity.ActivityController{}
	err := repo.query(func(tx *gorm.DB) error {
		return tx.Where("activity_id = ? AND user_id = ?", activityId, userId).First(activityController).Error
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrActivityControllerNotFound
	}

	return activityController, err
}

// Save 保存活动管制员信息
//
// 该函数用于保存活动管制员实体到数据库。根据实体的ID判断执行创建还是更新操作：
// 如果ID为0，则执行创建操作；否则执行更新操作。
//
// 参数:
//   - entity: 要保存的活动管制员实体指针
//
// 返回值:
//   - error: 保存过程中可能发生的错误，如实体为空或数据库操作错误
func (repo *ActivityControllerRepository) Save(entity *entity.ActivityController) error {
	return repo.save(entity)
}

// Delete 删除活动管制员信息
//
// 该函数用于从数据库中删除指定的活动管制员记录。
// 删除前会验证实体是否有效以及是否包含有效的ID。
//
// 参数:
//   - entity: 要删除的活动管制员实体指针
//
// 返回值:
//   - error: 删除过程中可能发生的错误，如实体为空、参数无效或数据库操作错误
func (repo *ActivityControllerRepository) Delete(entity *entity.ActivityController) error {
	return repo.delete(entity)
}

// Update 更新活动管制员信息
//
// 该函数用于更新活动管制员实体的部分字段。通过传入更新映射，
// 可以只更新指定的字段，而不是整个实体。
//
// 参数:
//   - entity: 要更新的活动管制员实体指针
//   - updates: 包含要更新的字段和值的映射
//
// 返回值:
//   - error: 更新过程中可能发生的错误，如实体为空、参数无效或数据库操作错误
func (repo *ActivityControllerRepository) Update(entity *entity.ActivityController, updates map[string]interface{}) error {
	return repo.update(entity, updates)
}

// JoinActivity 允许用户作为管制员加入活动席位
//
// 该函数处理用户申请成为特定活动席位管制员的完整业务流程，包括权限验证、
// 并发控制和数据一致性检查。
//
// 参数:
//   - activity: 活动实体指针
//   - activityFacility: 要加入的活动席位实体指针
//   - user: 申请加入的用户实体指针
//
// 返回值:
//   - error: 操作过程中可能发生的错误
func (repo *ActivityControllerRepository) JoinActivity(
	activity *entity.Activity,
	activityFacility *entity.ActivityFacility,
	user *entity.User,
) error {
	if activity == nil || activityFacility == nil || user == nil {
		return repository.ErrArgument
	}

	if activity.DeletedAt.Valid {
		return repository.ErrActivityDeleted
	}

	if activity.Status == repository.ActivityStatusEnded.Value {
		return repository.ErrActivityEnded
	}

	if activity.ID != activityFacility.ActivityId {
		return repository.ErrDataConflicts
	}

	if user.Rating < activityFacility.MinRating || (activityFacility.Tier2Tower && !user.Tier2) {
		return repository.ErrRatingNotAllowed
	}

	processFunc := func(tx *gorm.DB) error {
		activityFacility, err := repo.facilityRepo.GetById(activityFacility.ID)
		if err != nil {
			return err
		}

		if activityFacility.Controller != nil {
			if activityFacility.Controller.UserId == user.ID {
				return repository.ErrFacilityYouSigned
			}
			return repository.ErrFacilityOtherSigned
		}

		_, err = repo.GetByActivityIdAndUserId(activityFacility.ActivityId, user.ID)
		if err == nil {
			return repository.ErrControllerAlreadySign
		}
		if !errors.Is(err, repository.ErrActivityControllerNotFound) {
			return err
		}

		_, err = repo.pilotRepo.GetByActivityIdAndUserId(activityFacility.ActivityId, user.ID)
		if err == nil {
			return repository.ErrPilotAlreadySigned
		}
		if !errors.Is(err, repository.ErrActivityPilotNotFound) {
			return err
		}

		activityController := repo.New(activityFacility.ActivityId, activityFacility.ID, user.ID)
		if activityController == nil {
			return repository.ErrArgument
		}

		err = tx.Create(activityController).Error

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return repository.ErrFacilityOtherSigned
		}

		return nil
	}

	return repo.queryWithLock(processFunc, clause.Locking{Strength: "UPDATE"})
}

// LeaveActivity 处理用户作为管制员离开活动席位的操作
//
// 该函数允许已签约席位的用户取消其管制员身份，完整实现了退出流程，
// 包括权限验证、数据一致性检查和记录删除操作。
//
// 参数:
//   - activityFacility: 要离开的活动席位实体指针
//   - user: 申请离开的用户实体指针
//
// 返回值:
//   - error: 操作过程中可能发生的错误
func (repo *ActivityControllerRepository) LeaveActivity(
	activity *entity.Activity,
	activityFacility *entity.ActivityFacility,
	user *entity.User,
) error {
	if activity == nil || activityFacility == nil || user == nil {
		return repository.ErrArgument
	}

	deleteFunc := func(tx *gorm.DB) error {
		activityFacility, err := repo.facilityRepo.GetById(activityFacility.ID)
		if err != nil {
			return err
		}

		if activityFacility.Controller == nil {
			return repository.ErrFacilityNotSigned
		}

		if activityFacility.Controller.UserId != user.ID {
			return repository.ErrFacilityNotYourSign
		}

		activityController := activityFacility.Controller

		if activityFacility.ID != activityController.FacilityId ||
			activityFacility.ActivityId != activityController.ActivityId {
			return repository.ErrDataConflicts
		}

		return tx.Delete(activityController).Error
	}

	return repo.queryWithLock(deleteFunc, clause.Locking{Strength: "UPDATE"})
}
