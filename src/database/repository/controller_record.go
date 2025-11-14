// Package repository
package repository

import (
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"github.com/half-nothing/simple-fsd/src/interfaces/database/repository"
	"github.com/half-nothing/simple-fsd/src/interfaces/logger"
	"gorm.io/gorm"
)

type ControllerRecordRepository struct {
	*BaseRepository[*entity.ControllerRecord]
	pageReq PageableInterface[*entity.ControllerRecord]
}

func NewControllerRecordRepository(
	lg logger.Interface,
	db *gorm.DB,
	queryTimeout time.Duration,
) *ControllerRecordRepository {
	return &ControllerRecordRepository{
		BaseRepository: NewBaseRepository[*entity.ControllerRecord](lg, "ControllerRecordRepository", db, queryTimeout),
		pageReq:        NewPageRequest[*entity.ControllerRecord](db),
	}
}

func (repo *ControllerRecordRepository) New(
	uid uint,
	operatorCid int,
	recordType repository.ControllerRecordType,
	content string,
) *entity.ControllerRecord {
	if uid <= 0 || operatorCid <= 0 || content == "" {
		return nil
	}

	return &entity.ControllerRecord{
		UserId:      uid,
		OperatorCid: operatorCid,
		Type:        recordType.Value,
		Content:     content,
	}
}

func (repo *ControllerRecordRepository) GetPage(
	uid uint,
	pageNumber int,
	pageSize int,
) (records []*entity.ControllerRecord, total int64, err error) {
	records = make([]*entity.ControllerRecord, 0, pageSize)
	page := NewPage(pageNumber, pageSize, records, &entity.ControllerRecord{}, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("user_id = ?", uid)
	})
	total, err = repo.queryWithPagination(repo.pageReq, page)
	return
}

func (repo *ControllerRecordRepository) GetByIdAndUserId(id uint, uid uint) (*entity.ControllerRecord, error) {
	record := &entity.ControllerRecord{ID: id}
	err := repo.query(func(tx *gorm.DB) error {
		return tx.Where("user_id = ?", uid).First(record).Error
	})
	return record, err
}

func (repo *ControllerRecordRepository) GetById(id uint) (*entity.ControllerRecord, error) {
	record := &entity.ControllerRecord{ID: id}
	err := repo.query(func(tx *gorm.DB) error {
		return tx.First(record).Error
	})
	return record, err
}

func (repo *ControllerRecordRepository) Save(entity *entity.ControllerRecord) error {
	return repo.save(entity)
}

func (repo *ControllerRecordRepository) Delete(entity *entity.ControllerRecord) error {
	return repo.delete(entity)
}

func (repo *ControllerRecordRepository) Update(entity *entity.ControllerRecord, updates map[string]interface{}) error {
	return repo.update(entity, updates)
}
