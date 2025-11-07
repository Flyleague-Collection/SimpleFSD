// Package controller
package controller

import (
	. "github.com/half-nothing/simple-fsd/src/interfaces/http/service"
	"github.com/half-nothing/simple-fsd/src/interfaces/log"
	"github.com/labstack/echo/v4"
)

type AuditLogControllerInterface interface {
	GetAuditLogs(ctx echo.Context) error
	LogUnlawfulOverreach(ctx echo.Context) error
}

type AuditLogController struct {
	logger       log.LoggerInterface
	auditService AuditServiceInterface
}

func NewAuditLogController(logger log.LoggerInterface, auditService AuditServiceInterface) *AuditLogController {
	return &AuditLogController{
		logger:       log.NewLoggerAdapter(logger, "AuditLogController"),
		auditService: auditService,
	}
}

func (controller *AuditLogController) GetAuditLogs(ctx echo.Context) error {
	data := &RequestGetAuditLog{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetAuditLogs bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("GetAuditLogs jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.auditService.GetAuditLogPage(data).Response(ctx)
}

func (controller *AuditLogController) LogUnlawfulOverreach(ctx echo.Context) error {
	data := &RequestLogUnlawfulOverreach{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("LogUnlawfulOverreach bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("LogUnlawfulOverreach jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.auditService.LogUnlawfulOverreach(data).Response(ctx)
}
