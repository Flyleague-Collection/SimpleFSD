// Package repository
package repository

import "github.com/half-nothing/simple-fsd/src/interfaces/database/entity"

type UserHistory struct {
	Pilots      []*entity.History `json:"pilots"`
	Controllers []*entity.History `json:"controllers"`
}

// HistoryInterface 联飞记录操作接口定义
type HistoryInterface interface {
	// NewHistory 创建新联飞记录
	NewHistory(cid int, callsign string, isAtc bool) (history *entity.History)
	// SaveHistory 保存联飞记录到数据库, 当err为nil时保存成功
	SaveHistory(history *entity.History) (err error)
	EndRecord(history *entity.History)
	// EndRecordAndSaveHistory 结束联飞记录并保存到数据库, 当err为nil时保存成功
	EndRecordAndSaveHistory(history *entity.History) (err error)
	// GetUserHistory 获取用户最近十次的连线记录, 当err为nil时返回值userHistory有效
	GetUserHistory(cid int) (userHistory *UserHistory, err error)
}
