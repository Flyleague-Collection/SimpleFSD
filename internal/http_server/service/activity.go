// Package service
// 存放 ActivityServiceInterface 的实现
package service

import (
	"encoding/json"
	"errors"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"strconv"
	"time"
)

type ActivityService struct {
	logger            log.LoggerInterface
	config            *config.HttpServerConfig
	messageQueue      queue.MessageQueueInterface
	userOperation     operation.UserOperationInterface
	activityOperation operation.ActivityOperationInterface
	storeService      StoreServiceInterface
	auditLogOperation operation.AuditLogOperationInterface
}

func NewActivityService(
	logger log.LoggerInterface,
	config *config.HttpServerConfig,
	messageQueue queue.MessageQueueInterface,
	userOperation operation.UserOperationInterface,
	activityOperation operation.ActivityOperationInterface,
	auditLogOperation operation.AuditLogOperationInterface,
	storeService StoreServiceInterface,
) *ActivityService {
	return &ActivityService{
		logger:            logger,
		config:            config,
		messageQueue:      messageQueue,
		userOperation:     userOperation,
		activityOperation: activityOperation,
		storeService:      storeService,
		auditLogOperation: auditLogOperation,
	}
}

var (
	ErrParseTime         = NewApiStatus("TIME_FORMAT_ERROR", "格式错误", BadRequest)
	SuccessGetActivities = NewApiStatus("GET_ACTIVITIES", "成功获取活动", Ok)
)

func (activityService *ActivityService) GetActivities(req *RequestGetActivities) *ApiResponse[ResponseGetActivities] {
	targetMonth, err := time.Parse("2006-01", req.Time)
	if err != nil {
		return NewApiResponse[ResponseGetActivities](ErrParseTime, nil)
	}
	firstDay := targetMonth.AddDate(0, -1, 0)
	lastDay := targetMonth.AddDate(0, 2, 0).Add(-time.Second)
	activities, res := CallDBFunc[[]*operation.Activity, ResponseGetActivities](func() ([]*operation.Activity, error) {
		return activityService.activityOperation.GetActivities(firstDay, lastDay)
	})
	if res != nil {
		return res
	}
	data := ResponseGetActivities(activities)
	return NewApiResponse[ResponseGetActivities](SuccessGetActivities, &data)
}

var SuccessGetActivitiesPage = NewApiStatus("GET_ACTIVITIES_PAGE", "成功获取活动分页", Ok)

func (activityService *ActivityService) GetActivitiesPage(req *RequestGetActivitiesPage) *ApiResponse[ResponseGetActivitiesPage] {
	if req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseGetActivitiesPage](ErrIllegalParam, nil)
	}
	if res := CheckPermission[ResponseGetActivitiesPage](req.Permission, operation.ActivityShowList); res != nil {
		return res
	}
	activities, total, err := activityService.activityOperation.GetActivitiesPage(req.Page, req.PageSize)
	if res := CheckDatabaseError[ResponseGetActivitiesPage](err); res != nil {
		return res
	}
	return NewApiResponse(SuccessGetActivitiesPage, &ResponseGetActivitiesPage{
		Items:    activities,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	})
}

var SuccessGetActivityInfo = NewApiStatus("GET_ACTIVITY_INFO", "成功获取活动信息", Ok)

func (activityService *ActivityService) GetActivityInfo(req *RequestActivityInfo) *ApiResponse[ResponseActivityInfo] {
	if req.ActivityId <= 0 {
		return NewApiResponse[ResponseActivityInfo](ErrIllegalParam, nil)
	}
	activity, res := CallDBFunc[*operation.Activity, ResponseActivityInfo](func() (*operation.Activity, error) {
		return activityService.activityOperation.GetActivityById(req.ActivityId)
	})
	if res != nil {
		return res
	}
	return NewApiResponse(SuccessGetActivityInfo, (*ResponseActivityInfo)(activity))
}

var SuccessAddActivity = NewApiStatus("ADD_ACTIVITY", "成功添加活动", Ok)

