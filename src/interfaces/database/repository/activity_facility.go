// Package repository
package repository

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

type ActivityFacilityInterface interface {
	Base[*entity.ActivityFacility]
	New(activity *entity.Activity, minRating int, callsign string, frequency float64, tier2Tower bool) *entity.ActivityFacility
}
