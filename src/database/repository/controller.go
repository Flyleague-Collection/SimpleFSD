// Package repository
package repository

import (
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"github.com/half-nothing/simple-fsd/src/interfaces/logger"
	"gorm.io/gorm"
)

type ControllerRepository struct {
	*BaseRepository[*entity.User]
	pageReq PageableInterface[*entity.User]
}

func NewControllerRepository(
	lg logger.Interface,
	db *gorm.DB,
	queryTimeout time.Duration,
) *ControllerRepository {
	return &ControllerRepository{
		BaseRepository: NewBaseRepository[*entity.User](lg, "ControllerRepository", db, queryTimeout),
		pageReq:        NewPageRequest[*entity.User](db),
	}
}

func (repo *ControllerRepository) GetTotal() (total int64, err error) {
	// TODO: FSD权限重构完成后修改此处的权限
	err = repo.query(func(tx *gorm.DB) error {
		return tx.Model(&entity.User{}).Where("rating > ?", 1).Count(&total).Error
	})
	return
}

func (repo *ControllerRepository) GetPage(pageNumber int, pageSize int) (users []*entity.User, total int64, err error) {
	users = make([]*entity.User, 0, pageSize)
	total, err = repo.queryWithPagination(repo.pageReq, NewPage(pageNumber, pageSize, users, &entity.User{}, nil))
	return
}

func (repo *ControllerRepository) SetRating(user *entity.User, updateInfo map[string]interface{}) error {
	return repo.update(user, updateInfo)
}
