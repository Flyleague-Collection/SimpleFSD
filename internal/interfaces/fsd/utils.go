// Package fsd
package fsd

import "bytes"

var (
	SplitSign    = []byte("\r\n")
	SplitSignLen = len(SplitSign)
)

const (
	FrequencyMin = 118000
	FrequencyMax = 136975
)

func MakePacketWithoutSign(command ClientCommand, parts ...string) []byte {
	return bytes.TrimRight(MakePacket(command, parts...), string(SplitSign))
}

func MakePacket(command ClientCommand, parts ...string) []byte {
	totalLen := len(command)
	if len(parts) > 0 {
		for _, part := range parts {
			totalLen += len(part)
		}
		totalLen += len(parts) - 1
	}

	totalLen += SplitSignLen

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

	copy(result[pos:], SplitSign)

	return result
}

func FrequencyValid(frequency int) bool {
	return FrequencyMin <= frequency && frequency <= FrequencyMax
}
