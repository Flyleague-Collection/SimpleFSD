// Package entity
package entity

import "time"

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
