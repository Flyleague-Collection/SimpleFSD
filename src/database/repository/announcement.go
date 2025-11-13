// Package repository 包含公告相关的数据库操作
package repository

import (
	"errors"
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"github.com/half-nothing/simple-fsd/src/interfaces/database/repository"
	"github.com/half-nothing/simple-fsd/src/interfaces/logger"
	"gorm.io/gorm"
)

// AnnouncementRepository 公告仓库，用于处理公告相关的数据库操作
type AnnouncementRepository struct {
	*BaseRepository[*entity.Announcement]
	pageReq PageableInterface[*entity.Announcement]
}

// NewAnnouncementRepository 创建一个新的公告仓库实例
//
// 参数:
//   - lg: 日志记录器接口
//   - db: GORM数据库连接
//   - queryTimeout: 查询超时时间
//
// 返回值:
//   - *AnnouncementRepository: 返回一个新的公告仓库实例
func NewAnnouncementRepository(
	lg logger.Interface,
	db *gorm.DB,
	queryTimeout time.Duration,
) *AnnouncementRepository {
	return &AnnouncementRepository{
		BaseRepository: NewBaseRepository[*entity.Announcement](lg, "AnnouncementRepository", db, queryTimeout),
		pageReq:        NewPageRequest[*entity.Announcement](db),
	}
}

// New 创建一个新的公告实体
//
// 参数:
//   - user: 发布公告的用户实体指针，如果为nil或ID无效则返回nil
//   - content: 公告的内容字符串，不能为空
//   - announcementType: 公告类型，使用 [repository.AnnouncementType] 枚举
//   - important: 是否标记为重要公告
//   - forceShow: 是否强制显示该公告
//
// 返回值:
//   - *entity.Announcement: 新创建的公告实体指针，如果参数校验失败则返回nil
func (repo *AnnouncementRepository) New(
	user *entity.User,
	content string,
	announcementType repository.AnnouncementType,
	important bool,
	forceShow bool,
) *entity.Announcement {
	if user == nil || user.ID <= 0 || content == "" {
		return nil
	}

	return &entity.Announcement{
		PublisherId: user.ID,
		Content:     content,
		Type:        announcementType.Index,
		Important:   important,
		ForceShow:   forceShow,
	}
}

// GetPage 获取公告分页数据
//
// 参数:
//   - pageNumber: 页码
//   - pageSize: 每页大小
//
// 返回值:
//   - []*entity.Announcement: 公告列表
//   - int64: 总数
//   - error: 可能的错误
func (repo *AnnouncementRepository) GetPage(
	pageNumber int,
	pageSize int,
) (announcements []*entity.Announcement, total int64, err error) {
	announcements = make([]*entity.Announcement, 0, pageSize)
	total, err = repo.queryWithPagination(repo.pageReq, NewPage(pageNumber, pageSize, announcements, &entity.Announcement{}, nil))
	return
}

// GetById 根据ID获取公告
//
// 参数:
//   - id: 公告ID
//
// 返回值:
//   - *entity.Announcement: 对应的公告实体
//   - error: 可能的错误，如果ID无效或公告不存在，返回相应的错误
func (repo *AnnouncementRepository) GetById(id uint) (*entity.Announcement, error) {
	if id <= 0 {
		return nil, repository.ErrArgument
	}

	announcement := &entity.Announcement{ID: id}
	err := repo.query(func(tx *gorm.DB) error {
		return tx.Preload("User").First(announcement).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrAnnouncementNotFound
	}

	return announcement, err
}

// Save 保存公告实体到数据库
//
// 参数:
//   - entity: 要保存的公告实体
//
// 返回值:
//   - error: 可能的错误
func (repo *AnnouncementRepository) Save(entity *entity.Announcement) error {
	return repo.save(entity)
}

// Delete 从数据库删除公告实体
//
// 参数:
//   - entity: 要删除的公告实体
//
// 返回值:
//   - error: 可能的错误
func (repo *AnnouncementRepository) Delete(entity *entity.Announcement) error {
	return repo.delete(entity)
}

// Update 更新公告实体的部分字段
//
// 参数:
//   - entity: 要更新的公告实体
//   - updates: 包含要更新的字段和值的映射
//
// 返回值:
//   - error: 可能的错误
func (repo *AnnouncementRepository) Update(entity *entity.Announcement, updates map[string]interface{}) error {
	return repo.update(entity, updates)
}
