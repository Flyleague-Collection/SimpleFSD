// Package repository
package repository

import (
	"errors"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"github.com/half-nothing/simple-fsd/src/interfaces/enum"
)

type ControllerRecordType *enum.Enum[int]

var (
	ControllerRecordInterview    ControllerRecordType = enum.New(0, "面试")
	ControllerRecordSimulator    ControllerRecordType = enum.New(1, "模拟机")
	ControllerRecordRatingChange ControllerRecordType = enum.New(2, "权限变动")
	ControllerRecordTraining     ControllerRecordType = enum.New(3, "训练内容")
	ControllerRecordUnderMonitor ControllerRecordType = enum.New(4, "UM权限变动")
	ControllerRecordSolo         ControllerRecordType = enum.New(5, "Solo权限变动")
	ControllerRecordGuest        ControllerRecordType = enum.New(6, "客座权限变动")
	ControllerRecordApplication  ControllerRecordType = enum.New(7, "管制员申请")
	ControllerRecordOther        ControllerRecordType = enum.New(8, "其他未定义内容")
)

var ControllerRecordManager = enum.NewManager(
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
