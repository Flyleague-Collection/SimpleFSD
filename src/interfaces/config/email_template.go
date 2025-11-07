// Package config
package config

import (
	"errors"
	"fmt"
	"html/template"
	"net/url"

	"github.com/half-nothing/simple-fsd/src/interfaces/global"
	"github.com/half-nothing/simple-fsd/src/interfaces/log"
	"golang.org/x/sync/errgroup"
)

type EmailTemplateConfig struct {
	FilePath   string             `json:"file_path"`
	EmailTitle string             `json:"email_title"`
	Template   *template.Template `json:"-"`
	Enable     bool               `json:"enable"`
}

type EmailTemplateConfigs struct {
	VerifyCodeEmail            *EmailTemplateConfig `json:"verify_code_email"`
	ATCRatingChangeEmail       *EmailTemplateConfig `json:"atc_rating_change_email"`
	PermissionChangeEmail      *EmailTemplateConfig `json:"permission_change_email"`
	KickedFromServerEmail      *EmailTemplateConfig `json:"kicked_from_server_email"`
	PasswordChangeEmail        *EmailTemplateConfig `json:"password_change_email"`
	PasswordResetEmail         *EmailTemplateConfig `json:"password_reset_email"`
	ApplicationPassedEmail     *EmailTemplateConfig `json:"application_passed_email"`
	ApplicationRejectedEmail   *EmailTemplateConfig `json:"application_rejected_email"`
	ApplicationProcessingEmail *EmailTemplateConfig `json:"application_processing_email"`
	TicketReplyEmail           *EmailTemplateConfig `json:"ticket_reply_email"`
}

func defaultEmailTemplateConfig() *EmailTemplateConfigs {
	return &EmailTemplateConfigs{
		VerifyCodeEmail: &EmailTemplateConfig{
			FilePath:   "template/email_verify.template",
			EmailTitle: "邮箱验证码",
			Enable:     true,
		},
		ATCRatingChangeEmail: &EmailTemplateConfig{
			FilePath:   "template/atc_rating_change.template",
			EmailTitle: "管制权限变更通知",
			Enable:     true,
		},
		PermissionChangeEmail: &EmailTemplateConfig{
			FilePath:   "template/permission_change.template",
			EmailTitle: "管理权限变更通知",
			Enable:     true,
		},
		KickedFromServerEmail: &EmailTemplateConfig{
			FilePath:   "template/kicked_from_server.template",
			EmailTitle: "踢出服务器通知",
			Enable:     true,
		},
		PasswordChangeEmail: &EmailTemplateConfig{
			FilePath:   "template/password_change.template",
			EmailTitle: "飞控密码更改通知",
			Enable:     true,
		},
		PasswordResetEmail: &EmailTemplateConfig{
			FilePath:   "template/password_reset.template",
			EmailTitle: "飞控密码重置通知",
			Enable:     true,
		},
		ApplicationPassedEmail: &EmailTemplateConfig{
			FilePath:   "template/application_passed.template",
			EmailTitle: "管制员申请通过",
			Enable:     true,
		},
		ApplicationRejectedEmail: &EmailTemplateConfig{
			FilePath:   "template/application_rejected.template",
			EmailTitle: "管制员申请被拒",
			Enable:     true,
		},
		ApplicationProcessingEmail: &EmailTemplateConfig{
			FilePath:   "template/application_processing.template",
			EmailTitle: "管制员申请进度通知",
			Enable:     true,
		},
		TicketReplyEmail: &EmailTemplateConfig{
			FilePath:   "template/ticket_reply.template",
			EmailTitle: "工单回复通知",
			Enable:     true,
		},
	}
}

