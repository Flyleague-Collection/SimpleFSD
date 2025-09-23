// Package service
package service

type MetarServiceInterface interface {
	QueryMetar(req *RequestQueryMetar) *ApiResponse[ResponseQueryMetar]
}

type RequestQueryMetar struct {
	ICAO string `query:"icao"`
	Raw  bool   `query:"raw"`
}

type ResponseQueryMetar []string
