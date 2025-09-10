// Package database
package database

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"golang.org/x/crypto/bcrypt"
)

type UserOperation struct {
	logger log.LoggerInterface
	config *config.Config
}

func NewUserOperation(logger log.LoggerInterface, config *config.Config) *UserOperation {
	return &UserOperation{
		logger: logger,
		config: config,
	}
}

func (userOperation *UserOperation) GetUserByCid(cid string) (user *User, err error) {
	user, exist := data[cid]
	if !exist {
		return nil, ErrUserNotFound
	}
	return
}

func (userOperation *UserOperation) VerifyUserPassword(user *User, password string) (pass bool) {
	switch userOperation.config.EncryptionType {
	case 0:
		// 明文
		return user.Password == password
	case 1:
		// MD5
		hashValue := md5.Sum([]byte(password))
		return bytes.Equal(hashValue[:], []byte(user.Password))
	case 2:
		// SHA256
		hashValue := sha256.Sum256([]byte(password))
		return bytes.Equal(hashValue[:], []byte(user.Password))
	case 3:
		// bcrypt
		return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
	default:
		return false
	}
}
