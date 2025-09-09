// Package operation
package operation

type EventType string

const (
	UserInformationEdit  EventType = "UserInformationEdit"
	UserPermissionGrant  EventType = "UserPermissionGrant"
	UserPermissionRevoke EventType = "UserPermissionRevoke"
	UserRatingChange     EventType = "UserRatingChange"
	ActivityCreated      EventType = "ActivityCreated"
	ActivityDeleted      EventType = "ActivityDeleted"
	ActivityUpdated      EventType = "ActivityUpdated"
	TicketCreated        EventType = "TicketCreated"
	TicketReply          EventType = "TicketReply"
	ClientKicked         EventType = "ClientKicked"
	ClientMessage        EventType = "ClientMessage"
	UnlawfulOverreach    EventType = "UnlawfulOverreach"
)

type AuditLogOperationInterface interface {
	NewAuditLog(eventType EventType, subject int, object, ip, userAgent string, changeDetails *ChangeDetail) (auditLog *AuditLog)
	SaveAuditLog(auditLog *AuditLog) (err error)
	SaveAuditLogs(auditLogs []*AuditLog) (err error)
	GetAuditLogs(page, pageSize int) (auditLogs []*AuditLog, total int64, err error)
}
