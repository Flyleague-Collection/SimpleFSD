// Package database
package database

import (
	"context"
	"errors"
	"time"

	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"gorm.io/gorm"
)

type ControllerRecordOperation struct {
	logger       log.LoggerInterface
	db           *gorm.DB
	queryTimeout time.Duration
}

func NewControllerRecordOperation(logger log.LoggerInterface, db *gorm.DB, queryTimeout time.Duration) *ControllerRecordOperation {
	return &ControllerRecordOperation{
		logger:       logger,
		db:           db,
		queryTimeout: queryTimeout,
	}
}

func (controllerRecordOperation *ControllerRecordOperation) NewControllerRecord(uid uint, operatorCid int, recordType ControllerRecordType, content string) (record *ControllerRecord) {
	return &ControllerRecord{
		UserId:      uid,
		OperatorCid: operatorCid,
		Type:        int(recordType),
		Content:     content,
	}
}

func (controllerRecordOperation *ControllerRecordOperation) SaveControllerRecord(record *ControllerRecord) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), controllerRecordOperation.queryTimeout)
	defer cancel()
	if record.ID == 0 {
		return controllerRecordOperation.db.WithContext(ctx).Create(record).Error
	}
	return controllerRecordOperation.db.WithContext(ctx).Save(record).Error
}

func (controllerRecordOperation *ControllerRecordOperation) GetControllerRecords(uid uint, page, pageSize int) (records []*ControllerRecord, total int64, err error) {
	records = make([]*ControllerRecord, 0, pageSize)
	ctx, cancel := context.WithTimeout(context.Background(), controllerRecordOperation.queryTimeout)
	defer cancel()
	controllerRecordOperation.db.WithContext(ctx).Model(&ControllerRecord{}).Select("id").Where("user_id = ?", uid).Count(&total)
	err = controllerRecordOperation.db.WithContext(ctx).Offset((page-1)*pageSize).Where("user_id = ?", uid).Order("created_at desc").Limit(pageSize).Find(&records).Error
	return
}

func (controllerRecordOperation *ControllerRecordOperation) GetControllerRecord(id uint, uid uint) (record *ControllerRecord, err error) {
	record = &ControllerRecord{}
	ctx, cancel := context.WithTimeout(context.Background(), controllerRecordOperation.queryTimeout)
	defer cancel()
	err = controllerRecordOperation.db.WithContext(ctx).First(record, "id = ? and user_id = ?", id, uid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrControllerRecordNotFound
	}
	return
}

func (controllerRecordOperation *ControllerRecordOperation) DeleteControllerRecord(id uint) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), controllerRecordOperation.queryTimeout)
	defer cancel()
	return controllerRecordOperation.db.WithContext(ctx).Delete(&ControllerRecord{}, id).Error
}
