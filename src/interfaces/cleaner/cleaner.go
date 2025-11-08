// Package cleaner
package cleaner

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/global"
)

type Interface interface {
	Init()
	Add(callable global.Callable)
	Clean()
}
