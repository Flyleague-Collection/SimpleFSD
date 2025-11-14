// Package repository
package repository

import (
	"errors"
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"github.com/half-nothing/simple-fsd/src/interfaces/database/repository"
	"github.com/half-nothing/simple-fsd/src/interfaces/logger"
	"gorm.io/gorm"
)

type ControllerApplicationRepository struct {
	*BaseRepository[*entity.ControllerApplication]
	pageReq PageableInterface[*entity.ControllerApplication]
}

func NewControllerApplicationRepository(
	lg logger.Interface,
	db *gorm.DB,
	queryTimeout time.Duration,
) *ControllerApplicationRepository {
	return &ControllerApplicationRepository{
		BaseRepository: NewBaseRepository[*entity.ControllerApplication](lg, "ControllerApplicationRepository", db, queryTimeout),
		pageReq:        NewPageRequest[*entity.ControllerApplication](db),
	}
}

func (repo *ControllerApplicationRepository) New(
	user *entity.User,
	reason string,
	record string,
	guset bool,
	platform string,
	evidence string,
) *entity.ControllerApplication {
	if user == nil || user.ID <= 0 || reason == "" || record == "" || (guset && (platform == "" || evidence == "")) {
		return nil
	}

	return &entity.ControllerApplication{
		UserId:                user.ID,
		WhyWantToBeController: reason,
		ControllerRecord:      record,
		IsGuest:               guset,
		Platform:              platform,
		Evidence:              evidence,
		Status:                repository.ApplicationStatusSubmitted.Value,
	}
}

func (repo *ControllerApplicationRepository) GetByUserId(userId uint) (*entity.ControllerApplication, error) {
	if userId <= 0 {
		return nil, repository.ErrArgument
	}
	controllerApplication := &entity.ControllerApplication{}
	err := repo.query(func(tx *gorm.DB) error {
		return tx.Where("user_id = ?", userId).First(controllerApplication).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = repository.ErrApplicationNotFound
	}
	return controllerApplication, err
}

func (repo *ControllerApplicationRepository) GetPage(
	pageNumber int,
	pageSize int,
) (applications []*entity.ControllerApplication, total int64, err error) {
	applications = make([]*entity.ControllerApplication, 0, pageSize)
	total, err = repo.queryWithPagination(repo.pageReq, NewPage(pageNumber, pageSize, applications, &entity.ControllerApplication{}, nil))
	return
}

func (repo *ControllerApplicationRepository) GetById(id uint) (*entity.ControllerApplication, error) {
	if id <= 0 {
		return nil, repository.ErrArgument
	}

	application := &entity.ControllerApplication{ID: id}
	err := repo.query(func(tx *gorm.DB) error {
		return tx.First(application).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = repository.ErrApplicationNotFound
	}

	return application, err
}

func (repo *ControllerApplicationRepository) Save(entity *entity.ControllerApplication) error {
	return repo.save(entity)
}

func (repo *ControllerApplicationRepository) Delete(entity *entity.ControllerApplication) error {
	return repo.delete(entity)
}

func (repo *ControllerApplicationRepository) Update(entity *entity.ControllerApplication, updates map[string]interface{}) error {
	return repo.update(entity, updates)
}

func (repo *ControllerApplicationRepository) UpdateStatus(application *entity.ControllerApplication, status repository.ApplicationStatus, message string) error {
	return repo.update(application, map[string]interface{}{
		"status":  status.Value,
		"message": message,
	})
}
