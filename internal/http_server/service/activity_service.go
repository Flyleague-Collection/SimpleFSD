// Package service
package service

import (
	"encoding/json"
	"errors"
	c "github.com/half-nothing/simple-fsd/internal/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"strconv"
	"time"
)

type ActivityService struct {
	config            *c.HttpServerConfig
	userOperation     operation.UserOperationInterface
	activityOperation operation.ActivityOperationInterface
	storeService      StoreServiceInterface
	auditLogOperation operation.AuditLogOperationInterface
}

func NewActivityService(
	config *c.HttpServerConfig,
	userOperation operation.UserOperationInterface,
	activityOperation operation.ActivityOperationInterface,
	storeService StoreServiceInterface,
	auditLogOperation operation.AuditLogOperationInterface,
) *ActivityService {
	return &ActivityService{
		config:            config,
		userOperation:     userOperation,
		activityOperation: activityOperation,
		storeService:      storeService,
		auditLogOperation: auditLogOperation,
	}
}

var SuccessGetActivities = ApiStatus{StatusName: "GET_ACTIVITIES", Description: "成功获取活动", HttpCode: Ok}

func (activityService *ActivityService) GetActivities(req *RequestGetActivities) *ApiResponse[ResponseGetActivities] {
	targetMonth, _ := time.Parse("2006-01", req.Time)
	firstDay := targetMonth.AddDate(0, -1, 0)
	lastDay := targetMonth.AddDate(0, 2, 0).Add(-time.Second)
	activities, err := activityService.activityOperation.GetActivities(firstDay, lastDay)
	if err != nil {
		return NewApiResponse[ResponseGetActivities](&ErrDatabaseFail, Unsatisfied, nil)
	}
	return NewApiResponse[ResponseGetActivities](&SuccessGetActivities, Unsatisfied, &ResponseGetActivities{Items: activities})
}

var SuccessGetActivitiesPage = ApiStatus{StatusName: "GET_ACTIVITIES_PAGE", Description: "成功获取活动分页", HttpCode: Ok}

func (activityService *ActivityService) GetActivitiesPage(req *RequestGetActivitiesPage) *ApiResponse[ResponseGetActivitiesPage] {
	if req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseGetActivitiesPage](&ErrIllegalParam, Unsatisfied, nil)
	}
	if req.Permission <= 0 {
		return NewApiResponse[ResponseGetActivitiesPage](&ErrNoPermission, Unsatisfied, nil)
	}
	permission := operation.Permission(req.Permission)
	if !permission.HasPermission(operation.ActivityShowList) {
		return NewApiResponse[ResponseGetActivitiesPage](&ErrNoPermission, Unsatisfied, nil)
	}
	activities, total, err := activityService.activityOperation.GetActivitiesPage(req.Page, req.PageSize)
	if err != nil {
		return NewApiResponse[ResponseGetActivitiesPage](&ErrDatabaseFail, Unsatisfied, nil)
	}
	return NewApiResponse(&SuccessGetActivitiesPage, Unsatisfied, &ResponseGetActivitiesPage{
		Items:    activities,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	})
}

var SuccessGetActivityInfo = ApiStatus{StatusName: "GET_ACTIVITY_INFO", Description: "成功获取活动信息", HttpCode: Ok}

func (activityService *ActivityService) GetActivityInfo(req *RequestActivityInfo) *ApiResponse[ResponseActivityInfo] {
	if req.ActivityId <= 0 {
		return NewApiResponse[ResponseActivityInfo](&ErrIllegalParam, Unsatisfied, nil)
	}
	activity, res := CallDBFuncAndCheckError[operation.Activity, ResponseActivityInfo](func() (*operation.Activity, error) {
		return activityService.activityOperation.GetActivityById(req.ActivityId)
	})
	if res != nil {
		return res
	}
	return NewApiResponse(&SuccessGetActivityInfo, Unsatisfied, (*ResponseActivityInfo)(activity))
}

var SuccessAddActivity = ApiStatus{StatusName: "ADD_ACTIVITY", Description: "成功添加活动", HttpCode: Ok}

