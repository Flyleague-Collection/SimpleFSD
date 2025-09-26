// Package operation
package operation

import (
	"errors"
	"time"
)

type ControllerApplication struct {
	ID                    uint      `gorm:"primarykey" json:"id"`
	UserId                uint      `gorm:"index;not null" json:"user_id"`
	User                  *User     `gorm:"foreignKey:UserId;references:ID;constraint:OnUpdate:cascade,OnDelete:cascade;" json:"user"`
	WhyWantToBeController string    `gorm:"type:text;not null" json:"why_want_to_be_controller"`
	ControllerRecord      string    `gorm:"type:text;not null" json:"controller_record"`
	IsGuest               bool      `gorm:"not null" json:"is_guest"`
	Platform              string    `gorm:"not null" json:"platform"`
	Evidence              string    `gorm:"not null" json:"evidence"`
	Status                int       `gorm:"not null" json:"status"`
	Message               string    `gorm:"type:text;not null" json:"message"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"-"`
}

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

func IsValidApplicationStatus(s int) bool {
	return int(Submitted) <= s && s <= int(Rejected)
}

var (
	ErrApplicationNotFound      = errors.New("application does not exist")
	ErrApplicationAlreadyExists = errors.New("application already exists")
)

type ControllerApplicationOperationInterface interface {
	NewApplication(userId uint, reason string, record string, guset bool, platform string, evidence string) *ControllerApplication
	GetApplicationByUserId(userId uint) (*ControllerApplication, error)
	GetApplicationById(id uint) (application *ControllerApplication, err error)
	GetApplications(page, pageSize int) ([]*ControllerApplication, int64, error)
	SaveApplication(application *ControllerApplication) error
	ConfirmApplicationUnderProcessing(application *ControllerApplication) error
	UpdateApplicationStatus(application *ControllerApplication, status ControllerApplicationStatus, message string) error
	CancelApplication(application *ControllerApplication) error
}
