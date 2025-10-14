// Package fsd
package fsd

import "errors"

var (
	ErrCidNotFound        = errors.New("target cid not found")
	ErrConnectionNotFound = errors.New("connection not found")
)

type ConnectionManagerInterface interface {
	AddConnection(client ClientInterface)
	RemoveConnection(client ClientInterface) error
	GetConnections(cid int) ([]ClientInterface, error)
}
