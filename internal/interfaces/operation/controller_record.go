// Package operation
package operation

import (
	"errors"
	"time"
)

type ControllerRecord struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Type        int       `gorm:"not null" json:"type"`
	UserId      uint      `gorm:"index:Uid;not null" json:"uid"`
	OperatorCid int       `gorm:"index:OperatorCid;not null" json:"operator_cid"`
	Content     string    `gorm:"not null" json:"content"`
	CreatedAt   time.Time `gorm:"not null" json:"time"`
}

type ControllerRecordType int

const (
	Interview    ControllerRecordType = iota // 面试
	Simulator                                // 模拟机
	RatingChange                             // 权限变动
	Training                                 // 训练内容
	UnderMonitor                             // UM权限授予
	Solo                                     // Solo权限授予
	Guest                                    // 客座权限变动
	Other                                    // 其他未定义内容
)

func IsValidControllerRecordType(s int) bool {
	return int(Interview) <= s && s <= int(Other)
}

var (
	ErrControllerRecordNotFound = errors.New("controller record does not exist")
)

type ControllerRecordOperationInterface interface {
	NewControllerRecord(uid uint, operatorCid int, recordType ControllerRecordType, content string) (record *ControllerRecord)
	SaveControllerRecord(record *ControllerRecord) (err error)
	GetControllerRecords(uid uint, page, pageSize int) (records []*ControllerRecord, total int64, err error)
	GetControllerRecord(id uint, uid uint) (record *ControllerRecord, err error)
	DeleteControllerRecord(id uint) (err error)
}
