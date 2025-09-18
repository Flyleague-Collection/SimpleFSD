// Package controller
package controller

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/labstack/echo/v4"
)

type FlightPlanControllerInterface interface {
	SubmitFlightPlan(ctx echo.Context) error
	GetFlightPlan(ctx echo.Context) error
	GetFlightPlans(ctx echo.Context) error
	DeleteFlightPlan(ctx echo.Context) error
	LockFlightPlan(ctx echo.Context) error
	UnlockFlightPlan(ctx echo.Context) error
}

type FlightPlanController struct {
	logger            log.LoggerInterface
	flightPlanService FlightPlanServiceInterface
}

func NewFlightPlanController(
	logger log.LoggerInterface,
	flightPlanService FlightPlanServiceInterface,
) *FlightPlanController {
	return &FlightPlanController{
		logger:            log.NewLoggerAdapter(logger, "FlightPlanController"),
		flightPlanService: flightPlanService,
	}
}

func (controller *FlightPlanController) SubmitFlightPlan(ctx echo.Context) error {
	data := &RequestSubmitFlightPlan{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("SubmitFlightPlan bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.JwtHeader.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	return controller.flightPlanService.SubmitFlightPlan(data).Response(ctx)
}

func (controller *FlightPlanController) GetFlightPlan(ctx echo.Context) error {
	data := &RequestGetFlightPlan{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetFlightPlan bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.JwtHeader.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	return controller.flightPlanService.GetFlightPlan(data).Response(ctx)
}

func (controller *FlightPlanController) GetFlightPlans(ctx echo.Context) error {
	data := &RequestGetFlightPlans{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("GetFlightPlans bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.JwtHeader.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	return controller.flightPlanService.GetFlightPlans(data).Response(ctx)
}

func (controller *FlightPlanController) DeleteFlightPlan(ctx echo.Context) error {
	data := &RequestDeleteFlightPlan{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("DeleteFlightPlan bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.JwtHeader.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	return controller.flightPlanService.DeleteFlightPlan(data).Response(ctx)
}

func (controller *FlightPlanController) LockFlightPlan(ctx echo.Context) error {
	data := &RequestLockFlightPlan{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("LockFlightPlan bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.JwtHeader.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	data.Lock = true
	return controller.flightPlanService.LockFlightPlan(data).Response(ctx)
}

func (controller *FlightPlanController) UnlockFlightPlan(ctx echo.Context) error {
	data := &RequestLockFlightPlan{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("UnlockFlightPlan bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.JwtHeader.Cid = claim.Cid
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	data.Lock = false
	return controller.flightPlanService.LockFlightPlan(data).Response(ctx)
}
