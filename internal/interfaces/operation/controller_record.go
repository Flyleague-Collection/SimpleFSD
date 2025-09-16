// Package operation
package operation

import (
	"time"
)

type ControllerRecord struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Type      int       `gorm:"not null" json:"type"`
	Cid       int       `gorm:"index:Cid;not null" json:"cid"`
	Operator  int       `gorm:"index:Operator;not null" json:"operator"`
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `gorm:"not null" json:"time"`
}

type ControllerRecordType int

const (
	Application  ControllerRecordType = iota // 管制员申请
	Interview                                // 面试
	Simulator                                // 模拟机
	RatingChange                             // 权限变动
	Training                                 // 训练内容
	UnderMonitor                             // UM权限授予
	Solo                                     // Solo权限授予
	Guest                                    // 客座权限变动
	Other                                    // 其他未定义内容
)

func IsValidControllerRecordType(s int) bool {
	return int(Application) <= s && s <= int(Other)
}

func ToControllerRecordType(s int) ControllerRecordType {
	if !IsValidControllerRecordType(s) {
		return Other
	}
	return ControllerRecordType(s)
}

type ControllerRecordOperationInterface interface {
	NewControllerRecord(cid, operator int, recordType ControllerRecordType, content string) (record *ControllerRecord)
	SaveControllerRecord(record *ControllerRecord) (err error)
	GetControllerRecords(cid, page, pageSize int) (records []*ControllerRecord, total int64, err error)
	GetControllerRecord(id uint) (record *ControllerRecord, err error)
	DeleteControllerRecord(id uint) (err error)
}
