// Package entity
package entity

import "time"

type ActivityController struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	ActivityId uint      `gorm:"uniqueIndex:index_activity_controller;not null" json:"activity_id"`
	FacilityId uint      `gorm:"uniqueIndex:index_activity_controller;not null" json:"facility_id"`
	UserId     uint      `gorm:"not null" json:"uid"`
	User       *User     `gorm:"foreignKey:UserId;references:ID" json:"user"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

func (entity *ActivityController) GetId() uint {
	return entity.ID
}
