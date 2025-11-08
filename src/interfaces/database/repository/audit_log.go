// Package repository
package repository

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

const ValueNotAvailable = "NOT AVAILABLE"

type AuditEventType string

const (
	UserInformationEdit             AuditEventType = "UserInformationEdit"
	UserPermissionGrant             AuditEventType = "UserPermissionGrant"
	UserPermissionRevoke            AuditEventType = "UserPermissionRevoke"
	ActivityCreated                 AuditEventType = "ActivityCreated"
	ActivityDeleted                 AuditEventType = "ActivityDeleted"
	ActivityUpdated                 AuditEventType = "ActivityUpdated"
	ClientKickedFsd                 AuditEventType = "ClientKickedFromFsd"
	ClientKicked                    AuditEventType = "ClientKickedFromWeb"
	ClientMessage                   AuditEventType = "ClientMessage"
	ClientBroadcastMessage          AuditEventType = "ClientBroadcastMessage"
	UnlawfulOverreach               AuditEventType = "UnlawfulOverreach"
	TicketOpen                      AuditEventType = "TicketOpen"
	TicketClose                     AuditEventType = "TicketClose"
	TicketDeleted                   AuditEventType = "TicketDeleted"
	ControllerRecordCreated         AuditEventType = "ControllerRecordCreated"
	ControllerRecordDeleted         AuditEventType = "ControllerRecordDeleted"
	ControllerRatingChange          AuditEventType = "ControllerRatingChange"
	ControllerApplicationSubmit     AuditEventType = "ControllerApplicationSubmit"
	ControllerApplicationCancel     AuditEventType = "ControllerApplicationCancel"
	ControllerApplicationPassed     AuditEventType = "ControllerApplicationPassed"
	ControllerApplicationProcessing AuditEventType = "ControllerApplicationProcessing"
	ControllerApplicationRejected   AuditEventType = "ControllerApplicationRejected"
	FlightPlanDeleted               AuditEventType = "FlightPlanDeleted"
	FlightPlanSelfDeleted           AuditEventType = "FlightPlanSelfDeleted"
	FlightPlanLock                  AuditEventType = "FlightPlanLock"
	FlightPlanUnlock                AuditEventType = "FlightPlanUnlock"
	FileUpload                      AuditEventType = "FileUpload"
	AnnouncementPublished           AuditEventType = "AnnouncementPublished"
	AnnouncementUpdated             AuditEventType = "AnnouncementUpdated"
	AnnouncementDeleted             AuditEventType = "AnnouncementDeleted"
)

type AuditLogInterface interface {
	NewAuditLog(eventType AuditEventType, subject int, object, ip, userAgent string, changeDetails *entity.ChangeDetail) (auditLog *entity.AuditLog)
	SaveAuditLog(auditLog *entity.AuditLog) (err error)
	SaveAuditLogs(auditLogs []*entity.AuditLog) (err error)
	GetAuditLogs(page, pageSize int) (auditLogs []*entity.AuditLog, total int64, err error)
}
