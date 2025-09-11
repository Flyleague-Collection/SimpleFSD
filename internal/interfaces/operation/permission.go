// Package operation
package operation

type Permission int64

// 权限节点上限是64, 超过64需要使用切片
const (
	AdminEntry Permission = 1 << iota
	UserShowList
	UserGetProfile
	UserSetPassword
	UserEditBaseInfo
	UserEditPermission
	UserEditRating
	ActivityPublish
	ActivityShowList
	ActivityEdit
	ActivityEditState
	ActivityEditPilotState
	ActivityDelete
	AuditLogShow
	ClientManagerEntry
	ClientSendMessage
	ClientKill
)

var PermissionMap = map[string]Permission{
	"AdminEntry":             AdminEntry,
	"UserShowList":           UserShowList,
	"UserGetProfile":         UserGetProfile,
	"UserSetPassword":        UserSetPassword,
	"UserEditBaseInfo":       UserEditBaseInfo,
	"UserEditPermission":     UserEditPermission,
	"UserEditRating":         UserEditRating,
	"ActivityPublish":        ActivityPublish,
	"ActivityShowList":       ActivityShowList,
	"ActivityEdit":           ActivityEdit,
	"ActivityEditState":      ActivityEditState,
	"ActivityEditPilotState": ActivityEditPilotState,
	"ActivityDelete":         ActivityDelete,
	"AuditLogShow":           AuditLogShow,
	"ClientManagerEntry":     ClientManagerEntry,
	"ClientSendMessage":      ClientSendMessage,
	"ClientKill":             ClientKill,
}

func (p *Permission) IsValid() bool {
	maxPerm := ClientKill<<1 - 1 // 计算最大有效位
	return *p >= 0 && *p <= maxPerm
}

func (p *Permission) HasPermission(perm Permission) bool {
	return *p&perm != 0
}

func (p *Permission) Grant(perm Permission) {
	*p |= perm
}

func (p *Permission) Revoke(perm Permission) {
	*p &^= perm
}
