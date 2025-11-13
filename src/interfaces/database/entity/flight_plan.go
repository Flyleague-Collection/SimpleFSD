// Package entity
package entity

import (
	"time"
)

type FlightPlan struct {
	ID               uint      `gorm:"primarykey" json:"id"`
	Cid              int       `gorm:"index;not null" json:"cid"`
	User             *User     `gorm:"foreignKey:Cid;references:Cid;constraint:OnUpdate:cascade,OnDelete:cascade;" json:"user"`
	Callsign         string    `gorm:"size:16;uniqueIndex;not null" json:"callsign"`
	FlightType       string    `gorm:"size:4;not null" json:"flight_rules"`
	AircraftType     string    `gorm:"size:128;not null" json:"aircraft"`
	Tas              int       `gorm:"not null" json:"cruise_tas"`
	DepartureAirport string    `gorm:"size:4;not null" json:"departure"`
	DepartureTime    int       `gorm:"not null" json:"departure_time"`
	AtcDepartureTime int       `gorm:"not null" json:"-"`
	CruiseAltitude   string    `gorm:"size:8;not null" json:"altitude"`
	ArrivalAirport   string    `gorm:"size:4;not null" json:"arrival"`
	RouteTimeHour    string    `gorm:"size:2;not null" json:"route_time_hour"`
	RouteTimeMinute  string    `gorm:"size:2;not null" json:"route_time_minute"`
	FuelTimeHour     string    `gorm:"size:2;not null" json:"fuel_time_hour"`
	FuelTimeMinute   string    `gorm:"size:2;not null" json:"fuel_time_minute"`
	AlternateAirport string    `gorm:"size:4;not null" json:"alternate"`
	Remarks          string    `gorm:"type:text;not null" json:"remarks"`
	Route            string    `gorm:"type:text;not null" json:"route"`
	Locked           bool      `gorm:"default:0;not null" json:"locked"`
	FromWeb          bool      `gorm:"default:0;not null" json:"from_web"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
}
