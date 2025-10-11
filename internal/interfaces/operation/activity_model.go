// Package operation
package operation

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

type ActivityFacility struct {
	ID         uint         `gorm:"primarykey" json:"id"`
	ActivityId uint         `gorm:"index;not null" json:"activity_id"`
	Tier2Tower bool         `gorm:"default:false;not null" json:"tier2_tower"`
	MinRating  int          `gorm:"default:2;not null" json:"min_rating"`
	Callsign   string       `gorm:"size:16;not null" json:"callsign"`
	Frequency  string       `gorm:"size:16;not null" json:"frequency"`
	Controller *ActivityATC `gorm:"foreignKey:FacilityId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"controller"`
	CreatedAt  time.Time    `json:"-"`
	UpdatedAt  time.Time    `json:"-"`
}

func (facility *ActivityFacility) Equal(other *ActivityFacility) bool {
	return facility.ID == other.ID && facility.ActivityId == other.ActivityId && facility.MinRating == other.MinRating &&
		facility.Callsign == other.Callsign && facility.Frequency == other.Frequency && facility.Tier2Tower == other.Tier2Tower
}

func (facility *ActivityFacility) Diff(other *ActivityFacility) map[string]interface{} {
	result := make(map[string]interface{})
	if facility.MinRating != other.MinRating {
		result["min_rating"] = facility.MinRating
	}
	if facility.Tier2Tower != other.Tier2Tower {
		result["tier2_tower"] = facility.Tier2Tower
	}
	if facility.Callsign != other.Callsign {
		result["callsign"] = facility.Callsign
	}
	if facility.Frequency != other.Frequency {
		result["frequency"] = facility.Frequency
	}
	return result
}

type ActivityATC struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	ActivityId uint      `gorm:"uniqueIndex:index_activity_controller;not null" json:"activity_id"`
	FacilityId uint      `gorm:"uniqueIndex:index_activity_controller;not null" json:"facility_id"`
	UserId     uint      `gorm:"not null" json:"uid"`
	User       *User     `gorm:"foreignKey:UserId;references:ID" json:"user"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

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
type ActivityStatus int

const (
	Open     ActivityStatus = iota // 报名中
	InActive                       // 活动中
	Closed                         // 已结束
)

type ActivityPilotStatus int

const (
	Signed    ActivityPilotStatus = iota // 已报名
	Clearance                            // 已放行
	Takeoff                              // 已起飞
	Landing                              // 已落地
)
