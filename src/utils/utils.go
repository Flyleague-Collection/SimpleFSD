// Package utils
package utils

import "fmt"

func FormatCid(cid int) string {
	return fmt.Sprintf("%04d", cid)
}
