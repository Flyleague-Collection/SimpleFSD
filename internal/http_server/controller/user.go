// Package controller
package controller

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/labstack/echo/v4"
)

type UserControllerInterface interface {
	UserRegister(ctx echo.Context) error
	UserLogin(ctx echo.Context) error
	CheckUserAvailability(ctx echo.Context) error
	GetCurrentUserProfile(ctx echo.Context) error
	EditCurrentProfile(ctx echo.Context) error
	GetUserProfile(ctx echo.Context) error
	EditProfile(ctx echo.Context) error
	GetUsers(ctx echo.Context) error
	EditUserPermission(ctx echo.Context) error
	GetUserHistory(ctx echo.Context) error
	GetToken(ctx echo.Context) error
}

type UserController struct {
	logger  log.LoggerInterface
	service UserServiceInterface
}

func NewUserHandler(logger log.LoggerInterface, service UserServiceInterface) *UserController {
	return &UserController{
		logger:  logger,
		service: service,
	}
}

func (controller *UserController) UserRegister(ctx echo.Context) error {
	data := &RequestUserRegister{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("UserController.UserRegister bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	return controller.service.UserRegister(data).Response(ctx)
}

func (controller *UserController) UserLogin(ctx echo.Context) error {
	data := &RequestUserLogin{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("UserController.UserLogin bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	return controller.service.UserLogin(data).Response(ctx)
}

func (controller *UserController) CheckUserAvailability(ctx echo.Context) error {
	data := &RequestUserAvailability{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("UserController.CheckUserAvailability bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	return controller.service.CheckAvailability(data).Response(ctx)
}

func (controller *UserController) GetCurrentUserProfile(ctx echo.Context) error {
	data := &RequestUserCurrentProfile{}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	return controller.service.GetCurrentProfile(data).Response(ctx)
}

func (controller *UserController) EditCurrentProfile(ctx echo.Context) error {
	data := &RequestUserEditCurrentProfile{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("UserController.EditCurrentProfile bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.ID = claim.Uid
	data.Cid = claim.Cid
	return controller.service.EditCurrentProfile(data).Response(ctx)
}

func (controller *UserController) GetUserProfile(ctx echo.Context) error {
	data := &RequestUserProfile{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("UserController.GetUserProfile bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	return controller.service.GetUserProfile(data).Response(ctx)
}

func (controller *UserController) EditProfile(ctx echo.Context) error {
	data := &RequestUserEditProfile{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("UserController.EditProfile bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.JwtHeader.Cid = claim.Cid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	return controller.service.EditUserProfile(data).Response(ctx)
}

func (controller *UserController) GetUsers(ctx echo.Context) error {
	data := &RequestUserList{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("UserController.GetUsers bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	return controller.service.GetUserList(data).Response(ctx)
}

func (controller *UserController) EditUserPermission(ctx echo.Context) error {
	data := &RequestUserEditPermission{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("UserController.EditUserPermission bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Cid = claim.Cid
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	return controller.service.EditUserPermission(data).Response(ctx)
}

func (controller *UserController) GetUserHistory(ctx echo.Context) error {
	data := &RequestGetUserHistory{}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Cid = claim.Cid
	return controller.service.GetUserHistory(data).Response(ctx)
}

func (controller *UserController) GetToken(ctx echo.Context) error {
	data := &RequestGetToken{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("UserController.GetToken bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Claims = claim
	return controller.service.GetTokenWithFlushToken(data).Response(ctx)
}
