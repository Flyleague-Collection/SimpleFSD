// Package controller
package controller

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/labstack/echo/v4"
)

type TicketControllerInterface interface {
	GetTickets(ctx echo.Context) error
	GetUserTickets(ctx echo.Context) error
	CreateTicket(ctx echo.Context) error
	CloseTicket(ctx echo.Context) error
	DeleteTicket(ctx echo.Context) error
}

type TicketController struct {
	logger        log.LoggerInterface
	ticketService TicketServiceInterface
}

func NewTicketController(
	logger log.LoggerInterface,
	ticketService TicketServiceInterface,
) *TicketController {
	return &TicketController{
		logger:        log.NewLoggerAdapter(logger, "TicketController"),
		ticketService: ticketService,
	}
}

func (controller *TicketController) GetTickets(ctx echo.Context) error {
	data := &RequestGetTickets{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetTickets bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("GetTickets jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.ticketService.GetTickets(data).Response(ctx)
}

func (controller *TicketController) GetUserTickets(ctx echo.Context) error {
	data := &RequestGetUserTickets{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetUserTickets bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfo(data, ctx); err != nil {
		controller.logger.ErrorF("GetUserTickets jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.ticketService.GetUserTickets(data).Response(ctx)
}

func (controller *TicketController) CreateTicket(ctx echo.Context) error {
	data := &RequestCreateTicket{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("CreateTicket bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("CreateTicket jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.ticketService.CreateTicket(data).Response(ctx)
}

func (controller *TicketController) CloseTicket(ctx echo.Context) error {
	data := &RequestCloseTicket{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("CloseTicket bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("CloseTicket jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.ticketService.CloseTicket(data).Response(ctx)
}

func (controller *TicketController) DeleteTicket(ctx echo.Context) error {
	data := &RequestDeleteTicket{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("DeleteTicket bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("DeleteTicket jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.ticketService.DeleteTicket(data).Response(ctx)
}
