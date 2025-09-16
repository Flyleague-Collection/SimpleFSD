// Package controller
package controller

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ClientControllerInterface interface {
	GetOnlineClients(ctx echo.Context) error
	GetClientPath(ctx echo.Context) error
	SendMessageToClient(ctx echo.Context) error
	KillClient(ctx echo.Context) error
}

type ClientController struct {
	logger        log.LoggerInterface
	clientService ClientServiceInterface
}

func NewClientController(logger log.LoggerInterface, clientService ClientServiceInterface) *ClientController {
	return &ClientController{
		logger:        logger,
		clientService: clientService,
	}
}

func (controller *ClientController) GetOnlineClients(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, controller.clientService.GetOnlineClient())
}

func (controller *ClientController) GetClientPath(ctx echo.Context) error {
	data := &RequestClientPath{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ClientController.GetClientPath bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	return controller.clientService.GetClientPath(data).Response(ctx)
}

func (controller *ClientController) SendMessageToClient(ctx echo.Context) error {
	data := &RequestSendMessageToClient{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ClientController.sendMessageToClient bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Cid = claim.Cid
	data.Permission = claim.Permission
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	return controller.clientService.SendMessageToClient(data).Response(ctx)
}

func (controller *ClientController) KillClient(ctx echo.Context) error {
	data := &RequestKillClient{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("ClientController.killClient bind error: %v", err)
		return NewErrorResponse(ctx, ErrLackParam)
	}
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*Claims)
	data.Uid = claim.Uid
	data.Permission = claim.Permission
	data.Cid = claim.Cid
	data.Ip = ctx.RealIP()
	data.UserAgent = ctx.Request().UserAgent()
	return controller.clientService.KillClient(data).Response(ctx)
}
