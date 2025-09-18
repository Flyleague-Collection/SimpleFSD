// Package service
// 存放 ControllerServiceInterface 的实现
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
		logger:                    log.NewLoggerAdapter(logger, "ControllerService"),
		config:                    config,
		messageQueue:              messageQueue,
		userOperation:             userOperation,
		controllerOperation:       controllerOperation,
		controllerRecordOperation: controllerRecordOperation,
		auditLogOperation:         auditLogOperation,
	}
}

var SuccessGetControllers = NewApiStatus("GET_CONTROLLER_PAGE", "获取管制员信息分页成功", Ok)

func (controllerService *ControllerService) GetControllerList(req *RequestControllerList) *ApiResponse[ResponseControllerList] {
	if req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseControllerList](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseControllerList](req.Permission, operation.ControllerShowList); res != nil {
		return res
	}

	users, total, err := controllerService.controllerOperation.GetControllers(req.Page, req.PageSize)
	if res := CheckDatabaseError[ResponseControllerList](err); res != nil {
		return res
	}

	return NewApiResponse(SuccessGetControllers, &ResponseControllerList{
		Items:    users,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	})
}

var SuccessGetCurrentControllerRecord = NewApiStatus("GET_CURRENT_CONTROLLER_RECORD", "获取管制员履历成功", Ok)

