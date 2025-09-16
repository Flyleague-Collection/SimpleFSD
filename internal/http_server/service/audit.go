// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
)

type AudioLogService struct {
	logger         log.LoggerInterface
	auditOperation operation.AuditLogOperationInterface
}

func NewAuditService(
	logger log.LoggerInterface,
	auditOperation operation.AuditLogOperationInterface,
) *AudioLogService {
	return &AudioLogService{
		logger:         logger,
		auditOperation: auditOperation,
	}
}

func (auditLogService *AudioLogService) HandleAuditLogMessage(message *queue.Message) error {
	if val, ok := message.Data.(*operation.AuditLog); ok {
		return auditLogService.auditOperation.SaveAuditLog(val)
	}
	return queue.ErrMessageDataType
}

func (auditLogService *AudioLogService) HandleAuditLogsMessage(message *queue.Message) error {
	if val, ok := message.Data.([]*operation.AuditLog); ok {
		return auditLogService.auditOperation.SaveAuditLogs(val)
	}
	return queue.ErrMessageDataType
}

var SuccessGetAuditLog = ApiStatus{StatusName: "GET_AUDIT_LOG", Description: "成功获取审计日志", HttpCode: Ok}

func (auditLogService *AudioLogService) GetAuditLogPage(req *RequestGetAuditLog) *ApiResponse[ResponseGetAuditLog] {
	if req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseGetAuditLog](ErrIllegalParam, nil)
	}
	if req.Permission <= 0 {
		return NewApiResponse[ResponseGetAuditLog](ErrNoPermission, nil)
	}
	permission := operation.Permission(req.Permission)
	if !permission.HasPermission(operation.AuditLogShow) {
		return NewApiResponse[ResponseGetAuditLog](ErrNoPermission, nil)
	}
	auditLogs, total, err := auditLogService.auditOperation.GetAuditLogs(req.Page, req.PageSize)
	if err != nil {
		return NewApiResponse[ResponseGetAuditLog](ErrDatabaseFail, nil)
	}
	return NewApiResponse(&SuccessGetAuditLog, &ResponseGetAuditLog{
		Items:    auditLogs,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	})
}

var SuccessLogUnlawfulOverreach = ApiStatus{StatusName: "LOG_UNLAWFUL_OVERREACH", Description: "成功记录非法访问", HttpCode: Ok}

func (auditLogService *AudioLogService) LogUnlawfulOverreach(req *RequestLogUnlawfulOverreach) *ApiResponse[ResponseLogUnlawfulOverreach] {
	auditLog := auditLogService.auditOperation.NewAuditLog(operation.UnlawfulOverreach, req.Cid, req.AccessPath,
		req.Ip, req.UserAgent, nil)
	err := auditLogService.auditOperation.SaveAuditLog(auditLog)
	if err != nil {
		auditLogService.logger.ErrorF("Fail to create audit log for unlawful_overreach, detail: %v", err)
	}
	data := ResponseLogUnlawfulOverreach(true)
	return NewApiResponse(&SuccessLogUnlawfulOverreach, &data)
}
