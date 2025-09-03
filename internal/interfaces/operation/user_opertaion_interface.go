// Package operation
package operation

import (
	"errors"
)

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user does not exist")
)

// UserOperationInterface 用户操作接口定义
type UserOperationInterface interface {
	// GetUserByCid 通过Cid获取用户, 当err为nil时返回值user有效
	GetUserByCid(cid string) (user *User, err error)
	// VerifyUserPassword 验证用户密码是否正确, pass为true表示验证通过
	VerifyUserPassword(user *User, password string) (pass bool)
}
