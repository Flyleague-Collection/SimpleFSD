// Package controller
package controller

import (
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
		logger:          log.NewLoggerAdapter(logger, "ActivityController"),
		activityService: activityService,
	}
}

func (controller *ActivityController) GetActivities(ctx echo.Context) error {
	data := &RequestGetActivities{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetActivities bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.activityService.GetActivities(data).Response(ctx)
}

func (controller *ActivityController) GetActivitiesPage(ctx echo.Context) error {
	data := &RequestGetActivitiesPage{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetActivitiesPage bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("GetActivitiesPage jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.activityService.GetActivitiesPage(data).Response(ctx)
}

func (controller *ActivityController) GetActivityInfo(ctx echo.Context) error {
	data := &RequestActivityInfo{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetActivityInfo bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.activityService.GetActivityInfo(data).Response(ctx)
}

func (controller *ActivityController) AddActivity(ctx echo.Context) error {
	data := &RequestAddActivity{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("AddActivity bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("AddActivity jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.activityService.AddActivity(data).Response(ctx)
}

func (controller *ActivityController) DeleteActivity(ctx echo.Context) error {
	data := &RequestDeleteActivity{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("DeleteActivity bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("DeleteActivity jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.activityService.DeleteActivity(data).Response(ctx)
}

func (controller *ActivityController) ControllerJoin(ctx echo.Context) error {
	data := &RequestControllerJoin{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ControllerJoin bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("ControllerJoin jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.activityService.ControllerJoin(data).Response(ctx)
}

func (controller *ActivityController) ControllerLeave(ctx echo.Context) error {
	data := &RequestControllerLeave{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ControllerLeave bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("ControllerLeave jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.activityService.ControllerLeave(data).Response(ctx)
}

func (controller *ActivityController) PilotJoin(ctx echo.Context) error {
	data := &RequestPilotJoin{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("PilotJoin bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("PilotJoin jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.activityService.PilotJoin(data).Response(ctx)
}

func (controller *ActivityController) PilotLeave(ctx echo.Context) error {
	data := &RequestPilotLeave{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("PilotLeave bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("PilotLeave jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.activityService.PilotLeave(data).Response(ctx)
}

func (controller *ActivityController) EditActivity(ctx echo.Context) error {
	data := &RequestEditActivity{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("EditActivity bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("EditActivity jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.activityService.EditActivity(data).Response(ctx)
}

func (controller *ActivityController) EditActivityStatus(ctx echo.Context) error {
	data := &RequestEditActivityStatus{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("EditActivityStatus bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("EditActivityStatus jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.activityService.EditActivityStatus(data).Response(ctx)
}

func (controller *ActivityController) EditPilotStatus(ctx echo.Context) error {
	data := &RequestEditPilotStatus{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("EditPilotStatus bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("EditPilotStatus jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.activityService.EditPilotStatus(data).Response(ctx)
}
