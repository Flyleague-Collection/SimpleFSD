package database

import (
	"context"
	"time"

	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"gorm.io/gorm"
)

type HistoryOperation struct {
	logger       log.LoggerInterface
	db           *gorm.DB
	queryTimeout time.Duration
}

func NewHistoryOperation(logger log.LoggerInterface, db *gorm.DB, queryTimeout time.Duration) *HistoryOperation {
	return &HistoryOperation{logger: logger, db: db, queryTimeout: queryTimeout}
}

func (historyOperation *HistoryOperation) NewHistory(cid int, callsign string, isAtc bool) (history *History) {
	return &History{
		Cid:        cid,
		Callsign:   callsign,
		StartTime:  time.Now(),
		EndTime:    time.Now(),
		OnlineTime: 0,
		IsAtc:      isAtc,
	}
}

func (historyOperation *HistoryOperation) SaveHistory(history *History) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), historyOperation.queryTimeout)
	defer cancel()
	if history.ID == 0 {
		return historyOperation.db.WithContext(ctx).Create(history).Error
	}
	return historyOperation.db.WithContext(ctx).Save(history).Error
}

func (historyOperation *HistoryOperation) EndRecord(history *History) {
	history.EndTime = time.Now()
	history.OnlineTime = int(history.EndTime.Sub(history.StartTime).Seconds())
}

func (historyOperation *HistoryOperation) EndRecordAndSaveHistory(history *History) (err error) {
	history.EndTime = time.Now()
	history.OnlineTime = int(history.EndTime.Sub(history.StartTime).Seconds())
	return historyOperation.SaveHistory(history)
}

func (historyOperation *HistoryOperation) GetUserHistory(cid int) (userHistory *UserHistory, err error) {
	userHistory = &UserHistory{
		Pilots:      make([]History, 0, 10),
		Controllers: make([]History, 0, 10),
	}
	ctx, cancel := context.WithTimeout(context.Background(), historyOperation.queryTimeout)
	defer cancel()
	err = historyOperation.db.WithContext(ctx).Order("id desc").Where("cid = ? and is_atc = ?", cid, false).Limit(10).Find(&userHistory.Pilots).Error
	if err != nil {
		return
	}
	err = historyOperation.db.WithContext(ctx).Order("id desc").Where("cid = ? and is_atc = ?", cid, true).Limit(10).Find(&userHistory.Controllers).Error
	if err != nil {
		return
	}
	return
}
