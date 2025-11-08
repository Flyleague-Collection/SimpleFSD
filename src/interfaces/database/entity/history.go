// Package entity
package entity

import "time"

type History struct {
	ID         uint      `gorm:"primarykey" json:"-"`
	Cid        int       `gorm:"index;not null" json:"-"`
	Callsign   string    `gorm:"size:16;index;not null" json:"callsign"`
	StartTime  time.Time `gorm:"not null" json:"start_time"`
	EndTime    time.Time `gorm:"not null" json:"end_time"`
	OnlineTime int       `gorm:"default:0;not null" json:"online_time"`
	IsAtc      bool      `gorm:"default:0;not null" json:"-"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}
