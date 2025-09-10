// Package fsd
package fsd

import (
	"math"
)

type BroadcastTarget string

var (
	AllClient BroadcastTarget = "*"
	AllATC    BroadcastTarget = "*A"
	AllSup    BroadcastTarget = "*S"
)

func (b BroadcastTarget) String() string {
	return string(b)
}

func (b BroadcastTarget) Index() int {
	return 0
}

type BroadcastFilter func(toClient, fromClient ClientInterface) bool

func BroadcastToAll(_, _ ClientInterface) bool {
	return true
}

func BroadcastToPilot(toClient, _ ClientInterface) bool {
	return !toClient.IsAtc()
}

func BroadcastToAtc(toClient, _ ClientInterface) bool {
	return toClient.IsAtc()
}

func BroadcastToSup(toClient, _ ClientInterface) bool {
	if !toClient.IsAtc() {
		return false
	}
	return toClient.Rating() >= Supervisor
}

func BroadcastToClientInRange(toClient, fromClient ClientInterface) bool {
	if fromClient == nil {
		return true
	}
	distance := FindNearestDistance(toClient.Position(), fromClient.Position())
	var threshold float64 = 0
	switch {
	case toClient.IsAtc() && fromClient.IsAtc():
		threshold = math.Max(toClient.VisualRange(), fromClient.VisualRange())
	case toClient.IsAtc():
		threshold = toClient.VisualRange()
	case fromClient.IsAtc():
		threshold = fromClient.VisualRange()
	default:
		threshold = toClient.VisualRange() + fromClient.VisualRange()
	}
	return distance <= threshold
}

func CombineBroadcastFilter(filters ...BroadcastFilter) BroadcastFilter {
	return func(toClient, fromClient ClientInterface) bool {
		for _, f := range filters {
			if f == nil {
				continue
			}
			if !f(toClient, fromClient) {
				return false
			}
		}
		return true
	}
}
