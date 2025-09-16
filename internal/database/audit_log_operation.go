// Package database
package database

import (
	"context"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"gorm.io/gorm"
	"time"
)

type AuditLogOperation struct {
	logger       log.LoggerInterface
	db           *gorm.DB
	queryTimeout time.Duration
}

func NewAuditLogOperation(logger log.LoggerInterface, db *gorm.DB, queryTimeout time.Duration) *AuditLogOperation {
	return &AuditLogOperation{logger: logger, db: db, queryTimeout: queryTimeout}
}

func (auditLogOperation *AuditLogOperation) NewAuditLog(eventType AuditEventType, subject int, object, ip, userAgent string, changeDetails *ChangeDetail) (auditLog *AuditLog) {
	return &AuditLog{
		EventType:     string(eventType),
		Subject:       subject,
		Object:        object,
		Ip:            ip,
		UserAgent:     userAgent,
		ChangeDetails: changeDetails,
	}
}

func (auditLogOperation *AuditLogOperation) GetAuditLogs(page, pageSize int) (auditLogs []*AuditLog, total int64, err error) {
	auditLogs = make([]*AuditLog, 0, pageSize)
	ctx, cancel := context.WithTimeout(context.Background(), auditLogOperation.queryTimeout)
	defer cancel()
	auditLogOperation.db.WithContext(ctx).Model(&AuditLog{}).Select("id").Count(&total)
	err = auditLogOperation.db.WithContext(ctx).Offset((page - 1) * pageSize).Order("created_at desc").Limit(pageSize).Find(&auditLogs).Error
	return
}

func (auditLogOperation *AuditLogOperation) SaveAuditLog(auditLog *AuditLog) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), auditLogOperation.queryTimeout)
	defer cancel()
	return auditLogOperation.db.WithContext(ctx).Create(auditLog).Error
}

func (auditLogOperation *AuditLogOperation) SaveAuditLogs(auditLogs []*AuditLog) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), auditLogOperation.queryTimeout)
	defer cancel()
	return auditLogOperation.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(auditLogs).Error
	})
}
