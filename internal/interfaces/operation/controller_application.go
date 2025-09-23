// Package operation
package operation

import (
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
	CreatedAt             time.Time `json:"-"`
	UpdatedAt             time.Time `json:"-"`
}

type ControllerApplicationStatus int

const (
	Submitted ControllerApplicationStatus = iota
	UnderProcessing
	Passed
	Rejected
)

func IsValidApplicationStatus(s int) bool {
	return int(Submitted) <= s && s <= int(Rejected)
}

type ControllerApplicationOperationInterface interface {
	NewApplication(userId uint, reason string, record string, guset bool, platform string, evidence string, status int) *ControllerApplication
	GetApplicationByUserId(userId uint) (*ControllerApplication, error)
	SaveApplication(application *ControllerApplication) error
	ConfirmApplicationUnderProcessing(application *ControllerApplication) error
	UpdateApplicationStatus(application *ControllerApplication, status ControllerApplicationStatus, message string) error
	CancelApplication(application *ControllerApplication) error
}
