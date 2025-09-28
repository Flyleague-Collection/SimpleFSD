// Package database
package database

import (
	"context"
	"errors"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"gorm.io/gorm"
	"time"
)

type AnnouncementOperation struct {
	logger       log.LoggerInterface
	db           *gorm.DB
	queryTimeout time.Duration
}

func NewAnnouncementOperation(
	logger log.LoggerInterface,
	db *gorm.DB,
	queryTimeout time.Duration,
) *AnnouncementOperation {
	return &AnnouncementOperation{
		logger:       logger,
		db:           db,
		queryTimeout: queryTimeout,
	}
}

func (operation *AnnouncementOperation) NewAnnouncement(uid uint, content string, announcementType AnnouncementType, important bool, forceShow bool) *Announcement {
	return &Announcement{
		PublisherId: uid,
		Content:     content,
		Type:        int(announcementType),
		Important:   important,
		ForceShow:   forceShow,
	}
}

func (operation *AnnouncementOperation) SaveAnnouncement(announcement *Announcement) error {
	ctx, cancel := context.WithTimeout(context.Background(), operation.queryTimeout)
	defer cancel()
	if announcement.ID == 0 {
		return operation.db.WithContext(ctx).Create(announcement).Error
	}
	return operation.db.WithContext(ctx).Save(announcement).Error
}

func (operation *AnnouncementOperation) GetAnnouncementById(id uint) (announcement *Announcement, err error) {
	announcement = &Announcement{}
	ctx, cancel := context.WithTimeout(context.Background(), operation.queryTimeout)
	defer cancel()
	err = operation.db.WithContext(ctx).Preload("User").First(announcement, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrAnnouncementNotFound
	}
	return
}

func (operation *AnnouncementOperation) GetAnnouncements(page, pageSize int) (announcements []*UserAnnouncement, total int64, err error) {
	announcements = make([]*UserAnnouncement, 0, pageSize)
	ctx, cancel := context.WithTimeout(context.Background(), operation.queryTimeout)
	defer cancel()
	operation.db.WithContext(ctx).Model(&Announcement{}).Select("id").Count(&total)
	err = operation.db.WithContext(ctx).Model(&Announcement{}).Order("important desc, updated_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&announcements).Error
	return
}

func (operation *AnnouncementOperation) GetDetailAnnouncements(page, pageSize int) (announcements []*Announcement, total int64, err error) {
	announcements = make([]*Announcement, 0, pageSize)
	ctx, cancel := context.WithTimeout(context.Background(), operation.queryTimeout)
	defer cancel()
	operation.db.WithContext(ctx).Model(&Announcement{}).Select("id").Count(&total)
	err = operation.db.WithContext(ctx).Preload("User").Offset((page - 1) * pageSize).Limit(pageSize).Find(&announcements).Error
	return
}

func (operation *AnnouncementOperation) DeleteAnnouncement(announcement *Announcement) error {
	ctx, cancel := context.WithTimeout(context.Background(), operation.queryTimeout)
	defer cancel()
	result := operation.db.WithContext(ctx).Delete(announcement)
	if result.RowsAffected == 0 {
		return ErrAnnouncementNotFound
	}
	return result.Error
}

func (operation *AnnouncementOperation) UpdateAnnouncement(announcement *Announcement, updates map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), operation.queryTimeout)
	defer cancel()
	result := operation.db.WithContext(ctx).Model(announcement).Updates(updates)
	if result.RowsAffected == 0 {
		return ErrAnnouncementNotFound
	}
	return result.Error
}
