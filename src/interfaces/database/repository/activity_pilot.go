// Package repository
package repository

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"github.com/half-nothing/simple-fsd/src/interfaces/enum"
)

type ActivityPilotStatus *enum.Enum[int]

var (
	ActivityPilotStatusSigned    ActivityPilotStatus = enum.New(0, "报名")
	ActivityPilotStatusClearance ActivityPilotStatus = enum.New(1, "放行")
	ActivityPilotStatusTakeoff   ActivityPilotStatus = enum.New(2, "起飞")
	ActivityPilotStatusLanding   ActivityPilotStatus = enum.New(3, "着陆")
)

var ActivityPilotManager = enum.NewManager(
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
