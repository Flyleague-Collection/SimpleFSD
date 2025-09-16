// Package service
package service

import (
	"encoding/json"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"strconv"
	"time"
)

type ControllerService struct {
	logger                    log.LoggerInterface
	config                    *config.HttpServerConfig
	messageQueue              queue.MessageQueueInterface
	userOperation             operation.UserOperationInterface
	controllerOperation       operation.ControllerOperationInterface
	controllerRecordOperation operation.ControllerRecordOperationInterface
	auditLogOperation         operation.AuditLogOperationInterface
}

func NewControllerService(
	logger log.LoggerInterface,
	config *config.HttpServerConfig,
	messageQueue queue.MessageQueueInterface,
	userOperation operation.UserOperationInterface,
	controllerOperation operation.ControllerOperationInterface,
	controllerRecordOperation operation.ControllerRecordOperationInterface,
	auditLogOperation operation.AuditLogOperationInterface,
) *ControllerService {
	return &ControllerService{
		logger:                    logger,
		config:                    config,
		messageQueue:              messageQueue,
		userOperation:             userOperation,
		controllerOperation:       controllerOperation,
		controllerRecordOperation: controllerRecordOperation,
		auditLogOperation:         auditLogOperation,
	}
}

var SuccessGetControllers = ApiStatus{StatusName: "GET_CONTROLLER_PAGE", Description: "获取管制员信息分页成功", HttpCode: Ok}

func (controllerService *ControllerService) GetControllerList(req *RequestControllerList) *ApiResponse[ResponseControllerList] {
	if req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseControllerList](ErrIllegalParam, nil)
	}
	if req.Permission <= 0 {
		return NewApiResponse[ResponseControllerList](ErrNoPermission, nil)
	}
	permission := operation.Permission(req.Permission)
	if !permission.HasPermission(operation.UserShowList) {
		return NewApiResponse[ResponseControllerList](ErrNoPermission, nil)
	}
	users, total, err := controllerService.controllerOperation.GetControllers(req.Page, req.PageSize)
	if err != nil {
		return NewApiResponse[ResponseControllerList](ErrDatabaseFail, nil)
	}
	return NewApiResponse(&SuccessGetControllers, &ResponseControllerList{
		Items:    users,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	})
}

var SuccessGetCurrentControllerRecord = ApiStatus{StatusName: "GET_CURRENT_CONTROLLER_RECORD", Description: "获取管制员履历成功", HttpCode: Ok}

