// Package service
package service

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"io"
	"net/http"
	"strings"
	"time"
)

type CachedMetar struct {
	Response       *ApiResponse[ResponseQueryMetar]
	ExpirationTime time.Time
}

type MetarService struct {
	logger     log.LoggerInterface
	metarCache map[string]*CachedMetar
}

func NewMetarService(logger log.LoggerInterface) *MetarService {
	return &MetarService{
		logger:     log.NewLoggerAdapter(logger, "MetarService"),
		metarCache: make(map[string]*CachedMetar),
	}
}

var (
	ErrConnectionFail = NewApiStatus("CONNECTION_FAIL", "查询Metar失败", ServerInternalError)
	ErrQueryMetarFail = NewApiStatus("QUERY_METAR_FAIL", "查询Metar失败", ServerInternalError)
	ErrMetarNotFound  = NewApiStatus("METAR_NOT_FOUND", "未找到Metar信息", NotFound)
	SuccessGetMetar   = NewApiStatus("GET_METAR", "成功获取Metar", Ok)
)

func (metarService *MetarService) cacheMetar(icao string, response *ApiResponse[ResponseQueryMetar]) *ApiResponse[ResponseQueryMetar] {
	currentTime := time.Now()
	minute := currentTime.Minute()
	var addMinutes int
	if minute < 30 {
		addMinutes = 30 - minute
	} else {
		addMinutes = 60 - minute
	}
	next := currentTime.Add(time.Duration(addMinutes) * time.Minute)
	expirationTime := time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())

	metar := &CachedMetar{Response: response, ExpirationTime: expirationTime}
	metarService.metarCache[icao] = metar

	return response
}

func (metarService *MetarService) QueryMetar(req *RequestQueryMetar) *ApiResponse[ResponseQueryMetar] {
	if req.ICAO == "" || len(req.ICAO) != 4 {
		return NewApiResponse[ResponseQueryMetar](ErrIllegalParam, nil)
	}

	if metar, ok := metarService.metarCache[req.ICAO]; ok {
		if metar.ExpirationTime.After(time.Now()) {
			return metar.Response
		}
		delete(metarService.metarCache, req.ICAO)
	}

	url := fmt.Sprintf("https://aviationweather.gov/api/data/metar?ids=%s", req.ICAO)
	metarService.logger.InfoF("Get metar from url %s", url)
	response, err := http.Get(url)
	if err != nil {
		return NewApiResponse[ResponseQueryMetar](ErrConnectionFail, nil)
	}

	if response.StatusCode == 204 {
		url := fmt.Sprintf("https://xmairavt7.xiamenair.com/WarningPage/AirportReports?arp4code=%s/1", req.ICAO)
		metarService.logger.InfoF("Get metar from url %s", url)
		response, err = http.Get(url)
		if err != nil {
			return NewApiResponse[ResponseQueryMetar](ErrConnectionFail, nil)
		}

		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			return NewApiResponse[ResponseQueryMetar](ErrQueryMetarFail, nil)
		}

		metars := strings.Split(doc.Find("pre").Get(0).LastChild.Data, "\n")

		if len(metars) == 0 {
			return metarService.cacheMetar(req.ICAO, NewApiResponse[ResponseQueryMetar](ErrMetarNotFound, nil))
		}

		for _, metar := range metars {
			if strings.HasPrefix(metar, "METAR") {
				data := ResponseQueryMetar(metar)
				return metarService.cacheMetar(req.ICAO, NewApiResponse(SuccessGetMetar, &data))
			}
		}

		return metarService.cacheMetar(req.ICAO, NewApiResponse[ResponseQueryMetar](ErrMetarNotFound, nil))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return NewApiResponse[ResponseQueryMetar](ErrQueryMetarFail, nil)
	}
	_ = response.Body.Close()

	data := ResponseQueryMetar(body)
	return metarService.cacheMetar(req.ICAO, NewApiResponse(SuccessGetMetar, &data))
}
