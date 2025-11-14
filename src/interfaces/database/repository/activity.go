// Package repository
package repository

import (
	"errors"
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"github.com/half-nothing/simple-fsd/src/interfaces/enum"
)

var (
	ErrActivityNotFound = errors.New("activity not found")
	ErrActivityDeleted  = errors.New("activity has been deleted")
	ErrActivityEnded    = errors.New("activity has closed")

	ErrActivityPilotNotFound = errors.New("activity pilot not found")
	ErrPilotAlreadySigned    = errors.New("you have already signed up for the activity")
	ErrPilotUnsigned         = errors.New("you have not signed up for the activity yet")
	ErrCallsignAlreadyUsed   = errors.New("callsign already used")

	ErrFacilityNotFound    = errors.New("facility not found")
	ErrFacilityOtherSigned = errors.New("facility signed")
	ErrFacilityYouSigned   = errors.New("you have already signed up for the activity")
	ErrFacilityNotSigned   = errors.New("facility not signed")
	ErrFacilityNotYourSign = errors.New("you can not cancel other's facility sign")

	ErrActivityControllerNotFound = errors.New("activity controller not found")
	ErrRatingNotAllowed           = errors.New("rating not allowed")
	ErrControllerAlreadySign      = errors.New("you can not sign more than one facility")
)

type ActivityStatus *enum.Enum[int]

var (
	ActivityStatusRegistering ActivityStatus = enum.New(0, "报名中")
	ActivityStatusInTheEvent  ActivityStatus = enum.New(1, "活动中")
	ActivityStatusEnded       ActivityStatus = enum.New(2, "已结束")
)

var ActivityStatusManager = enum.NewManager(
	ActivityStatusRegistering,
	ActivityStatusInTheEvent,
	ActivityStatusEnded,
)

type ActivityType *enum.Enum[int]

var (
	ActivityTypeOneWay     ActivityType = enum.New(0, "单向单站")
	ActivityTypeBothWay    ActivityType = enum.New(1, "双向双站")
	ActivityTypeFIROpenDay ActivityType = enum.New(2, "空域开放日")
)

var ActivityTypeManager = enum.NewManager(
	ActivityTypeOneWay,
	ActivityTypeBothWay,
	ActivityTypeFIROpenDay,
)

// ActivityInterface 联飞活动操作接口定义
type ActivityInterface interface {
	Base[*entity.Activity]
	NewBuilder(user *entity.User, title string, image *entity.Image, activeTime time.Time, notams string) *ActivityBuilder
	NewOneWay(builder *ActivityBuilder, dep string, arr string, route string, distance int) *entity.Activity
	NewBothWay(builder *ActivityBuilder, dep string, arr string, route string, distance int, route2 string, distance2 int) *entity.Activity
	NewFIROpenDay(builder *ActivityBuilder, firs ...string) *entity.Activity
	GetNumber() (int64, error)
	GetBetween(startDay time.Time, endDay time.Time) ([]*entity.Activity, error)
	GetPage(pageNumber int, pageSize int) ([]*entity.Activity, int64, error)
	UpdateStatus(activityId uint, status ActivityStatus) error
	UpdateInfo(oldActivity *entity.Activity, newActivity *entity.Activity) error
	GetPilotRepository() ActivityPilotInterface
	GetControllerRepository() ActivityControllerInterface
	GetFacilityRepository() ActivityFacilityInterface
}
