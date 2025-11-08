// Package email
package email

type EmailSenderInterface interface {
	SendApplicationPassedEmail(data *ApplicationPassedEmailData) error
	SendApplicationProcessingEmail(data *ApplicationProcessingEmailData) error
	SendApplicationRejectedEmail(data *ApplicationRejectedEmailData) error
	SendAtcRatingChangeEmail(data *AtcRatingChangeEmailData) error
	SendEmailVerifyEmail(data *EmailVerifyEmailData) error
	SendKickedFromServerEmail(data *KickedFromServerEmailData) error
	SendPasswordChangeEmail(data *PasswordChangeEmailData) error
	SendPasswordResetEmail(data *PasswordResetEmailData) error
	SendPermissionChangeEmail(data *PermissionChangeEmailData) error
	SendTicketReplyEmail(data *TicketReplyEmailData) error
}
