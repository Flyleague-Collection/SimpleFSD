// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/half-nothing/simple-fsd/internal/utils"
	"strings"
)

type MetarService struct {
	logger       log.LoggerInterface
	metarManager interfaces.MetarManagerInterface
}

func NewMetarService(
	logger log.LoggerInterface,
	metarManager interfaces.MetarManagerInterface,
) *MetarService {
	return &MetarService{
		logger:       log.NewLoggerAdapter(logger, "MetarService"),
		metarManager: metarManager,
	}
}

var (
	ErrMetarNotFound = NewApiStatus("METAR_NOT_FOUND", "未找到Metar信息", NotFound)
	SuccessGetMetar  = NewApiStatus("GET_METAR", "成功获取Metar", Ok)
)

func (metarService *MetarService) QueryMetar(req *RequestQueryMetar) *ApiResponse[ResponseQueryMetar] {
	if req.ICAO == "" {
		return NewApiResponse[ResponseQueryMetar](ErrIllegalParam, nil)
	}

	icaos := strings.Split(req.ICAO, ",")
	if len(icaos) < 1 {
		return NewApiResponse[ResponseQueryMetar](ErrIllegalParam, nil)
	}

	utils.Map[string](icaos, func(element *string) { *element = strings.TrimSpace(*element) })
	metars := metarService.metarManager.QueryMetars(icaos)
	if len(metars) < 1 {
		return NewApiResponse[ResponseQueryMetar](ErrMetarNotFound, nil)
	}
	data := ResponseQueryMetar(metars)
	return NewApiResponse(SuccessGetMetar, &data)
}
