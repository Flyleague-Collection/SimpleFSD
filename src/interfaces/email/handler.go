// Package email
package email

import "github.com/half-nothing/simple-fsd/src/interfaces/queue"

type EmailMessageHandlerInterface interface {
	HandleSendApplicationPassedEmailMessage(message *queue.Message) error
	HandleSendApplicationProcessingEmailMessage(message *queue.Message) error
	HandleSendApplicationRejectedEmailMessage(message *queue.Message) error
	HandleSendAtcRatingChangeEmailMessage(message *queue.Message) error
	HandleSendEmailVerifyEmailMessage(message *queue.Message) error
	HandleSendKickedFromServerEmailMessage(message *queue.Message) error
	HandleSendPasswordChangeEmailMessage(message *queue.Message) error
	HandleSendPasswordResetEmailMessage(message *queue.Message) error
	HandleSendPermissionChangeEmailMessage(message *queue.Message) error
	HandleSendTicketReplyEmailMessage(message *queue.Message) error
}
