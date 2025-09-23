// Package email
package email

import (
	"fmt"
	. "github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"gopkg.in/gomail.v2"
	"html/template"
	"strings"
	"time"
)

type EmailSender struct {
	logger log.LoggerInterface
	config *config.EmailConfig
}

func NewEmailSender(
	logger log.LoggerInterface,
	config *config.EmailConfig,
) *EmailSender {
	return &EmailSender{
		logger: log.NewLoggerAdapter(logger, "EmailSender"),
		config: config,
	}
}

func formatCid(cid int) string {
	return fmt.Sprintf("%04d", cid)
}
func (sender *EmailSender) renderTemplate(template *template.Template, data interface{}) (string, error) {
	if template == nil {
		return "", ErrTemplateNotInitialized
	}

	var sb strings.Builder
	if err := template.Execute(&sb, data); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (sender *EmailSender) generateEmail(email string, title string, content string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", sender.config.Username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)

	return m
}

func (sender *EmailSender) SendApplicationPassedEmail(data *ApplicationPassedEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}

	email := strings.ToLower(data.User.Email)
	d := &ApplicationPassedEmail{
		Cid:      formatCid(data.User.Cid),
		Operator: formatCid(data.Operator.Cid),
		Rating:   data.Rating,
		Contact:  data.Operator.Email,
	}

	message, err := sender.renderTemplate(sender.config.Template.ApplicationPassedTemplate, d)
	if err != nil {
		sender.logger.WarnF("Error rendering application passed email template: %v", err)
		return ErrRenderingTemplate
	}

	m := sender.generateEmail(email, "管制员申请通过", message)

	sender.logger.InfoF("Sending application passed email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendApplicationProcessingEmail(data *ApplicationProcessingEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}

	email := strings.ToLower(data.User.Email)
	times := make([]string, len(data.AvailableTimes))
	for _, availableTime := range data.AvailableTimes {
		times = append(times, availableTime.Format("2006-01-02 15:04:05 MST"))
	}
	d := &ApplicationProcessingEmail{
		Cid:  formatCid(data.User.Cid),
		Time: strings.Join(times, ", "),
	}

	message, err := sender.renderTemplate(sender.config.Template.ApplicationProcessingTemplate, d)
	if err != nil {
		sender.logger.WarnF("Error rendering application processing email template: %v", err)
		return ErrRenderingTemplate
	}

	m := sender.generateEmail(email, "管制员申请进度通知", message)

	sender.logger.InfoF("Sending password application processing email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendApplicationRejectedEmail(data *ApplicationRejectedEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}

	email := strings.ToLower(data.User.Email)
	d := &ApplicationRejectedEmail{
		Cid:      formatCid(data.User.Cid),
		Operator: formatCid(data.Operator.Cid),
		Reason:   data.Reason,
		Contact:  data.Operator.Email,
	}

	message, err := sender.renderTemplate(sender.config.Template.ApplicationRejectedTemplate, d)
	if err != nil {
		sender.logger.WarnF("Error rendering application rejected email template: %v", err)
		return ErrRenderingTemplate
	}

	m := sender.generateEmail(email, "飞控密码更改通知", message)

	sender.logger.InfoF("Sending application rejected email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendAtcRatingChangeEmail(data *AtcRatingChangeEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}

	email := strings.ToLower(data.User.Email)
	d := &AtcRatingChangeEmail{
		Cid:      fmt.Sprintf("%04d", data.User.Cid),
		OldValue: data.OldRating,
		NewValue: data.NewRating,
		Operator: fmt.Sprintf("%04d", data.Operator.Cid),
		Contact:  data.Operator.Email,
	}

	message, err := sender.renderTemplate(sender.config.Template.ATCRatingChangeTemplate, d)
	if err != nil {
		sender.logger.WarnF("Error rendering rating change email template: %v", err)
		return ErrRenderingTemplate
	}

	m := sender.generateEmail(email, "管制权限变更通知", message)

	sender.logger.InfoF("Sending rating change email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendEmailVerifyEmail(data *EmailVerifyEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}

	email := strings.ToLower(data.Email)

	d := &EmailVerifyEmail{
		Cid:     fmt.Sprintf("%04d", data.Cid),
		Code:    fmt.Sprintf("%06d", data.Code),
		Expired: fmt.Sprintf("%.0f", sender.config.VerifyExpiredDuration.Minutes()),
	}

	message, err := sender.renderTemplate(sender.config.Template.EmailVerifyTemplate, d)
	if err != nil {
		sender.logger.WarnF("Error rendering email verification template: %v", err)
		return ErrRenderingTemplate
	}

	m := sender.generateEmail(email, "您的验证码", message)

	sender.logger.InfoF("Sending email verification code(%d) to %s(%d)", data.Code, email, data.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendKickedFromServerEmail(data *KickedFromServerEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}

	email := strings.ToLower(data.User.Email)
	d := &KickedFromServerEmail{
		Cid:      fmt.Sprintf("%04d", data.User.Cid),
		Time:     time.Now().Format(time.DateTime),
		Reason:   data.Reason,
		Operator: fmt.Sprintf("%04d", data.Operator.Cid),
		Contact:  data.Operator.Email,
	}

	message, err := sender.renderTemplate(sender.config.Template.KickedFromServerTemplate, d)
	if err != nil {
		sender.logger.WarnF("Error rendering kick message email template: %v", err)
		return ErrRenderingTemplate
	}

	m := sender.generateEmail(email, "踢出服务器通知", message)

	sender.logger.InfoF("Sending kick message email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendPasswordChangeEmail(data *PasswordChangeEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}

	email := strings.ToLower(data.User.Email)
	d := &PasswordChangeEmail{
		Cid:       formatCid(data.User.Cid),
		IP:        data.Ip,
		UserAgent: data.UserAgent,
		Time:      time.Now().Format(time.DateTime),
	}

	message, err := sender.renderTemplate(sender.config.Template.PasswordChangeTemplate, d)
	if err != nil {
		sender.logger.WarnF("Error rendering password change email template: %v", err)
		return ErrRenderingTemplate
	}

	m := sender.generateEmail(email, "飞控密码更改通知", message)

	sender.logger.InfoF("Sending password change email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendPermissionChangeEmail(data *PermissionChangeEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}

	email := strings.ToLower(data.User.Email)
	d := &PermissionChangeEmail{
		Cid:         fmt.Sprintf("%04d", data.User.Cid),
		Operator:    fmt.Sprintf("%04d", data.Operator.Cid),
		Contact:     data.Operator.Email,
		Permissions: strings.Join(data.Permissions, ", "),
	}

	message, err := sender.renderTemplate(sender.config.Template.PermissionChangeTemplate, d)
	if err != nil {
		sender.logger.WarnF("Error rendering permission change email template: %v", err)
		return ErrRenderingTemplate
	}

	m := sender.generateEmail(email, "管理权限变更通知", message)

	sender.logger.InfoF("Sending permission change email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}

func (sender *EmailSender) SendTicketReplyEmail(data *TicketReplyEmailData) error {
	if sender.config.EmailServer == nil {
		return nil
	}

	email := strings.ToLower(data.User.Email)
	d := &TicketReplyEmail{
		Cid:   formatCid(data.User.Cid),
		Title: data.Title,
		Reply: data.Reply,
	}

	message, err := sender.renderTemplate(sender.config.Template.TicketReplyTemplate, d)
	if err != nil {
		sender.logger.WarnF("Error rendering ticket reply email template: %v", err)
		return ErrRenderingTemplate
	}

	m := sender.generateEmail(email, "工单回复通知", message)

	sender.logger.InfoF("Sending ticket reply email to %s(%d)", email, data.User.Cid)

	return sender.config.EmailServer.DialAndSend(m)
}
