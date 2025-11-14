// Package repository
package repository

import (
	"errors"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"github.com/half-nothing/simple-fsd/src/interfaces/enum"
)

type AnnouncementType *enum.Enum[int]

var (
	AnnouncementTypeNormal     AnnouncementType = enum.New(0, "普通公告")
	AnnouncementTypeController AnnouncementType = enum.New(1, "空管中心公告")
	AnnouncementTypeTechnical  AnnouncementType = enum.New(2, "技术组公告")
)

var AnnouncementTypeManager = enum.NewManager(
	AnnouncementTypeNormal,
	AnnouncementTypeController,
	AnnouncementTypeTechnical,
)

var (
	ErrAnnouncementNotFound = errors.New("announcement not found")
)

type AnnouncementInterface interface {
	Base[*entity.Announcement]
	New(user *entity.User, content string, announcementType AnnouncementType, important bool, forceShow bool) *entity.Announcement
	GetPage(pageNumber int, pageSize int) ([]*entity.Announcement, int64, error)
}
