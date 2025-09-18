// Package controller
package controller

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/labstack/echo/v4"
)

type FileControllerInterface interface {
	UploadImages(ctx echo.Context) error
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

func (controller *FileController) UploadImages(ctx echo.Context) error {
	if file, err := ctx.FormFile("file"); err != nil {
		controller.logger.ErrorF("UploadImages form file error: %v", err)
		return NewErrorResponse(ctx, ErrParseParam)
	} else {
		data := &RequestUploadFile{File: file}
		token := ctx.Get("user").(*jwt.Token)
		claim := token.Claims.(*Claims)
		data.Cid = claim.Cid
		data.Uid = claim.Uid
		data.Permission = claim.Permission
		data.Ip = ctx.RealIP()
		data.UserAgent = ctx.Request().UserAgent()
		return controller.storeService.SaveUploadImages(data).Response(ctx)
	}
}
