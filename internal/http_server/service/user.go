// Package service
// 存放 UserServiceInterface 的实现
package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/half-nothing/simple-fsd/internal/utils"
	"strings"
	"time"
)

type UserService struct {
	logger            log.LoggerInterface
	config            *config.HttpServerConfig
	messageQueue      queue.MessageQueueInterface
	emailService      EmailServiceInterface
	userOperation     operation.UserOperationInterface
	historyOperation  operation.HistoryOperationInterface
	storeService      StoreServiceInterface
	auditLogOperation operation.AuditLogOperationInterface
}

func NewUserService(
	logger log.LoggerInterface,
	config *config.HttpServerConfig,
	messageQueue queue.MessageQueueInterface,
	userOperation operation.UserOperationInterface,
	historyOperation operation.HistoryOperationInterface,
	auditLogOperation operation.AuditLogOperationInterface,
	storeService StoreServiceInterface,
	emailService EmailServiceInterface,
) *UserService {
	return &UserService{
		logger:            logger,
		messageQueue:      messageQueue,
		emailService:      emailService,
		config:            config,
		userOperation:     userOperation,
		historyOperation:  historyOperation,
		storeService:      storeService,
		auditLogOperation: auditLogOperation,
	}
}

var (
	ErrEmailNotFound    = NewApiStatus("EMAIL_CODE_NOT_FOUND", "未向该邮箱发送验证码", BadRequest)
	ErrCidNotMatch      = NewApiStatus("CID_NOT_MATCH", "注册cid与验证码发送时的cid不一致", BadRequest)
	ErrEmailExpired     = NewApiStatus("EMAIL_CODE_EXPIRED", "验证码已过期", BadRequest)
	ErrEmailCodeInvalid = NewApiStatus("EMAIL_CODE_INVALID", "邮箱验证码错误", BadRequest)
	SuccessRegister     = NewApiStatus("REGISTER_SUCCESS", "注册成功", Ok)
)

func (userService *UserService) verifyEmailCode(email string, emailCode, cid int) *ApiStatus {
	err := userService.emailService.VerifyEmailCode(email, emailCode, cid)
	switch {
	case errors.Is(err, ErrEmailCodeNotFound):
		return ErrEmailNotFound
	case errors.Is(err, ErrEmailCodeExpired):
		return ErrEmailExpired
	case errors.Is(err, ErrInvalidEmailCode):
		return ErrEmailCodeInvalid
	case errors.Is(err, ErrCidMismatch):
		return ErrCidNotMatch
	default:
		return nil
	}
}

func (userService *UserService) UserRegister(req *RequestUserRegister) *ApiResponse[ResponseUserRegister] {
	if req.Username == "" || req.Email == "" || req.Password == "" || req.Cid <= 0 || req.EmailCode <= 0 {
		return NewApiResponse[ResponseUserRegister](ErrIllegalParam, nil)
	}

	if err := usernameValidator.CheckString(req.Username); err != nil {
		return NewApiResponse[ResponseUserRegister](err, nil)
	}

	if err := emailValidator.CheckString(req.Email); err != nil {
		return NewApiResponse[ResponseUserRegister](err, nil)
	}

	if err := passwordValidator.CheckString(req.Password); err != nil {
		return NewApiResponse[ResponseUserRegister](err, nil)
	}

	if err := cidValidator.CheckInt(req.Cid); err != nil {
		return NewApiResponse[ResponseUserRegister](err, nil)
	}

	if res := userService.verifyEmailCode(req.Email, req.EmailCode, req.Cid); res != nil {
		return NewApiResponse[ResponseUserRegister](res, nil)
	}

	user, err := userService.userOperation.NewUser(req.Username, req.Email, req.Cid, req.Password)
	if res := CheckDatabaseError[ResponseUserRegister](err); res != nil {
		return res
	}

	if res := CallDBFuncWithoutRet[ResponseUserRegister](func() error {
		return userService.userOperation.AddUser(user)
	}); res != nil {
		return res
	}

	data := ResponseUserRegister(true)
	return NewApiResponse(SuccessRegister, &data)
}

