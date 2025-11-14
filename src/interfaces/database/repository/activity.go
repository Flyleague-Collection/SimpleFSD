// Package repository
package repository

import (
	"errors"
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
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

type ActivityStatus *Enum[int]

var (
	ActivityStatusRegistering ActivityStatus = NewEnum(0, "报名中")
	ActivityStatusInTheEvent  ActivityStatus = NewEnum(1, "活动中")
	ActivityStatusEnded       ActivityStatus = NewEnum(2, "已结束")
)

var activityStatuses = []ActivityStatus{
	ActivityStatusRegistering,
	ActivityStatusInTheEvent,
	ActivityStatusEnded,
}

func IsValidActivityStatus(index int) bool {
	return 0 <= index && index < len(activityStatuses)
}

func GetActivityStatus(index int) ActivityStatus {
	if !IsValidActivityStatus(index) {
		return nil
	}
	return activityStatuses[index]
}

type ActivityType *Enum[int]

var (
	ActivityTypeOneWay     ActivityType = NewEnum(0, "单向单站")
	ActivityTypeBothWay    ActivityType = NewEnum(1, "双向双站")
	ActivityTypeFIROpenDay ActivityType = NewEnum(2, "空域开放日")
)

var activityTypes = []ActivityType{
	ActivityTypeOneWay,
	ActivityTypeBothWay,
	ActivityTypeFIROpenDay,
}

func IsValidActivityType(index int) bool {
	return 0 <= index && index < len(activityTypes)
}

func GetActivityType(index int) ActivityType {
	if !IsValidActivityType(index) {
		return nil
	}
	return activityTypes[index]
}

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
