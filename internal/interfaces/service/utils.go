// Package service
package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/labstack/echo/v4"
	"time"
)

type HttpCode int

const (
	Unsatisfied         HttpCode = 0
	Ok                  HttpCode = 200
	BadRequest          HttpCode = 400
	Unauthorized        HttpCode = 401
	PermissionDenied    HttpCode = 403
	NotFound            HttpCode = 404
	Conflict            HttpCode = 409
	ServerInternalError HttpCode = 500
)

func (hc HttpCode) Code() int {
	return int(hc)
}

type ApiStatus struct {
	StatusName  string
	Description string
	HttpCode    HttpCode
}

func NewApiStatus(statusName, description string, httpCode HttpCode) *ApiStatus {
	return &ApiStatus{
		StatusName:  statusName,
		Description: description,
		HttpCode:    httpCode,
	}
}

type ApiResponse[T any] struct {
	HttpCode int    `json:"-"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Data     *T     `json:"data"`
}

type Claims struct {
	Uid        uint   `json:"uid"`
	Cid        int    `json:"cid"`
	Username   string `json:"username"`
	Permission int64  `json:"permission"`
	Rating     int    `json:"rating"`
	FlushToken bool   `json:"flushToken"`
	config     *config.JWTConfig
	jwt.RegisteredClaims
}

type EchoContentHeader struct {
	Ip        string
	UserAgent string
}

type JwtHeader struct {
	Uid        uint
	Permission int64
	Cid        int
}

func NewClaims(config *config.JWTConfig, user *operation.User, flushToken bool) *Claims {
	expiredDuration := config.ExpiresDuration
	if flushToken {
		expiredDuration += config.RefreshDuration
	}
	return &Claims{
		Uid:        user.ID,
		Cid:        user.Cid,
		Username:   user.Username,
		Permission: user.Permission,
		Rating:     user.Rating,
		FlushToken: flushToken,
		config:     config,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "FsdHttpServer",
			Subject:   user.Username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiredDuration)),
		},
	}
}

func (claim *Claims) GenerateKey() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	tokenString, _ := token.SignedString([]byte(claim.config.Secret))
	return tokenString
}

func (res *ApiResponse[T]) Response(ctx echo.Context) error {
	return ctx.JSON(res.HttpCode, res)
}

var (
	ErrIllegalParam          = NewApiStatus("PARAM_ERROR", "参数不正确", BadRequest)
	ErrLackParam             = NewApiStatus("PARAM_LACK_ERROR", "缺少参数", BadRequest)
	ErrNoPermission          = NewApiStatus("NO_PERMISSION", "无权这么做", PermissionDenied)
	ErrDatabaseFail          = NewApiStatus("DATABASE_ERROR", "服务器内部错误", ServerInternalError)
	ErrUserNotFound          = NewApiStatus("USER_NOT_FOUND", "指定用户不存在", NotFound)
	ErrActivityNotFound      = NewApiStatus("ACTIVITY_NOT_FOUND", "活动不存在", NotFound)
	ErrFlightPlanNotFound    = NewApiStatus("FLIGHT_PLAN_NOT_FOUND", "飞行计划不存在", NotFound)
	ErrFlightPlanLocked      = NewApiStatus("FLIGHT_PLAN_LOCKED", "飞行计划已锁定", Conflict)
	ErrFacilityNotFound      = NewApiStatus("FACILITY_NOT_FOUND", "管制席位不存在", NotFound)
	ErrRegisterFail          = NewApiStatus("REGISTER_FAIL", "注册失败", ServerInternalError)
	ErrIdentifierTaken       = NewApiStatus("USER_EXISTS", "用户已存在", BadRequest)
	ErrMissingOrMalformedJwt = NewApiStatus("MISSING_OR_MALFORMED_JWT", "缺少JWT令牌或者令牌格式错误", BadRequest)
	ErrInvalidOrExpiredJwt   = NewApiStatus("INVALID_OR_EXPIRED_JWT", "无效或过期的JWT令牌", Unauthorized)
	ErrInvalidJwtType        = NewApiStatus("INVALID_JWT_TYPE", "非法的JWT令牌类型", Unauthorized)
	ErrUnknown               = NewApiStatus("UNKNOWN_JWT_ERROR", "未知的JWT解析错误", ServerInternalError)
)

func NewErrorResponse(ctx echo.Context, codeStatus *ApiStatus) error {
	return NewApiResponse[any](codeStatus, nil).Response(ctx)
}

func NewApiResponse[T any](codeStatus *ApiStatus, data *T) *ApiResponse[T] {
	return &ApiResponse[T]{
		HttpCode: codeStatus.HttpCode.Code(),
		Code:     codeStatus.StatusName,
		Message:  codeStatus.Description,
		Data:     data,
	}
}

func checkDatabaseError[T any](err error) *ApiResponse[T] {
	switch {
	case errors.Is(err, operation.ErrIdentifierCheck):
		return NewApiResponse[T](ErrRegisterFail, nil)
	case errors.Is(err, operation.ErrIdentifierTaken):
		return NewApiResponse[T](ErrIdentifierTaken, nil)
	case errors.Is(err, operation.ErrUserNotFound):
		return NewApiResponse[T](ErrUserNotFound, nil)
	case errors.Is(err, operation.ErrActivityNotFound):
		return NewApiResponse[T](ErrActivityNotFound, nil)
	case errors.Is(err, operation.ErrFlightPlanNotFound):
		return NewApiResponse[T](ErrFlightPlanNotFound, nil)
	case errors.Is(err, operation.ErrFacilityNotFound):
		return NewApiResponse[T](ErrFacilityNotFound, nil)
	case err != nil:
		return NewApiResponse[T](ErrDatabaseFail, nil)
	default:
		return nil
	}
}

func CheckPermission[T any](permission int64, perm operation.Permission) *ApiResponse[T] {
	userPermission := operation.Permission(permission)
	if userPermission.HasPermission(perm) {
		return NewApiResponse[T](ErrNoPermission, nil)
	}
	return nil
}

// CallDBFunc 调用数据库操作函数并处理错误
func CallDBFunc[R any, T any](fc func() (*R, error)) (result *R, response *ApiResponse[T]) {
	result, err := fc()
	response = checkDatabaseError[T](err)
	return
}

func CallDBFuncWithoutRet[T any](fc func() error) *ApiResponse[T] {
	err := fc()
	return checkDatabaseError[T](err)
}

// GetUsersAndCheckPermission 从数据库获取用户数据并检查权限
func GetUsersAndCheckPermission[T any](
	userOperation operation.UserOperationInterface,
	uid uint,
	targetUid uint,
	perm operation.Permission,
) (user *operation.User, targetUser *operation.User, response *ApiResponse[T]) {
	// 敏感操作获取实时数据
	user, response = CallDBFunc[operation.User, T](func() (*operation.User, error) { return userOperation.GetUserByUid(uid) })
	if response != nil {
		return
	}
	if response = CheckPermission[T](user.Permission, perm); response != nil {
		return
	}
	targetUser, response = CallDBFunc[operation.User, T](func() (*operation.User, error) { return userOperation.GetUserByUid(targetUid) })
	return
}

func GetUserAndCheckPermission[T any](
	userOperation operation.UserOperationInterface,
	permission int64,
	targetUid uint,
	perm operation.Permission,
) (*operation.User, *ApiResponse[T]) {
	if res := CheckPermission[T](permission, perm); res != nil {
		return nil, res
	}
	return CallDBFunc[operation.User, T](func() (*operation.User, error) { return userOperation.GetUserByUid(targetUid) })
}
