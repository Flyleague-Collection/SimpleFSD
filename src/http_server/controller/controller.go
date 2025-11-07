// Package controller
package controller

import (
	. "github.com/half-nothing/simple-fsd/src/interfaces/http/service"
	"github.com/half-nothing/simple-fsd/src/interfaces/log"
	"github.com/labstack/echo/v4"
)

type ATCControllerInterface interface {
	GetControllers(ctx echo.Context) error
	GetControllerRatings(ctx echo.Context) error
	GetCurrentControllerRecord(ctx echo.Context) error
	GetControllerRecord(ctx echo.Context) error
	UpdateControllerRating(ctx echo.Context) error
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
		logger:            log.NewLoggerAdapter(logger, "ATCController"),
		controllerService: controllerService,
	}
}

func (controller *ATCController) GetControllers(ctx echo.Context) error {
	data := &RequestControllerList{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetControllers bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("GetControllers jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.controllerService.GetControllerList(data).Response(ctx)
}

func (controller *ATCController) GetControllerRatings(ctx echo.Context) error {
	data := &RequestControllerRatingList{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetControllerRatings bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.controllerService.GetControllerRatings(data).Response(ctx)
}

func (controller *ATCController) UpdateControllerRating(ctx echo.Context) error {
	data := &RequestUpdateControllerRating{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("UpdateControllerRating bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("UpdateControllerRating jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.controllerService.UpdateControllerRating(data).Response(ctx)
}

func (controller *ATCController) GetCurrentControllerRecord(ctx echo.Context) error {
	data := &RequestGetCurrentControllerRecord{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetCurrentControllerRecord bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("GetCurrentControllerRecord jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.controllerService.GetCurrentControllerRecord(data).Response(ctx)
}

func (controller *ATCController) GetControllerRecord(ctx echo.Context) error {
	data := &RequestGetControllerRecord{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetControllerRecord bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("GetControllerRecord jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.controllerService.GetControllerRecord(data).Response(ctx)
}

func (controller *ATCController) AddControllerRecord(ctx echo.Context) error {
	data := &RequestAddControllerRecord{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("AddControllerRecord bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("AddControllerRecord jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.controllerService.AddControllerRecord(data).Response(ctx)
}

func (controller *ATCController) DeleteControllerRecord(ctx echo.Context) error {
	data := &RequestDeleteControllerRecord{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("DeleteControllerRecord bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("DeleteControllerRecord jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.controllerService.DeleteControllerRecord(data).Response(ctx)
}
