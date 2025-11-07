// Package email
package email

import (
	. "github.com/half-nothing/simple-fsd/src/interfaces"
	"github.com/half-nothing/simple-fsd/src/interfaces/queue"
)

type EmailMessageHandler struct {
	sender EmailSenderInterface
}

func NewEmailMessageHandler(
	sender EmailSenderInterface,
) *EmailMessageHandler {
	return &EmailMessageHandler{
		sender: sender,
	}
}

func (handler *EmailMessageHandler) HandleSendApplicationPassedEmailMessage(message *queue.Message) error {
	if val, ok := message.Data.(*ApplicationPassedEmailData); ok {
		return handler.sender.SendApplicationPassedEmail(val)
	}
	return queue.ErrMessageDataType
}

func (handler *EmailMessageHandler) HandleSendApplicationProcessingEmailMessage(message *queue.Message) error {
	if val, ok := message.Data.(*ApplicationProcessingEmailData); ok {
		return handler.sender.SendApplicationProcessingEmail(val)
	}
	return queue.ErrMessageDataType
}

func (handler *EmailMessageHandler) HandleSendApplicationRejectedEmailMessage(message *queue.Message) error {
	if val, ok := message.Data.(*ApplicationRejectedEmailData); ok {
		return handler.sender.SendApplicationRejectedEmail(val)
	}
	return queue.ErrMessageDataType
}

func (handler *EmailMessageHandler) HandleSendAtcRatingChangeEmailMessage(message *queue.Message) error {
	if val, ok := message.Data.(*AtcRatingChangeEmailData); ok {
		return handler.sender.SendAtcRatingChangeEmail(val)
	}
	return queue.ErrMessageDataType
}

func (handler *EmailMessageHandler) HandleSendEmailVerifyEmailMessage(message *queue.Message) error {
	if val, ok := message.Data.(*EmailVerifyEmailData); ok {
		return handler.sender.SendEmailVerifyEmail(val)
	}
	return queue.ErrMessageDataType
}

func (handler *EmailMessageHandler) HandleSendKickedFromServerEmailMessage(message *queue.Message) error {
	if val, ok := message.Data.(*KickedFromServerEmailData); ok {
		return handler.sender.SendKickedFromServerEmail(val)
	}
	return queue.ErrMessageDataType
}

func (handler *EmailMessageHandler) HandleSendPasswordChangeEmailMessage(message *queue.Message) error {
	if val, ok := message.Data.(*PasswordChangeEmailData); ok {
		return handler.sender.SendPasswordChangeEmail(val)
	}
	return queue.ErrMessageDataType
}

func (handler *EmailMessageHandler) HandleSendPasswordResetEmailMessage(message *queue.Message) error {
	if val, ok := message.Data.(*PasswordResetEmailData); ok {
		return handler.sender.SendPasswordResetEmail(val)
	}
	return queue.ErrMessageDataType
}

func (handler *EmailMessageHandler) HandleSendPermissionChangeEmailMessage(message *queue.Message) error {
	if val, ok := message.Data.(*PermissionChangeEmailData); ok {
		return handler.sender.SendPermissionChangeEmail(val)
	}
	return queue.ErrMessageDataType
}

func (handler *EmailMessageHandler) HandleSendTicketReplyEmailMessage(message *queue.Message) error {
	if val, ok := message.Data.(*TicketReplyEmailData); ok {
		return handler.sender.SendTicketReplyEmail(val)
	}
	return queue.ErrMessageDataType
}
