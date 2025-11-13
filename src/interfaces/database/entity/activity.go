// Package entity
package entity

import (
	"time"

	"gorm.io/gorm"
)

type Activity struct {
	ID               uint                  `gorm:"primarykey"`
	Type             int                   `gorm:"default:0;not null"`
	Publisher        int                   `gorm:"index;not null"`
	Title            string                `gorm:"size:128;not null"`
	ImageId          uint                  `gorm:"index;not null"`
	Image            *Image                `gorm:"foreignKey:ImageId;references:ID;constraint:OnUpdate:cascade,OnDelete:cascade;"`
	ActiveTime       time.Time             `gorm:"not null"`
	DepartureAirport string                `gorm:"size:64"`
	ArrivalAirport   string                `gorm:"size:64"`
	Route            string                `gorm:"type:text"`
	Distance         int                   `gorm:"default:0"`
	Route2           string                `gorm:"type:text"`
	Distance2        int                   `gorm:"default:0"`
	OpenFirs         string                `gorm:"size:128"`
	Status           int                   `gorm:"default:0;not null"`
	NOTAMS           string                `gorm:"type:text;not null"`
	Facilities       []*ActivityFacility   `gorm:"foreignKey:ActivityId;references:ID"`
	Controllers      []*ActivityController `gorm:"foreignKey:ActivityId;references:ID"`
	Pilots           []*ActivityPilot      `gorm:"foreignKey:ActivityId;references:ID"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}

func (activity *Activity) Equal(other *Activity) bool {
	if activity == nil || other == nil {
		return false
	}

	return activity.ID == other.ID &&
		activity.Type == other.Type &&
		activity.Publisher == other.Publisher &&
		activity.Title == other.Title &&
		activity.ImageId == other.ImageId &&
		activity.ActiveTime == other.ActiveTime &&
		activity.DepartureAirport == other.DepartureAirport &&
		activity.ArrivalAirport == other.ArrivalAirport &&
		activity.Route == other.Route &&
		activity.Distance == other.Distance &&
		activity.Route2 == other.Route2 &&
		activity.Distance2 == other.Distance2 &&
		activity.OpenFirs == other.OpenFirs &&
		activity.Status == other.Status &&
		activity.NOTAMS == other.NOTAMS
}

func (activity *Activity) Diff(other *Activity) map[string]interface{} {
	if activity == nil || other == nil {
		return nil
	}

	result := make(map[string]interface{})
	if other.Type >= 0 && activity.Type != other.Type {
		result["type"] = other.Type
	}
	if other.Publisher > 0 && activity.Publisher != other.Publisher {
		result["publisher"] = other.Publisher
	}
	if other.Title != "" && activity.Title != other.Title {
		result["title"] = other.Title
	}
	if other.ImageId > 0 && activity.ImageId != other.ImageId {
		result["image_id"] = other.ImageId
	}
	if other.ActiveTime.IsZero() && other.ActiveTime != other.ActiveTime {
		result["active_time"] = other.ActiveTime
	}
	if other.DepartureAirport != "" && activity.DepartureAirport != other.DepartureAirport {
		result["departure_airport"] = other.DepartureAirport
	}
	if other.ArrivalAirport != "" && activity.ArrivalAirport != other.ArrivalAirport {
		result["arrival_airport"] = other.ArrivalAirport
	}
	if other.Route != "" && activity.Route != other.Route {
		result["route"] = other.Route
	}
	if other.Distance >= 0 && activity.Distance != other.Distance {
		result["distance"] = other.Distance
	}
	if other.Route2 != "" && activity.Route2 != other.Route2 {
		result["route2"] = other.Route2
	}
	if other.Distance2 >= 0 && activity.Distance2 != other.Distance2 {
		result["distance2"] = other.Distance2
	}
	if other.OpenFirs != "" && activity.OpenFirs != other.OpenFirs {
		result["open_firs"] = other.OpenFirs
	}
	if other.Status >= 0 && activity.Status != other.Status {
		result["status"] = other.Status
	}
	if activity.NOTAMS != other.NOTAMS {
		result["NOTAMS"] = other.NOTAMS
	}
	return result
}

func (activity *Activity) GetId() uint {
	return activity.ID
}