var (
	ErrWrongUsernameOrPassword = NewApiStatus("WRONG_USERNAME_OR_PASSWORD", "用户名或密码错误", NotFound)
	SuccessLogin               = NewApiStatus("LOGIN_SUCCESS", "登陆成功", Ok)
)

func (userService *UserService) UserLogin(req *RequestUserLogin) *ApiResponse[ResponseUserLogin] {
	if req.Username == "" || req.Password == "" {
		return NewApiResponse[ResponseUserLogin](ErrIllegalParam, nil)
	}

	userId := operation.GetUserId(req.Username)

	user, res := CallDBFunc[*operation.User, ResponseUserLogin](func() (*operation.User, error) {
		return userId.GetUser(userService.userOperation)
	})
	if res != nil {
		return res
	}

	if pass := userService.userOperation.VerifyUserPassword(user, req.Password); !pass {
		return NewApiResponse[ResponseUserLogin](ErrWrongUsernameOrPassword, nil)
	}

	token := NewClaims(userService.config.JWT, user, false)
	flushToken := NewClaims(userService.config.JWT, user, true)
	return NewApiResponse(SuccessLogin, &ResponseUserLogin{
		User:       user,
		Token:      token.GenerateKey(),
		FlushToken: flushToken.GenerateKey(),
	})
}

var (
	NameNotAvailability = NewApiStatus("INFO_NOT_AVAILABILITY", "用户信息不可用", Ok)
	NameAvailability    = NewApiStatus("INFO_AVAILABILITY", "用户信息可用", Ok)
)

func (userService *UserService) CheckAvailability(req *RequestUserAvailability) *ApiResponse[ResponseUserAvailability] {
	if req.Username == "" && req.Email == "" && req.Cid == "" {
		return NewApiResponse[ResponseUserAvailability](ErrIllegalParam, nil)
	}

	exist, err := userService.userOperation.IsUserIdentifierTaken(nil, utils.StrToInt(req.Cid, 0), req.Username, req.Email)
	if res := CheckDatabaseError[ResponseUserAvailability](err); res != nil {
		return res
	}

	data := ResponseUserAvailability(!exist)
	if exist {
		return NewApiResponse(NameNotAvailability, &data)
	}
	return NewApiResponse(NameAvailability, &data)
}

var SuccessGetCurrentProfile = NewApiStatus("GET_CURRENT_PROFILE_SUCCESS", "获取当前用户信息成功", Ok)

func (userService *UserService) GetCurrentProfile(req *RequestUserCurrentProfile) *ApiResponse[ResponseUserCurrentProfile] {
	user, res := CallDBFunc[*operation.User, ResponseUserCurrentProfile](func() (*operation.User, error) {
		return userService.userOperation.GetUserByUid(req.Uid)
	})
	if res != nil {
		return res
	}

	data := ResponseUserCurrentProfile(user)
	return NewApiResponse(SuccessGetCurrentProfile, &data)
}

var (
	ErrOriginPasswordRequired = NewApiStatus("ORIGIN_PASSWORD_REQUIRED", "未提供原始密码", BadRequest)
	ErrNewPasswordRequired    = NewApiStatus("NEW_PASSWORD_REQUIRED", "未提供新密码", BadRequest)
	ErrWrongOriginPassword    = NewApiStatus("WRONG_ORIGIN_PASSWORD_ERROR", "原始密码不正确", BadRequest)
	ErrQQInvalid              = NewApiStatus("QQ_INVALID", "qq号不正确", BadRequest)
	SuccessEditCurrentProfile = NewApiStatus("SUCCESS_EDIT_CURRENT_PROFILE", "编辑用户信息成功", Ok)
)

func checkQQ(qq int) *ApiStatus {
	// QQ 号码应当在 10000 - 100000000000之间
	if 1e4 <= qq && qq < 1e11 {
		return nil
	}
	return ErrQQInvalid
}

