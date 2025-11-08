// Package entity
package entity

import (
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/DTO"
)

type Announcement struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	PublisherId uint      `gorm:"index;not null" json:"publisher_id"`
	User        *User     `gorm:"foreignKey:PublisherId;references:ID;constraint:OnUpdate:cascade,OnDelete:cascade;" json:"user"`
	Title       string    `gorm:"type:text;not null" json:"title"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	Type        int       `gorm:"index;not null" json:"type"`
	Important   bool      `gorm:"type:bool;default:false;not null" json:"important"`
	ForceShow   bool      `gorm:"type:bool;default:false;not null" json:"force_show"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (announcement *Announcement) ToUserAnnouncementDTO() *DTO.UserAnnouncement {
	return &DTO.UserAnnouncement{
		ID:        announcement.ID,
		Title:     announcement.Title,
		Content:   announcement.Content,
		Type:      announcement.Type,
		Important: announcement.Important,
		ForceShow: announcement.ForceShow,
		CreatedAt: announcement.CreatedAt,
		UpdatedAt: announcement.UpdatedAt,
	}
}
