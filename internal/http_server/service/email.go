// Package service
// 存放 EmailServiceInterface 的实现
package service

import (
	"errors"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/half-nothing/simple-fsd/internal/utils"
	"math/rand"
	"strings"
	"time"
)

type EmailService struct {
	logger            log.LoggerInterface
	config            *config.EmailConfig
	emailCodeCache    interfaces.CacheInterface[*EmailCode]
	lastSendTimeCache interfaces.CacheInterface[time.Time]
	messageQueue      queue.MessageQueueInterface
}

func NewEmailService(
	logger log.LoggerInterface,
	config *config.EmailConfig,
	emailCodeCache interfaces.CacheInterface[*EmailCode],
	lastSendTimeCache interfaces.CacheInterface[time.Time],
	messageQueue queue.MessageQueueInterface,
) *EmailService {
	return &EmailService{
		logger:            log.NewLoggerAdapter(logger, "EmailService"),
		config:            config,
		emailCodeCache:    emailCodeCache,
		lastSendTimeCache: lastSendTimeCache,
		messageQueue:      messageQueue,
	}
}

var (
	ErrRenderingTemplate = errors.New("error rendering template")
	ErrEmailCodeExpired  = errors.New("email code expired")
	ErrEmailCodeIllegal  = errors.New("email code illegal")
	ErrInvalidEmailCode  = errors.New("invalid email code")
	ErrCidMismatch       = errors.New("cid mismatch")
)

func (emailService *EmailService) VerifyEmailCode(email string, code string, cid int) error {
	if emailService.config.EmailServer == nil {
		return nil
	}

	realEmailCode := utils.StrToInt(code, -1)
	if realEmailCode == -1 {
		return ErrEmailCodeIllegal
	}

	email = strings.ToLower(email)
	emailCode, ok := emailService.emailCodeCache.Get(email)
	if !ok {
		return ErrEmailCodeExpired
	}

	if emailCode.Code != realEmailCode {
		return ErrInvalidEmailCode
	}

	if emailCode.Cid != cid {
		return ErrCidMismatch
	}

	emailService.emailCodeCache.Del(email)
	return nil
}

func (emailService *EmailService) SendEmailVerifyCode(req *RequestEmailVerifyCode) *ApiResponse[ResponseEmailVerifyCode] {
	if emailService.config.EmailServer == nil {
		return NewApiResponse(SendEmailSuccess, &ResponseEmailVerifyCode{Email: req.Email})
	}

	if req.Email == "" || req.Cid <= 0 {
		return NewApiResponse[ResponseEmailVerifyCode](ErrIllegalParam, nil)
	}

	if val, ok := emailService.lastSendTimeCache.Get(req.Email); ok {
		return NewApiResponse[ResponseEmailVerifyCode](NewApiStatus(
			"EMAIL_SEND_INTERVAL",
			fmt.Sprintf("邮件已发送, 请在%.0f秒后重试", time.Now().Sub(val).Seconds()),
			BadRequest,
		), nil)
	}

	code := rand.Intn(1e6)
	if err := emailService.messageQueue.SyncPublish(&queue.Message{
		Type: queue.SendEmailVerifyEmail,
		Data: &interfaces.EmailVerifyEmailData{
			Email: req.Email,
			Cid:   req.Cid,
			Code:  code,
		},
	}); err != nil {
		if errors.Is(err, ErrRenderingTemplate) {
			return NewApiResponse[ResponseEmailVerifyCode](ErrRenderTemplate, nil)
		}
		return NewApiResponse[ResponseEmailVerifyCode](ErrSendEmail, nil)
	}

	emailService.emailCodeCache.SetWithTTL(req.Email, &EmailCode{Code: code, Cid: req.Cid}, emailService.config.VerifyExpiredDuration)
	emailService.lastSendTimeCache.SetWithTTL(req.Email, time.Now(), emailService.config.SendDuration)

	return NewApiResponse(SendEmailSuccess, &ResponseEmailVerifyCode{Email: req.Email})
}
