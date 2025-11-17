// Package pdu
package pdu

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/interfaces/global"
)

type WeatherResponse struct {
	*Base
	Metar string
}

func NewWeatherResponse(to string, metar string) *WeatherResponse {
	return &WeatherResponse{
		Base:  NewBase(fsd.ClientCommandWeatherResponse, global.FSDServerName, to),
		Metar: metar,
	}
}

func (c *WeatherResponse) Build() []byte {
	return MakeProtocolDataUnitPacket(
		c.GetType(),
		global.FSDServerName,
		c.To,
		c.Metar,
	)
}