func (activityService *ActivityService) AddActivity(req *RequestAddActivity) *ApiResponse[ResponseAddActivity] {
	if req.Permission <= 0 {
		return NewApiResponse[ResponseAddActivity](&ErrNoPermission, Unsatisfied, nil)
	}
	permission := operation.Permission(req.Permission)
	if !permission.HasPermission(operation.ActivityPublish) {
		return NewApiResponse[ResponseAddActivity](&ErrNoPermission, Unsatisfied, nil)
	}
	req.Activity.Publisher = req.Cid
	err := activityService.activityOperation.SaveActivity(req.Activity)
	if err != nil {
		c.ErrorF("Error adding activity: %v", err)
		return NewApiResponse[ResponseAddActivity](&ErrDatabaseFail, Unsatisfied, nil)
	}

	go func() {
		newValue, _ := json.Marshal(req.Activity)
		changeDetail := &operation.ChangeDetail{
			OldValue: "",
			NewValue: string(newValue),
		}
		auditLog := activityService.auditLogOperation.NewAuditLog(operation.ActivityCreated, req.Cid,
			strconv.Itoa(int(req.Activity.ID)), req.Ip, req.UserAgent, changeDetail)
		err := activityService.auditLogOperation.SaveAuditLog(auditLog)
		if err != nil {
			c.ErrorF("Fail to create audit log for activity_created, detail: %v", err)
		}
	}()

	return NewApiResponse[ResponseAddActivity](&SuccessAddActivity, Unsatisfied, &ResponseAddActivity{Activity: req.Activity})
}

var (
	SuccessDeleteActivity = ApiStatus{StatusName: "DELETE_ACTIVITY", Description: "成功删除活动", HttpCode: Ok}
)

func (activityService *ActivityService) DeleteActivity(req *RequestDeleteActivity) *ApiResponse[ResponseDeleteActivity] {
	if req.Permission <= 0 {
		return NewApiResponse[ResponseDeleteActivity](&ErrNoPermission, Unsatisfied, nil)
	}
	permission := operation.Permission(req.Permission)
	if !permission.HasPermission(operation.ActivityDelete) {
		return NewApiResponse[ResponseDeleteActivity](&ErrNoPermission, Unsatisfied, nil)
	}
	activity, err := activityService.activityOperation.GetActivityById(req.ActivityId)
	if err != nil {
		if errors.Is(err, operation.ErrActivityNotFound) {
			return NewApiResponse[ResponseDeleteActivity](&ErrActivityNotFound, Unsatisfied, nil)
		}
		return NewApiResponse[ResponseDeleteActivity](&ErrDatabaseFail, Unsatisfied, nil)
	}
	if err := activityService.activityOperation.DeleteActivity(activity); err != nil {
		return NewApiResponse[ResponseDeleteActivity](&ErrDatabaseFail, Unsatisfied, nil)
	}

	go func() {
		auditLog := activityService.auditLogOperation.NewAuditLog(operation.ActivityDeleted, req.Cid,
			strconv.Itoa(int(activity.ID)), req.Ip, req.UserAgent, nil)
		err := activityService.auditLogOperation.SaveAuditLog(auditLog)
		if err != nil {
			c.ErrorF("Fail to create audit log for activity_deleted, detail: %v", err)
		}
	}()

	data := ResponseDeleteActivity(true)
	return NewApiResponse(&SuccessDeleteActivity, Unsatisfied, &data)
}

var (
	ErrRatingTooLow          = ApiStatus{StatusName: "RATING_TOO_LOW", Description: "管制权限不够", HttpCode: PermissionDenied}
	ErrFacilityAlreadyExist  = ApiStatus{StatusName: "FACILITY_ALREADY_EXIST", Description: "你不能同时报名两个以上的席位", HttpCode: Conflict}
	ErrFacilityAlreadySigned = ApiStatus{StatusName: "FACILITY_ALREADY_SIGNED", Description: "已有其他管制员报名", HttpCode: Conflict}
	SuccessSignFacility      = ApiStatus{StatusName: "SIGNED_FACILITY", Description: "报名成功", HttpCode: Ok}
)

func (activityService *ActivityService) ControllerJoin(req *RequestControllerJoin) *ApiResponse[ResponseControllerJoin] {
	if req.Rating <= fsd.Observer.Index() {
		return NewApiResponse[ResponseControllerJoin](&ErrRatingTooLow, Unsatisfied, nil)
	}
	user, res := CallDBFuncAndCheckError[operation.User, ResponseControllerJoin](func() (*operation.User, error) {
		return activityService.userOperation.GetUserByUid(req.Uid)
	})
	if res != nil {
		return res
	}
	facility, res := CallDBFuncAndCheckError[operation.ActivityFacility, ResponseControllerJoin](func() (*operation.ActivityFacility, error) {
		return activityService.activityOperation.GetFacilityById(req.FacilityId)
	})
	if res != nil {
		return res
	}
	err := activityService.activityOperation.SignFacilityController(facility, user)
	if err != nil {
		if errors.Is(err, operation.ErrRatingNotAllowed) {
			return NewApiResponse[ResponseControllerJoin](&ErrRatingTooLow, Unsatisfied, nil)
		}
		if errors.Is(err, operation.ErrFacilityAlreadyExists) {
			return NewApiResponse[ResponseControllerJoin](&ErrFacilityAlreadyExist, Unsatisfied, nil)
		}
		if errors.Is(err, operation.ErrFacilitySigned) {
			return NewApiResponse[ResponseControllerJoin](&ErrFacilityAlreadySigned, Unsatisfied, nil)
		}
		return NewApiResponse[ResponseControllerJoin](&ErrDatabaseFail, Unsatisfied, nil)
	}
	data := ResponseControllerJoin(true)
	return NewApiResponse(&SuccessSignFacility, Unsatisfied, &data)
}

