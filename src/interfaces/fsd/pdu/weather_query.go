// Package pdu
package pdu

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
)

type WeatherQuery struct {
	*Base
	Station string
}

func (c *WeatherQuery) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// ZSHA_CTR SERVER METAR ZSSS
	// [   0  ] [  1 ] [ 2 ] [ 3]
	if r := c.CheckLength(fsd.ClientCommandWeatherQuery, len(data)); r != nil {
		return nil, r
	}
	command := &WeatherQuery{
		Base:    NewBase(fsd.ClientCommandWeatherQuery, data[0], data[1]),
		Station: data[3],
	}
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
