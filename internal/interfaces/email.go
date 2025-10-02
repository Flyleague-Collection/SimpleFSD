// Package interfaces
package interfaces

import (
	"errors"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	"time"
)

var (
	ErrTemplateNotInitialized = errors.New("error template not initialized")
	ErrRenderingTemplate      = errors.New("error rendering template")
)

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

type ApplicationPassedEmailData struct {
	User     *operation.User
	Operator *operation.User
	Message  string
}

// ApplicationPassedEmail 管制员申请通过
type ApplicationPassedEmail struct {
	Cid      string // 用户CID
	Operator string // 操作者CID
	Message  string // 通过信息
	Contact  string // 操作者邮箱
}

type ApplicationProcessingEmailData struct {
	User           *operation.User
	Operator       *operation.User
	AvailableTimes []time.Time
}

// ApplicationProcessingEmail 管制员申请进度通知
type ApplicationProcessingEmail struct {
	Cid     string // 申请者CID
	Time    string // 可用时间, 例: 2025-09-24 12:00:00 CST, 025-09-25 12:00:00 CST
	Contact string // 回复邮件
}

type ApplicationRejectedEmailData struct {
	User     *operation.User
	Operator *operation.User
	Reason   string
}

// ApplicationRejectedEmail 管制员申请拒绝通知
type ApplicationRejectedEmail struct {
	Cid      string // 申请者CID
	Operator string // 操作者CID
	Reason   string // 拒绝理由
	Contact  string // 操作者邮箱
}

type AtcRatingChangeEmailData struct {
	User      *operation.User
	Operator  *operation.User
	NewRating string
	OldRating string
}

// AtcRatingChangeEmail 管制权限变更
type AtcRatingChangeEmail struct {
	Cid      string // 用户CID
	NewValue string // 操作者CID
	OldValue string // 原管制权限
	Operator string // 新管制权限
	Contact  string // 操作者邮箱
}

type EmailVerifyEmailData struct {
	Email string
	Cid   int
	Code  int
}

// EmailVerifyEmail 邮箱验证码
type EmailVerifyEmail struct {
	Cid       string // 用户注册CID
	Code      string // 验证码
	Email     string // 用户注册邮箱
	Expired   string // 过期时间, 单位为分钟, 比如5
	ExpiredAt string // 过期时间, 时间点, 比如2025-09-25 12:00:00 CST
}

type KickedFromServerEmailData struct {
	User     *operation.User
	Operator *operation.User
	Reason   string
}

// KickedFromServerEmail 踢出服务器通知
type KickedFromServerEmail struct {
	Cid      string // 用户CID
	Time     string // 时间
	Operator string // 操作者CID
	Reason   string // 理由
	Contact  string // 操作者邮箱
}

type PasswordChangeEmailData struct {
	User      *operation.User
	Ip        string
	UserAgent string
}

// PasswordChangeEmail 密码修改通知
type PasswordChangeEmail struct {
	Cid       string // 用户CID
	IP        string // 用户IP
	UserAgent string // 用户UA
	Time      string // 修改时间
}

type PermissionChangeEmailData struct {
	User        *operation.User
	Operator    *operation.User
	Permissions []string
}

// PasswordResetEmail 密码修改通知
type PasswordResetEmail struct {
	Cid       string // 用户CID
	IP        string // 用户IP
	UserAgent string // 用户UA
	Time      string // 修改时间
}

type PasswordResetEmailData struct {
	User      *operation.User
	Ip        string
	UserAgent string
}

// PermissionChangeEmail 飞控权限修改通知
type PermissionChangeEmail struct {
	Cid         string // 用户CID
	Operator    string // 操作者CID
	Permissions string // 受影响权限, 例: ControllerCreateRecord, ControllerChangeGuest
	Contact     string // 操作者邮箱
}

type TicketReplyEmailData struct {
	User  *operation.User
	Title string
	Reply string
}

// TicketReplyEmail 工单回复通知
type TicketReplyEmail struct {
	Cid   string // 用户CID
	Title string // 工单标题
	Reply string // 工单回复内容
}
