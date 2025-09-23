// Package interfaces
package interfaces

import (
	"errors"
)

var (
	ErrICAOInvalid   = errors.New("invalid ICAO value")
	ErrMetarNotFound = errors.New("metar not found")
)

type MetarManagerInterface interface {
	QueryMetar(icao string) (metar string, err error)
	QueryMetars(icaos []string) (metars []string)
}

type MetarGetterInterface interface {
	GetMetar(icao string) (string, error)
}
