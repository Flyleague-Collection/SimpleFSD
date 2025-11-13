// Package entity
package entity

import (
	"time"
)

type User struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	Username       string    `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Email          string    `gorm:"size:128;uniqueIndex;not null" json:"email"`
	Cid            int       `gorm:"uniqueIndex;not null" json:"cid"`
	Password       string    `gorm:"size:128;not null" json:"-"`
	AvatarUrl      string    `gorm:"size:128;not null;default:''" json:"avatar_url"`
	ImageId        uint      `gorm:"default:0;not null" json:"avatar_id"`
	Image          *Image    `gorm:"foreignKey:AvatarId;references:ID;constraint:OnUpdate:cascade,OnDelete:cascade;" json:"avatar"`
	QQ             int       `gorm:"default:0" json:"qq"`
	Rating         int       `gorm:"default:0" json:"rating"`
	Guest          bool      `gorm:"default:false" json:"guest"`
	UnderMonitor   bool      `gorm:"default:false;not null" json:"under_monitor"`
	UnderSolo      bool      `gorm:"default:false;not null" json:"under_solo"`
	Tier2          bool      `gorm:"default:false;not null" json:"tier2"`
	SoloUntil      time.Time `gorm:"default:null" json:"solo_until"`
	Permission     uint64    `gorm:"default:0" json:"permission"`
	TotalPilotTime int       `gorm:"default:0" json:"total_pilot_time"`
	TotalAtcTime   int       `gorm:"default:0" json:"total_atc_time"`
	CreatedAt      time.Time `json:"register_time"`
	UpdatedAt      time.Time `json:"-"`
}
