// Package entity
package entity

import (
	"time"
)

type User struct {
	ID             uint      `gorm:"primarykey"`
	Username       string    `gorm:"size:64;uniqueIndex;not null"`
	Email          string    `gorm:"size:128;uniqueIndex;not null"`
	Cid            int       `gorm:"uniqueIndex;not null"`
	Password       string    `gorm:"size:128;not null"`
	AvatarUrl      string    `gorm:"size:128;not null;default:''"`
	ImageId        uint      `gorm:"default:0;not null"`
	Image          *Image    `gorm:"foreignKey:AvatarId;references:ID;constraint:OnUpdate:cascade,OnDelete:cascade;"`
	QQ             int       `gorm:"default:0;not null"`
	Rating         int       `gorm:"default:0;not null"`
	Guest          bool      `gorm:"default:false;not null"`
	UnderMonitor   bool      `gorm:"default:false;not null"`
	UnderSolo      bool      `gorm:"default:false;not null"`
	Tier2          bool      `gorm:"default:false;not null"`
	SoloUntil      time.Time `gorm:"default:null"`
	Permission     uint64    `gorm:"default:0;not null"`
	TotalPilotTime int       `gorm:"default:0;not null"`
	TotalAtcTime   int       `gorm:"default:0;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (user *User) GetId() uint {
	return user.ID
}
