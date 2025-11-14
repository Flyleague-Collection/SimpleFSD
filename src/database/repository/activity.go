package repository

import (
	"errors"
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"github.com/half-nothing/simple-fsd/src/interfaces/database/repository"
	"github.com/half-nothing/simple-fsd/src/interfaces/logger"
	"github.com/half-nothing/simple-fsd/src/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ActivityRepository struct {
	*BaseRepository[*entity.Activity]
	pageReq        PageableInterface[*entity.Activity]
	controllerRepo repository.ActivityControllerInterface
	pilotRepo      repository.ActivityPilotInterface
	facilityRepo   repository.ActivityFacilityInterface
}

func NewActivityRepository(
	lg logger.Interface,
	db *gorm.DB,
	queryTimeout time.Duration,
	controllerRepo repository.ActivityControllerInterface,
	pilotRepo repository.ActivityPilotInterface,
	facilityRepo repository.ActivityFacilityInterface,
) *ActivityRepository {
	return &ActivityRepository{
		BaseRepository: NewBaseRepository[*entity.Activity](lg, "ActivityRepository", db, queryTimeout),
		pageReq:        NewPageRequest[*entity.Activity](db),
		controllerRepo: controllerRepo,
		pilotRepo:      pilotRepo,
		facilityRepo:   facilityRepo,
	}
}

func (repo *ActivityRepository) NewBuilder(user *entity.User, title string, image *entity.Image, activeTime time.Time, notams string) *repository.ActivityBuilder {
	if user == nil || user.ID <= 0 || title == "" || image == nil || image.ID <= 0 || activeTime.IsZero() {
		return nil
	}

	return repository.NewActivityBuilder().
		SetTitle(title).
		SetImage(image).
		SetPublisher(user.Cid).
		SetActiveTime(activeTime).
		SetNOTAMS(notams)
}

func (repo *ActivityRepository) NewOneWay(builder *repository.ActivityBuilder, dep string, arr string, route string, distance int) *entity.Activity {
	if builder == nil || dep == "" || arr == "" || route == "" || distance <= 0 {
		return nil
	}

	return builder.SetType(repository.ActivityTypeOneWay).
		SetDepartureAirport(dep).
		SetArrivalAirport(arr).
		SetRoute(route).
		SetDistance(distance).
		Build()
}

func (repo *ActivityRepository) NewBothWay(builder *repository.ActivityBuilder, dep string, arr string, route string, distance int, route2 string, distance2 int) *entity.Activity {
	if builder == nil || dep == "" || arr == "" || route == "" || distance <= 0 || route2 == "" || distance2 <= 0 {
		return nil
	}

	return builder.SetType(repository.ActivityTypeBothWay).
		SetDepartureAirport(dep).
		SetArrivalAirport(arr).
		SetRoute(route).
		SetDistance(distance).
		SetRoute2(route2).
		SetDistance2(distance2).
		Build()
}

func (repo *ActivityRepository) NewFIROpenDay(builder *repository.ActivityBuilder, firs ...string) *entity.Activity {
	if builder == nil || len(firs) == 0 {
		return nil
	}

	return builder.SetType(repository.ActivityTypeFIROpenDay).
		SetOpenFir(firs).
		Build()
}

// GetById 根据ID获取活动信息
// 该方法会预加载活动相关的设施、飞行员和控制器等关联信息
// 参数:
//   - id: 活动的唯一标识符
//
// 返回值:
//   - *entity.Activity: 对应ID的活动实体指针
//   - error: 查询过程中可能发生的错误，如果记录不存在则返回ErrActivityNotFound
func (repo *ActivityRepository) GetById(id uint) (*entity.Activity, error) {
	activity := &entity.Activity{ID: id}

	queryBuilder := NewQueryBuilder[*entity.Activity]()
	queryBuilder.Preload("Image", func(db *gorm.DB) *gorm.DB {
		return db.Select("url")
	})
	queryBuilder.Preload("Facilities", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_index DESC")
	})
	queryBuilder.Preload("Facilities.Controller")
	queryBuilder.Preload("Facilities.Controller.User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, cid, avatar_url")
	})
	queryBuilder.Preload("Pilots.User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, cid, avatar_url")
	})
	queryBuilder.Preload("Controllers")

	err := repo.queryEntityWithBuilder(queryBuilder, activity)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = repository.ErrActivityNotFound
	}

	return activity, err
}

func (repo *ActivityRepository) Save(entity *entity.Activity) error {
	return repo.saveWithLock(entity, nil)
}

func (repo *ActivityRepository) Delete(entity *entity.Activity) error {
	return repo.deleteWithLock(entity, nil)
}

func (repo *ActivityRepository) Update(entity *entity.Activity, updates map[string]interface{}) error {
	return repo.updateWithLock(entity, updates, nil)
}

func (repo *ActivityRepository) GetNumber() (total int64, err error) {
	err = repo.query(func(tx *gorm.DB) error {
		return tx.Model(&entity.Activity{}).Select("id").Count(&total).Error
	})
	return
}

func (repo *ActivityRepository) GetBetween(
	startDay time.Time,
	endDay time.Time,
) (activities []*entity.Activity, err error) {
	activities = make([]*entity.Activity, 0)
	err = repo.query(func(tx *gorm.DB) error {
		return tx.Where("active_time between ? and ?", startDay, endDay).Find(activities).Error
	})
	return
}

