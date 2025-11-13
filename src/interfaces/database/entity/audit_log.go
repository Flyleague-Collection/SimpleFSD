// Package entity
package entity

import "time"

type AuditLog struct {
	ID            uint          `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time     `gorm:"not null" json:"time"`
	EventType     string        `gorm:"index:eventType;not null" json:"event_type"`
	Subject       int           `gorm:"index:Subject;not null" json:"subject"`
	User          *User         `gorm:"foreignKey:Subject;references:Cid;constraint:OnUpdate:cascade,OnDelete:cascade;" json:"user"`
	Object        string        `gorm:"index:Object;not null" json:"object"`
	Ip            string        `gorm:"not null" json:"ip"`
	UserAgent     string        `gorm:"not null" json:"user_agent"`
	ChangeDetails *ChangeDetail `gorm:"type:text;serializer:json" json:"change_details"`
}

type ChangeDetail struct {
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}