var (
	ErrFacilityUnSigned   = ApiStatus{StatusName: "FACILITY_UNSIGNED", Description: "该席位尚未有人报名", HttpCode: Conflict}
	SuccessUnsignFacility = ApiStatus{StatusName: "UNSIGNED_FACILITY", Description: "成功取消报名", HttpCode: Ok}
)

func (activityService *ActivityService) ControllerLeave(req *RequestControllerLeave) *ApiResponse[ResponseControllerLeave] {
	facility, res := CallDBFuncAndCheckError[operation.ActivityFacility, ResponseControllerLeave](func() (*operation.ActivityFacility, error) {
		return activityService.activityOperation.GetFacilityById(req.FacilityId)
	})
	if res != nil {
		return res
	}
	err := activityService.activityOperation.UnsignFacilityController(facility, req.Cid)
	if err != nil {
		if errors.Is(err, operation.ErrFacilityNotSigned) {
			return NewApiResponse[ResponseControllerLeave](&ErrFacilityUnSigned, Unsatisfied, nil)
		}
		return NewApiResponse[ResponseControllerLeave](&ErrDatabaseFail, Unsatisfied, nil)
	}
	data := ResponseControllerLeave(true)
	return NewApiResponse(&SuccessUnsignFacility, Unsatisfied, &data)
}

var (
	ErrAlreadySigned      = ApiStatus{StatusName: "ALREADY_SIGNED", Description: "你已经报名该活动了", HttpCode: Conflict}
	ErrCallsignUsed       = ApiStatus{StatusName: "CALLSIGN_USED", Description: "呼号已被占用", HttpCode: Conflict}
	SuccessSignedActivity = ApiStatus{StatusName: "SIGNED_ACTIVITY", Description: "报名成功", HttpCode: Ok}
)

func (activityService *ActivityService) PilotJoin(req *RequestPilotJoin) *ApiResponse[ResponsePilotJoin] {
	err := activityService.activityOperation.SignActivityPilot(req.ActivityId, req.Cid, req.Callsign, req.AircraftType)
	if err != nil {
		if errors.Is(err, operation.ErrActivityAlreadySigned) {
			return NewApiResponse[ResponsePilotJoin](&ErrAlreadySigned, Unsatisfied, nil)
		}
		if errors.Is(err, operation.ErrCallsignAlreadyUsed) {
			return NewApiResponse[ResponsePilotJoin](&ErrCallsignUsed, Unsatisfied, nil)
		}
		return NewApiResponse[ResponsePilotJoin](&ErrDatabaseFail, Unsatisfied, nil)
	}
	data := ResponsePilotJoin(true)
	return NewApiResponse(&SuccessSignedActivity, Unsatisfied, &data)
}

var (
	ErrNoSigned             = ApiStatus{StatusName: "NO_SIGNED", Description: "你还没有报名该活动", HttpCode: Conflict}
	SuccessUnsignedActivity = ApiStatus{StatusName: "UNSIGNED_ACTIVITY", Description: "取消报名成功", HttpCode: Ok}
)

func (activityService *ActivityService) PilotLeave(req *RequestPilotLeave) *ApiResponse[ResponsePilotLeave] {
	err := activityService.activityOperation.UnsignActivityPilot(req.ActivityId, req.Cid)
	if err != nil {
		if errors.Is(err, operation.ErrActivityUnsigned) {
			return NewApiResponse[ResponsePilotLeave](&ErrNoSigned, Unsatisfied, nil)
		}
		return NewApiResponse[ResponsePilotLeave](&ErrDatabaseFail, Unsatisfied, nil)
	}
	data := ResponsePilotLeave(true)
	return NewApiResponse(&SuccessUnsignedActivity, Unsatisfied, &data)
}

var SuccessEditActivity = ApiStatus{StatusName: "EDIT_ACTIVITY", Description: "修改活动成功", HttpCode: Ok}