func (repo *ActivityRepository) GetPage(
	pageNumber int, pageSize int,
) (activities []*entity.Activity, total int64, err error) {
	activities = make([]*entity.Activity, 0, pageSize)
	total, err = repo.queryWithPagination(repo.pageReq, NewPage(pageNumber, pageSize, activities, &entity.Activity{}, nil))
	return
}

func (repo *ActivityRepository) UpdateStatus(activityId uint, status repository.ActivityStatus) error {
	processFunc := func(tx *gorm.DB) error {
		result := tx.Model(&entity.Activity{ID: activityId}).Update("status", status.Value)
		if result.RowsAffected == 0 {
			return repository.ErrActivityNotFound
		}
		return result.Error
	}

	return repo.queryWithLock(processFunc, nil)
}

// UpdateInfo 更新活动信息
// 比较新旧活动信息，计算出需要更新、插入和删除的设施席位，并执行相应的数据库操作
// 参数:
//   - oldActivity: 原始活动信息
//   - newActivity: 新的活动信息
//
// 返回值:
//   - error: 更新过程中可能出现的错误
func (repo *ActivityRepository) UpdateInfo(oldActivity *entity.Activity, newActivity *entity.Activity) error {
	deleteFacilities, insertFacilities, updateFacilities, err := repo.preprocessActivityData(oldActivity, newActivity)
	if err != nil {
		return err
	}

	// 定义数据库操作过程
	processFunc := func(tx *gorm.DB) error {
		// 获取最新的活动信息
		dbActivity, err := repo.GetById(oldActivity.ID)
		if err != nil {
			return err
		}

		// 计算活动基本信息的更新内容
		updateInfo := dbActivity.Diff(newActivity)
		if updateInfo != nil {
			if err := tx.Model(dbActivity).Updates(updateInfo).Error; err != nil {
				return err
			}
		}

		// 外键已配置级联删除, 这里直接批量删除即可
		if len(deleteFacilities) > 0 {
			if err := tx.Delete(deleteFacilities).Error; err != nil {
				return err
			}
		}

		if len(insertFacilities) > 0 {
			if err := tx.Create(insertFacilities).Error; err != nil {
				return err
			}
		}

		if len(updateFacilities) > 0 {
			if err := tx.Model(&entity.ActivityFacility{}).Updates(updateFacilities).Error; err != nil {
				return err
			}
		}
		return nil
	}

	// 执行带锁的数据库操作
	return repo.queryWithLock(processFunc, clause.Locking{Strength: "UPDATE"})
}

func (repo *ActivityRepository) preprocessActivityData(
	oldActivity *entity.Activity,
	newActivity *entity.Activity,
) ([]*entity.ActivityFacility, []*entity.ActivityFacility, []map[string]interface{}, error) {
	if oldActivity == nil || newActivity == nil ||
		oldActivity.ID == 0 || newActivity.ID != oldActivity.ID ||
		oldActivity.Facilities == nil || newActivity.Facilities == nil {
		return nil, nil, nil, repository.ErrArgument
	}

	oldFacilities := utils.Filter(oldActivity.Facilities, utils.FilterNotNull)
	newFacilities := utils.Filter(newActivity.Facilities, utils.FilterNotNull)

	// 设置席位的排序索引
	index := len(newFacilities)
	for _, facility := range newFacilities {
		index--
		facility.SortIndex = index
	}

	// 初始化需要操作的设施席位列表
	deleteFacilities := make([]*entity.ActivityFacility, 0)
	insertFacilities := make([]*entity.ActivityFacility, 0)
	updateFacilities := make([]map[string]interface{}, 0)

	// 创建旧席位的映射，便于快速查找
	oldFacilityMap := make(map[uint]*entity.ActivityFacility)
	utils.ForEach(oldFacilities, func(index int, facility *entity.ActivityFacility) {
		oldFacilityMap[facility.ID] = facility
	})

	// 处理新席位
	for _, newFacility := range newFacilities {
		if newFacility.ID == 0 {
			// 新增席位
			newFacility.ActivityId = newActivity.ID
			insertFacilities = append(insertFacilities, newFacility)
		} else if oldFacility, exists := oldFacilityMap[newFacility.ID]; exists {
			// 检查是否需要更新
			if !oldFacility.Equal(newFacility) {
				updateData := oldFacility.Diff(newFacility)
				if updateData != nil && len(updateData) > 0 {
					updateData["id"] = oldFacility.ID
					updateFacilities = append(updateFacilities, updateData)
				}
			}
			// 从映射中移除
			delete(oldFacilityMap, newFacility.ID)
		} else {
			// 注意：如果newFacility.ID不为0但在oldFacilityMap中不存在
			// 这可能表示数据不一致，可能是脏读导致的数据错误，也可能是有人恶意注入，直接报错退出就行
			return nil, nil, nil, repository.ErrDataConflicts
		}
	}

	// 剩余在oldFacilityMap中的席位是需要删除的
	for _, facility := range oldFacilityMap {
		deleteFacilities = append(deleteFacilities, facility)
	}

	return deleteFacilities, insertFacilities, updateFacilities, nil
}

func (repo *ActivityRepository) GetPilotRepository() repository.ActivityPilotInterface {
	return repo.pilotRepo
}

func (repo *ActivityRepository) GetControllerRepository() repository.ActivityControllerInterface {
	return repo.controllerRepo
}

func (repo *ActivityRepository) GetFacilityRepository() repository.ActivityFacilityInterface {
	return repo.facilityRepo
}
