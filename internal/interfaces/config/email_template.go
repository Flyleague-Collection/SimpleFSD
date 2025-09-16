// Package config
package config

import (
	"errors"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"golang.org/x/sync/errgroup"
	"html/template"
	"net/url"
)

type EmailTemplateConfig struct {
	EmailVerifyTemplateFile      string             `json:"email_verify_template_file"`
	EmailVerifyTemplate          *template.Template `json:"-"`
	ATCRatingChangeTemplateFile  string             `json:"atc_rating_change_template_file"`
	ATCRatingChangeTemplate      *template.Template `json:"-"`
	EnableRatingChangeEmail      bool               `json:"enable_rating_change_email"`
	PermissionChangeTemplateFile string             `json:"permission_change_template_file"`
	PermissionChangeTemplate     *template.Template `json:"-"`
	EnablePermissionChangeEmail  bool               `json:"enable_permission_change_email"`
	KickedFromServerTemplateFile string             `json:"kicked_from_server_template_file"`
	KickedFromServerTemplate     *template.Template `json:"-"`
	EnableKickedFromServerEmail  bool               `json:"enable_kicked_from_server_email"`
	PasswordChangeTemplateFile   string             `json:"password_change_template_file"`
	PasswordChangeTemplate       *template.Template `json:"-"`
	EnablePasswordChangeEmail    bool               `json:"enable_password_change_email"`
}

func defaultEmailTemplateConfig() *EmailTemplateConfig {
	return &EmailTemplateConfig{
		EmailVerifyTemplateFile:      "template/email_verify.template",
		ATCRatingChangeTemplateFile:  "template/atc_rating_change.template",
		EnableRatingChangeEmail:      true,
		PermissionChangeTemplateFile: "template/permission_change.template",
		EnablePermissionChangeEmail:  true,
		KickedFromServerTemplateFile: "template/kicked_from_server.template",
		EnableKickedFromServerEmail:  true,
		PasswordChangeTemplateFile:   "template/password_change.template",
		EnablePasswordChangeEmail:    true,
	}
}

func validateTemplate(
	logger log.LoggerInterface,
	enable bool,
	filePath, urlPath string,
	tplName string,
	setter func(*template.Template),
	errMsgLoad, errMsgParse string,
) error {
	if !enable {
		return nil
	}

	fileUrl, err := url.JoinPath(*global.DownloadPrefix, urlPath)
	if err != nil {
		return ValidFailWith(fmt.Errorf("fail to parse url %s", *global.DownloadPrefix), err)
	}

	bytes, err := cachedContent(logger, filePath, fileUrl)
	if err != nil {
		return ValidFailWith(errors.New(errMsgLoad), err)
	}

	parsed, err := template.New(tplName).Parse(string(bytes))
	if err != nil {
		return ValidFailWith(errors.New(errMsgParse), err)
	}

	setter(parsed)
	return nil
}

func (config *EmailTemplateConfig) checkValid(logger log.LoggerInterface) *ValidResult {
	var eg errgroup.Group

	eg.Go(func() error {
		return validateTemplate(
			logger,
			true,
			config.EmailVerifyTemplateFile,
			global.EmailVerifyTemplateFilePath,
			"email_verify",
			func(t *template.Template) { config.EmailVerifyTemplate = t },
			"fail to load email_verify_template_file",
			"fail to parse email_verify_template",
		)
	})

	eg.Go(func() error {
		return validateTemplate(
			logger,
			config.EnableRatingChangeEmail,
			config.ATCRatingChangeTemplateFile,
			global.ATCRatingChangeTemplateFilePath,
			"atc_rating_change",
			func(t *template.Template) { config.ATCRatingChangeTemplate = t },
			"fail to load atc_rating_change_template_file",
			"fail to parse atc_rating_change_template",
		)
	})

	eg.Go(func() error {
		return validateTemplate(
			logger,
			config.EnablePermissionChangeEmail,
			config.PermissionChangeTemplateFile,
			global.PermissionChangeTemplateFilePath,
			"permission_change",
			func(t *template.Template) { config.PermissionChangeTemplate = t },
			"fail to load permission_change_template_file",
			"fail to parse permission_change_template",
		)
	})

	eg.Go(func() error {
		return validateTemplate(
			logger,
			config.EnableKickedFromServerEmail,
			config.KickedFromServerTemplateFile,
			global.KickedFromServerTemplateFilePath,
			"kicked_from_server",
			func(t *template.Template) { config.KickedFromServerTemplate = t },
			"fail to load kicked_from_server_template",
			"fail to parse kicked_from_server_template",
		)
	})

	eg.Go(func() error {
		return validateTemplate(
			logger,
			config.EnablePasswordChangeEmail,
			config.PasswordChangeTemplateFile,
			global.PasswordChangeTemplateFilePath,
			"password_change",
			func(t *template.Template) { config.PasswordChangeTemplate = t },
			"fail to load password_change_template",
			"fail to parse password_change_template",
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
