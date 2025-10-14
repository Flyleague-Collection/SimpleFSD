// Package queue
package queue

import (
	"errors"

	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
)

type MessageType int

const (
	SendApplicationPassedEmail MessageType = iota
	SendApplicationProcessingEmail
	SendApplicationRejectedEmail
	SendAtcRatingChangeEmail
	SendEmailVerifyEmail
	SendKickedFromServerEmail
	SendPasswordChangeEmail
	SendPasswordResetEmail
	SendPermissionChangeEmail
	SendTicketReplyEmail
	SendMessageToClient
	DeleteVerifyCode
	KickClientFromServer
	BroadcastMessage
	FlushFlightPlan
	ChangeFlightPlanLockStatus
	AuditLog
	AuditLogs
	FsdMessageReceived
)

var messageTypes = []string{
	"SendApplicationPassedEmail",
	"SendApplicationProcessingEmail",
	"SendApplicationRejectedEmail",
	"SendAtcRatingChangeEmail",
	"SendEmailVerifyEmail",
	"SendKickedFromServerEmail",
	"SendPasswordChangeEmail",
	"SendPasswordResetEmail",
	"SendPermissionChangeEmail",
	"SendTicketReplyEmail",
	"SendMessageToClient",
	"DeleteVerifyCode",
	"KickClientFromServer",
	"BroadcastMessage",
	"FlushFlightPlan",
	"ChangeFlightPlanLockStatus",
	"AuditLog",
	"AuditLogs",
	"FsdMessageReceived",
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