func validateTemplate(
	logger log.LoggerInterface,
	emailTemplate *EmailTemplateConfig,
	urlPath string,
	tplName string,
	errMsgLoad, errMsgParse string,
) error {
	if !emailTemplate.Enable {
		return nil
	}

	fileUrl, err := url.JoinPath(*global.DownloadPrefix, urlPath)
	if err != nil {
		return ValidFailWith(fmt.Errorf("fail to parse url %s", *global.DownloadPrefix), err)
	}

	bytes, err := cachedContent(logger, emailTemplate.FilePath, fileUrl)
	if err != nil {
		return ValidFailWith(errors.New(errMsgLoad), err)
	}

	parsed, err := template.New(tplName).Parse(string(bytes))
	if err != nil {
		return ValidFailWith(errors.New(errMsgParse), err)
	}

	emailTemplate.Template = parsed
	return nil
}

func (config *EmailTemplateConfigs) checkValid(logger log.LoggerInterface) *ValidResult {
	if !config.VerifyCodeEmail.Enable {
		return ValidFail(errors.New("verify code email can not be disabled"))
	}

	var eg errgroup.Group

	eg.Go(func() error {
		return validateTemplate(
			logger,
			config.VerifyCodeEmail,
			global.EmailVerifyTemplateFilePath,
			"email_verify",
			"fail to load email_verify_template_file",
			"fail to parse email_verify_template",
		)
	})

	eg.Go(func() error {
		return validateTemplate(
			logger,
			config.ATCRatingChangeEmail,
			global.ATCRatingChangeTemplateFilePath,
			"atc_rating_change",
			"fail to load atc_rating_change_template_file",
			"fail to parse atc_rating_change_template",
		)
	})

	eg.Go(func() error {
		return validateTemplate(
			logger,
			config.PermissionChangeEmail,
			global.PermissionChangeTemplateFilePath,
			"permission_change",
			"fail to load permission_change_template_file",
			"fail to parse permission_change_template",
		)
	})

	eg.Go(func() error {
		return validateTemplate(
			logger,
			config.KickedFromServerEmail,
			global.KickedFromServerTemplateFilePath,
			"kicked_from_server",
			"fail to load kicked_from_server_template",
			"fail to parse kicked_from_server_template",
		)
	})

	eg.Go(func() error {
		return validateTemplate(
			logger,
			config.PasswordChangeEmail,
			global.PasswordChangeTemplateFilePath,
			"password_change",
			"fail to load password_change_template",
			"fail to parse password_change_template",
		)
	})

	eg.Go(func() error {
		return validateTemplate(
			logger,
			config.PasswordResetEmail,
			global.PasswordResetTemplateFilePath,
			"password_reset",
			"fail to load password_reset_template",
			"fail to parse password_reset_template",
		)
	})

	eg.Go(func() error {
		return validateTemplate(
			logger,
			config.ApplicationPassedEmail,
			global.ApplicationPassedTemplateFilePath,
			"application_passed",
			"fail to load application_passed_template",
			"fail to parse application_passed_template",
		)
	})

	eg.Go(func() error {
		return validateTemplate(
			logger,
			config.ApplicationRejectedEmail,
			global.ApplicationRejectedTemplateFilePath,
			"application_rejected",
			"fail to load application_rejected_template",
			"fail to parse application_rejected_template",
		)
	})

	eg.Go(func() error {
		return validateTemplate(
			logger,
			config.ApplicationProcessingEmail,
			global.ApplicationProcessingTemplateFilePath,
			"application_processing",
			"fail to load application_processing_template",
			"fail to parse application_processing_template",
		)
	})

	eg.Go(func() error {
		return validateTemplate(
			logger,
			config.TicketReplyEmail,
			global.TicketReplyTemplateFilePath,
			"ticket_reply",
			"fail to load ticket_reply_template",
			"fail to parse ticket_reply_template",
		)
	})

	if err := eg.Wait(); err != nil {
		// 我们这里很确定只会有ValidResult类型的错误
		// 不可能有其他类型的错误, 代码里根本没有返回其他错误
		// 所以这里强制类型转换是安全的
		return err.(*ValidResult)
	}
	return ValidPass()
}
