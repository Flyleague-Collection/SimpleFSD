// Package repository
package repository

import (
	"errors"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"github.com/half-nothing/simple-fsd/src/interfaces/enum"
)

const ValueNotAvailable = "NOT AVAILABLE"

type AuditEvent *enum.Enum[string]

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
	AuditEventUserInformationEdit  AuditEvent = enum.New("UserInformationEdit", "管理员修改用户信息")
	AuditEventUserPermissionGrant  AuditEvent = enum.New("UserPermissionGrant", "管理员授予用户权限")
	AuditEventUserPermissionRevoke AuditEvent = enum.New("UserPermissionRevoke", "管理员撤销用户权限")
	AuditEventUserRoleGrant        AuditEvent = enum.New("UserRoleGrant", "管理员授予用户角色")
	AuditEventUserRoleRevoke       AuditEvent = enum.New("UserRoleRevoke", "管理员撤销用户角色")

	// 活动相关
	AuditEventActivityCreated           AuditEvent = enum.New("ActivityCreated", "管理员创建活动")
	AuditEventActivityDeleted           AuditEvent = enum.New("ActivityDeleted", "管理员删除活动")
	AuditEventActivityUpdated           AuditEvent = enum.New("ActivityUpdated", "管理员修改活动信息")
	AuditEventActivityPilotSign         AuditEvent = enum.New("ActivityPilotSign", "飞行员报名活动")
	AuditEventActivityPilotLeave        AuditEvent = enum.New("ActivityPilotLeave", "飞行员退出活动")
	AuditEventActivityPilotStatusChange AuditEvent = enum.New("ActivityPilotStatusChange", "管理员修改飞行员活动状态")
	AuditEventActivityControllerJoin    AuditEvent = enum.New("AuditEventActivityControllerJoin", "管制员加入活动")
	AuditEventActivityControllerLeave   AuditEvent = enum.New("AuditEventActivityControllerLeave", "管制员退出活动")
	AuditEventActivityStatusChange      AuditEvent = enum.New("ActivityStatusChange", "管理员修改活动状态")

	// 在线管理相关
	AuditEventClientKickedFsd        AuditEvent = enum.New("ClientKickedFromFsd", "管理员在FSD中踢出用户")
	AuditEventClientKicked           AuditEvent = enum.New("ClientKickedFromWeb", "管理员在WEB中踢出用户")
	AuditEventClientMessage          AuditEvent = enum.New("ClientMessage", "管理员发送消息给用户")
	AuditEventClientBroadcastMessage AuditEvent = enum.New("ClientBroadcastMessage", "管理员广播消息给用户")

	// 异常访问
	AuditEventUnlawfulOverreach AuditEvent = enum.New("UnlawfulOverreach", "用户发生非法越权访问")

	// 工单相关
	AuditEventTicketOpen    AuditEvent = enum.New("TicketOpen", "用户创建工单")
	AuditEventTicketClose   AuditEvent = enum.New("TicketClose", "用户或管理员关闭工单")
	AuditEventTicketDeleted AuditEvent = enum.New("TicketDeleted", "用户或管理员删除工单")

	// 管制员相关
	AuditEventControllerRecordCreated         AuditEvent = enum.New("ControllerRecordCreated", "管理员创建管制员履历")
	AuditEventControllerRecordDeleted         AuditEvent = enum.New("ControllerRecordDeleted", "管理员删除管制员履历")
	AuditEventControllerRatingChange          AuditEvent = enum.New("ControllerRatingChange", "管理员修改管制员权限")
	AuditEventControllerApplicationSubmit     AuditEvent = enum.New("ControllerApplicationSubmit", "用户提交管制员申请")
	AuditEventControllerApplicationCancel     AuditEvent = enum.New("ControllerApplicationCancel", "用户取消管制员申请")
	AuditEventControllerApplicationPassed     AuditEvent = enum.New("ControllerApplicationPassed", "管理员通过管制员申请")
	AuditEventControllerApplicationProcessing AuditEvent = enum.New("ControllerApplicationProcessing", "管理员正在处理管制员申请")
	AuditEventControllerApplicationRejected   AuditEvent = enum.New("ControllerApplicationRejected", "管理员拒绝管制员申请")

	// 飞行计划相关
	AuditEventFlightPlanDeleted     AuditEvent = enum.New("FlightPlanDeleted", "管理员删除用户飞行计划")
	AuditEventFlightPlanSelfDeleted AuditEvent = enum.New("FlightPlanSelfDeleted", "用户删除自己的飞行计划")
	AuditEventFlightPlanLock        AuditEvent = enum.New("FlightPlanLock", "管制员锁定飞行计划")
	AuditEventFlightPlanUnlock      AuditEvent = enum.New("FlightPlanUnlock", "管制员解锁飞行计划")

	// 文件相关
	AuditEventFileUpload AuditEvent = enum.New("FileUpload", "用户上传文件")

	// 公告相关
	AuditEventAnnouncementPublished AuditEvent = enum.New("AnnouncementPublished", "管理员发布公告")
	AuditEventAnnouncementUpdated   AuditEvent = enum.New("AnnouncementUpdated", "管理员修改公告")
	AuditEventAnnouncementDeleted   AuditEvent = enum.New("AnnouncementDeleted", "管理员删除公告")
)

var AuditEventManager = enum.NewManager(
	AuditEventUserInformationEdit,
	AuditEventUserPermissionGrant,
	AuditEventUserPermissionRevoke,
	AuditEventUserRoleGrant,
	AuditEventUserRoleRevoke,
	AuditEventActivityCreated,
	AuditEventActivityDeleted,
	AuditEventActivityUpdated,
	AuditEventActivityPilotSign,
	AuditEventActivityPilotLeave,
	AuditEventActivityPilotStatusChange,
	AuditEventActivityControllerJoin,
	AuditEventActivityControllerLeave,
	AuditEventActivityStatusChange,
	AuditEventClientKickedFsd,
	AuditEventClientKicked,
	AuditEventClientMessage,
	AuditEventClientBroadcastMessage,
	AuditEventUnlawfulOverreach,
	AuditEventTicketOpen,
	AuditEventTicketClose,
	AuditEventTicketDeleted,
	AuditEventControllerRecordCreated,
	AuditEventControllerRecordDeleted,
	AuditEventControllerRatingChange,
	AuditEventControllerApplicationSubmit,
	AuditEventControllerApplicationCancel,
	AuditEventControllerApplicationPassed,
	AuditEventControllerApplicationProcessing,
	AuditEventControllerApplicationRejected,
	AuditEventFlightPlanDeleted,
	AuditEventFlightPlanSelfDeleted,
	AuditEventFlightPlanLock,
	AuditEventFlightPlanUnlock,
	AuditEventFileUpload,
	AuditEventAnnouncementPublished,
	AuditEventAnnouncementUpdated,
	AuditEventAnnouncementDeleted,
)
