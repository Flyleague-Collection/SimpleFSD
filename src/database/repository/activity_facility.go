// Package repository
package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"github.com/half-nothing/simple-fsd/src/interfaces/database/repository"
	"github.com/half-nothing/simple-fsd/src/interfaces/logger"
	"gorm.io/gorm"
)

type ActivityFacilityRepository struct {
	*BaseRepository[*entity.ActivityFacility]
	activityControllerRepo repository.ActivityControllerInterface
}

func NewActivityFacilityRepository(
	lg logger.Interface,
	db *gorm.DB,
	queryTimeout time.Duration,
	activityControllerRepo repository.ActivityControllerInterface,
) *ActivityFacilityRepository {
	return &ActivityFacilityRepository{
		BaseRepository:         NewBaseRepository[*entity.ActivityFacility](lg, "ActivityFacilityRepository", db, queryTimeout),
		activityControllerRepo: activityControllerRepo,
	}
}

func (repo *ActivityFacilityRepository) New(activity *entity.Activity, minRating int, callsign string, frequency float64, tier2Tower bool) *entity.ActivityFacility {
	return &entity.ActivityFacility{
		ActivityId: activity.ID,
		MinRating:  minRating,
		Callsign:   callsign,
		Frequency:  fmt.Sprintf("%.3f", frequency),
		Tier2Tower: tier2Tower,
	}
}

func (repo *ActivityFacilityRepository) GetById(id uint) (*entity.ActivityFacility, error) {
	if id <= 0 {
		return nil, repository.ErrArgument
	}

	activityFacility := &entity.ActivityFacility{ID: id}
	err := repo.query(func(tx *gorm.DB) error {
		return tx.Preload("Controller").First(activityFacility).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = repository.ErrFacilityNotFound
	}

	return activityFacility, err
}

func (repo *ActivityFacilityRepository) Save(entity *entity.ActivityFacility) error {
	return repo.save(entity)
}

func (repo *ActivityFacilityRepository) Delete(entity *entity.ActivityFacility) error {
	return repo.delete(entity)
}

func (repo *ActivityFacilityRepository) Update(entity *entity.ActivityFacility, updates map[string]interface{}) error {
	return repo.update(entity, updates)
}
