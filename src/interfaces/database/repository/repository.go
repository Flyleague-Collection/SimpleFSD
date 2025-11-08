// Package repository
package repository

type DatabaseInterface interface {
	GetUserRepository() UserInterface
	GetTicketRepository() TicketInterface
	GetActivityRepository() ActivityInterface
	GetHistoryRepository() HistoryInterface
	GetFlightPlanRepository() FlightPlanInterface
	GetControllerRecordRepository() ControllerRecordInterface
	GetControllerApplicationRepository() ControllerApplicationInterface
	GetAuditLogRepository() AuditLogInterface
	GetControllerRepository() ControllerInterface
	GetAnnouncementRepository() AnnouncementInterface
}
