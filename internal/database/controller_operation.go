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

func (controllerOperation *ControllerOperation) SetControllerRating(user *User, rating int) (err error) {
	user.Rating = rating
	ctx, cancel := context.WithTimeout(context.Background(), controllerOperation.queryTimeout)
	defer cancel()
	return controllerOperation.db.Clauses(clause.Locking{Strength: "UPDATE"}).WithContext(ctx).Model(user).Updates(&User{Rating: rating}).Error
}

func (controllerOperation *ControllerOperation) SetControllerSolo(user *User, untilTime time.Time) (err error) {
	user.UnderSolo = true
	user.SoloUntil = untilTime
	ctx, cancel := context.WithTimeout(context.Background(), controllerOperation.queryTimeout)
	defer cancel()
	return controllerOperation.db.Clauses(clause.Locking{Strength: "UPDATE"}).WithContext(ctx).Model(user).Updates(&User{UnderSolo: true, SoloUntil: untilTime}).Error
}

func (controllerOperation *ControllerOperation) UnsetControllerSolo(user *User) (err error) {
	user.UnderSolo = false
	ctx, cancel := context.WithTimeout(context.Background(), controllerOperation.queryTimeout)
	defer cancel()
	return controllerOperation.db.Clauses(clause.Locking{Strength: "UPDATE"}).WithContext(ctx).Model(user).Updates(&User{UnderSolo: false}).Error
}

func (controllerOperation *ControllerOperation) SetControllerUnderMonitor(user *User, underMonitor bool) (err error) {
	user.UnderMonitor = underMonitor
	ctx, cancel := context.WithTimeout(context.Background(), controllerOperation.queryTimeout)
	defer cancel()
	return controllerOperation.db.Clauses(clause.Locking{Strength: "UPDATE"}).WithContext(ctx).Model(user).Updates(&User{UnderMonitor: underMonitor}).Error
}

func (controllerOperation *ControllerOperation) SetControllerGuest(user *User, guest bool) (err error) {
	user.Guest = guest
	ctx, cancel := context.WithTimeout(context.Background(), controllerOperation.queryTimeout)
	defer cancel()
	return controllerOperation.db.Clauses(clause.Locking{Strength: "UPDATE"}).WithContext(ctx).Model(user).Updates(&User{Guest: guest}).Error
}

func (controllerOperation *ControllerOperation) SetControllerGuestRating(user *User, rating int) (err error) {
	user.Guest = true
	user.Rating = rating
	ctx, cancel := context.WithTimeout(context.Background(), controllerOperation.queryTimeout)
	defer cancel()
	return controllerOperation.db.Clauses(clause.Locking{Strength: "UPDATE"}).WithContext(ctx).Model(user).Updates(&User{Guest: true, Rating: rating}).Error
}
