// Package config
package config

import (
	"errors"

	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
)

type AirportData struct {
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	Alt          float64 `json:"alt"`
	AirportRange int     `json:"airport_range"`
}

func (config *AirportData) checkValid(_ log.LoggerInterface) *ValidResult {
	if config.AirportRange <= 0 {
		return ValidFail(errors.New("airport_range must be greater than zero"))
	}
	return ValidPass()
}
