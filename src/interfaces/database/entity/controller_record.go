// Package entity
package entity

import (
	"time"
)

type ControllerRecord struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Type        int       `gorm:"not null" json:"type"`
	UserId      uint      `gorm:"index:Uid;not null" json:"uid"`
	OperatorCid int       `gorm:"index:OperatorCid;not null" json:"operator_cid"`
	Content     string    `gorm:"not null" json:"content"`
	CreatedAt   time.Time `gorm:"not null" json:"time"`
}
