// Package service
// 存放 EmailServiceInterface 的实现
package service

import (
	"errors"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"gopkg.in/gomail.v2"
	"html/template"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type EmailService struct {
	logger       log.LoggerInterface
	emailCodes   map[string]EmailCode
	lastSendTime map[string]time.Time
	config       *config.EmailConfig
}

type EmailCode struct {
	code     int
	cid      int
	sendTime time.Time
}

type EmailVerifyTemplateData struct {
	Cid     string
	Code    string
	Expired string
}

type EmailPermissionChangeData struct {
	Cid      string
	Operator string
	Contact  string
}

type EmailRatingChangeData struct {
	Cid      string
	NewValue string
	OldValue string
	Operator string
	Contact  string
}

type EmailKickedFromServerData struct {
	Cid      string
	Time     string
	Reason   string
	Operator string
	Contact  string
}

func NewEmailService(logger log.LoggerInterface, config *config.EmailConfig) *EmailService {
	return &EmailService{
		logger:       log.NewLoggerAdapter(logger, "EmailService"),
		config:       config,
		emailCodes:   make(map[string]EmailCode),
		lastSendTime: make(map[string]time.Time),
	}
}

var (
	ErrEmailSendInterval      = errors.New("email send interval")
	ErrRenderingTemplate      = errors.New("error rendering template")
	ErrTemplateNotInitialized = errors.New("error template not initialized")
	ErrEmailCodeNotFound      = errors.New("email code not found")
	ErrEmailCodeExpired       = errors.New("email code expired")
	ErrInvalidEmailCode       = errors.New("invalid email code")
	ErrCidMismatch            = errors.New("cid mismatch")
)

func (emailService *EmailService) renderTemplate(template *template.Template, data interface{}) (string, error) {
	if template == nil {
		return "", ErrTemplateNotInitialized
	}

	var sb strings.Builder
	if err := template.Execute(&sb, data); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (emailService *EmailService) VerifyEmailCode(email string, code int, cid int) error {
	if emailService.config.EmailServer == nil {
		return nil
	}

	email = strings.ToLower(email)
	emailCode, ok := emailService.emailCodes[email]
	if !ok {
		return ErrEmailCodeNotFound
	}

	if time.Since(emailCode.sendTime) > emailService.config.VerifyExpiredDuration {
		return ErrEmailCodeExpired
	}

	if emailCode.code != code {
		return ErrInvalidEmailCode
	}

	if emailCode.cid != cid {
		return ErrCidMismatch
	}

	delete(emailService.emailCodes, email)
	return nil
}

func (emailService *EmailService) HandleSendKickedFromServerEmailMessage(message *queue.Message) error {
	if val, ok := message.Data.(*KickedFromServerEmailData); ok {
		return emailService.sendKickedFromServerEmail(val)
	}
	return queue.ErrMessageDataType
}

func (emailService *EmailService) HandleSendPermissionChangeEmailMessage(message *queue.Message) error {
	if val, ok := message.Data.(*PermissionChangeEmailData); ok {
		return emailService.sendPermissionChangeEmail(val)
	}
	return queue.ErrMessageDataType
}

func (emailService *EmailService) HandleSendRatingChangeEmailMessage(message *queue.Message) error {
	if val, ok := message.Data.(*RatingChangeEmailData); ok {
		return emailService.sendRatingChangeEmail(val)
	}
	return queue.ErrMessageDataType
}

func (emailService *EmailService) HandleSendVerifyEmailMessage(message *queue.Message) error {
	if val, ok := message.Data.(*VerifyCodeEmailData); ok {
		err, _ := emailService.sendEmailCode(val)
		return err
	}
	return queue.ErrMessageDataType
}

func (emailService *EmailService) sendEmailCode(data *VerifyCodeEmailData) (error, time.Duration) {
	if emailService.config.EmailServer == nil {
		return nil, 0
	}

	email := strings.ToLower(data.Email)
	if lastSendTime, ok := emailService.lastSendTime[email]; ok {
		now := time.Now()
		timeRemain := lastSendTime.Add(emailService.config.SendDuration).Sub(now)
		if timeRemain > 0 {
			return ErrEmailSendInterval, timeRemain
		}
	}

	code := rand.Intn(1e6)
	emailCode := EmailCode{code: code, cid: data.Cid, sendTime: time.Now()}
	d := &EmailVerifyTemplateData{
		Cid:     fmt.Sprintf("%04d", data.Cid),
		Code:    strconv.Itoa(code),
		Expired: strconv.Itoa(int(emailService.config.VerifyExpiredDuration.Minutes())),
	}

	message, err := emailService.renderTemplate(emailService.config.Template.EmailVerifyTemplate, d)
	if err != nil {
		emailService.logger.WarnF("Error rendering email verification template: %v", err)
		return ErrRenderingTemplate, 0
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailService.config.Username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "您的验证码")
	m.SetBody("text/html", message)

	emailService.emailCodes[email] = emailCode
	emailService.lastSendTime[email] = time.Now()

	emailService.logger.InfoF("Sending email verification code(%d) to %s(%d)", code, email, data.Cid)

	return emailService.config.EmailServer.DialAndSend(m), 0
}

func (emailService *EmailService) sendPermissionChangeEmail(data *PermissionChangeEmailData) error {
	if emailService.config.EmailServer == nil {
		return nil
	}

	email := strings.ToLower(data.User.Email)
	d := &EmailPermissionChangeData{
		Cid:      fmt.Sprintf("%04d", data.User.Cid),
		Operator: fmt.Sprintf("%04d", data.Operator.Cid),
		Contact:  data.Operator.Email,
	}

	message, err := emailService.renderTemplate(emailService.config.Template.PermissionChangeTemplate, d)
	if err != nil {
		emailService.logger.WarnF("Error rendering email verification template: %v", err)
		return ErrRenderingTemplate
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailService.config.Username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "管理权限变更通知")
	m.SetBody("text/html", message)

	emailService.logger.InfoF("Sending permission change email to %s(%d)", email, data.User.Cid)

	return emailService.config.EmailServer.DialAndSend(m)
}

func (emailService *EmailService) sendRatingChangeEmail(data *RatingChangeEmailData) error {
	if emailService.config.EmailServer == nil {
		return nil
	}

	email := strings.ToLower(data.User.Email)
	d := &EmailRatingChangeData{
		Cid:      fmt.Sprintf("%04d", data.User.Cid),
		OldValue: data.OldRating.String(),
		NewValue: data.NewRating.String(),
		Operator: fmt.Sprintf("%04d", data.Operator.Cid),
		Contact:  data.Operator.Email,
	}

	message, err := emailService.renderTemplate(emailService.config.Template.ATCRatingChangeTemplate, d)
	if err != nil {
		emailService.logger.WarnF("Error rendering email verification template: %v", err)
		return ErrRenderingTemplate
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailService.config.Username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "管制权限变更通知")
	m.SetBody("text/html", message)

	emailService.logger.InfoF("Sending rating change email to %s(%d)", email, data.User.Cid)

	return emailService.config.EmailServer.DialAndSend(m)
}

func (emailService *EmailService) sendKickedFromServerEmail(data *KickedFromServerEmailData) error {
	if emailService.config.EmailServer == nil {
		return nil
	}

	email := strings.ToLower(data.User.Email)
	d := &EmailKickedFromServerData{
		Cid:      fmt.Sprintf("%04d", data.User.Cid),
		Time:     time.Now().Format(time.DateTime),
		Reason:   data.Reason,
		Operator: fmt.Sprintf("%04d", data.Operator.Cid),
		Contact:  data.Operator.Email,
	}

	message, err := emailService.renderTemplate(emailService.config.Template.KickedFromServerTemplate, d)
	if err != nil {
		emailService.logger.WarnF("Error rendering email verification template: %v", err)
		return ErrRenderingTemplate
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailService.config.Username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "踢出服务器通知")
	m.SetBody("text/html", message)

	emailService.logger.InfoF("Sending kick message email to %s(%d)", email, data.User.Cid)

	return emailService.config.EmailServer.DialAndSend(m)
}

var (
	ErrRenderTemplate = NewApiStatus("RENDER_TEMPLATE_ERROR", "邮件发送失败", ServerInternalError)
	ErrSendEmail      = NewApiStatus("EMAIL_SEND_ERROR", "发送失败", ServerInternalError)
	SendEmailSuccess  = NewApiStatus("SEND_EMAIL_SUCCESS", "邮件发送成功", Ok)
)

func (emailService *EmailService) SendEmailVerifyCode(req *RequestEmailVerifyCode) *ApiResponse[ResponseEmailVerifyCode] {
	if emailService.config.EmailServer == nil {
		return NewApiResponse(SendEmailSuccess, &ResponseEmailVerifyCode{Email: req.Email})
	}

	if req.Email == "" || req.Cid <= 0 {
		return NewApiResponse[ResponseEmailVerifyCode](ErrIllegalParam, nil)
	}

	err, remainTime := emailService.sendEmailCode(&VerifyCodeEmailData{Email: req.Email, Cid: req.Cid})
	if err == nil {
		return NewApiResponse(SendEmailSuccess, &ResponseEmailVerifyCode{Email: req.Email})
	}

	if errors.Is(err, ErrEmailSendInterval) {
		return NewApiResponse[ResponseEmailVerifyCode](NewApiStatus(
			"EMAIL_SEND_INTERVAL",
			fmt.Sprintf("邮件已发送, 请在%.0f秒后重试", remainTime.Seconds()),
			BadRequest,
		), nil)
	}

	if errors.Is(err, ErrRenderingTemplate) {
		return NewApiResponse[ResponseEmailVerifyCode](ErrRenderTemplate, nil)
	}

	return NewApiResponse[ResponseEmailVerifyCode](ErrSendEmail, nil)
}
