// Package controller
package controller

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/labstack/echo/v4"
)

type ActivityControllerInterface interface {
	GetActivities(ctx echo.Context) error
	GetActivitiesPage(ctx echo.Context) error
	GetActivityInfo(ctx echo.Context) error
	AddActivity(ctx echo.Context) error
	DeleteActivity(ctx echo.Context) error
	ControllerJoin(ctx echo.Context) error
	ControllerLeave(ctx echo.Context) error
	PilotJoin(ctx echo.Context) error
	PilotLeave(ctx echo.Context) error
	EditActivity(ctx echo.Context) error
	EditActivityStatus(ctx echo.Context) error
	EditPilotStatus(ctx echo.Context) error
}

type ActivityController struct {
	logger          log.LoggerInterface
	activityService ActivityServiceInterface
}

func NewActivityController(logger log.LoggerInterface, activityService ActivityServiceInterface) *ActivityController {
	return &ActivityController{
		logger:          logger,
		activityService: activityService,
	}
}

func (controller *ActivityController) GetActivities(ctx echo.Context) error {
	data := &RequestGetActivities{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ActivityController.GetActivities bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	return controller.activityService.GetActivities(data).Response(ctx)
}

func (controller *ActivityController) GetActivitiesPage(ctx echo.Context) error {
	data := &RequestGetActivitiesPage{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ActivityController.GetActivitiesPage bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	return controller.activityService.GetActivitiesPage(data).Response(ctx)
}

func (controller *ActivityController) GetActivityInfo(ctx echo.Context) error {
	data := &RequestActivityInfo{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ActivityController.GetActivityInfo bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	return controller.activityService.GetActivityInfo(data).Response(ctx)
}

func (controller *ActivityController) AddActivity(ctx echo.Context) error {
	data := &RequestAddActivity{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ActivityController.AddActivity bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Cid = claim.Cid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	return controller.activityService.AddActivity(data).Response(ctx)
}

func (controller *ActivityController) DeleteActivity(ctx echo.Context) error {
	data := &RequestDeleteActivity{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ActivityController.DeleteActivity bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Cid = claim.Cid
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	return controller.activityService.DeleteActivity(data).Response(ctx)
}

func (controller *ActivityController) ControllerJoin(ctx echo.Context) error {
	data := &RequestControllerJoin{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ActivityController.ControllerJoin bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Rating = claim.Rating
	return controller.activityService.ControllerJoin(data).Response(ctx)
}

func (controller *ActivityController) ControllerLeave(ctx echo.Context) error {
	data := &RequestControllerLeave{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ActivityController.ControllerLeave bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Cid = claim.Cid
	return controller.activityService.ControllerLeave(data).Response(ctx)
}

func (controller *ActivityController) PilotJoin(ctx echo.Context) error {
	data := &RequestPilotJoin{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ActivityController.PilotJoin bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Cid = claim.Cid
	return controller.activityService.PilotJoin(data).Response(ctx)
}

func (controller *ActivityController) PilotLeave(ctx echo.Context) error {
	data := &RequestPilotLeave{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ActivityController.PilotLeave bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Cid = claim.Cid
	return controller.activityService.PilotLeave(data).Response(ctx)
}

func (controller *ActivityController) EditActivity(ctx echo.Context) error {
	data := &RequestEditActivity{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ActivityController.EditActivity bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Cid = claim.Cid
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	return controller.activityService.EditActivity(data).Response(ctx)
}

func (controller *ActivityController) EditActivityStatus(ctx echo.Context) error {
	data := &RequestEditActivityStatus{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ActivityController.EditActivityStatus bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	return controller.activityService.EditActivityStatus(data).Response(ctx)
}

func (controller *ActivityController) EditPilotStatus(ctx echo.Context) error {
	data := &RequestEditPilotStatus{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ActivityController.EditPilotStatus bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Cid = claim.Cid
	return controller.activityService.EditPilotStatus(data).Response(ctx)
}
