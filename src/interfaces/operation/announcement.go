// Package operation
package operation

import (
	"errors"
	"time"
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

type UserAnnouncement struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Title     string    `gorm:"type:text;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Type      int       `gorm:"index;not null" json:"type"`
	Important bool      `gorm:"type:bool;default:false;not null" json:"important"`
	ForceShow bool      `gorm:"type:bool;default:false;not null" json:"force_show"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AnnouncementType int

const (
	Normal     AnnouncementType = iota // 普通公告
	Controller                         // 空管中心公告
	Technical                          // 技术组公告
)

func IsValidAnnouncementType(val int) bool {
	return int(Normal) <= val && val <= int(Technical)
}

var (
	ErrAnnouncementNotFound = errors.New("announcement not found")
)

type AnnouncementOperationInterface interface {
	NewAnnouncement(uid uint, content string, announcementType AnnouncementType, important bool, forceShow bool) *Announcement
	SaveAnnouncement(announcement *Announcement) error
	GetAnnouncementById(id uint) (announcement *Announcement, err error)
	GetAnnouncements(page, pageSize int) (announcements []*UserAnnouncement, total int64, err error)
	GetDetailAnnouncements(page, pageSize int) (announcements []*Announcement, total int64, err error)
	DeleteAnnouncement(announcement *Announcement) error
	UpdateAnnouncement(announcement *Announcement, updates map[string]interface{}) error
}
