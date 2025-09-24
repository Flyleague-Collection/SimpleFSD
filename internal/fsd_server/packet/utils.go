// Package packet
package packet

import (
	"bufio"
	"bytes"
	"errors"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"net"
	"strings"
)

func parserCommandLine(line []byte, possibleCommands [][]byte) (ClientCommand, []string) {
	for _, prefix := range possibleCommands {
		if bytes.HasPrefix(line, prefix) {
			decodeLine := string(line[len(prefix):])
			return ClientCommand(prefix), strings.Split(decodeLine, ":")
		}
	}
	return TempData, nil
}

func isNetClosedError(err error) bool {
	if errors.Is(err, net.ErrClosed) {
		return true
	}
	var opErr *net.OpError
	ok := errors.As(err, &opErr)
	return ok && opErr.Timeout()
}

func createSplitFunc(sep []byte) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.Index(data, sep); i >= 0 {
			return i + len(sep), data[0:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	}
}
