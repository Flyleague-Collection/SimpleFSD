// Package entity
package entity

import (
	"time"
)

type ControllerApplication struct {
	ID                    uint   `gorm:"primarykey"`
	UserId                uint   `gorm:"index;not null"`
	User                  *User  `gorm:"foreignKey:UserId;references:ID;constraint:OnUpdate:cascade,OnDelete:cascade"`
	WhyWantToBeController string `gorm:"type:text;not null"`
	ControllerRecord      string `gorm:"type:text;not null"`
	IsGuest               bool   `gorm:"default:false;not null"`
	Platform              string `gorm:"not null"`
	Evidence              string `gorm:"not null"`
	Status                int    `gorm:"default:0;not null"`
	Message               string `gorm:"type:text;not null"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func (application *ControllerApplication) GetId() uint {
	return application.ID
}
