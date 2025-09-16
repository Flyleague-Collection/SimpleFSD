// Package service
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
	messageQueue      queue.MessageQueueInterface
	emailService      EmailServiceInterface
	config            *config.HttpServerConfig
	userOperation     operation.UserOperationInterface
	historyOperation  operation.HistoryOperationInterface
	storeService      StoreServiceInterface
	auditLogOperation operation.AuditLogOperationInterface
}

func NewUserService(
	logger log.LoggerInterface,
	config *config.HttpServerConfig,
	userOperation operation.UserOperationInterface,
	historyOperation operation.HistoryOperationInterface,
	auditLogOperation operation.AuditLogOperationInterface,
	storeService StoreServiceInterface,
	emailService EmailServiceInterface,
	messageQueue queue.MessageQueueInterface,
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
	ErrEmailNotFound    = ApiStatus{StatusName: "EMAIL_CODE_NOT_FOUND", Description: "未向该邮箱发送验证码", HttpCode: BadRequest}
	ErrCidNotMatch      = ApiStatus{StatusName: "CID_NOT_MATCH", Description: "注册cid与验证码发送时的cid不一致", HttpCode: BadRequest}
	ErrEmailExpired     = ApiStatus{StatusName: "EMAIL_CODE_EXPIRED", Description: "验证码已过期", HttpCode: BadRequest}
	ErrEmailCodeInvalid = ApiStatus{StatusName: "EMAIL_CODE_INVALID", Description: "邮箱验证码错误", HttpCode: BadRequest}
	SuccessRegister     = ApiStatus{StatusName: "REGISTER_SUCCESS", Description: "注册成功", HttpCode: Ok}
)

func (userService *UserService) verifyEmailCode(email string, emailCode, cid int) *ApiStatus {
	err := userService.emailService.VerifyCode(email, emailCode, cid)
	switch {
	case errors.Is(err, ErrEmailCodeNotFound):
		return &ErrEmailNotFound
	case errors.Is(err, ErrEmailCodeExpired):
		return &ErrEmailExpired
	case errors.Is(err, ErrInvalidEmailCode):
		return &ErrEmailCodeInvalid
	case errors.Is(err, ErrCidMismatch):
		return &ErrCidNotMatch
	default:
		return nil
	}
}

