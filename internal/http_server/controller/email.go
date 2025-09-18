// Package controller
package controller

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/labstack/echo/v4"
)

type EmailControllerInterface interface {
	SendVerifyEmail(ctx echo.Context) error
}

type EmailController struct {
	logger       log.LoggerInterface
	emailService EmailServiceInterface
}

func NewEmailController(logger log.LoggerInterface, emailService EmailServiceInterface) *EmailController {
	return &EmailController{
		logger:       log.NewLoggerAdapter(logger, "EmailController"),
		emailService: emailService,
	}
}

func (controller *EmailController) SendVerifyEmail(ctx echo.Context) error {
	data := &RequestEmailVerifyCode{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("SendVerifyEmail bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.emailService.SendEmailVerifyCode(data).Response(ctx)
}
