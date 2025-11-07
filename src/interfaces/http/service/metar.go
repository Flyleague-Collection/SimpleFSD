// Package service
package service

var (
	ErrMetarNotFound = NewApiStatus("METAR_NOT_FOUND", "未找到Metar信息", NotFound)
	SuccessGetMetar  = NewApiStatus("GET_METAR", "成功获取Metar", Ok)
)

type MetarServiceInterface interface {
	QueryMetar(req *RequestQueryMetar) *ApiResponse[ResponseQueryMetar]
}

type RequestQueryMetar struct {
	ICAO string `query:"icao"`
	Raw  bool   `query:"raw"`
}

type ResponseQueryMetar []string
