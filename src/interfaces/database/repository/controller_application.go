// Package repository
package repository

import (
	"errors"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

type ApplicationStatus *Enum[int]

var (
	ApplicationStatusSubmitted       ApplicationStatus = NewEnum(0, "已提交")
	ApplicationStatusUnderProcessing ApplicationStatus = NewEnum(1, "处理中")
	ApplicationStatusPassed          ApplicationStatus = NewEnum(2, "已通过")
	ApplicationStatusRejected        ApplicationStatus = NewEnum(3, "已拒绝")
)

var ApplicationStatusManager = NewEnumManager(
	ApplicationStatusSubmitted,
	ApplicationStatusUnderProcessing,
	ApplicationStatusPassed,
	ApplicationStatusRejected,
)

var ApplicationStatusTransformMap = map[ApplicationStatus][]ApplicationStatus{
	ApplicationStatusSubmitted:       {ApplicationStatusUnderProcessing, ApplicationStatusPassed, ApplicationStatusRejected},
	ApplicationStatusUnderProcessing: {ApplicationStatusPassed, ApplicationStatusRejected},
	ApplicationStatusPassed:          {},
	ApplicationStatusRejected:        {},
}

var (
	ErrApplicationNotFound      = errors.New("application does not exist")
	ErrApplicationAlreadyExists = errors.New("application already exists")
)

type ControllerApplicationInterface interface {
	Base[*entity.ControllerApplication]
	New(user *entity.User, reason string, record string, guset bool, platform string, evidence string) *entity.ControllerApplication
	GetByUserId(userId uint) (*entity.ControllerApplication, error)
	GetPage(pageNumber int, pageSize int) ([]*entity.ControllerApplication, int64, error)
	UpdateStatus(application *entity.ControllerApplication, status ApplicationStatus, message string) error
}
