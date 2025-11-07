// Package interfaces
package interfaces

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/global"
)

type CleanerInterface interface {
	Init()
	Add(callable global.Callable)
	Clean()
}