func (activityService *ActivityService) EditActivity(req *RequestEditActivity) *ApiResponse[ResponseEditActivity] {
	if req.Activity == nil {
		return NewApiResponse[ResponseEditActivity](&ErrIllegalParam, Unsatisfied, nil)
	}
	if req.Permission <= 0 {
		return NewApiResponse[ResponseEditActivity](&ErrNoPermission, Unsatisfied, nil)
	}
	permission := operation.Permission(req.Permission)
	if !permission.HasPermission(operation.ActivityEdit) {
		return NewApiResponse[ResponseEditActivity](&ErrNoPermission, Unsatisfied, nil)
	}
	activity, res := CallDBFuncAndCheckError[operation.Activity, ResponseEditActivity](func() (*operation.Activity, error) {
		return activityService.activityOperation.GetActivityById(req.ID)
	})
	oldValue, _ := json.Marshal(activity)
	if res != nil {
		return res
	}
	if req.ImageUrl != "" && req.ImageUrl != activity.ImageUrl && activity.ImageUrl != "" {
		_, err := activityService.storeService.DeleteImageFile(activity.ImageUrl)
		if err != nil {
			c.ErrorF("err while delete old activity image, %v", err)
		}
	}
	updateInfo := req.Activity.Diff(activity)
	err := activityService.activityOperation.UpdateActivityInfo(activity, req.Activity, updateInfo)
	if err != nil {
		return NewApiResponse[ResponseEditActivity](&ErrDatabaseFail, Unsatisfied, nil)
	}

	go func() {
		newValue, _ := json.Marshal(req.Activity)
		changeDetail := &operation.ChangeDetail{
			OldValue: string(oldValue),
			NewValue: string(newValue),
		}
		auditLog := activityService.auditLogOperation.NewAuditLog(operation.ActivityUpdated, req.Cid,
			strconv.Itoa(int(req.Activity.ID)), req.Ip, req.UserAgent, changeDetail)
		err := activityService.auditLogOperation.SaveAuditLog(auditLog)
		if err != nil {
			c.ErrorF("Fail to create audit log for activity_created, detail: %v", err)
		}
	}()

	data := ResponseEditActivity(true)
	return NewApiResponse(&SuccessEditActivity, Unsatisfied, &data)
}

var SuccessEditActivityStatus = ApiStatus{StatusName: "EDIT_ACTIVITY_STATUS", Description: "成功修改活动状态", HttpCode: Ok}

func (activityService *ActivityService) EditActivityStatus(req *RequestEditActivityStatus) *ApiResponse[ResponseEditActivityStatus] {
	if req.Status < int(operation.Open) || req.Status > int(operation.Closed) {
		return NewApiResponse[ResponseEditActivityStatus](&ErrIllegalParam, Unsatisfied, nil)
	}

	status := operation.ActivityStatus(req.Status)

	if req.Permission <= 0 {
		return NewApiResponse[ResponseEditActivityStatus](&ErrNoPermission, Unsatisfied, nil)
	}
	permission := operation.Permission(req.Permission)
	if !permission.HasPermission(operation.ActivityEditState) {
		return NewApiResponse[ResponseEditActivityStatus](&ErrNoPermission, Unsatisfied, nil)
	}
	err := activityService.activityOperation.SetActivityStatus(req.ActivityId, status)
	if err != nil {
		return NewApiResponse[ResponseEditActivityStatus](&ErrDatabaseFail, Unsatisfied, nil)
	}
	data := ResponseEditActivityStatus(true)
	return NewApiResponse(&SuccessEditActivityStatus, Unsatisfied, &data)
}

var SuccessEditPilotsStatus = ApiStatus{StatusName: "EDIT_PILOTS_STATUS", Description: "成功修改活动机组状态", HttpCode: Ok}

func (activityService *ActivityService) EditPilotStatus(req *RequestEditPilotStatus) *ApiResponse[ResponseEditPilotStatus] {
	if req.Status < int(operation.Signed) || req.Status > int(operation.Landing) {
		return NewApiResponse[ResponseEditPilotStatus](&ErrIllegalParam, Unsatisfied, nil)
	}

	status := operation.ActivityPilotStatus(req.Status)

	if req.Permission <= 0 {
		return NewApiResponse[ResponseEditPilotStatus](&ErrNoPermission, Unsatisfied, nil)
	}
	permission := operation.Permission(req.Permission)
	if !permission.HasPermission(operation.ActivityEditState) {
		return NewApiResponse[ResponseEditPilotStatus](&ErrNoPermission, Unsatisfied, nil)
	}
	pilot, res := CallDBFuncAndCheckError[operation.ActivityPilot, ResponseEditPilotStatus](func() (*operation.ActivityPilot, error) {
		return activityService.activityOperation.GetActivityPilotById(req.ActivityId, req.Cid)
	})
	if res != nil {
		return res
	}
	err := activityService.activityOperation.SetActivityPilotStatus(pilot, status)
	if err != nil {
		return NewApiResponse[ResponseEditPilotStatus](&ErrDatabaseFail, Unsatisfied, nil)
	}
	data := ResponseEditPilotStatus(true)
	return NewApiResponse(&SuccessEditPilotsStatus, Unsatisfied, &data)
}