func (activityService *ActivityService) AddActivity(req *RequestAddActivity) *ApiResponse[ResponseAddActivity] {
	if req.Activity == nil {
		return NewApiResponse[ResponseAddActivity](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseAddActivity](req.Permission, operation.ActivityPublish); res != nil {
		return res
	}

	req.Activity.ID = 0
	req.Activity.Publisher = req.Cid

	if res := CallDBFuncWithoutRet[ResponseAddActivity](func() error {
		return activityService.activityOperation.SaveActivity(req.Activity)
	}); res != nil {
		return res
	}

	newValue, _ := json.Marshal(req.Activity)
	activityService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: activityService.auditLogOperation.NewAuditLog(
			operation.ActivityCreated,
			req.Cid,
			strconv.Itoa(int(req.Activity.ID)),
			req.Ip,
			req.UserAgent,
			&operation.ChangeDetail{
				OldValue: operation.ValueNotAvailable,
				NewValue: string(newValue),
			}),
	})

	data := ResponseAddActivity(true)
	return NewApiResponse[ResponseAddActivity](SuccessAddActivity, &data)
}

var SuccessDeleteActivity = NewApiStatus("DELETE_ACTIVITY", "成功删除活动", Ok)

func (activityService *ActivityService) DeleteActivity(req *RequestDeleteActivity) *ApiResponse[ResponseDeleteActivity] {
	if req.ActivityId <= 0 {
		return NewApiResponse[ResponseDeleteActivity](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseDeleteActivity](req.Permission, operation.ActivityDelete); res != nil {
		return res
	}

	if res := CallDBFuncWithoutRet[ResponseDeleteActivity](func() error {
		return activityService.activityOperation.DeleteActivity(req.ActivityId)
	}); res != nil {
		return res
	}

	activityService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: activityService.auditLogOperation.NewAuditLog(
			operation.ActivityDeleted,
			req.Cid,
			strconv.Itoa(int(req.ActivityId)),
			req.Ip,
			req.UserAgent,
			nil,
		),
	})

	data := ResponseDeleteActivity(true)
	return NewApiResponse(SuccessDeleteActivity, &data)
}

var (
	ErrRatingTooLow          = NewApiStatus("RATING_TOO_LOW", "管制权限不够", PermissionDenied)
	ErrFacilityAlreadyExist  = NewApiStatus("FACILITY_ALREADY_EXIST", "你不能同时报名两个以上的席位", Conflict)
	ErrFacilityAlreadySigned = NewApiStatus("FACILITY_ALREADY_SIGNED", "已有其他管制员报名", Conflict)
	SuccessSignFacility      = NewApiStatus("SIGNED_FACILITY", "报名成功", Ok)
)

func (activityService *ActivityService) ControllerJoin(req *RequestControllerJoin) *ApiResponse[ResponseControllerJoin] {
	if req.ActivityId <= 0 || req.FacilityId <= 0 {
		return NewApiResponse[ResponseControllerJoin](ErrIllegalParam, nil)
	}

	if req.Rating <= fsd.Observer.Index() {
		return NewApiResponse[ResponseControllerJoin](ErrRatingTooLow, nil)
	}

	if res := WithErrorHandlerWithoutRet[ResponseControllerJoin](func(err error) *ApiResponse[ResponseControllerJoin] {
		if errors.Is(err, operation.ErrRatingNotAllowed) {
			return NewApiResponse[ResponseControllerJoin](ErrRatingTooLow, nil)
		}
		if errors.Is(err, operation.ErrFacilityAlreadyExists) {
			return NewApiResponse[ResponseControllerJoin](ErrFacilityAlreadyExist, nil)
		}
		if errors.Is(err, operation.ErrFacilitySigned) {
			return NewApiResponse[ResponseControllerJoin](ErrFacilityAlreadySigned, nil)
		}
		return nil
	}).CallDBFuncWithoutRet(func() error {
		activity, err := activityService.activityOperation.GetActivityById(req.ActivityId)
		if err != nil {
			return err
		}
		if activity.Status >= int(operation.InActive) {
			return operation.ErrActivityHasClosed
		}
		user, err := activityService.userOperation.GetUserByUid(req.Uid)
		if err != nil {
			return err
		}
		facility, err := activityService.activityOperation.GetFacilityById(req.FacilityId)
		if err != nil {
			return err
		}
		if facility.ActivityId != req.ActivityId {
			return operation.ErrActivityIdMismatch
		}
		return activityService.activityOperation.SignFacilityController(facility, user)
	}); res != nil {
		return res
	}

	data := ResponseControllerJoin(true)
	return NewApiResponse(SuccessSignFacility, &data)
}

