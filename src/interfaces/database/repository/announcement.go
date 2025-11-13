// Package repository
package repository

import (
	"errors"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

type AnnouncementType *Enum

var (
	AnnouncementTypeNormal     AnnouncementType = NewEnum(0, "普通公告")
	AnnouncementTypeController AnnouncementType = NewEnum(1, "空管中心公告")
	AnnouncementTypeTechnical  AnnouncementType = NewEnum(2, "技术组公告")
)

var announcementTypes = []AnnouncementType{
	AnnouncementTypeNormal,
	AnnouncementTypeController,
	AnnouncementTypeTechnical,
}

func IsValidAnnouncementType(val int) bool {
	return 0 <= val && val < len(announcementTypes)
}

func GetAnnouncementType(val int) AnnouncementType {
	if !IsValidAnnouncementType(val) {
		return nil
	}
	return announcementTypes[val]
}

var (
	ErrAnnouncementNotFound = errors.New("announcement not found")
)

type AnnouncementInterface interface {
	Base[*entity.Announcement]
	New(user *entity.User, content string, announcementType AnnouncementType, important bool, forceShow bool) *entity.Announcement
	GetPage(pageNumber int, pageSize int) ([]*entity.Announcement, int64, error)
}
