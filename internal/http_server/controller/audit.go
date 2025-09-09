// Package controller
package controller

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
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
		logger:       logger,
		auditService: auditService,
	}
}

func (controller *AuditLogController) GetAuditLogs(ctx echo.Context) error {
	data := &RequestGetAuditLog{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("AuditLogController.GetAuditLogs bind error: %v", err)
		return NewErrorResponse(ctx, &ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	return controller.auditService.GetAuditLogPage(data).Response(ctx)
}

func (controller *AuditLogController) LogUnlawfulOverreach(ctx echo.Context) error {
	data := &RequestLogUnlawfulOverreach{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("AuditLogController.LogUnlawfulOverreach bind error: %v", err)
		return NewErrorResponse(ctx, &ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Cid = claim.Cid
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	return controller.auditService.LogUnlawfulOverreach(data).Response(ctx)
}
