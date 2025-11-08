// Package repository
package repository

import (
	"errors"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

type ControllerApplicationStatus int

const (
	Submitted ControllerApplicationStatus = iota
	UnderProcessing
	Passed
	Rejected
)

var AllowedStatusMap = map[ControllerApplicationStatus][]ControllerApplicationStatus{
	Submitted:       {UnderProcessing, Passed, Rejected},
	UnderProcessing: {Passed, Rejected},
	Passed:          {},
	Rejected:        {},
}

func IsValidApplicationStatus(val int) bool {
	return int(Submitted) <= val && val <= int(Rejected)
}

var (
	ErrApplicationNotFound      = errors.New("application does not exist")
	ErrApplicationAlreadyExists = errors.New("application already exists")
)

type ControllerApplicationInterface interface {
	NewApplication(userId uint, reason string, record string, guset bool, platform string, evidence string) *entity.ControllerApplication
	GetApplicationByUserId(userId uint) (*entity.ControllerApplication, error)
	GetApplicationById(id uint) (application *entity.ControllerApplication, err error)
	GetApplications(page, pageSize int) ([]*entity.ControllerApplication, int64, error)
	SaveApplication(application *entity.ControllerApplication) error
	ConfirmApplicationUnderProcessing(application *entity.ControllerApplication) error
	UpdateApplicationStatus(application *entity.ControllerApplication, status ControllerApplicationStatus, message string) error
	CancelApplication(application *entity.ControllerApplication) error
}
