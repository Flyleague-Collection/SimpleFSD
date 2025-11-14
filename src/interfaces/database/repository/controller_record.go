// Package repository
package repository

import (
	"errors"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

type ControllerRecordType *Enum[int]

var (
	ControllerRecordInterview    ControllerRecordType = NewEnum(0, "面试")
	ControllerRecordSimulator    ControllerRecordType = NewEnum(1, "模拟机")
	ControllerRecordRatingChange ControllerRecordType = NewEnum(2, "权限变动")
	ControllerRecordTraining     ControllerRecordType = NewEnum(3, "训练内容")
	ControllerRecordUnderMonitor ControllerRecordType = NewEnum(4, "UM权限变动")
	ControllerRecordSolo         ControllerRecordType = NewEnum(5, "Solo权限变动")
	ControllerRecordGuest        ControllerRecordType = NewEnum(6, "客座权限变动")
	ControllerRecordApplication  ControllerRecordType = NewEnum(7, "管制员申请")
	ControllerRecordOther        ControllerRecordType = NewEnum(8, "其他未定义内容")
)

var ControllerRecordManager = NewEnumManager(
	ControllerRecordInterview,
	ControllerRecordSimulator,
	ControllerRecordRatingChange,
	ControllerRecordTraining,
	ControllerRecordUnderMonitor,
	ControllerRecordSolo,
	ControllerRecordGuest,
	ControllerRecordApplication,
	ControllerRecordOther,
)

var (
	ErrControllerRecordNotFound = errors.New("controller record does not exist")
)

type ControllerRecordInterface interface {
	Base[*entity.ControllerRecord]
	New(uid uint, operatorCid int, recordType ControllerRecordType, content string) *entity.ControllerRecord
	GetPage(uid uint, pageNumber, pageSize int) ([]*entity.ControllerRecord, int64, error)
	GetByIdAndUserId(id uint, uid uint) (*entity.ControllerRecord, error)
}
