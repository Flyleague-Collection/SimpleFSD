// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
)

type VerifyCodeEmailData struct {
	Email string
	Cid   int
}

type PermissionChangeEmailData struct {
	User     *operation.User
	Operator *operation.User
}

type RatingChangeEmailData struct {
	User      *operation.User
	Operator  *operation.User
	OldRating string
	NewRating string
}

type KickedFromServerEmailData struct {
	User     *operation.User
	Operator *operation.User
	Reason   string
}

type EmailServiceInterface interface {
	VerifyEmailCode(email string, code string, cid int) error
	SendEmailVerifyCode(req *RequestEmailVerifyCode) *ApiResponse[ResponseEmailVerifyCode]
}

type RequestEmailVerifyCode struct {
	Email string `json:"email"`
	Cid   int    `json:"cid"`
}

type ResponseEmailVerifyCode struct {
	Email string `json:"email"`
}
