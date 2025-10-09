// Package websocket
package websocket

const (
	TokenInvalid int = 4000 + iota
	TokenExpired
	TokenFormatError
)

type SendMessage struct {
	Target string `json:"target"`
	Data   string `json:"data"`
}

type ReceiveMessage struct {
	From string `json:"from"`
	Data string `json:"data"`
}
