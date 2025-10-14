// Package controller
package controller

import (
	. "github.com/half-nothing/simple-fsd/internal/interfaces/http/service"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/labstack/echo/v4"
)

type ControllerApplicationControllerInterface interface {
	GetSelfApplication(ctx echo.Context) error
	GetApplications(ctx echo.Context) error
	SubmitApplication(ctx echo.Context) error
	CancelSelfApplication(ctx echo.Context) error
	UpdateApplication(ctx echo.Context) error
}

type ControllerApplicationController struct {
	logger             log.LoggerInterface
	applicationService ControllerApplicationServiceInterface
}

func NewControllerApplicationController(
	logger log.LoggerInterface,
	applicationService ControllerApplicationServiceInterface,
) *ControllerApplicationController {
	return &ControllerApplicationController{
		logger:             log.NewLoggerAdapter(logger, "ApplicationController"),
		applicationService: applicationService,
	}
}

func (controller *ControllerApplicationController) GetSelfApplication(ctx echo.Context) error {
	data := &RequestGetSelfApplication{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetSelfApplication bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("GetSelfApplication jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.applicationService.GetSelfApplication(data).Response(ctx)
}

func (controller *ControllerApplicationController) GetApplications(ctx echo.Context) error {
	data := &RequestGetApplications{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetApplications bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("GetApplications jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.applicationService.GetApplications(data).Response(ctx)
}

func (controller *ControllerApplicationController) SubmitApplication(ctx echo.Context) error {
	data := &RequestSubmitControllerApplication{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("SubmitApplication bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("SubmitApplication jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.applicationService.SubmitControllerApplication(data).Response(ctx)
}

func (controller *ControllerApplicationController) CancelSelfApplication(ctx echo.Context) error {
	data := &RequestCancelSelfApplication{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("CancelSelfApplication bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("CancelSelfApplication jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.applicationService.CancelSelfApplication(data).Response(ctx)
}

func (controller *ControllerApplicationController) UpdateApplication(ctx echo.Context) error {
	data := &RequestUpdateApplicationStatus{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("UpdateApplication bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("UpdateApplication jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.applicationService.UpdateApplicationStatus(data).Response(ctx)
}
