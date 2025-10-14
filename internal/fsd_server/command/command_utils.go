// Package command
package command

import "strings"

const (
	CallsignMinLen = 3
	CallsignMaxLen = 12
	ForbiddenChars = "!@#$%*:& \t"
)

var validSuffix = [6]string{"DEL", "GND", "TWR", "APP", "CTR", "FSS"}

func isValidAtc(callsign string) bool {
	if !callsignValid(callsign) {
		return false
	}
	for _, prefix := range validSuffix {
		if strings.HasSuffix(callsign, prefix) {
			return true
		}
	}
	return false
}

func callsignValid(callsign string) bool {
	if len(callsign) < CallsignMinLen || len(callsign) >= CallsignMaxLen {
		return false
	}

	if strings.ContainsAny(callsign, ForbiddenChars) {
		return false
	}

	for _, r := range callsign {
		if r > 127 {
			return false
		}
	}

	return true
}
