// Package service
// 存放工具函数
package service

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
)

type FieldValidator struct {
	Min, Max          int
	ErrShort, ErrLong *ApiStatus
}

func (v *FieldValidator) CheckString(value string) *ApiStatus {
	length := len(value)
	if length > v.Max {
		return v.ErrLong
	}
	if length < v.Min {
		return v.ErrShort
	}
	return nil
}

func (v *FieldValidator) CheckInt(value int) *ApiStatus {
	if value > v.Max {
		return v.ErrLong
	}
	if value < v.Min {
		return v.ErrShort
	}
	return nil
}

var (
	usernameValidator *FieldValidator
	passwordValidator *FieldValidator
	emailValidator    *FieldValidator
	cidValidator      *FieldValidator
)

func InitValidator(config *config.HttpServerLimit) {
	usernameValidator = &FieldValidator{
		Min:      config.UsernameLengthMin,
		Max:      config.UsernameLengthMax,
		ErrShort: NewApiStatus("USERNAME_TOO_SHORT", "用户名过短", BadRequest),
		ErrLong:  NewApiStatus("USERNAME_TOO_LONG", "用户名过长", BadRequest),
	}
	passwordValidator = &FieldValidator{
		Min:      config.PasswordLengthMin,
		Max:      config.PasswordLengthMax,
		ErrShort: NewApiStatus("PASSWORD_TOO_SHORT", "密码长度过短", BadRequest),
		ErrLong:  NewApiStatus("PASSWORD_TOO_LONG", "密码长度过长", BadRequest),
	}
	emailValidator = &FieldValidator{
		Min:      config.EmailLengthMin,
		Max:      config.EmailLengthMax,
		ErrShort: NewApiStatus("EMAIL_TOO_SHORT", "邮箱过短", BadRequest),
		ErrLong:  NewApiStatus("EMAIL_TOO_LONG", "邮箱过长", BadRequest),
	}
	cidValidator = &FieldValidator{
		Min:      config.CidMin,
		Max:      config.CidMax,
		ErrShort: NewApiStatus("CID_TOO_SHORT", "cid过短", BadRequest),
		ErrLong:  NewApiStatus("CID_TOO_LONG", "cid过长", BadRequest),
	}
}
