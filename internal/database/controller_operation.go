// Package database
package database

import (
	"context"
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ControllerOperation struct {
	logger       log.LoggerInterface
	db           *gorm.DB
	queryTimeout time.Duration
}

func NewControllerOperation(logger log.LoggerInterface, db *gorm.DB, queryTimeout time.Duration) *ControllerOperation {
	return &ControllerOperation{
		logger:       logger,
		db:           db,
		queryTimeout: queryTimeout,
	}
}

func (controllerOperation *ControllerOperation) GetTotalControllers() (total int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), controllerOperation.queryTimeout)
	defer cancel()
	err = controllerOperation.db.WithContext(ctx).Model(&User{}).Select("id").Where("rating > ?", fsd.Normal).Count(&total).Error
	return
}

func (controllerOperation *ControllerOperation) GetControllers(page, pageSize int) (users []*User, total int64, err error) {
	users = make([]*User, 0, pageSize)
	ctx, cancel := context.WithTimeout(context.Background(), controllerOperation.queryTimeout)
	defer cancel()
	controllerOperation.db.WithContext(ctx).Model(&User{}).Select("id").Where("rating > ?", fsd.Normal).Count(&total)
	err = controllerOperation.db.WithContext(ctx).Offset((page-1)*pageSize).Order("cid").Where("rating > ?", fsd.Normal).Limit(pageSize).Find(&users).Error
	return
}

func (controllerOperation *ControllerOperation) SetControllerRating(user *User, updateInfo map[string]interface{}) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), controllerOperation.queryTimeout)
	defer cancel()
	return controllerOperation.db.Clauses(clause.Locking{Strength: "UPDATE"}).WithContext(ctx).Model(user).Updates(updateInfo).Error
}
