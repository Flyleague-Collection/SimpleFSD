// Package entity
package entity

import (
	"time"
)

type ControllerRecord struct {
	ID          uint      `gorm:"primarykey"`
	Type        int       `gorm:"default:0;not null"`
	UserId      uint      `gorm:"index:Uid;not null"`
	User        *User     `gorm:"foreignKey:UserId;references:ID;constraint:OnUpdate:cascade,OnDelete:cascade"`
	OperatorCid int       `gorm:"index:OperatorCid;not null"`
	Content     string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`
}

func (record *ControllerRecord) GetId() uint {
	return record.ID
}
