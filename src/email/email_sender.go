// Package email
package email

import (
	"fmt"
	"strings"
	"time"

	. "github.com/half-nothing/simple-fsd/src/interfaces"
	"github.com/half-nothing/simple-fsd/src/interfaces/config"
	"github.com/half-nothing/simple-fsd/src/interfaces/log"
	"github.com/half-nothing/simple-fsd/src/utils"
	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	logger         log.LoggerInterface
	config         *config.EmailConfig
	templateConfig *config.EmailTemplateConfigs
}

func NewEmailSender(
	logger log.LoggerInterface,
	config *config.EmailConfig,
) *EmailSender {
	return &EmailSender{
		logger:         log.NewLoggerAdapter(logger, "EmailSender"),
		config:         config,
		templateConfig: config.Template,
	}
}

func (sender *EmailSender) renderTemplate(template *config.EmailTemplateConfig, data interface{}) (string, error) {
	if template.Template == nil {
		return "", ErrTemplateNotInitialized
	}

	var sb strings.Builder
	if err := template.Template.Execute(&sb, data); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (sender *EmailSender) generateEmail(email string, config *config.EmailTemplateConfig, data interface{}) (*gomail.Message, error) {
	content, err := sender.renderTemplate(config, data)
	if err != nil {
		return nil, err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", sender.config.Username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", config.EmailTitle)
	m.SetBody("text/html", content)

	return m, nil
}

func (sender *EmailSender) SendApplicationPassedEmail(data *ApplicationPassedEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}
	if !sender.templateConfig.ApplicationPassedEmail.Enable {
		return nil
	}

	email := strings.ToLower(data.User.Email)

	m, err := sender.generateEmail(email, sender.templateConfig.ApplicationPassedEmail, &ApplicationPassedEmail{
		Cid:      utils.FormatCid(data.User.Cid),
		Operator: utils.FormatCid(data.Operator.Cid),
		Contact:  data.Operator.Email,
		Message:  data.Message,
	})
	if err != nil {
		sender.logger.WarnF("Error rendering application passed email template: %v", err)
		return ErrRenderingTemplate
	}

	sender.logger.InfoF("Sending application passed email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendApplicationProcessingEmail(data *ApplicationProcessingEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}
	if !sender.templateConfig.ApplicationProcessingEmail.Enable {
		return nil
	}

	email := strings.ToLower(data.User.Email)
	times := make([]string, 0, len(data.AvailableTimes))
	for _, availableTime := range data.AvailableTimes {
		times = append(times, availableTime.Local().Format("2006-01-02 15:04:05 MST"))
	}

	m, err := sender.generateEmail(email, sender.templateConfig.ApplicationProcessingEmail, &ApplicationProcessingEmail{
		Cid:     utils.FormatCid(data.User.Cid),
		Time:    strings.Join(times, ", "),
		Contact: data.Operator.Email,
	})
	if err != nil {
		sender.logger.WarnF("Error rendering application processing email template: %v", err)
		return ErrRenderingTemplate
	}

	sender.logger.InfoF("Sending password application processing email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendApplicationRejectedEmail(data *ApplicationRejectedEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}
	if !sender.templateConfig.ApplicationRejectedEmail.Enable {
		return nil
	}

	email := strings.ToLower(data.User.Email)

	m, err := sender.generateEmail(email, sender.templateConfig.ApplicationRejectedEmail, &ApplicationRejectedEmail{
		Cid:      utils.FormatCid(data.User.Cid),
		Operator: utils.FormatCid(data.Operator.Cid),
		Reason:   data.Reason,
		Contact:  data.Operator.Email,
	})
	if err != nil {
		sender.logger.WarnF("Error rendering application rejected email template: %v", err)
		return ErrRenderingTemplate
	}

	sender.logger.InfoF("Sending application rejected email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendAtcRatingChangeEmail(data *AtcRatingChangeEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}
	if !sender.templateConfig.ATCRatingChangeEmail.Enable {
		return nil
	}

	email := strings.ToLower(data.User.Email)

	m, err := sender.generateEmail(email, sender.templateConfig.ATCRatingChangeEmail, &AtcRatingChangeEmail{
		Cid:      fmt.Sprintf("%04d", data.User.Cid),
		OldValue: data.OldRating,
		NewValue: data.NewRating,
		Operator: fmt.Sprintf("%04d", data.Operator.Cid),
		Contact:  data.Operator.Email,
	})
	if err != nil {
		sender.logger.WarnF("Error rendering rating change email template: %v", err)
		return ErrRenderingTemplate
	}

	sender.logger.InfoF("Sending rating change email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendEmailVerifyEmail(data *EmailVerifyEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}

	email := strings.ToLower(data.Email)

	m, err := sender.generateEmail(email, sender.templateConfig.VerifyCodeEmail, &EmailVerifyEmail{
		Cid:       fmt.Sprintf("%04d", data.Cid),
		Code:      fmt.Sprintf("%06d", data.Code),
		Expired:   fmt.Sprintf("%.0f", sender.config.VerifyExpiredDuration.Minutes()),
		ExpiredAt: time.Now().Add(sender.config.VerifyExpiredDuration).Format("2006-01-02 15:04:05 MST"),
	})
	if err != nil {
		sender.logger.WarnF("Error rendering email verification template: %v", err)
		return ErrRenderingTemplate
	}

	sender.logger.InfoF("Sending email verification code(%d) to %s(%d)", data.Code, email, data.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendKickedFromServerEmail(data *KickedFromServerEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}
	if !sender.templateConfig.KickedFromServerEmail.Enable {
		return nil
	}

	email := strings.ToLower(data.User.Email)

	m, err := sender.generateEmail(email, sender.templateConfig.KickedFromServerEmail, &KickedFromServerEmail{
		Cid:      fmt.Sprintf("%04d", data.User.Cid),
		Time:     time.Now().Format(time.DateTime),
		Reason:   data.Reason,
		Operator: fmt.Sprintf("%04d", data.Operator.Cid),
		Contact:  data.Operator.Email,
	})
	if err != nil {
		sender.logger.WarnF("Error rendering kick message email template: %v", err)
		return ErrRenderingTemplate
	}

	sender.logger.InfoF("Sending kick message email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendPasswordChangeEmail(data *PasswordChangeEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}
	if !sender.templateConfig.PasswordChangeEmail.Enable {
		return nil
	}

	email := strings.ToLower(data.User.Email)

	m, err := sender.generateEmail(email, sender.templateConfig.PasswordChangeEmail, &PasswordChangeEmail{
		Cid:       utils.FormatCid(data.User.Cid),
		IP:        data.Ip,
		UserAgent: data.UserAgent,
		Time:      time.Now().Format(time.DateTime),
	})
	if err != nil {
		sender.logger.WarnF("Error rendering password change email template: %v", err)
		return ErrRenderingTemplate
	}

	sender.logger.InfoF("Sending password change email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendPasswordResetEmail(data *PasswordResetEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}
	if !sender.templateConfig.PasswordChangeEmail.Enable {
		return nil
	}

	email := strings.ToLower(data.User.Email)

	m, err := sender.generateEmail(email, sender.templateConfig.PasswordResetEmail, &PasswordResetEmail{
		Cid:       utils.FormatCid(data.User.Cid),
		IP:        data.Ip,
		UserAgent: data.UserAgent,
		Time:      time.Now().Format(time.DateTime),
	})
	if err != nil {
		sender.logger.WarnF("Error rendering password reset email template: %v", err)
		return ErrRenderingTemplate
	}

	sender.logger.InfoF("Sending password reset email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendPermissionChangeEmail(data *PermissionChangeEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}
	if !sender.templateConfig.PermissionChangeEmail.Enable {
		return nil
	}

	email := strings.ToLower(data.User.Email)

	m, err := sender.generateEmail(email, sender.templateConfig.PermissionChangeEmail, &PermissionChangeEmail{
		Cid:         fmt.Sprintf("%04d", data.User.Cid),
		Operator:    fmt.Sprintf("%04d", data.Operator.Cid),
		Contact:     data.Operator.Email,
		Permissions: strings.Join(data.Permissions, ", "),
	})
	if err != nil {
		sender.logger.WarnF("Error rendering permission change email template: %v", err)
		return ErrRenderingTemplate
	}

	sender.logger.InfoF("Sending permission change email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendTicketReplyEmail(data *TicketReplyEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}
	if !sender.templateConfig.TicketReplyEmail.Enable {
		return nil
	}

	email := strings.ToLower(data.User.Email)

	m, err := sender.generateEmail(email, sender.templateConfig.TicketReplyEmail, &TicketReplyEmail{
		Cid:   utils.FormatCid(data.User.Cid),
		Title: data.Title,
		Reply: data.Reply,
	})
	if err != nil {
		sender.logger.WarnF("Error rendering ticket reply email template: %v", err)
		return ErrRenderingTemplate
	}

	sender.logger.InfoF("Sending ticket reply email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}
