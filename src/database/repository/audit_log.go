// Package repository
package repository

import (
	"errors"
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"github.com/half-nothing/simple-fsd/src/interfaces/database/repository"
	"github.com/half-nothing/simple-fsd/src/interfaces/logger"
	"github.com/half-nothing/simple-fsd/src/utils"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	*BaseRepository[*entity.AuditLog]
	pageReq PageableInterface[*entity.AuditLog]
}

func NewAuditLogRepository(
	lg logger.Interface,
	db *gorm.DB,
	queryTimeout time.Duration,
) *AuditLogRepository {
	return &AuditLogRepository{
		BaseRepository: NewBaseRepository[*entity.AuditLog](lg, "AuditLogRepository", db, queryTimeout),
		pageReq:        NewPageRequest[*entity.AuditLog](db),
	}
}

func (repo *AuditLogRepository) New(
	eventType repository.AuditEvent,
	subject int,
	object string,
	ip string,
	userAgent string,
	changeDetails *entity.ChangeDetail,
) *entity.AuditLog {
	if eventType == nil || subject <= 0 || object == "" || ip == "" || userAgent == "" {
		return nil
	}

	return &entity.AuditLog{
		EventType:     eventType.Value,
		Subject:       subject,
		Object:        object,
		Ip:            ip,
		UserAgent:     userAgent,
		ChangeDetails: changeDetails,
	}
}

func (repo *AuditLogRepository) GetById(id uint) (*entity.AuditLog, error) {
	auditLog := &entity.AuditLog{ID: id}
	err := repo.query(func(tx *gorm.DB) error {
		return tx.First(auditLog).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = repository.ErrAuditLogNotFound
	}
	return auditLog, err
}

func (repo *AuditLogRepository) GetPage(pageNumber int, pageSize int) (auditLogs []*entity.AuditLog, total int64, err error) {
	auditLogs = make([]*entity.AuditLog, 0, pageSize)
	total, err = repo.queryWithPagination(repo.pageReq, NewPage(pageNumber, pageSize, auditLogs, &entity.AuditLog{}, nil))
	return
}

func (repo *AuditLogRepository) Save(entity *entity.AuditLog) error {
	return repo.save(entity)
}

func (repo *AuditLogRepository) BatchCreate(auditLogs []*entity.AuditLog) error {
	auditLogs = utils.Filter(auditLogs, func(auditLog *entity.AuditLog) bool {
		return auditLog != nil && auditLog.ID == 0
	})
	if len(auditLogs) == 0 {
		return nil
	}
	return repo.queryWithTransaction(func(tx *gorm.DB) error {
		return tx.Create(auditLogs).Error
	})
}

func (repo *AuditLogRepository) Delete(entity *entity.AuditLog) error {
	return repo.delete(entity)
}

func (repo *AuditLogRepository) Update(entity *entity.AuditLog, updates map[string]interface{}) error {
	return repo.update(entity, updates)
}
