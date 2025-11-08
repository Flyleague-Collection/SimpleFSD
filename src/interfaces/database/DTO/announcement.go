// Package DTO
package DTO

import "time"

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
