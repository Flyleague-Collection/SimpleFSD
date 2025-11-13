// Package repository
package repository

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

type ActivityControllerInterface interface {
	Base[*entity.ActivityController]
	New(activityId uint, facilityId uint, userId uint) *entity.ActivityController
	GetByActivityIdAndFacilityIdAndUserId(activityId uint, facilityId uint, userId uint) (*entity.ActivityController, error)
	GetByActivityIdAndUserId(activityId uint, userId uint) (*entity.ActivityController, error)
	JoinActivity(activity *entity.Activity, activityFacility *entity.ActivityFacility, user *entity.User) error
	LeaveActivity(activity *entity.Activity, activityFacility *entity.ActivityFacility, user *entity.User) error
}
