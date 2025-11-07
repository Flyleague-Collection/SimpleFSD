// Package controller
package controller

import (
	. "github.com/half-nothing/simple-fsd/src/interfaces/http/service"
	"github.com/half-nothing/simple-fsd/src/interfaces/log"
	"github.com/labstack/echo/v4"
)

type FileControllerInterface interface {
	UploadImage(ctx echo.Context) error
	UploadFile(ctx echo.Context) error
}

type FileController struct {
	logger       log.LoggerInterface
	storeService StoreServiceInterface
}

func NewFileController(logger log.LoggerInterface, storeService StoreServiceInterface) *FileController {
	return &FileController{
		logger:       log.NewLoggerAdapter(logger, "FileController"),
		storeService: storeService,
	}
}

func (controller *FileController) UploadImage(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		controller.logger.ErrorF("UploadImage form file error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	data := &RequestUploadImage{File: file}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("UploadImage jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.storeService.SaveUploadImage(data).Response(ctx)
}

func (controller *FileController) UploadFile(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		controller.logger.ErrorF("UploadFile form file error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	data := &RequestUploadFile{File: file}
	if err := SetJwtInfoAndEchoContent(data, ctx); err != nil {
		controller.logger.ErrorF("UploadFile jwt token parse error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	}
	return controller.storeService.SaveUploadFile(data).Response(ctx)
}
