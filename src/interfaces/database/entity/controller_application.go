// Package entity
package entity

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
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"-"`
}