var (
	ErrFacilityUnSigned    = NewApiStatus("FACILITY_UNSIGNED", "该席位尚未有人报名", Conflict)
	ErrFacilityNotYourSign = NewApiStatus("FACILITY_NOT_YOUR_SIGN", "这不是你报名的席位", Conflict)
	SuccessUnsignFacility  = NewApiStatus("UNSIGNED_FACILITY", "成功取消报名", Ok)
)

func (activityService *ActivityService) ControllerLeave(req *RequestControllerLeave) *ApiResponse[ResponseControllerLeave] {
	if req.ActivityId <= 0 || req.FacilityId <= 0 {
		return NewApiResponse[ResponseControllerLeave](ErrIllegalParam, nil)
	}

	if res := WithErrorHandlerWithoutRet[ResponseControllerLeave](func(err error) *ApiResponse[ResponseControllerLeave] {
		if errors.Is(err, operation.ErrFacilityNotSigned) {
			return NewApiResponse[ResponseControllerLeave](ErrFacilityUnSigned, nil)
		}
		if errors.Is(err, operation.ErrFacilityNotYourSign) {
			return NewApiResponse[ResponseControllerLeave](ErrFacilityNotYourSign, nil)
		}
		return nil
	}).CallDBFuncWithoutRet(func() error {
		activity, err := activityService.activityOperation.GetActivityById(req.ActivityId)
		if err != nil {
			return err
		}
		if activity.Status >= int(operation.InActive) {
			return operation.ErrActivityHasClosed
		}
		facility, err := activityService.activityOperation.GetFacilityById(req.FacilityId)
		if err != nil {
			return err
		}
		if facility.ActivityId != req.ActivityId {
			return operation.ErrActivityIdMismatch
		}
		return activityService.activityOperation.UnsignFacilityController(facility, req.Uid)
	}); res != nil {
		return res
	}

	data := ResponseControllerLeave(true)
	return NewApiResponse(SuccessUnsignFacility, &data)
}

var (
	ErrAlreadySigned      = NewApiStatus("ALREADY_SIGNED", "你已经报名该活动了", Conflict)
	ErrCallsignUsed       = NewApiStatus("CALLSIGN_USED", "呼号已被占用", Conflict)
	SuccessSignedActivity = NewApiStatus("SIGNED_ACTIVITY", "报名成功", Ok)
)

func (activityService *ActivityService) PilotJoin(req *RequestPilotJoin) *ApiResponse[ResponsePilotJoin] {
	if req.ActivityId <= 0 || req.Callsign == "" || req.AircraftType == "" {
		return NewApiResponse[ResponsePilotJoin](ErrIllegalParam, nil)
	}

	if res := WithErrorHandlerWithoutRet[ResponsePilotJoin](func(err error) *ApiResponse[ResponsePilotJoin] {
		if errors.Is(err, operation.ErrActivityAlreadySigned) {
			return NewApiResponse[ResponsePilotJoin](ErrAlreadySigned, nil)
		}
		if errors.Is(err, operation.ErrCallsignAlreadyUsed) {
			return NewApiResponse[ResponsePilotJoin](ErrCallsignUsed, nil)
		}
		return nil
	}).CallDBFuncWithoutRet(func() error {
		activity, err := activityService.activityOperation.GetActivityById(req.ActivityId)
		if err != nil {
			return err
		}
		if activity.Status >= int(operation.InActive) {
			return operation.ErrActivityHasClosed
		}
		return activityService.activityOperation.SignActivityPilot(req.ActivityId, req.Uid, req.Callsign, req.AircraftType)
	}); res != nil {
		return res
	}

	data := ResponsePilotJoin(true)
	return NewApiResponse(SuccessSignedActivity, &data)
}

var (
	ErrNoSigned             = NewApiStatus("NO_SIGNED", "你还没有报名该活动", Conflict)
	SuccessUnsignedActivity = NewApiStatus("UNSIGNED_ACTIVITY", "取消报名成功", Ok)
)

func (activityService *ActivityService) PilotLeave(req *RequestPilotLeave) *ApiResponse[ResponsePilotLeave] {
	if req.ActivityId <= 0 {
		return NewApiResponse[ResponsePilotLeave](ErrIllegalParam, nil)
	}

	if res := WithErrorHandlerWithoutRet[ResponsePilotLeave](func(err error) *ApiResponse[ResponsePilotLeave] {
		if errors.Is(err, operation.ErrActivityUnsigned) {
			return NewApiResponse[ResponsePilotLeave](ErrNoSigned, nil)
		}
		return nil
	}).CallDBFuncWithoutRet(func() error {
		activity, err := activityService.activityOperation.GetActivityById(req.ActivityId)
		if err != nil {
			return err
		}
		if activity.Status >= int(operation.InActive) {
			return operation.ErrActivityHasClosed
		}
		return activityService.activityOperation.UnsignActivityPilot(req.ActivityId, req.Uid)
	}); res != nil {
		return res
	}

	data := ResponsePilotLeave(true)
	return NewApiResponse(SuccessUnsignedActivity, &data)
}

