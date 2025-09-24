// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/labstack/echo/v4"
)

type UserServiceInterface interface {
	UserRegister(req *RequestUserRegister) *ApiResponse[ResponseUserRegister]
	UserLogin(req *RequestUserLogin) *ApiResponse[ResponseUserLogin]
	CheckAvailability(req *RequestUserAvailability) *ApiResponse[ResponseUserAvailability]
	GetCurrentProfile(req *RequestUserCurrentProfile) *ApiResponse[ResponseUserCurrentProfile]
	EditCurrentProfile(req *RequestUserEditCurrentProfile) *ApiResponse[ResponseUserEditCurrentProfile]
	GetUserProfile(req *RequestUserProfile) *ApiResponse[ResponseUserProfile]
	EditUserProfile(req *RequestUserEditProfile) *ApiResponse[ResponseUserEditProfile]
	GetUserList(req *RequestUserList) *ApiResponse[ResponseUserList]
	EditUserPermission(req *RequestUserEditPermission) *ApiResponse[ResponseUserEditPermission]
	GetUserHistory(req *RequestGetUserHistory) *ApiResponse[ResponseGetUserHistory]
	GetTokenWithFlushToken(req *RequestGetToken) *ApiResponse[ResponseGetToken]
}

type RequestUserRegister struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Cid       int    `json:"cid"`
	EmailCode string `json:"email_code"`
}

type ResponseUserRegister bool

type RequestUserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseUserLogin struct {
	User       *operation.User `json:"user"`
	Token      string          `json:"token"`
	FlushToken string          `json:"flush_token"`
}

type RequestUserAvailability struct {
	Username string `query:"username"`
	Email    string `query:"email"`
	Cid      string `query:"cid"`
}

type ResponseUserAvailability bool

type RequestUserCurrentProfile struct {
	JwtHeader
}

type ResponseUserCurrentProfile *operation.User

type RequestUserEditCurrentProfile struct {
	JwtHeader
	ID             uint   `json:"id"`
	Cid            int    `json:"cid"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	EmailCode      string `json:"email_code"`
	AvatarUrl      string `json:"avatar_url"`
	QQ             int    `json:"qq"`
	OriginPassword string `json:"origin_password"`
	NewPassword    string `json:"new_password"`
}

type ResponseUserEditCurrentProfile bool

type RequestUserProfile struct {
	JwtHeader
	TargetUid uint `param:"uid"`
}

type ResponseUserProfile *operation.User

type RequestUserList struct {
	JwtHeader
	Page     int `query:"page_number"`
	PageSize int `query:"page_size"`
}

type ResponseUserList struct {
	Items    []*operation.User `json:"items"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Total    int64             `json:"total"`
}

type RequestUserEditProfile struct {
	JwtHeader
	EchoContentHeader
	TargetUid uint `param:"uid"`
	RequestUserEditCurrentProfile
}

type ResponseUserEditProfile bool

type RequestUserEditPermission struct {
	JwtHeader
	EchoContentHeader
	TargetUid   uint     `param:"uid"`
	Permissions echo.Map `json:"permissions"`
}

type ResponseUserEditPermission bool

type RequestGetUserHistory struct {
	JwtHeader
}

type ResponseGetUserHistory struct {
	*operation.UserHistory
	TotalAtcTime   int `json:"total_atc_time"`
	TotalPilotTime int `json:"total_pilot_time"`
}

type RequestGetToken struct {
	*Claims
	FirstTime bool `query:"first"`
}

type ResponseGetToken struct {
	User       *operation.User `json:"user"`
	Token      string          `json:"token"`
	FlushToken string          `json:"flush_token"`
}
