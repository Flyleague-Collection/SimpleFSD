// Package database
package database

import "github.com/half-nothing/simple-fsd/src/interfaces/database/repository"

type DatabaseRepositoryBuilder struct {
	repo *DatabaseRepository
}

func NewDatabaseRepositoryBuilder() *DatabaseRepositoryBuilder {
	return &DatabaseRepositoryBuilder{
		repo: &DatabaseRepository{},
	}
}

func (builder *DatabaseRepositoryBuilder) SetUserRepository(repo repository.UserInterface) *DatabaseRepositoryBuilder {
	builder.repo.userRepo = repo
	return builder
}

func (builder *DatabaseRepositoryBuilder) SetTicketRepository(repo repository.TicketInterface) *DatabaseRepositoryBuilder {
	builder.repo.ticketRepo = repo
	return builder
}

func (builder *DatabaseRepositoryBuilder) SetActivityRepository(repo repository.ActivityInterface) *DatabaseRepositoryBuilder {
	builder.repo.activityRepo = repo
	return builder
}

func (builder *DatabaseRepositoryBuilder) SetHistoryRepository(repo repository.HistoryInterface) *DatabaseRepositoryBuilder {
	builder.repo.historyRepo = repo
	return builder
}

func (builder *DatabaseRepositoryBuilder) SetFlightPlanRepository(repo repository.FlightPlanInterface) *DatabaseRepositoryBuilder {
	builder.repo.flightPlanRepo = repo
	return builder
}

func (builder *DatabaseRepositoryBuilder) SetControllerRecordRepository(repo repository.ControllerRecordInterface) *DatabaseRepositoryBuilder {
	builder.repo.controllerRecordRepo = repo
	return builder
}

func (builder *DatabaseRepositoryBuilder) SetControllerApplicationRepository(repo repository.ControllerApplicationInterface) *DatabaseRepositoryBuilder {
	builder.repo.controllerApplicationRepo = repo
	return builder
}

func (builder *DatabaseRepositoryBuilder) SetAuditLogRepository(repo repository.AuditLogInterface) *DatabaseRepositoryBuilder {
	builder.repo.auditLogRepo = repo
	return builder
}

func (builder *DatabaseRepositoryBuilder) SetControllerRepository(repo repository.ControllerInterface) *DatabaseRepositoryBuilder {
	builder.repo.controllerRepo = repo
	return builder
}

func (builder *DatabaseRepositoryBuilder) SetAnnouncementRepository(repo repository.AnnouncementInterface) *DatabaseRepositoryBuilder {
	builder.repo.announcementRepo = repo
	return builder
}

func (builder *DatabaseRepositoryBuilder) Build() *DatabaseRepository {
	return builder.repo
}
