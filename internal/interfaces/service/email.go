// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"html/template"
	"time"
)

type SendEmailCodeData struct {
	Email string
	Cid   int
}

type SendPermissionChangeData struct {
	User     *operation.User
	Operator *operation.User
}

type SendRatingChangeData struct {
	User      *operation.User
	Operator  *operation.User
	OldRating fsd.Rating
	NewRating fsd.Rating
}

type SendKickedFromServerData struct {
	User     *operation.User
	Operator *operation.User
	Reason   string
}

type EmailServiceInterface interface {
	RenderTemplate(template *template.Template, data interface{}) (string, error)
	VerifyCode(email string, code int, cid int) error
	SendEmailCode(data *SendEmailCodeData) (error, time.Duration)
	SendPermissionChangeEmail(data *SendPermissionChangeData) error
	SendRatingChangeEmail(data *SendRatingChangeData) error
	SendKickedFromServerEmail(data *SendKickedFromServerData) error
	SendEmailVerifyCode(req *RequestEmailVerifyCode) *ApiResponse[ResponseEmailVerifyCode]
}

type RequestEmailVerifyCode struct {
	Email string `json:"email"`
	Cid   int    `json:"cid"`
}

type ResponseEmailVerifyCode struct {
	Email string `json:"email"`
}
