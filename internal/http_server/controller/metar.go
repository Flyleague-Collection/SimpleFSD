// Package controller
package controller

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/labstack/echo/v4"
)

type MetarControllerInterface interface {
	QueryMetar(ctx echo.Context) error
}

type MetarController struct {
	logger       log.LoggerInterface
	metarService MetarServiceInterface
}

func NewMetarServiceController(
	logger log.LoggerInterface,
	metarService MetarServiceInterface,
) *MetarController {
	return &MetarController{
		logger:       log.NewLoggerAdapter(logger, "MetarController"),
		metarService: metarService,
	}
}

func (controller *MetarController) QueryMetar(ctx echo.Context) error {
	data := &RequestQueryMetar{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.ErrorF("QueryMetar bind error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.metarService.QueryMetar(data).Response(ctx)
}