func (controllerService *ControllerService) GetCurrentControllerRecord(req *RequestGetCurrentControllerRecord) *ApiResponse[ResponseGetCurrentControllerRecord] {
	if req.Uid <= 0 {
		return NewApiResponse[ResponseGetCurrentControllerRecord](ErrIllegalParam, nil)
	}

	records, total, err := controllerService.controllerRecordOperation.GetControllerRecords(req.Cid, req.Page, req.PageSize)
	if err != nil {
		return NewApiResponse[ResponseGetCurrentControllerRecord](ErrDatabaseFail, nil)
	}

	return NewApiResponse(&SuccessGetCurrentControllerRecord, &ResponseGetCurrentControllerRecord{
		Items:    records,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

func (controllerService *ControllerService) GetControllerRecord(req *RequestGetControllerRecord) *ApiResponse[ResponseGetControllerRecord] {
	if req.Uid <= 0 || req.TargetCid <= 0 {
		return NewApiResponse[ResponseGetControllerRecord](ErrIllegalParam, nil)
	}

	permission := operation.Permission(req.Permission)
	if !permission.HasPermission(operation.ControllerShowRecord) {
		return NewApiResponse[ResponseGetControllerRecord](ErrNoPermission, nil)
	}

	records, total, err := controllerService.controllerRecordOperation.GetControllerRecords(req.TargetCid, req.Page, req.PageSize)
	if err != nil {
		return NewApiResponse[ResponseGetControllerRecord](ErrDatabaseFail, nil)
	}

	return NewApiResponse(&SuccessGetCurrentControllerRecord, &ResponseGetControllerRecord{
		Items:    records,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

var (
	ErrSameRating                 = ApiStatus{StatusName: "SAME_RATING", Description: "用户已是该权限", HttpCode: BadRequest}
	SuccessUpdateControllerRating = ApiStatus{StatusName: "UPDATE_CONTROLLER_RATING", Description: "编辑用户管制权限成功", HttpCode: Ok}
)

func (controllerService *ControllerService) UpdateControllerRating(req *RequestUpdateControllerRating) *ApiResponse[ResponseUpdateControllerRating] {
	if req.Uid <= 0 || req.TargetUid < 0 || req.Rating < fsd.Ban.Index() || req.Rating > fsd.Administrator.Index() {
		return NewApiResponse[ResponseUpdateControllerRating](ErrIllegalParam, nil)
	}
	user, targetUser, res := GetUsersAndCheckPermission[ResponseUpdateControllerRating](controllerService.userOperation, req.Uid, req.TargetUid, operation.ControllerEditRating)
	if res != nil {
		return res
	}
	oldRating := fsd.Rating(targetUser.Rating)
	newRating := fsd.Rating(req.Rating)
	if oldRating == newRating {
		return NewApiResponse[ResponseUpdateControllerRating](&ErrSameRating, nil)
	}

	if _, res := CallDBFunc[interface{}, ResponseUpdateControllerRating](func() (*interface{}, error) {
		return nil, controllerService.controllerOperation.SetControllerRating(targetUser, newRating.Index())
	}); res != nil {
		return res
	}

	if controllerService.config.Email.Template.EnableRatingChangeEmail {
		controllerService.messageQueue.Publish(&queue.Message{
			Type: queue.SendRatingChangeEmail,
			Data: &SendRatingChangeData{
				User:      targetUser,
				Operator:  user,
				OldRating: oldRating,
				NewRating: newRating,
			},
		})
	}

	controllerService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: controllerService.auditLogOperation.NewAuditLog(operation.ControllerRatingChange, req.Cid,
			strconv.Itoa(targetUser.Cid), req.Ip, req.UserAgent, &operation.ChangeDetail{
				OldValue: oldRating.String(),
				NewValue: newRating.String(),
			},
		),
	})

	data := ResponseUpdateControllerRating(true)
	return NewApiResponse(&SuccessUpdateControllerRating, &data)
}

var (
	ErrNoChangeRequired       = NewApiStatus("NO_CHANGE_REQUIRED", "已经处于UnderMonitor", Conflict)
	SuccessChangeUnderMonitor = NewApiStatus("CHANGE_UNDER_MONITOR", "修改状态成功", Ok)
)

func (controllerService *ControllerService) UpdateControllerUnderMonitor(req *RequestUpdateControllerUnderMonitor) *ApiResponse[ResponseUpdateControllerUnderMonitor] {
	if req.Uid <= 0 || req.TargetUid < 0 {
		return NewApiResponse[ResponseUpdateControllerUnderMonitor](ErrIllegalParam, nil)
	}

	targetUser, res := GetUserAndCheckPermission[ResponseUpdateControllerUnderMonitor](controllerService.userOperation, req.Permission, req.TargetUid, operation.ControllerChangeUnderMonitor)
	if res != nil {
		return res
	}

	if targetUser.UnderMonitor == req.UnderMonitor {
		return NewApiResponse[ResponseUpdateControllerUnderMonitor](ErrNoChangeRequired, nil)
	}

	if res := CallDBFuncWithoutRet[ResponseUpdateControllerUnderMonitor](func() error {
		return controllerService.controllerOperation.SetControllerUnderMonitor(targetUser, req.UnderMonitor)
	}); res != nil {
		return res
	}

	controllerService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: controllerService.auditLogOperation.NewAuditLog(
			operation.ControllerUMChange,
			req.Cid,
			fmt.Sprintf("%04d", targetUser.Cid),
			req.Ip,
			req.UserAgent,
			&operation.ChangeDetail{
				OldValue: strconv.FormatBool(!req.UnderMonitor),
				NewValue: strconv.FormatBool(req.UnderMonitor),
			},
		),
	})

	data := ResponseUpdateControllerUnderMonitor(true)
	return NewApiResponse(SuccessChangeUnderMonitor, &data)
}

var SuccessUpdateControllerSolo = NewApiStatus("UPDATE_CONTROLLER_SOLO", "修改SOLO状态成功", Ok)

func (controllerService *ControllerService) UpdateControllerUnderSolo(req *RequestUpdateControllerUnderSolo) *ApiResponse[ResponseUpdateControllerUnderSolo] {
	if req.Uid <= 0 || req.TargetUid < 0 {
		return NewApiResponse[ResponseUpdateControllerUnderSolo](ErrIllegalParam, nil)
	}

	targetUser, res := GetUserAndCheckPermission[ResponseUpdateControllerUnderSolo](controllerService.userOperation, req.Permission, req.TargetUid, operation.ControllerChangeSolo)
	if res != nil {
		return res
	}

	if targetUser.UnderSolo == req.Solo && (!targetUser.UnderMonitor || targetUser.SoloUntil.Equal(req.EndTime)) {
		return NewApiResponse[ResponseUpdateControllerUnderSolo](ErrNoChangeRequired, nil)
	}

	details := &operation.ChangeDetail{}

	if targetUser.UnderSolo {
		details.OldValue = fmt.Sprintf("true(%s)", targetUser.SoloUntil.Format(time.DateTime))
	} else {
		details.OldValue = "false"
	}

	if res := CallDBFuncWithoutRet[ResponseUpdateControllerUnderSolo](func() error {
		if req.Solo {
			return controllerService.controllerOperation.SetControllerSolo(targetUser, req.EndTime)
		} else {
			return controllerService.controllerOperation.UnsetControllerSolo(targetUser)
		}
	}); res != nil {
		return res
	}

	if req.Solo {
		details.NewValue = fmt.Sprintf("true(%s)", req.EndTime.Format(time.DateTime))
	} else {
		details.NewValue = "false"
	}

	controllerService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: controllerService.auditLogOperation.NewAuditLog(
			operation.ControllerSoloChange,
			req.Cid,
			fmt.Sprintf("%04d", targetUser.Cid),
			req.Ip,
			req.UserAgent,
			details,
		),
	})

	data := ResponseUpdateControllerUnderSolo(true)
	return NewApiResponse(SuccessUpdateControllerSolo, &data)
}

