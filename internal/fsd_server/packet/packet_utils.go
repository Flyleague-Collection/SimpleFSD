package packet

import (
	"bytes"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"strings"
)

const (
	CallsignMinLen = 3
	CallsignMaxLen = 12
	ForbiddenChars = "!@#$%*:& \t"
	FrequencyMin   = 18000
	FrequencyMax   = 36975
)

func parserCommandLine(line []byte) (ClientCommand, []string) {
	for _, prefix := range PossibleClientCommands {
		if bytes.HasPrefix(line, prefix) {
			decodeLine := string(line[len(prefix):])
			return ClientCommand(prefix), strings.Split(decodeLine, ":")
		}
	}
	return TempData, nil
}

func makePacket(command ClientCommand, parts ...string) []byte {
	totalLen := len(command)
	if len(parts) > 0 {
		for _, part := range parts {
			totalLen += len(part)
		}
		totalLen += len(parts) - 1
	}

	totalLen += splitSignLen

	result := make([]byte, totalLen)
	pos := 0

	pos += copy(result[pos:], command)

	for i, part := range parts {
		if i > 0 {
			result[pos] = ':'
			pos++
		}
		pos += copy(result[pos:], part)
	}

	copy(result[pos:], splitSign)

	return result
}

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

func frequencyValid(frequency int) bool {
	return FrequencyMin <= frequency && frequency <= FrequencyMax
}
