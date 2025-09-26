// Package operation
package operation

import "time"

type History struct {
	ID         uint      `gorm:"primarykey" json:"-"`
	Cid        int       `gorm:"index;not null" json:"-"`
	Callsign   string    `gorm:"size:16;index;not null" json:"callsign"`
	StartTime  time.Time `gorm:"not null" json:"start_time"`
	EndTime    time.Time `gorm:"not null" json:"end_time"`
	OnlineTime int       `gorm:"default:0;not null" json:"online_time"`
	IsAtc      bool      `gorm:"default:0;not null" json:"-"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

// HistoryOperationInterface 联飞记录操作接口定义
type HistoryOperationInterface interface {
	// NewHistory 创建新联飞记录
	NewHistory(cid int, callsign string, isAtc bool) (history *History)
	// SaveHistory 保存联飞记录到数据库, 当err为nil时保存成功
	SaveHistory(history *History) (err error)
	EndRecord(history *History)
	// EndRecordAndSaveHistory 结束联飞记录并保存到数据库, 当err为nil时保存成功
	EndRecordAndSaveHistory(history *History) (err error)
	// GetUserHistory 获取用户最近十次的连线记录, 当err为nil时返回值userHistory有效
	GetUserHistory(cid int) (userHistory *UserHistory, err error)
}

type UserHistory struct {
	Pilots      []History `json:"pilots"`
	Controllers []History `json:"controllers"`
}