var SuccessUpdateControllerGuest = NewApiStatus("UPDATE_CONTROLLER_GUEST", "修改客座状态成功", Ok)

func (controllerService *ControllerService) UpdateControllerGuest(req *RequestUpdateControllerGuest) *ApiResponse[ResponseUpdateControllerGuest] {
	if req.Uid <= 0 || req.TargetUid < 0 {
		return NewApiResponse[ResponseUpdateControllerGuest](ErrIllegalParam, nil)
	}

	targetUser, res := GetUserAndCheckPermission[ResponseUpdateControllerGuest](controllerService.userOperation, req.Permission, req.TargetUid, operation.ControllerChangeGuest)
	if res != nil {
		return res
	}

	if targetUser.Guest == req.Guest && (!targetUser.Guest || targetUser.Rating == req.Rating) {
		return NewApiResponse[ResponseUpdateControllerGuest](ErrNoChangeRequired, nil)
	}

	details := &operation.ChangeDetail{}

	if targetUser.Guest {
		details.OldValue = fmt.Sprintf("true(%s)", fsd.Rating(targetUser.Rating).String())
	} else {
		details.OldValue = "false"
	}

	if res := CallDBFuncWithoutRet[ResponseUpdateControllerGuest](func() error {
		if req.Guest {
			return controllerService.controllerOperation.SetControllerGuestRating(targetUser, req.Rating)
		} else {
			return controllerService.controllerOperation.SetControllerGuest(targetUser, false)
		}
	}); res != nil {
		return res
	}

	if req.Guest {
		details.NewValue = fmt.Sprintf("true(%s)", fsd.Rating(req.Rating).String())
	} else {
		details.NewValue = "false"
	}

	controllerService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: controllerService.auditLogOperation.NewAuditLog(
			operation.ControllerGuestChange,
			req.Cid,
			fmt.Sprintf("%04d", targetUser.Cid),
			req.Ip,
			req.UserAgent,
			details,
		),
	})

	data := ResponseUpdateControllerGuest(true)
	return NewApiResponse(SuccessUpdateControllerGuest, &data)
}

