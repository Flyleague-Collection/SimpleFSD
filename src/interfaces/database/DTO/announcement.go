// Package DTO
package DTO

import (
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

type UserAnnouncement struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Type      int       `json:"type"`
	Important bool      `json:"important"`
	ForceShow bool      `json:"force_show"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (announcement *UserAnnouncement) FromAnnouncementEntity(
	dbo *entity.Announcement,
) *UserAnnouncement {
	return &UserAnnouncement{
		ID:        dbo.ID,
		Title:     dbo.Title,
		Content:   dbo.Content,
		Type:      dbo.Type,
		Important: dbo.Important,
		ForceShow: dbo.ForceShow,
		CreatedAt: dbo.CreatedAt,
		UpdatedAt: dbo.UpdatedAt,
	}
}

type AdminAnnouncement struct {
	ID        uint         `json:"id"`
	Publisher *entity.User `json:"publisher"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Type      int          `json:"type"`
	Important bool         `json:"important"`
	ForceShow bool         `json:"force_show"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (announcement *AdminAnnouncement) FromAnnouncementEntity(
	dbo *entity.Announcement,
) *AdminAnnouncement {
	return &AdminAnnouncement{
		ID:        dbo.ID,
		Publisher: dbo.User,
		Title:     dbo.Title,
		Content:   dbo.Content,
		Type:      dbo.Type,
		Important: dbo.Important,
		ForceShow: dbo.ForceShow,
		CreatedAt: dbo.CreatedAt,
	}
}