func (userService *UserService) editUserProfile(req *RequestUserEditCurrentProfile, skipEmailVerify bool, skipPasswordVerify bool) (*ApiStatus, *operation.User, string) {
	if req.Username == "" && req.Email == "" && req.QQ <= 0 && req.OriginPassword == "" && req.NewPassword == "" && req.AvatarUrl == "" {
		return ErrIllegalParam, nil, ""
	}

	if req.OriginPassword != "" && req.NewPassword != "" {
		if err := passwordValidator.CheckString(req.NewPassword); err != nil {
			return err, nil, ""
		}
	} else if req.OriginPassword != "" && req.NewPassword == "" {
		return ErrNewPasswordRequired, nil, ""
	} else if req.OriginPassword == "" && req.NewPassword != "" && !skipPasswordVerify {
		return ErrOriginPasswordRequired, nil, ""
	}

	if req.Username != "" {
		if err := usernameValidator.CheckString(req.Username); err != nil {
			return err, nil, ""
		}
	}

	if req.Email != "" {
		if err := emailValidator.CheckString(req.Email); err != nil {
			return err, nil, ""
		}
		if !skipEmailVerify {
			if req.EmailCode <= 0 {
				return ErrIllegalParam, nil, ""
			}
			if res := userService.verifyEmailCode(req.Email, req.EmailCode, req.Cid); res != nil {
				return res, nil, ""
			}
		}
	}

	if req.QQ > 0 {
		if err := checkQQ(req.QQ); err != nil {
			return err, nil, ""
		}
	}

	user, err := userService.userOperation.GetUserByUid(req.ID)
	if errors.Is(err, operation.ErrUserNotFound) {
		return ErrUserNotFound, nil, ""
	} else if err != nil {
		return ErrDatabaseFail, nil, ""
	}

	updateInfo := &operation.User{}

	oldValue, _ := json.Marshal(user)

	if req.Username != "" || req.Email != "" {
		exist, _ := userService.userOperation.IsUserIdentifierTaken(nil, 0, req.Username, req.Email)
		if exist {
			return ErrIdentifierTaken, nil, ""
		}

		if req.Username != "" && req.Username != user.Username {
			user.Username = req.Username
			updateInfo.Username = req.Username
		}

		if req.Email != "" && req.Email != user.Email {
			user.Email = req.Email
			updateInfo.Email = req.Email
		}
	}

	if req.QQ > 0 && req.QQ != user.QQ {
		user.QQ = req.QQ
		updateInfo.QQ = req.QQ
		if req.AvatarUrl == "" && (user.AvatarUrl == "" || strings.HasPrefix(user.AvatarUrl, "https://q2.qlogo.cn/")) {
			user.AvatarUrl = fmt.Sprintf("https://q2.qlogo.cn/headimg_dl?dst_uin=%d&spec=100", user.QQ)
			updateInfo.AvatarUrl = user.AvatarUrl
		}
	}

	if req.AvatarUrl != "" {
		if user.AvatarUrl != "" && !strings.HasPrefix(user.AvatarUrl, "https://q2.qlogo.cn/") {
			_, err = userService.storeService.DeleteImageFile(user.AvatarUrl)
			if err != nil {
				userService.logger.ErrorF("err while delete user old avatar, %v", err)
			}
		}
		user.AvatarUrl = req.AvatarUrl
		updateInfo.AvatarUrl = user.AvatarUrl
	}

	if req.OriginPassword != "" || (skipPasswordVerify && req.NewPassword != "") {
		password, err := userService.userOperation.UpdateUserPassword(user, req.OriginPassword, req.NewPassword, skipPasswordVerify)
		if errors.Is(err, operation.ErrPasswordEncode) {
			return ErrUnknownServerError, nil, ""
		} else if errors.Is(err, operation.ErrOldPassword) {
			return ErrWrongOriginPassword, nil, ""
		} else if err != nil {
			return ErrDatabaseFail, nil, ""
		}
		updateInfo.Password = string(password)
	}

	if err := userService.userOperation.UpdateUserInfo(user, updateInfo); err != nil {
		if errors.Is(err, operation.ErrUserNotFound) {
			return ErrUserNotFound, nil, ""
		} else {
			return ErrDatabaseFail, nil, ""
		}
	}

	return nil, user, string(oldValue)
}

