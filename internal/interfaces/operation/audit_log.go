// Package operation
package operation

import "time"

type AuditLog struct {
	ID            uint          `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time     `gorm:"not null" json:"time"`
	EventType     string        `gorm:"index:eventType;not null" json:"event_type"`
	Subject       int           `gorm:"index:Subject;not null" json:"subject"`
	Object        string        `gorm:"index:Object;not null" json:"object"`
	Ip            string        `gorm:"not null" json:"ip"`
	UserAgent     string        `gorm:"not null" json:"user_agent"`
	ChangeDetails *ChangeDetail `gorm:"type:text;serializer:json" json:"change_details"`
}

type ChangeDetail struct {
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}
type AuditEventType string

const (
	UserInformationEdit     AuditEventType = "UserInformationEdit"
	UserPermissionGrant     AuditEventType = "UserPermissionGrant"
	UserPermissionRevoke    AuditEventType = "UserPermissionRevoke"
	ActivityCreated         AuditEventType = "ActivityCreated"
	ActivityDeleted         AuditEventType = "ActivityDeleted"
	ActivityUpdated         AuditEventType = "ActivityUpdated"
	ClientKickedFsd         AuditEventType = "ClientKickedFromFsd"
	ClientKicked            AuditEventType = "ClientKickedFromWeb"
	ClientMessage           AuditEventType = "ClientMessage"
	UnlawfulOverreach       AuditEventType = "UnlawfulOverreach"
	TicketOpen              AuditEventType = "TicketOpen"
	TicketClose             AuditEventType = "TicketClose"
	TicketDeleted           AuditEventType = "TicketDeleted"
	ControllerRecordCreated AuditEventType = "ControllerRecordCreated"
	ControllerRecordDeleted AuditEventType = "ControllerRecordUpdated"
	ControllerRatingChange  AuditEventType = "ControllerRatingChange"
	ControllerUMChange      AuditEventType = "ControllerUMChange"
	ControllerSoloChange    AuditEventType = "ControllerSoloChange"
	ControllerGuestChange   AuditEventType = "ControllerGuestChange"
	FlightPlanDeleted       AuditEventType = "FlightPlanDeleted"
	FlightPlanLock          AuditEventType = "FlightPlanLock"
	FlightPlanUnlock        AuditEventType = "FlightPlanUnlock"
)

type AuditLogOperationInterface interface {
	NewAuditLog(eventType AuditEventType, subject int, object, ip, userAgent string, changeDetails *ChangeDetail) (auditLog *AuditLog)
	SaveAuditLog(auditLog *AuditLog) (err error)
	SaveAuditLogs(auditLogs []*AuditLog) (err error)
	GetAuditLogs(page, pageSize int) (auditLogs []*AuditLog, total int64, err error)
}
