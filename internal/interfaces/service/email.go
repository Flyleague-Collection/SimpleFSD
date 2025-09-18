// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
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
	OldRating fsd.Rating
	NewRating fsd.Rating
}

type KickedFromServerEmailData struct {
	User     *operation.User
	Operator *operation.User
	Reason   string
}

type EmailServiceInterface interface {
	VerifyEmailCode(email string, code int, cid int) error
	SendEmailVerifyCode(req *RequestEmailVerifyCode) *ApiResponse[ResponseEmailVerifyCode]
}

type RequestEmailVerifyCode struct {
	Email string `json:"email"`
	Cid   int    `json:"cid"`
}

type ResponseEmailVerifyCode struct {
	Email string `json:"email"`
}
