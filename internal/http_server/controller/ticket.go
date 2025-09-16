// Package controller
package controller

import (
	"github.com/golang-jwt/jwt/v5"
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
		logger:        logger,
		ticketService: ticketService,
	}
}

func (ticketController *TicketController) GetTickets(ctx echo.Context) error {
	data := &RequestGetTicket{}
	if err := ctx.Bind(data); err != nil {
		ticketController.logger.ErrorF("TicketController.GetTickets bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	return ticketController.ticketService.GetTickets(data).Response(ctx)
}

func (ticketController *TicketController) GetUserTickets(ctx echo.Context) error {
	data := &RequestGetUserTicket{}
	if err := ctx.Bind(data); err != nil {
		ticketController.logger.ErrorF("TicketController.GetUserTickets bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	return ticketController.ticketService.GetUserTicket(data).Response(ctx)
}

func (ticketController *TicketController) CreateTicket(ctx echo.Context) error {
	data := &RequestCreateTicket{}
	if err := ctx.Bind(data); err != nil {
		ticketController.logger.ErrorF("TicketController.CreateTicket bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	return ticketController.ticketService.CreateTicket(data).Response(ctx)
}

func (ticketController *TicketController) CloseTicket(ctx echo.Context) error {
	data := &RequestCloseTicket{}
	if err := ctx.Bind(data); err != nil {
		ticketController.logger.ErrorF("TicketController.CloseTicket bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	return ticketController.ticketService.CloseTicket(data).Response(ctx)
}

func (ticketController *TicketController) DeleteTicket(ctx echo.Context) error {
	data := &RequestDeleteTicket{}
	if err := ctx.Bind(data); err != nil {
		ticketController.logger.ErrorF("TicketController.DeleteTicket bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	return ticketController.ticketService.DeleteTicket(data).Response(ctx)
}
