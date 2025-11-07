// Package global
package global

import "context"

type Callable interface {
	Invoke(ctx context.Context) error
}
