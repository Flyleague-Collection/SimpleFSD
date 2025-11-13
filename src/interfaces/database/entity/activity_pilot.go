// Package entity
package entity

import "time"

type ActivityPilot struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	ActivityId   uint      `gorm:"uniqueIndex:index_activity_pilot;not null" json:"activity_id"`
	UserId       uint      `gorm:"uniqueIndex:index_activity_pilot;not null" json:"uid"`
	User         *User     `gorm:"foreignKey:UserId;references:ID" json:"user"`
	Callsign     string    `gorm:"size:32;not null" json:"callsign"`
	AircraftType string    `gorm:"size:32;not null" json:"aircraft_type"`
	Status       int       `gorm:"default:0;not null" json:"status"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

func (entity *ActivityPilot) GetId() uint {
	return entity.ID
}
