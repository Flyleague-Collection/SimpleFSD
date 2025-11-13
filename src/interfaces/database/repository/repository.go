// Package repository
package repository

import "errors"

var (
	ErrArgument      = errors.New("argument error")
	ErrDataConflicts = errors.New("data conflicts")
)

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
