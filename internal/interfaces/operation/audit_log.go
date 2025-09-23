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

const ValueNotAvailable = "NOT AVAILABLE"

type AuditEventType string

const (
	UserInformationEdit           AuditEventType = "UserInformationEdit"
	UserPermissionGrant           AuditEventType = "UserPermissionGrant"
	UserPermissionRevoke          AuditEventType = "UserPermissionRevoke"
	ActivityCreated               AuditEventType = "ActivityCreated"
	ActivityDeleted               AuditEventType = "ActivityDeleted"
	ActivityUpdated               AuditEventType = "ActivityUpdated"
	ClientKickedFsd               AuditEventType = "ClientKickedFromFsd"
	ClientKicked                  AuditEventType = "ClientKickedFromWeb"
	ClientMessage                 AuditEventType = "ClientMessage"
	ClientBroadcastMessage        AuditEventType = "ClientBroadcastMessage"
	UnlawfulOverreach             AuditEventType = "UnlawfulOverreach"
	TicketOpen                    AuditEventType = "TicketOpen"
	TicketClose                   AuditEventType = "TicketClose"
	TicketDeleted                 AuditEventType = "TicketDeleted"
	ControllerRecordCreated       AuditEventType = "ControllerRecordCreated"
	ControllerRecordDeleted       AuditEventType = "ControllerRecordDeleted"
	ControllerRatingChange        AuditEventType = "ControllerRatingChange"
	ControllerApplicationPaused   AuditEventType = "ControllerApplicationPaused"
	ControllerApplicationRejected AuditEventType = "ControllerApplicationRejected"
	FlightPlanDeleted             AuditEventType = "FlightPlanDeleted"
	FlightPlanLock                AuditEventType = "FlightPlanLock"
	FlightPlanUnlock              AuditEventType = "FlightPlanUnlock"
	FileUpload                    AuditEventType = "FileUpload"
)

type AuditLogOperationInterface interface {
	NewAuditLog(eventType AuditEventType, subject int, object, ip, userAgent string, changeDetails *ChangeDetail) (auditLog *AuditLog)
	SaveAuditLog(auditLog *AuditLog) (err error)
	SaveAuditLogs(auditLogs []*AuditLog) (err error)
	GetAuditLogs(page, pageSize int) (auditLogs []*AuditLog, total int64, err error)
}
