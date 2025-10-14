// Package controller
package controller

import (
	. "github.com/half-nothing/simple-fsd/internal/interfaces/http/service"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/labstack/echo/v4"
)

type AnnouncementControllerInterface interface {
	GetAnnouncements(ctx echo.Context) error
	GetDetailAnnouncements(ctx echo.Context) error
	CreateAnnouncement(ctx echo.Context) error
	UpdateAnnouncement(ctx echo.Context) error
	DeleteAnnouncement(ctx echo.Context) error
}

type AnnouncementController struct {
	logger  log.LoggerInterface
	service AnnouncementServiceInterface
}

func NewAnnouncementController(
	logger log.LoggerInterface,
	service AnnouncementServiceInterface,
) *AnnouncementController {
	return &AnnouncementController{
		logger:  log.NewLoggerAdapter(logger, "AnnouncementController"),
		service: service,
	}
}

func (controller *AnnouncementController) GetAnnouncements(ctx echo.Context) error {
	data := &RequestGetAnnouncements{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetAnnouncements bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("GetAnnouncements jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.service.GetAnnouncements(data).Response(ctx)
}

func (controller *AnnouncementController) GetDetailAnnouncements(ctx echo.Context) error {
	data := &RequestGetDetailAnnouncements{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetDetailAnnouncements bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("GetDetailAnnouncements jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.service.GetDetailAnnouncements(data).Response(ctx)
}

func (controller *AnnouncementController) CreateAnnouncement(ctx echo.Context) error {
	data := &RequestPublishAnnouncement{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("CreateAnnouncement bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("CreateAnnouncement jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.service.PublishAnnouncement(data).Response(ctx)
}

func (controller *AnnouncementController) UpdateAnnouncement(ctx echo.Context) error {
	data := &RequestEditAnnouncement{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("UpdateAnnouncement bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("UpdateAnnouncement jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.service.EditAnnouncement(data).Response(ctx)
}

func (controller *AnnouncementController) DeleteAnnouncement(ctx echo.Context) error {
	data := &RequestDeleteAnnouncement{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("DeleteAnnouncement bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("DeleteAnnouncement jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.service.DeleteAnnouncement(data).Response(ctx)
}
