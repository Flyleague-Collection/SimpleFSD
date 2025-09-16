// Package queue
package queue

import (
	"errors"
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
)

type MessageType int

const (
	SendVerifyEmail MessageType = iota
	SendRatingChangeEmail
	SendPermissionChangeEmail
	SendPasswordChangeEmail
	SendKickedFromServerEmail
	SendMessageToClient
	KickClientFromServer
	AuditLog
	AuditLogs
)

var messageTypes = []string{
	"SendVerifyEmail", "SendRatingChangeEmail", "SendPermissionChangeEmail", "SendPasswordChangeEmail",
	"SendKickedFromServerEmail", "SendMessageToClient", "KickClientFromServer", "AuditLog", "AuditLogs",
}

func (messageType MessageType) String() string {
	return messageTypes[messageType]
}

type Message struct {
	Type MessageType
	Data interface{}
}

type Subscriber func(message *Message) error

var ErrMessageDataType = errors.New("message data type not correct")

type MessageQueueInterface interface {
	Start()
	Stop()
	ShutdownCallback() global.Callable
	Publish(message *Message)
	SyncPublish(message *Message) error
	Subscribe(messageType MessageType, handler Subscriber)
}
