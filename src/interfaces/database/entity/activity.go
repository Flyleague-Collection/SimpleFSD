// Package entity
package entity

import (
	"time"

	"gorm.io/gorm"
)

type Activity struct {
	ID               uint                `gorm:"primarykey" json:"id"`
	Publisher        int                 `gorm:"index;not null" json:"publisher"`
	Title            string              `gorm:"size:128;not null" json:"title"`
	ImageUrl         string              `gorm:"size:128;not null" json:"image_url"`
	ActiveTime       time.Time           `gorm:"not null" json:"active_time"`
	DepartureAirport string              `gorm:"size:64;not null" json:"departure_airport"`
	ArrivalAirport   string              `gorm:"size:64;not null" json:"arrival_airport"`
	Route            string              `gorm:"type:text;not null" json:"route"`
	Distance         int                 `gorm:"default:0;not null" json:"distance"`
	Status           int                 `gorm:"default:0;not null" json:"status"`
	NOTAMS           string              `gorm:"type:text;not null" json:"NOTAMS"`
	Facilities       []*ActivityFacility `gorm:"foreignKey:ActivityId;references:ID" json:"facilities"`
	Controllers      []*ActivityATC      `gorm:"foreignKey:ActivityId;references:ID" json:"controllers"`
	Pilots           []*ActivityPilot    `gorm:"foreignKey:ActivityId;references:ID" json:"pilots"`
	CreatedAt        time.Time           `json:"-"`
	UpdatedAt        time.Time           `json:"-"`
	DeletedAt        gorm.DeletedAt      `json:"-"`
}

func (facility *Activity) Equal(other *Activity) bool {
	return facility.ID == other.ID && facility.Publisher == other.Publisher && facility.Title == other.Title &&
		facility.ImageUrl == other.ImageUrl && facility.ActiveTime == other.ActiveTime &&
		facility.DepartureAirport == other.DepartureAirport && facility.ArrivalAirport == other.ArrivalAirport &&
		facility.Route == other.Route && facility.Distance == other.Distance && facility.Status == other.Status &&
		facility.NOTAMS == other.NOTAMS
}

func (facility *Activity) Diff(other *Activity) map[string]interface{} {
	result := make(map[string]interface{})
	if facility.Publisher != 0 && facility.Publisher != other.Publisher {
		other.Publisher = facility.Publisher
		result["publisher"] = facility.Publisher
	}
	if facility.Title != "" && facility.Title != other.Title {
		other.Title = facility.Title
		result["title"] = facility.Title
	}
	if facility.ImageUrl != "" && facility.ImageUrl != other.ImageUrl {
		other.ImageUrl = facility.ImageUrl
		result["image_url"] = facility.ImageUrl
	}
	if facility.ActiveTime != other.ActiveTime {
		other.ActiveTime = facility.ActiveTime
		result["active_time"] = facility.ActiveTime
	}
	if facility.DepartureAirport != "" && facility.DepartureAirport != other.DepartureAirport {
		other.DepartureAirport = facility.DepartureAirport
		result["departure_airport"] = facility.DepartureAirport
	}
	if facility.ArrivalAirport != "" && facility.ArrivalAirport != other.ArrivalAirport {
		other.ArrivalAirport = facility.ArrivalAirport
		result["arrival_airport"] = facility.ArrivalAirport
	}
	if facility.Route != "" && facility.Route != other.Route {
		other.Route = facility.Route
		result["route"] = facility.Route
	}
	if facility.Distance != 0 && facility.Distance != other.Distance {
		other.Distance = facility.Distance
		result["distance"] = facility.Distance
	}
	if facility.Status != other.Status {
		other.Status = facility.Status
		result["status"] = facility.Status
	}
	if facility.NOTAMS != other.NOTAMS {
		other.NOTAMS = facility.NOTAMS
		result["NOTAMS"] = facility.NOTAMS
	}
	return result
}
