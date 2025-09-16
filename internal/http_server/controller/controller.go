// Package controller
package controller

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/labstack/echo/v4"
)

type ATCControllerInterface interface {
	GetControllers(ctx echo.Context) error
	GetCurrentControllerRecord(ctx echo.Context) error
	GetControllerRecord(ctx echo.Context) error
	UpdateControllerRating(ctx echo.Context) error
	SetControllerUnderMonitor(ctx echo.Context) error
	UnsetControllerUnderMonitor(ctx echo.Context) error
	SetControllerUnderSolo(ctx echo.Context) error
	UnsetControllerUnderSolo(ctx echo.Context) error
	SetControllerGuest(ctx echo.Context) error
	UnsetControllerGuest(ctx echo.Context) error
	AddControllerRecord(ctx echo.Context) error
	DeleteControllerRecord(ctx echo.Context) error
}

type ATCController struct {
	logger            log.LoggerInterface
	controllerService ControllerServiceInterface
}

func NewATCController(
	logger log.LoggerInterface,
	controllerService ControllerServiceInterface,
) *ATCController {
	return &ATCController{
		logger:            logger,
		controllerService: controllerService,
	}
}

func (controller *ATCController) GetControllers(ctx echo.Context) error {
	data := &RequestControllerList{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ATCController.GetControllers bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Cid = claim.Cid
	return controller.controllerService.GetControllerList(data).Response(ctx)
}

func (controller *ATCController) UpdateControllerRating(ctx echo.Context) error {
	data := &RequestUpdateControllerRating{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ATCController.UpdateControllerRating bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	return controller.controllerService.UpdateControllerRating(data).Response(ctx)
}

func (controller *ATCController) GetCurrentControllerRecord(ctx echo.Context) error {
	data := &RequestGetCurrentControllerRecord{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ATCController.GetCurrentControllerRecord bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	return controller.controllerService.GetCurrentControllerRecord(data).Response(ctx)
}

func (controller *ATCController) GetControllerRecord(ctx echo.Context) error {
	data := &RequestGetControllerRecord{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ATCController.GetControllerRecord bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	return controller.controllerService.GetControllerRecord(data).Response(ctx)
}

func (controller *ATCController) SetControllerUnderMonitor(ctx echo.Context) error {
	data := &RequestUpdateControllerUnderMonitor{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ATCController.SetControllerUnderMonitor bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	data.UnderMonitor = true
	return controller.controllerService.UpdateControllerUnderMonitor(data).Response(ctx)
}

func (controller *ATCController) UnsetControllerUnderMonitor(ctx echo.Context) error {
	data := &RequestUpdateControllerUnderMonitor{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ATCController.UnsetControllerUnderMonitor bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	data.UnderMonitor = false
	return controller.controllerService.UpdateControllerUnderMonitor(data).Response(ctx)
}

func (controller *ATCController) SetControllerUnderSolo(ctx echo.Context) error {
	data := &RequestUpdateControllerUnderSolo{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ATCController.SetControllerUnderSolo bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	if data.EndTime.IsZero() {
		controller.logger.Error("Missing required data `end_time`")
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	data.Solo = true
	return controller.controllerService.UpdateControllerUnderSolo(data).Response(ctx)
}

func (controller *ATCController) UnsetControllerUnderSolo(ctx echo.Context) error {
	data := &RequestUpdateControllerUnderSolo{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ATCController.UnsetControllerUnderSolo bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	data.Solo = false
	return controller.controllerService.UpdateControllerUnderSolo(data).Response(ctx)
}

func (controller *ATCController) SetControllerGuest(ctx echo.Context) error {
	data := &RequestUpdateControllerGuest{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ATCController.SetControllerGuest bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	data.Guest = true
	return controller.controllerService.UpdateControllerGuest(data).Response(ctx)
}

func (controller *ATCController) UnsetControllerGuest(ctx echo.Context) error {
	data := &RequestUpdateControllerGuest{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ATCController.UnsetControllerGuest bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	data.Guest = false
	return controller.controllerService.UpdateControllerGuest(data).Response(ctx)
}

func (controller *ATCController) AddControllerRecord(ctx echo.Context) error {
	data := &RequestAddControllerRecord{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ATCController.AddControllerRecord bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	return controller.controllerService.AddControllerRecord(data).Response(ctx)
}

func (controller *ATCController) DeleteControllerRecord(ctx echo.Context) error {
	data := &RequestDeleteControllerRecord{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ATCController.DeleteControllerRecord bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	return controller.controllerService.DeleteControllerRecord(data).Response(ctx)
}
