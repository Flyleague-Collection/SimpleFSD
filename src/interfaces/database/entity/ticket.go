// Package entity
package entity

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserId    uint           `gorm:"index:userId;not null" json:"creator"`
	User      *User          `gorm:"foreignKey:UserId;references:ID;constraint:OnUpdate:cascade,OnDelete:cascade;" json:"user"`
	Type      int            `gorm:"not null" json:"type"`
	Title     string         `gorm:"not null" json:"title"`
	Content   string         `gorm:"not null" json:"content"`
	Reply     string         `gorm:"not null" json:"reply"`
	Closer    int            `gorm:"index:closer;not null" json:"closer"`
	CreatedAt time.Time      `json:"open_at"`
	UpdatedAt time.Time      `json:"close_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
