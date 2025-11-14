// Package repository
package repository

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

type ActivityPilotStatus *Enum[int]

var (
	ActivityPilotStatusSigned    ActivityPilotStatus = NewEnum(0, "报名")
	ActivityPilotStatusClearance ActivityPilotStatus = NewEnum(1, "放行")
	ActivityPilotStatusTakeoff   ActivityPilotStatus = NewEnum(2, "起飞")
	ActivityPilotStatusLanding   ActivityPilotStatus = NewEnum(3, "着陆")
)

var ActivityPilotManager = NewEnumManager(
	ActivityPilotStatusSigned,
	ActivityPilotStatusClearance,
	ActivityPilotStatusTakeoff,
	ActivityPilotStatusLanding,
)

type ActivityPilotInterface interface {
	Base[*entity.ActivityPilot]
	New(activityId uint, userId uint, callsign string, aircraftType string) *entity.ActivityPilot
	GetByActivityIdAndUserId(activityId uint, userId uint) (*entity.ActivityPilot, error)
	GetByActivityIdAndCallsign(activityId uint, callsign string) (*entity.ActivityPilot, error)
	VerifyUserIdAndCallsign(activityId uint, userId uint, callsign string) (*entity.ActivityPilot, error)
	UpdateStatus(activityPilot *entity.ActivityPilot, status ActivityPilotStatus) error
	JoinActivity(activity *entity.Activity, user *entity.User, callsign string, aircraftType string) error
	LeaveActivity(activity *entity.Activity, user *entity.User) error
}
