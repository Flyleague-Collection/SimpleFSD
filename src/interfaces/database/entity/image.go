// Package entity
package entity

import (
	"time"

	"gorm.io/gorm"
)

type Image struct {
	ID        uint   `gorm:"primarykey"`
	UserId    uint   `gorm:"index;not null"`
	User      *User  `gorm:"foreignKey:UserId;references:ID;constraint:OnUpdate:cascade,OnDelete:cascade;"`
	HashCode  string `gorm:"size:128;uniqueIndex;not null"`
	FileName  string `gorm:"size:128;not null"`
	Url       string `gorm:"type:text;not null"`
	Size      int64  `gorm:"default:0;not null"`
	MimeType  string `gorm:"size:128;not null"`
	Comment   string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (image *Image) GetId() uint {
	return image.ID
}