func (userService *UserService) UserRegister(req *RequestUserRegister) *ApiResponse[ResponseUserRegister] {
	if req.Username == "" || req.Email == "" || req.Password == "" || req.Cid <= 0 || req.EmailCode < 1e5 {
		return NewApiResponse[ResponseUserRegister](ErrIllegalParam, nil)
	}
	if res := userService.verifyEmailCode(req.Email, req.EmailCode, req.Cid); res != nil {
		return NewApiResponse[ResponseUserRegister](res, nil)
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
	user, err := userService.userOperation.NewUser(req.Username, req.Email, req.Cid, req.Password)
	if err != nil {
		return NewApiResponse[ResponseUserRegister](ErrRegisterFail, nil)
	}
	if _, res := CallDBFunc[interface{}, ResponseUserRegister](func() (*interface{}, error) {
		return nil, userService.userOperation.AddUser(user)
	}); res != nil {
		return res
	}
	token := NewClaims(userService.config.JWT, user, false)
	flushToken := NewClaims(userService.config.JWT, user, true)
	return NewApiResponse(&SuccessRegister, &ResponseUserRegister{
		User:       user,
		Token:      token.GenerateKey(),
		FlushToken: flushToken.GenerateKey(),
	})
}

var (
	ErrUsernameOrPassword = ApiStatus{StatusName: "WRONG_USERNAME_OR_PASSWORD", Description: "用户名或密码错误", HttpCode: NotFound}
	SuccessLogin          = ApiStatus{StatusName: "LOGIN_SUCCESS", Description: "登陆成功", HttpCode: Ok}
)

func (userService *UserService) UserLogin(req *RequestUserLogin) *ApiResponse[ResponseUserLogin] {
	if req.Username == "" || req.Password == "" {
		return NewApiResponse[ResponseUserLogin](ErrIllegalParam, nil)
	}
	userId := operation.GetUserId(req.Username)

	user, res := CallDBFunc[operation.User, ResponseUserLogin](func() (*operation.User, error) {
		return userId.GetUser(userService.userOperation)
	})
	if res != nil {
		return res
	}

	if pass := userService.userOperation.VerifyUserPassword(user, req.Password); pass {
		token := NewClaims(userService.config.JWT, user, false)
		flushToken := NewClaims(userService.config.JWT, user, true)
		return NewApiResponse(&SuccessLogin, &ResponseUserLogin{
			User:       user,
			Token:      token.GenerateKey(),
			FlushToken: flushToken.GenerateKey(),
		})
	} else {
		return NewApiResponse[ResponseUserLogin](&ErrUsernameOrPassword, nil)
	}
}

var (
	NameNotAvailability = ApiStatus{StatusName: "INFO_NOT_AVAILABILITY", Description: "用户信息不可用", HttpCode: Ok}
	NameAvailability    = ApiStatus{StatusName: "INFO_AVAILABILITY", Description: "用户信息可用", HttpCode: Ok}
)

func (userService *UserService) CheckAvailability(req *RequestUserAvailability) *ApiResponse[ResponseUserAvailability] {
	if req.Username == "" && req.Email == "" && req.Cid == "" {
		return NewApiResponse[ResponseUserAvailability](ErrIllegalParam, nil)
	}
	exist, _ := userService.userOperation.IsUserIdentifierTaken(nil, utils.StrToInt(req.Cid, 0), req.Username, req.Email)
	data := ResponseUserAvailability(!exist)
	if exist {
		return NewApiResponse(&NameNotAvailability, &data)
	}
	return NewApiResponse(&NameAvailability, &data)
}

var (
	SuccessGetCurrentProfile = ApiStatus{StatusName: "GET_CURRENT_PROFILE_SUCCESS", Description: "获取当前用户信息成功", HttpCode: Ok}
)

func (userService *UserService) GetCurrentProfile(req *RequestUserCurrentProfile) *ApiResponse[ResponseUserCurrentProfile] {
	if user, err := userService.userOperation.GetUserByUid(req.Uid); errors.Is(err, operation.ErrUserNotFound) {
		return NewApiResponse[ResponseUserCurrentProfile](ErrUserNotFound, nil)
	} else if err != nil {
		return NewApiResponse[ResponseUserCurrentProfile](ErrDatabaseFail, nil)
	} else {
		return NewApiResponse(&SuccessGetCurrentProfile, (*ResponseUserCurrentProfile)(user))
	}
}

var (
	ErrOriginPasswordRequired = ApiStatus{StatusName: "ORIGIN_PASSWORD_REQUIRED", Description: "请输入原始密码", HttpCode: BadRequest}
	ErrNewPasswordRequired    = ApiStatus{StatusName: "NEW_PASSWORD_REQUIRED", Description: "请输入新密码", HttpCode: BadRequest}
	ErrOriginPassword         = ApiStatus{StatusName: "ORIGIN_PASSWORD_ERROR", Description: "原始密码不正确", HttpCode: BadRequest}
	ErrQQInvalid              = ApiStatus{StatusName: "QQ_INVALID", Description: "qq号不正确", HttpCode: BadRequest}
	SuccessEditCurrentProfile = ApiStatus{StatusName: "SUCCESS_EDIT_CURRENT_PROFILE", Description: "编辑用户信息成功", HttpCode: Ok}
)

func checkQQ(qq int) *ApiStatus {
	if 1e4 <= qq && qq < 1e11 {
		return nil
	}
	return &ErrQQInvalid
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
		return &ErrNewPasswordRequired, nil, ""
	} else if req.OriginPassword == "" && req.NewPassword != "" && !skipPasswordVerify {
		return &ErrOriginPasswordRequired, nil, ""
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

	updateInfo := make(map[string]interface{})

	oldValue, _ := json.Marshal(user)

	if req.Username != "" || req.Email != "" {
		exist, _ := userService.userOperation.IsUserIdentifierTaken(nil, 0, req.Username, req.Email)
		if exist {
			return ErrIdentifierTaken, nil, ""
		}
		if req.Username != "" && req.Username != user.Username {
			user.Username = req.Username
			updateInfo["username"] = req.Username
		}
		if req.Email != "" && req.Email != user.Email {
			user.Email = req.Email
			updateInfo["email"] = req.Email
		}
	}

	if req.QQ > 0 && req.QQ != user.QQ {
		user.QQ = req.QQ
		updateInfo["qq"] = req.QQ
		if req.AvatarUrl == "" && (user.AvatarUrl == "" || strings.HasPrefix(user.AvatarUrl, "https://q2.qlogo.cn/")) {
			user.AvatarUrl = fmt.Sprintf("https://q2.qlogo.cn/headimg_dl?dst_uin=%d&spec=100", user.QQ)
			updateInfo["avatar_url"] = user.AvatarUrl
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
		updateInfo["avatar_url"] = user.AvatarUrl
	}

	if req.OriginPassword != "" || skipPasswordVerify {
		password, err := userService.userOperation.UpdateUserPassword(user, req.OriginPassword, req.NewPassword, skipPasswordVerify)
		if errors.Is(err, operation.ErrUserNotFound) {
			return ErrUserNotFound, nil, ""
		} else if errors.Is(err, operation.ErrOldPassword) {
			return &ErrOriginPassword, nil, ""
		} else if err != nil {
			return ErrDatabaseFail, nil, ""
		}
		updateInfo["password"] = password
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
	if err, user, _ := userService.editUserProfile(req, false, false); err != nil {
		return NewApiResponse[ResponseUserEditCurrentProfile](err, nil)
	} else {
		return NewApiResponse(&SuccessEditCurrentProfile, (*ResponseUserEditCurrentProfile)(user))
	}
}

var (
	SuccessGetProfile = ApiStatus{StatusName: "GET_PROFILE_SUCCESS", Description: "获取用户信息成功", HttpCode: Ok}
)

func (userService *UserService) GetUserProfile(req *RequestUserProfile) *ApiResponse[ResponseUserProfile] {
	if req.TargetUid <= 0 {
		return NewApiResponse[ResponseUserProfile](ErrIllegalParam, nil)
	}
	if req.Permission <= 0 {
		return NewApiResponse[ResponseUserProfile](ErrNoPermission, nil)
	}
	permission := operation.Permission(req.Permission)
	if !permission.HasPermission(operation.UserGetProfile) {
		return NewApiResponse[ResponseUserProfile](ErrNoPermission, nil)
	}
	user, res := CallDBFunc[operation.User, ResponseUserProfile](func() (*operation.User, error) {
		return userService.userOperation.GetUserByUid(req.TargetUid)
	})
	if res != nil {
		return res
	}
	return NewApiResponse(&SuccessGetProfile, (*ResponseUserProfile)(user))
}

var (
	SuccessEditUserProfile = ApiStatus{StatusName: "EDIT_USER_PROFILE", Description: "修改用户信息成功", HttpCode: Ok}
)

func (userService *UserService) EditUserProfile(req *RequestUserEditProfile) *ApiResponse[ResponseUserEditProfile] {
	if req.Permission <= 0 {
		return NewApiResponse[ResponseUserEditProfile](ErrNoPermission, nil)
	}
	permission := operation.Permission(req.Permission)
	if !permission.HasPermission(operation.UserEditBaseInfo) {
		return NewApiResponse[ResponseUserEditProfile](ErrNoPermission, nil)
	}
	req.ID = req.TargetUid
	err, user, oldValue := userService.editUserProfile(&req.RequestUserEditCurrentProfile, true, permission.HasPermission(operation.UserSetPassword))
	if err != nil {
		return NewApiResponse[ResponseUserEditProfile](err, nil)
	}

	go func() {
		newValue, _ := json.Marshal(user)
		var object string
		if req.NewPassword == "" {
			object = fmt.Sprintf("%04d", user.Cid)
		} else {
			object = fmt.Sprintf("%04d(%s)", user.Cid, req.NewPassword)
		}
		userService.messageQueue.Publish(&queue.Message{
			Type: queue.AuditLog,
			Data: userService.auditLogOperation.NewAuditLog(operation.UserInformationEdit, req.Cid, object,
				req.Ip, req.UserAgent, &operation.ChangeDetail{
					OldValue: oldValue,
					NewValue: string(newValue),
				},
			),
		})
	}()

	return NewApiResponse(&SuccessEditUserProfile, (*ResponseUserEditProfile)(user))
}

var SuccessGetUsers = ApiStatus{StatusName: "GET_USER_PAGE", Description: "获取用户信息分页成功", HttpCode: Ok}

func (userService *UserService) GetUserList(req *RequestUserList) *ApiResponse[ResponseUserList] {
	if req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseUserList](ErrIllegalParam, nil)
	}
	if req.Permission <= 0 {
		return NewApiResponse[ResponseUserList](ErrNoPermission, nil)
	}
	permission := operation.Permission(req.Permission)
	if !permission.HasPermission(operation.UserShowList) {
		return NewApiResponse[ResponseUserList](ErrNoPermission, nil)
	}
	users, total, err := userService.userOperation.GetUsers(req.Page, req.PageSize)
	if err != nil {
		return NewApiResponse[ResponseUserList](ErrDatabaseFail, nil)
	}
	return NewApiResponse(&SuccessGetUsers, &ResponseUserList{
		Items:    users,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	})
}

var (
	ErrPermissionNodeNotExists = ApiStatus{StatusName: "PERMISSION_NODE_NOT_EXISTS", Description: "无效权限节点", HttpCode: BadRequest}
	SuccessEditUserPermission  = ApiStatus{StatusName: "EDIT_USER_PERMISSION", Description: "编辑用户权限成功", HttpCode: Ok}
)

func (userService *UserService) EditUserPermission(req *RequestUserEditPermission) *ApiResponse[ResponseUserEditPermission] {
	if req.Uid <= 0 {
		return NewApiResponse[ResponseUserEditPermission](ErrIllegalParam, nil)
	}
	user, targetUser, res := GetUsersAndCheckPermission[ResponseUserEditPermission](userService.userOperation, req.Uid, req.TargetUid, operation.UserEditPermission)
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
					auditLogs = append(auditLogs, userService.auditLogOperation.NewAuditLog(operation.UserPermissionGrant, req.Cid,
						fmt.Sprintf("%04d(%s)", targetUser.Cid, key), req.Ip, req.UserAgent, nil))
				} else {
					targetPermission.Revoke(per)
					auditLogs = append(auditLogs, userService.auditLogOperation.NewAuditLog(operation.UserPermissionRevoke, req.Cid,
						fmt.Sprintf("%04d(%s)", targetUser.Cid, key), req.Ip, req.UserAgent, nil))
				}
			} else {
				return NewApiResponse[ResponseUserEditPermission](ErrIllegalParam, nil)
			}
		} else {
			return NewApiResponse[ResponseUserEditPermission](&ErrPermissionNodeNotExists, nil)
		}
	}

	if _, res := CallDBFunc[interface{}, ResponseUserEditPermission](func() (*interface{}, error) {
		return nil, userService.userOperation.UpdateUserPermission(targetUser, targetPermission)
	}); res != nil {
		return res
	}

	if userService.config.Email.Template.EnablePermissionChangeEmail {
		userService.messageQueue.Publish(&queue.Message{
			Type: queue.SendPermissionChangeEmail,
			Data: &SendPermissionChangeData{
				User:     targetUser,
				Operator: user,
			},
		})
	}

	userService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLogs,
		Data: auditLogs,
	})

	return NewApiResponse(&SuccessEditUserPermission, (*ResponseUserEditPermission)(targetUser))
}

