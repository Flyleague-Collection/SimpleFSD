// Package repository
package repository

import (
	"errors"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

const ValueNotAvailable = "NOT AVAILABLE"

type AuditEvent *Enum[string]

var (
	ErrAuditLogNotFound = errors.New("audit log not found")
)

type AuditLogInterface interface {
	Base[*entity.AuditLog]
	New(eventType AuditEvent, subject int, object, ip, userAgent string, changeDetails *entity.ChangeDetail) *entity.AuditLog
	BatchCreate(auditLogs []*entity.AuditLog) error
	GetPage(pageNumber int, pageSize int) ([]*entity.AuditLog, int64, error)
}

//goland:noinspection GoCommentStart
var (
	// 用户相关
	AuditEventUserInformationEdit  AuditEvent = NewEnum("UserInformationEdit", "管理员修改用户信息")
	AuditEventUserPermissionGrant  AuditEvent = NewEnum("UserPermissionGrant", "管理员授予用户权限")
	AuditEventUserPermissionRevoke AuditEvent = NewEnum("UserPermissionRevoke", "管理员撤销用户权限")
	AuditEventUserRoleGrant        AuditEvent = NewEnum("UserRoleGrant", "管理员授予用户角色")
	AuditEventUserRoleRevoke       AuditEvent = NewEnum("UserRoleRevoke", "管理员撤销用户角色")

	// 活动相关
	AuditEventActivityCreated           AuditEvent = NewEnum("ActivityCreated", "管理员创建活动")
	AuditEventActivityDeleted           AuditEvent = NewEnum("ActivityDeleted", "管理员删除活动")
	AuditEventActivityUpdated           AuditEvent = NewEnum("ActivityUpdated", "管理员修改活动信息")
	AuditEventActivityPilotSign         AuditEvent = NewEnum("ActivityPilotSign", "飞行员报名活动")
	AuditEventActivityPilotLeave        AuditEvent = NewEnum("ActivityPilotLeave", "飞行员退出活动")
	AuditEventActivityPilotStatusChange AuditEvent = NewEnum("ActivityPilotStatusChange", "管理员修改飞行员活动状态")
	AuditEventActivityControllerJoin    AuditEvent = NewEnum("AuditEventActivityControllerJoin", "管制员加入活动")
	AuditEventActivityControllerLeave   AuditEvent = NewEnum("AuditEventActivityControllerLeave", "管制员退出活动")
	AuditEventActivityStatusChange      AuditEvent = NewEnum("ActivityStatusChange", "管理员修改活动状态")

	// 在线管理相关
	AuditEventClientKickedFsd        AuditEvent = NewEnum("ClientKickedFromFsd", "管理员在FSD中踢出用户")
	AuditEventClientKicked           AuditEvent = NewEnum("ClientKickedFromWeb", "管理员在WEB中踢出用户")
	AuditEventClientMessage          AuditEvent = NewEnum("ClientMessage", "管理员发送消息给用户")
	AuditEventClientBroadcastMessage AuditEvent = NewEnum("ClientBroadcastMessage", "管理员广播消息给用户")

	// 异常访问
	AuditEventUnlawfulOverreach AuditEvent = NewEnum("UnlawfulOverreach", "用户发生非法越权访问")

	// 工单相关
	AuditEventTicketOpen    AuditEvent = NewEnum("TicketOpen", "用户创建工单")
	AuditEventTicketClose   AuditEvent = NewEnum("TicketClose", "用户或管理员关闭工单")
	AuditEventTicketDeleted AuditEvent = NewEnum("TicketDeleted", "用户或管理员删除工单")

	// 管制员相关
	AuditEventControllerRecordCreated         AuditEvent = NewEnum("ControllerRecordCreated", "管理员创建管制员履历")
	AuditEventControllerRecordDeleted         AuditEvent = NewEnum("ControllerRecordDeleted", "管理员删除管制员履历")
	AuditEventControllerRatingChange          AuditEvent = NewEnum("ControllerRatingChange", "管理员修改管制员权限")
	AuditEventControllerApplicationSubmit     AuditEvent = NewEnum("ControllerApplicationSubmit", "用户提交管制员申请")
	AuditEventControllerApplicationCancel     AuditEvent = NewEnum("ControllerApplicationCancel", "用户取消管制员申请")
	AuditEventControllerApplicationPassed     AuditEvent = NewEnum("ControllerApplicationPassed", "管理员通过管制员申请")
	AuditEventControllerApplicationProcessing AuditEvent = NewEnum("ControllerApplicationProcessing", "管理员正在处理管制员申请")
	AuditEventControllerApplicationRejected   AuditEvent = NewEnum("ControllerApplicationRejected", "管理员拒绝管制员申请")

	// 飞行计划相关
	AuditEventFlightPlanDeleted     AuditEvent = NewEnum("FlightPlanDeleted", "管理员删除用户飞行计划")
	AuditEventFlightPlanSelfDeleted AuditEvent = NewEnum("FlightPlanSelfDeleted", "用户删除自己的飞行计划")
	AuditEventFlightPlanLock        AuditEvent = NewEnum("FlightPlanLock", "管制员锁定飞行计划")
	AuditEventFlightPlanUnlock      AuditEvent = NewEnum("FlightPlanUnlock", "管制员解锁飞行计划")

	// 文件相关
	AuditEventFileUpload AuditEvent = NewEnum("FileUpload", "用户上传文件")

	// 公告相关
	AuditEventAnnouncementPublished AuditEvent = NewEnum("AnnouncementPublished", "管理员发布公告")
	AuditEventAnnouncementUpdated   AuditEvent = NewEnum("AnnouncementUpdated", "管理员修改公告")
	AuditEventAnnouncementDeleted   AuditEvent = NewEnum("AnnouncementDeleted", "管理员删除公告")
)
