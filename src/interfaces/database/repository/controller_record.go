// Package repository
package repository

import (
	"errors"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

type ControllerRecordType int

const (
	Interview    ControllerRecordType = iota // 面试
	Simulator                                // 模拟机
	RatingChange                             // 权限变动
	Training                                 // 训练内容
	UnderMonitor                             // UM权限变动
	Solo                                     // Solo权限变动
	Guest                                    // 客座权限变动
	Application                              // 管制员申请
	Other                                    // 其他未定义内容
)

func IsValidControllerRecordType(val int) bool {
	return int(Interview) <= val && val <= int(Other)
}

var (
	ErrControllerRecordNotFound = errors.New("controller record does not exist")
)

type ControllerRecordInterface interface {
	NewControllerRecord(uid uint, operatorCid int, recordType ControllerRecordType, content string) (record *entity.ControllerRecord)
	SaveControllerRecord(record *entity.ControllerRecord) (err error)
	GetControllerRecords(uid uint, page, pageSize int) (records []*entity.ControllerRecord, total int64, err error)
	GetControllerRecord(id uint, uid uint) (record *entity.ControllerRecord, err error)
	DeleteControllerRecord(id uint) (err error)
}