func (userService *UserService) EditCurrentProfile(req *RequestUserEditCurrentProfile) *ApiResponse[ResponseUserEditCurrentProfile] {
	err, _, _ := userService.editUserProfile(req, false, false)
	if err != nil {
		return NewApiResponse[ResponseUserEditCurrentProfile](err, nil)
	}
	data := ResponseUserEditCurrentProfile(true)
	return NewApiResponse(SuccessEditCurrentProfile, &data)
}

var SuccessGetProfile = NewApiStatus("GET_PROFILE_SUCCESS", "获取用户信息成功", Ok)

func (userService *UserService) GetUserProfile(req *RequestUserProfile) *ApiResponse[ResponseUserProfile] {
	if req.TargetUid <= 0 {
		return NewApiResponse[ResponseUserProfile](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseUserProfile](req.Permission, operation.UserGetProfile); res != nil {
		return res
	}

	user, res := CallDBFunc[*operation.User, ResponseUserProfile](func() (*operation.User, error) {
		return userService.userOperation.GetUserByUid(req.TargetUid)
	})
	if res != nil {
		return res
	}

	data := ResponseUserProfile(user)
	return NewApiResponse(SuccessGetProfile, &data)
}

var SuccessEditUserProfile = NewApiStatus("EDIT_USER_PROFILE", "修改用户信息成功", Ok)

func (userService *UserService) EditUserProfile(req *RequestUserEditProfile) *ApiResponse[ResponseUserEditProfile] {
	if req.TargetUid <= 0 {
		return NewApiResponse[ResponseUserEditProfile](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseUserEditProfile](req.Permission, operation.UserEditBaseInfo); res != nil {
		return res
	}

	permission := operation.Permission(req.Permission)

	if req.NewPassword != "" && !permission.HasPermission(operation.UserSetPassword) {
		return NewApiResponse[ResponseUserEditProfile](ErrNoPermission, nil)
	}

	req.RequestUserEditCurrentProfile.ID = req.TargetUid
	err, user, oldValue := userService.editUserProfile(&req.RequestUserEditCurrentProfile, true, true)
	if err != nil {
		return NewApiResponse[ResponseUserEditProfile](err, nil)
	}

	newValue, _ := json.Marshal(user)
	object := fmt.Sprintf("%04d", user.Cid)
	if req.NewPassword != "" {
		object += fmt.Sprintf("(%s)", req.NewPassword)
	}
	userService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: userService.auditLogOperation.NewAuditLog(
			operation.UserInformationEdit,
			req.JwtHeader.Cid,
			object,
			req.Ip,
			req.UserAgent,
			&operation.ChangeDetail{
				OldValue: oldValue,
				NewValue: string(newValue),
			},
		),
	})

	data := ResponseUserEditProfile(true)
	return NewApiResponse(SuccessEditUserProfile, &data)
}

var SuccessGetUsers = NewApiStatus("GET_USER_PAGE", "获取用户信息分页成功", Ok)

func (userService *UserService) GetUserList(req *RequestUserList) *ApiResponse[ResponseUserList] {
	if req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseUserList](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseUserList](req.Permission, operation.UserShowList); res != nil {
		return res
	}

	users, total, err := userService.userOperation.GetUsers(req.Page, req.PageSize)
	if res := CheckDatabaseError[ResponseUserList](err); res != nil {
		return res
	}

	return NewApiResponse(SuccessGetUsers, &ResponseUserList{
		Items:    users,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	})
}

var (
	ErrPermissionNodeNotExists = NewApiStatus("PERMISSION_NODE_NOT_EXISTS", "无效权限节点", BadRequest)
	SuccessEditUserPermission  = NewApiStatus("EDIT_USER_PERMISSION", "编辑用户权限成功", Ok)
)

