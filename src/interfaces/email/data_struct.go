// Package email
package email

import (
	"errors"
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

var (
	ErrTemplateNotInitialized = errors.New("error template not initialized")
	ErrRenderingTemplate      = errors.New("error rendering template")
)

type ApplicationPassedEmailData struct {
	User     *entity.User
	Operator *entity.User
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
	User           *entity.User
	Operator       *entity.User
	AvailableTimes []time.Time
}

// ApplicationProcessingEmail 管制员申请进度通知
type ApplicationProcessingEmail struct {
	Cid     string // 申请者CID
	Time    string // 可用时间, 例: 2025-09-24 12:00:00 CST, 025-09-25 12:00:00 CST
	Contact string // 回复邮件
}

type ApplicationRejectedEmailData struct {
	User     *entity.User
	Operator *entity.User
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
	User      *entity.User
	Operator  *entity.User
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
	User     *entity.User
	Operator *entity.User
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
	User      *entity.User
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
	User        *entity.User
	Operator    *entity.User
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
	User      *entity.User
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
	User  *entity.User
	Title string
	Reply string
}

// TicketReplyEmail 工单回复通知
type TicketReplyEmail struct {
	Cid   string // 用户CID
	Title string // 工单标题
	Reply string // 工单回复内容
}