var SuccessGetUserHistory = ApiStatus{StatusName: "GET_USER_HISTORY", Description: "成功获取用户历史数据", HttpCode: Ok}

func (userService *UserService) GetUserHistory(req *RequestGetUserHistory) *ApiResponse[ResponseGetUserHistory] {
	if req.Cid <= 0 {
		return NewApiResponse[ResponseGetUserHistory](ErrIllegalParam, nil)
	}

	user, res := CallDBFunc[operation.User, ResponseGetUserHistory](func() (*operation.User, error) {
		return userService.userOperation.GetUserByCid(req.Cid)
	})
	if res != nil {
		return res
	}

	userHistory, res := CallDBFunc[operation.UserHistory, ResponseGetUserHistory](func() (*operation.UserHistory, error) {
		return userService.historyOperation.GetUserHistory(req.Cid)
	})
	if res != nil {
		return res
	}

	return NewApiResponse(&SuccessGetUserHistory, &ResponseGetUserHistory{
		TotalPilotTime: user.TotalPilotTime,
		TotalAtcTime:   user.TotalAtcTime,
		UserHistory:    userHistory,
	})
}

var SuccessGetToken = ApiStatus{StatusName: "GET_TOKEN", Description: "成功刷新秘钥", HttpCode: Ok}

func (userService *UserService) GetTokenWithFlushToken(req *RequestGetToken) *ApiResponse[ResponseGetToken] {
	if !req.FlushToken {
		return NewApiResponse[ResponseGetToken](ErrIllegalParam, nil)
	}

	user, res := CallDBFunc[operation.User, ResponseGetToken](func() (*operation.User, error) {
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
	return NewApiResponse(&SuccessGetToken, &ResponseGetToken{
		User:       user,
		Token:      token.GenerateKey(),
		FlushToken: flushToken,
	})
}
