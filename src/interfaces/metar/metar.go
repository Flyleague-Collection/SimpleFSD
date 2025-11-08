// Package metar
package metar

import "errors"

var (
	ErrICAOInvalid   = errors.New("invalid ICAO value")
	ErrMetarNotFound = errors.New("metar not found")
)

type ManagerInterface interface {
	QueryMetar(icao string) (metar string, err error)
	QueryMetars(icaos []string) (metars []string)
}

type GetterInterface interface {
	GetMetar(icao string) (string, error)
}
