// Package service
package service

type EmailCode struct {
	Code int
	Cid  int
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