func (controllerService *ControllerService) GetCurrentControllerRecord(req *RequestGetCurrentControllerRecord) *ApiResponse[ResponseGetCurrentControllerRecord] {
	if req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseGetCurrentControllerRecord](ErrIllegalParam, nil)
	}

	records, total, err := controllerService.controllerRecordOperation.GetControllerRecords(req.Uid, req.Page, req.PageSize)
	if res := CheckDatabaseError[ResponseGetCurrentControllerRecord](err); res != nil {
		return res
	}

	return NewApiResponse(SuccessGetCurrentControllerRecord, &ResponseGetCurrentControllerRecord{
		Items:    records,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

func (controllerService *ControllerService) GetControllerRecord(req *RequestGetControllerRecord) *ApiResponse[ResponseGetControllerRecord] {
	if req.TargetUid <= 0 || req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseGetControllerRecord](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseGetControllerRecord](req.Permission, operation.ControllerShowRecord); res != nil {
		return res
	}

	records, total, err := controllerService.controllerRecordOperation.GetControllerRecords(req.TargetUid, req.Page, req.PageSize)
	if res := CheckDatabaseError[ResponseGetControllerRecord](err); res != nil {
		return res
	}

	return NewApiResponse(SuccessGetCurrentControllerRecord, &ResponseGetControllerRecord{
		Items:    records,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

var (
	ErrSameRating                 = NewApiStatus("SAME_RATING", "用户已是该权限", BadRequest)
	SuccessUpdateControllerRating = NewApiStatus("UPDATE_CONTROLLER_RATING", "编辑用户管制权限成功", Ok)
)

func (controllerService *ControllerService) UpdateControllerRating(req *RequestUpdateControllerRating) *ApiResponse[ResponseUpdateControllerRating] {
	if req.TargetUid <= 0 || !fsd.IsValidRating(req.Rating) {
		return NewApiResponse[ResponseUpdateControllerRating](ErrIllegalParam, nil)
	}

	user, targetUser, res := GetTargetUserAndCheckPermissionFromDatabase[ResponseUpdateControllerRating](
		controllerService.userOperation,
		req.Uid,
		req.TargetUid,
		operation.ControllerEditRating,
	)
	if res != nil {
		return res
	}

	oldRating := fsd.Rating(targetUser.Rating)
	newRating := fsd.Rating(req.Rating)

	if oldRating == newRating {
		return NewApiResponse[ResponseUpdateControllerRating](ErrSameRating, nil)
	}

	if res := CallDBFuncWithoutRet[ResponseUpdateControllerRating](func() error {
		return controllerService.controllerOperation.SetControllerRating(targetUser, newRating.Index())
	}); res != nil {
		return res
	}

	if controllerService.config.Email.Template.EnableRatingChangeEmail {
		controllerService.messageQueue.Publish(&queue.Message{
			Type: queue.SendRatingChangeEmail,
			Data: &RatingChangeEmailData{
				User:      targetUser,
				Operator:  user,
				OldRating: oldRating,
				NewRating: newRating,
			},
		})
	}

	controllerService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: controllerService.auditLogOperation.NewAuditLog(
			operation.ControllerRatingChange,
			req.Cid,
			fmt.Sprintf("%04d", targetUser.Cid),
			req.Ip,
			req.UserAgent,
			&operation.ChangeDetail{
				OldValue: oldRating.String(),
				NewValue: newRating.String(),
			},
		),
	})

	data := ResponseUpdateControllerRating(true)
	return NewApiResponse(SuccessUpdateControllerRating, &data)
}

var (
	ErrNoChangeRequired       = NewApiStatus("NO_CHANGE_REQUIRED", "状态无需修改", Conflict)
	SuccessChangeUnderMonitor = NewApiStatus("CHANGE_UNDER_MONITOR", "修改状态成功", Ok)
)

func (controllerService *ControllerService) UpdateControllerUnderMonitor(req *RequestUpdateControllerUnderMonitor) *ApiResponse[ResponseUpdateControllerUnderMonitor] {
	if req.TargetUid <= 0 {
		return NewApiResponse[ResponseUpdateControllerUnderMonitor](ErrIllegalParam, nil)
	}

	targetUser, res := GetTargetUserAndCheckPermission[ResponseUpdateControllerUnderMonitor](
		controllerService.userOperation,
		req.Permission,
		req.TargetUid,
		operation.ControllerChangeUnderMonitor,
	)
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
	if req.TargetUid <= 0 || (req.Solo && (req.EndTime.IsZero() || req.EndTime.Before(time.Now()))) {
		return NewApiResponse[ResponseUpdateControllerUnderSolo](ErrIllegalParam, nil)
	}

	targetUser, res := GetTargetUserAndCheckPermission[ResponseUpdateControllerUnderSolo](
		controllerService.userOperation,
		req.Permission,
		req.TargetUid,
		operation.ControllerChangeSolo,
	)
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
	if req.TargetUid <= 0 || !fsd.IsValidRating(req.Rating) {
		return NewApiResponse[ResponseUpdateControllerGuest](ErrIllegalParam, nil)
	}

	targetUser, res := GetTargetUserAndCheckPermission[ResponseUpdateControllerGuest](
		controllerService.userOperation,
		req.Permission,
		req.TargetUid,
		operation.ControllerChangeGuest,
	)
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
	if req.TargetUid <= 0 || req.Content == "" || !operation.IsValidControllerRecordType(req.Type) {
		return NewApiResponse[ResponseAddControllerRecord](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseAddControllerRecord](req.Permission, operation.ControllerCreateRecord); res != nil {
		return res
	}

	targetUser, res := CallDBFunc[*operation.User, ResponseAddControllerRecord](func() (*operation.User, error) {
		return controllerService.userOperation.GetUserByUid(req.TargetUid)
	})
	if res != nil {
		return res
	}

	controllerRecordType := operation.ControllerRecordType(req.Type)

	record := controllerService.controllerRecordOperation.NewControllerRecord(req.TargetUid, req.Cid, controllerRecordType, req.Content)

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
			fmt.Sprintf("%04d", targetUser.Cid),
			req.Ip,
			req.UserAgent,
			&operation.ChangeDetail{
				OldValue: operation.ValueNotAvailable,
				NewValue: string(newValue),
			},
		),
	})

	data := ResponseAddControllerRecord(true)
	return NewApiResponse(SuccessAddControllerRecord, &data)
}

var SuccessDeleteControllerRecord = NewApiStatus("DELETE_CONTROLLER_RECORD", "删除管制员履历成功", Ok)

func (controllerService *ControllerService) DeleteControllerRecord(req *RequestDeleteControllerRecord) *ApiResponse[ResponseDeleteControllerRecord] {
	if req.TargetRecord <= 0 {
		return NewApiResponse[ResponseDeleteControllerRecord](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseDeleteControllerRecord](req.Permission, operation.ControllerDeleteRecord); res != nil {
		return res
	}

	record, res := CallDBFunc[*operation.ControllerRecord, ResponseDeleteControllerRecord](func() (*operation.ControllerRecord, error) {
		return controllerService.controllerRecordOperation.GetControllerRecord(req.TargetRecord, req.TargetUid)
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
				NewValue: operation.ValueNotAvailable,
			},
		),
	})

	data := ResponseDeleteControllerRecord(true)
	return NewApiResponse(SuccessDeleteControllerRecord, &data)
}
