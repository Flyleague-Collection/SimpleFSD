// Package repository
package repository

import (
	"errors"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/DTO"
	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

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

type AnnouncementInterface interface {
	NewAnnouncement(uid uint, content string, announcementType AnnouncementType, important bool, forceShow bool) *entity.Announcement
	SaveAnnouncement(announcement *entity.Announcement) error
	GetAnnouncementById(id uint) (announcement *entity.Announcement, err error)
	GetAnnouncements(page, pageSize int) (announcements []*DTO.UserAnnouncement, total int64, err error)
	GetDetailAnnouncements(page, pageSize int) (announcements []*entity.Announcement, total int64, err error)
	DeleteAnnouncement(announcement *entity.Announcement) error
	UpdateAnnouncement(announcement *entity.Announcement, updates map[string]interface{}) error
}