func (userService *UserService) EditUserPermission(req *RequestUserEditPermission) *ApiResponse[ResponseUserEditPermission] {
	if req.TargetUid <= 0 || len(req.Permissions) == 0 {
		return NewApiResponse[ResponseUserEditPermission](ErrIllegalParam, nil)
	}

	user, targetUser, res := GetTargetUserAndCheckPermissionFromDatabase[ResponseUserEditPermission](
		userService.userOperation,
		req.Uid,
		req.TargetUid,
		operation.UserEditPermission,
	)
	if res != nil {
		return res
	}

	permission := operation.Permission(user.Permission)
	targetPermission := operation.Permission(targetUser.Permission)
	auditLogs := make([]*operation.AuditLog, 0, len(req.Permissions))

	for key, value := range req.Permissions {
		if per, ok := operation.PermissionMap[key]; ok {
			if !permission.HasPermission(per) {
				return NewApiResponse[ResponseUserEditPermission](ErrNoPermission, nil)
			}
			if value, ok := value.(bool); ok {
				if value {
					targetPermission.Grant(per)
					auditLogs = append(auditLogs,
						userService.auditLogOperation.NewAuditLog(
							operation.UserPermissionGrant,
							req.Cid,
							fmt.Sprintf("%04d(%s)", targetUser.Cid, key),
							req.Ip,
							req.UserAgent,
							nil,
						))
				} else {
					targetPermission.Revoke(per)
					auditLogs = append(auditLogs,
						userService.auditLogOperation.NewAuditLog(
							operation.UserPermissionRevoke,
							req.Cid,
							fmt.Sprintf("%04d(%s)", targetUser.Cid, key),
							req.Ip,
							req.UserAgent,
							nil,
						))
				}
			} else {
				return NewApiResponse[ResponseUserEditPermission](ErrIllegalParam, nil)
			}
		} else {
			return NewApiResponse[ResponseUserEditPermission](ErrPermissionNodeNotExists, nil)
		}
	}

	if res := CallDBFuncWithoutRet[ResponseUserEditPermission](func() error {
		return userService.userOperation.UpdateUserPermission(targetUser, targetPermission)
	}); res != nil {
		return res
	}

	if userService.config.Email.Template.EnablePermissionChangeEmail {
		userService.messageQueue.Publish(&queue.Message{
			Type: queue.SendPermissionChangeEmail,
			Data: &PermissionChangeEmailData{
				User:     targetUser,
				Operator: user,
			},
		})
	}

	userService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLogs,
		Data: auditLogs,
	})

	data := ResponseUserEditPermission(true)
	return NewApiResponse(SuccessEditUserPermission, &data)
}

var SuccessGetUserHistory = NewApiStatus("GET_USER_HISTORY", "成功获取用户历史数据", Ok)

func (userService *UserService) GetUserHistory(req *RequestGetUserHistory) *ApiResponse[ResponseGetUserHistory] {
	user, res := CallDBFunc[*operation.User, ResponseGetUserHistory](func() (*operation.User, error) {
		return userService.userOperation.GetUserByCid(req.Cid)
	})
	if res != nil {
		return res
	}

	userHistory, res := CallDBFunc[*operation.UserHistory, ResponseGetUserHistory](func() (*operation.UserHistory, error) {
		return userService.historyOperation.GetUserHistory(req.Cid)
	})
	if res != nil {
		return res
	}

	return NewApiResponse(SuccessGetUserHistory, &ResponseGetUserHistory{
		TotalPilotTime: user.TotalPilotTime,
		TotalAtcTime:   user.TotalAtcTime,
		UserHistory:    userHistory,
	})
}

var SuccessGetToken = NewApiStatus("GET_TOKEN", "成功刷新秘钥", Ok)

func (userService *UserService) GetTokenWithFlushToken(req *RequestGetToken) *ApiResponse[ResponseGetToken] {
	if !req.FlushToken {
		return NewApiResponse[ResponseGetToken](ErrIllegalParam, nil)
	}

	user, res := CallDBFunc[*operation.User, ResponseGetToken](func() (*operation.User, error) {
		return userService.userOperation.GetUserByUid(req.Uid)
	})
	if res != nil {
		return res
	}

	var flushToken string
	if !req.FirstTime && req.ExpiresAt.Add(-2*userService.config.JWT.ExpiresDuration).After(time.Now()) {
		flushToken = ""
	} else {
		flushToken = NewClaims(userService.config.JWT, user, true).GenerateKey()
	}

	token := NewClaims(userService.config.JWT, user, false)
	return NewApiResponse(SuccessGetToken, &ResponseGetToken{
		Token:      token.GenerateKey(),
		FlushToken: flushToken,
	})
}