var SuccessAddControllerRecord = NewApiStatus("ADD_CONTROLLER_RECORD", "添加管制员履历成功", Ok)

func (controllerService *ControllerService) AddControllerRecord(req *RequestAddControllerRecord) *ApiResponse[ResponseAddControllerRecord] {
	if req.Uid <= 0 || req.TargetCid < 0 || req.Content == "" || !operation.IsValidControllerRecordType(req.Type) {
		return NewApiResponse[ResponseAddControllerRecord](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseAddControllerRecord](req.Permission, operation.ControllerCreateRecord); res != nil {
		return res
	}

	controllerRecordType := operation.ToControllerRecordType(req.Type)

	record := controllerService.controllerRecordOperation.NewControllerRecord(req.TargetCid, req.Cid, controllerRecordType, req.Content)

	if res := CallDBFuncWithoutRet[ResponseAddControllerRecord](func() error {
		return controllerService.controllerRecordOperation.SaveControllerRecord(record)
	}); res != nil {
		return res
	}

	newValue, _ := json.Marshal(record)
	controllerService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: controllerService.auditLogOperation.NewAuditLog(
			operation.ControllerRecordCreated,
			req.Cid,
			fmt.Sprintf("%04d", req.TargetCid),
			req.Ip,
			req.UserAgent,
			&operation.ChangeDetail{
				OldValue: "NOT AVAILABLE",
				NewValue: string(newValue),
			},
		),
	})

	data := ResponseAddControllerRecord(true)
	return NewApiResponse(SuccessAddControllerRecord, &data)
}

var SuccessDeleteControllerRecord = NewApiStatus("DELETE_CONTROLLER_RECORD", "删除管制员履历成功", Ok)

func (controllerService *ControllerService) DeleteControllerRecord(req *RequestDeleteControllerRecord) *ApiResponse[ResponseDeleteControllerRecord] {
	if req.Uid <= 0 || req.TargetRecord < 0 {
		return NewApiResponse[ResponseDeleteControllerRecord](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseDeleteControllerRecord](req.Permission, operation.ControllerDeleteRecord); res != nil {
		return res
	}

	record, res := CallDBFunc[operation.ControllerRecord, ResponseDeleteControllerRecord](func() (*operation.ControllerRecord, error) {
		return controllerService.controllerRecordOperation.GetControllerRecord(req.TargetRecord)
	})
	if res != nil {
		return res
	}

	if res := CallDBFuncWithoutRet[ResponseDeleteControllerRecord](func() error {
		return controllerService.controllerRecordOperation.DeleteControllerRecord(req.TargetRecord)
	}); res != nil {
		return res
	}

	oldValue, _ := json.Marshal(record)
	controllerService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: controllerService.auditLogOperation.NewAuditLog(
			operation.ControllerRecordDeleted,
			req.Cid,
			fmt.Sprintf("%d", req.TargetRecord),
			req.Ip,
			req.UserAgent,
			&operation.ChangeDetail{
				OldValue: string(oldValue),
				NewValue: "NOT AVAILABLE",
			},
		),
	})

	data := ResponseDeleteControllerRecord(true)
	return NewApiResponse(SuccessDeleteControllerRecord, &data)
}
