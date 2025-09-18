// Package service
package service

type MetarServiceInterface interface {
	QueryMetar(req *RequestQueryMetar) *ApiResponse[ResponseQueryMetar]
}

type RequestQueryMetar struct {
	ICAO string `query:"icao"`
}

type ResponseQueryMetar string