var SuccessEditActivity = NewApiStatus("EDIT_ACTIVITY", "修改活动成功", Ok)

func (activityService *ActivityService) EditActivity(req *RequestEditActivity) *ApiResponse[ResponseEditActivity] {
	if req.Activity == nil {
		return NewApiResponse[ResponseEditActivity](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseEditActivity](req.Permission, operation.ActivityEdit); res != nil {
		return res
	}

	activity, res := CallDBFunc[*operation.Activity, ResponseEditActivity](func() (*operation.Activity, error) {
		return activityService.activityOperation.GetActivityById(req.ID)
	})
	if res != nil {
		return res
	}

	oldValue, _ := json.Marshal(activity)

	if req.ImageUrl != "" && req.ImageUrl != activity.ImageUrl && activity.ImageUrl != "" {
		_, err := activityService.storeService.DeleteImageFile(activity.ImageUrl)
		if err != nil {
			activityService.logger.ErrorF("err while delete old activity image, %v", err)
		}
	}

	updateInfo := req.Activity.Diff(activity)

	if res := CallDBFuncWithoutRet[ResponseEditActivity](func() error {
		return activityService.activityOperation.UpdateActivityInfo(activity, req.Activity, updateInfo)
	}); res != nil {
		return res
	}

	newValue, _ := json.Marshal(req.Activity)
	activityService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: activityService.auditLogOperation.NewAuditLog(
			operation.ActivityUpdated,
			req.Cid,
			strconv.Itoa(int(req.Activity.ID)),
			req.Ip,
			req.UserAgent, &operation.ChangeDetail{
				OldValue: string(oldValue),
				NewValue: string(newValue),
			},
		),
	})

	data := ResponseEditActivity(true)
	return NewApiResponse(SuccessEditActivity, &data)
}

var SuccessEditActivityStatus = NewApiStatus("EDIT_ACTIVITY_STATUS", "成功修改活动状态", Ok)

func (activityService *ActivityService) EditActivityStatus(req *RequestEditActivityStatus) *ApiResponse[ResponseEditActivityStatus] {
	if req.ActivityId <= 0 || req.Status < int(operation.Open) || req.Status > int(operation.Closed) {
		return NewApiResponse[ResponseEditActivityStatus](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseEditActivityStatus](req.Permission, operation.ActivityEditState); res != nil {
		return res
	}

	status := operation.ActivityStatus(req.Status)

	if res := CallDBFuncWithoutRet[ResponseEditActivityStatus](func() error {
		return activityService.activityOperation.SetActivityStatus(req.ActivityId, status)
	}); res != nil {
		return res
	}

	data := ResponseEditActivityStatus(true)
	return NewApiResponse(SuccessEditActivityStatus, &data)
}

var SuccessEditPilotsStatus = NewApiStatus("EDIT_PILOTS_STATUS", "成功修改活动机组状态", Ok)

func (activityService *ActivityService) EditPilotStatus(req *RequestEditPilotStatus) *ApiResponse[ResponseEditPilotStatus] {
	if req.ActivityId <= 0 || req.UserId <= 0 || req.Status < int(operation.Signed) || req.Status > int(operation.Landing) {
		return NewApiResponse[ResponseEditPilotStatus](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseEditPilotStatus](req.Permission, operation.ActivityEditPilotState); res != nil {
		return res
	}

	status := operation.ActivityPilotStatus(req.Status)

	if res := CallDBFuncWithoutRet[ResponseEditPilotStatus](func() error {
		pilot, err := activityService.activityOperation.GetActivityPilotById(req.ActivityId, req.UserId)
		if err != nil {
			return err
		}
		return activityService.activityOperation.SetActivityPilotStatus(pilot, status)
	}); res != nil {
		return res
	}

	data := ResponseEditPilotStatus(true)
	return NewApiResponse(SuccessEditPilotsStatus, &data)
}
