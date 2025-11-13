// Package database
package database

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/database/repository"
)

type DatabaseRepository struct {
	userRepo                  repository.UserInterface
	ticketRepo                repository.TicketInterface
	activityRepo              repository.ActivityInterface
	historyRepo               repository.HistoryInterface
	flightPlanRepo            repository.FlightPlanInterface
	controllerRecordRepo      repository.ControllerRecordInterface
	controllerApplicationRepo repository.ControllerApplicationInterface
	auditLogRepo              repository.AuditLogInterface
	controllerRepo            repository.ControllerInterface
	announcementRepo          repository.AnnouncementInterface
}

func (repo *DatabaseRepository) GetUserRepository() repository.UserInterface {
	return repo.userRepo
}

func (repo *DatabaseRepository) GetTicketRepository() repository.TicketInterface {
	return repo.ticketRepo
}

func (repo *DatabaseRepository) GetActivityRepository() repository.ActivityInterface {
	return repo.activityRepo
}

func (repo *DatabaseRepository) GetHistoryRepository() repository.HistoryInterface {
	return repo.historyRepo
}

func (repo *DatabaseRepository) GetFlightPlanRepository() repository.FlightPlanInterface {
	return repo.flightPlanRepo
}

func (repo *DatabaseRepository) GetControllerRecordRepository() repository.ControllerRecordInterface {
	return repo.controllerRecordRepo
}

func (repo *DatabaseRepository) GetControllerApplicationRepository() repository.ControllerApplicationInterface {
	return repo.controllerApplicationRepo
}

func (repo *DatabaseRepository) GetAuditLogRepository() repository.AuditLogInterface {
	return repo.auditLogRepo
}

func (repo *DatabaseRepository) GetControllerRepository() repository.ControllerInterface {
	return repo.controllerRepo
}

func (repo *DatabaseRepository) GetAnnouncementRepository() repository.AnnouncementInterface {
	return repo.announcementRepo
}
