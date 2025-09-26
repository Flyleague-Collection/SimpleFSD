// Package database
package database

import (
	"context"
	"errors"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"gorm.io/gorm"
	"time"
)

type ControllerApplicationOperation struct {
	logger       log.LoggerInterface
	db           *gorm.DB
	queryTimeout time.Duration
}

func NewControllerApplicationOperation(
	logger log.LoggerInterface,
	db *gorm.DB,
	queryTimeout time.Duration,
) *ControllerApplicationOperation {
	return &ControllerApplicationOperation{
		logger:       logger,
		db:           db,
		queryTimeout: queryTimeout,
	}
}

func (operation *ControllerApplicationOperation) NewApplication(userId uint, reason string, record string, guset bool, platform string, evidence string) *ControllerApplication {
	return &ControllerApplication{
		UserId:                userId,
		WhyWantToBeController: reason,
		ControllerRecord:      record,
		IsGuest:               guset,
		Platform:              platform,
		Evidence:              evidence,
		Status:                int(Submitted),
	}
}

func (operation *ControllerApplicationOperation) GetApplicationByUserId(userId uint) (application *ControllerApplication, err error) {
	application = &ControllerApplication{}
	ctx, cancel := context.WithTimeout(context.Background(), operation.queryTimeout)
	defer cancel()
	err = operation.db.WithContext(ctx).Preload("User").Order("created_at desc").Where("user_id = ?", userId).First(application).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrApplicationNotFound
	}
	return
}

func (operation *ControllerApplicationOperation) GetApplicationById(id uint) (application *ControllerApplication, err error) {
	application = &ControllerApplication{}
	ctx, cancel := context.WithTimeout(context.Background(), operation.queryTimeout)
	defer cancel()
	err = operation.db.WithContext(ctx).Preload("User").Order("created_at desc").First(application, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrApplicationNotFound
	}
	return
}

func (operation *ControllerApplicationOperation) GetApplications(page, pageSize int) (applications []*ControllerApplication, total int64, err error) {
	applications = make([]*ControllerApplication, 0, pageSize)
	ctx, cancel := context.WithTimeout(context.Background(), operation.queryTimeout)
	defer cancel()
	operation.db.WithContext(ctx).Model(&ControllerApplication{}).Select("id").Count(&total)
	err = operation.db.WithContext(ctx).Preload("User").Offset((page - 1) * pageSize).Order("created_at desc").Limit(pageSize).Find(&applications).Error
	return
}

func (operation *ControllerApplicationOperation) SaveApplication(application *ControllerApplication) error {
	ctx, cancel := context.WithTimeout(context.Background(), operation.queryTimeout)
	defer cancel()
	if application.ID == 0 {
		applicationCheck, err := operation.GetApplicationByUserId(application.UserId)
		if err == nil && applicationCheck != nil && applicationCheck.Status != int(Passed) && applicationCheck.Status != int(Rejected) {
			return ErrApplicationAlreadyExists
		}
		return operation.db.WithContext(ctx).Create(application).Error
	}
	return operation.db.WithContext(ctx).Save(application).Error
}

func (operation *ControllerApplicationOperation) ConfirmApplicationUnderProcessing(application *ControllerApplication) error {
	ctx, cancel := context.WithTimeout(context.Background(), operation.queryTimeout)
	defer cancel()
	err := operation.db.WithContext(ctx).Model(application).Updates(&ControllerApplication{Status: int(UnderProcessing)}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrApplicationNotFound
	}
	return err
}

func (operation *ControllerApplicationOperation) UpdateApplicationStatus(application *ControllerApplication, status ControllerApplicationStatus, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), operation.queryTimeout)
	defer cancel()
	err := operation.db.WithContext(ctx).Model(application).Updates(&ControllerApplication{Status: int(status), Message: message}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrApplicationNotFound
	}
	return err
}

func (operation *ControllerApplicationOperation) CancelApplication(application *ControllerApplication) error {
	ctx, cancel := context.WithTimeout(context.Background(), operation.queryTimeout)
	defer cancel()
	err := operation.db.WithContext(ctx).Delete(application).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrApplicationNotFound
	}
	return err
}
