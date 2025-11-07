// Package fsd
package fsd

type Position struct {
	Latitude  float64
	Longitude float64
}

func (p *Position) PositionValid() bool {
	return p.Latitude != 0 && p.Longitude != 0
}
