// Package operation
package operation

import (
	"errors"
	"time"
)

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

var (
	ErrActivityNotFound      = errors.New("activity not found")
	ErrFacilityNotFound      = errors.New("facility not found")
	ErrRatingNotAllowed      = errors.New("rating not allowed")
	ErrFacilitySigned        = errors.New("facility signed")
	ErrFacilityNotSigned     = errors.New("facility not signed")
	ErrFacilityNotYourSign   = errors.New("you can not cancel other's facility sign")
	ErrFacilityAlreadyExists = errors.New("you can not sign more than one facility")
	ErrActivityAlreadySigned = errors.New("you have already signed up for the activity")
	ErrCallsignAlreadyUsed   = errors.New("callsign already used")
	ErrActivityUnsigned      = errors.New("you have not signed up for the activity yet")
)

// ActivityOperationInterface 联飞活动操作接口定义
type ActivityOperationInterface interface {
	// NewActivity 创建新活动
	NewActivity(user *User, title string, imageUrl string, activeTime time.Time, dep string, arr string, route string, distance int, notams string) (activity *Activity)
	// NewActivityFacility 创建新活动管制席位
	NewActivityFacility(activity *Activity, rating int, callsign string, frequency float64) (activityFacility *ActivityFacility)
	// NewActivityAtc 创建新参加活动的管制员
	NewActivityAtc(facility *ActivityFacility, user *User) (activityAtc *ActivityATC)
	// NewActivityPilot 创建新参加活动的飞行员
	NewActivityPilot(activity *Activity, cid int, callsign string, aircraftType string) (activityPilot *ActivityPilot)
	// GetActivities 获取指定日期内的所有活动, 当err为nil时返回值activities有效
	GetActivities(startDay, endDay time.Time) (activities []*Activity, err error)
	// GetActivityById 通过活动Id获取活动详细内容(包括外键关联的内容), 当err为nil时返回值activity有效
	GetActivityById(id uint) (activity *Activity, err error)
	// GetOnlyActivityById 通过活动Id获取活动详细内容(不包括外键关联的内容), 当err为nil时返回值activity有效
	GetOnlyActivityById(id uint) (activity *Activity, err error)
	// SaveActivity 保存活动到数据库, 当err为nil时保存成功
	SaveActivity(activity *Activity) (err error)
	// DeleteActivity 删除活动, 当err为nil时删除成功
	DeleteActivity(activity *Activity) (err error)
	// SetActivityStatus 设置活动状态, 当err为nil时设置成功
	SetActivityStatus(activity *Activity, status ActivityStatus) (err error)
	// SetActivityPilotStatus 设置参与活动的飞行员的状态, 当err为nil时设置成功
	SetActivityPilotStatus(activityPilot *ActivityPilot, status ActivityPilotStatus) (err error)
	// GetActivityPilotById 获取参与活动的指定机组, 当err为nil时返回值pilot有效
	GetActivityPilotById(activityId uint, cid int) (pilot *ActivityPilot, err error)
	// GetFacilityById 获取指定活动的指定席位, 当err为nil时返回值facility有效
	GetFacilityById(facilityId uint) (facility *ActivityFacility, err error)
	// SignFacilityController 设置报名席位的用户, 当err为nil时保存成功
	SignFacilityController(facility *ActivityFacility, user *User) (err error)
	// UnsignFacilityController 取消报名席位的用户, 当err为nil时取消成功
	UnsignFacilityController(facility *ActivityFacility, cid int) (err error)
	// SignActivityPilot 飞行员报名, 当err为nil时保存成功
	SignActivityPilot(activity *Activity, cid int, callsign string, aircraftType string) (err error)
	// UnsignActivityPilot 飞行员取消报名, 当err为nil时取消成功
	UnsignActivityPilot(activity *Activity, cid int) (err error)
	// UpdateActivityInfo 更新活动信息, 当err为nil时更新成功
	UpdateActivityInfo(activity *Activity, updateInfo map[string]interface{}) (err error)
}
