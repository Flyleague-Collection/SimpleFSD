// Package service
package service

import (
	"encoding/json"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"strconv"
)

type AnnouncementService struct {
	logger                log.LoggerInterface
	messageQueue          queue.MessageQueueInterface
	announcementOperation operation.AnnouncementOperationInterface
	auditLogOperation     operation.AuditLogOperationInterface
}

func NewAnnouncementService(
	logger log.LoggerInterface,
	messageQueue queue.MessageQueueInterface,
	announcementOperation operation.AnnouncementOperationInterface,
	auditLogOperation operation.AuditLogOperationInterface,
) *AnnouncementService {
	return &AnnouncementService{
		logger:                log.NewLoggerAdapter(logger, "AnnouncementService"),
		messageQueue:          messageQueue,
		announcementOperation: announcementOperation,
		auditLogOperation:     auditLogOperation,
	}
}

func (service *AnnouncementService) GetAnnouncements(req *RequestGetAnnouncements) *ApiResponse[ResponseGetAnnouncements] {
	if req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseGetAnnouncements](ErrIllegalParam, nil)
	}

	announcements, total, err := service.announcementOperation.GetAnnouncements(req.Page, req.PageSize)
	if res := CheckDatabaseError[ResponseGetAnnouncements](err); res != nil {
		return res
	}

	data := ResponseGetAnnouncements(&PageResponse[*operation.UserAnnouncement]{
		Items:    announcements,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	})
	return NewApiResponse(SuccessGetAnnouncements, &data)
}

func (service *AnnouncementService) GetDetailAnnouncements(req *RequestGetDetailAnnouncements) *ApiResponse[ResponseGetDetailAnnouncements] {
	if req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseGetDetailAnnouncements](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseGetDetailAnnouncements](req.Permission, operation.AnnouncementShowList); res != nil {
		return res
	}

	announcements, total, err := service.announcementOperation.GetDetailAnnouncements(req.Page, req.PageSize)
	if res := CheckDatabaseError[ResponseGetDetailAnnouncements](err); res != nil {
		return res
	}

	data := ResponseGetDetailAnnouncements(&PageResponse[*operation.Announcement]{
		Items:    announcements,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	})
	return NewApiResponse(SuccessGetDetailAnnouncements, &data)
}

func (service *AnnouncementService) PublishAnnouncement(req *RequestPublishAnnouncement) *ApiResponse[ResponsePublishAnnouncement] {
	if req.Announcement == nil || !operation.IsValidAnnouncementType(req.Announcement.Type) || req.Announcement.Content == "" {
		return NewApiResponse[ResponsePublishAnnouncement](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponsePublishAnnouncement](req.Permission, operation.AnnouncementPublish); res != nil {
		return res
	}

	req.Announcement.ID = 0
	req.Announcement.PublisherId = req.Uid

	if res := CallDBFuncWithoutRet[ResponsePublishAnnouncement](func() error {
		return service.announcementOperation.SaveAnnouncement(req.Announcement)
	}); res != nil {
		return res
	}

	newValue, _ := json.Marshal(req.Announcement)
	service.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: service.auditLogOperation.NewAuditLog(
			operation.AnnouncementPublished,
			req.Cid,
			strconv.Itoa(int(req.Announcement.ID)),
			req.Ip,
			req.UserAgent,
			&operation.ChangeDetail{
				OldValue: operation.ValueNotAvailable,
				NewValue: string(newValue),
			},
		),
	})

	data := ResponsePublishAnnouncement(true)
	return NewApiResponse(SuccessPublishAnnouncement, &data)
}

func (service *AnnouncementService) EditAnnouncement(req *RequestEditAnnouncement) *ApiResponse[ResponseEditAnnouncement] {
	if req.Announcement == nil || req.AnnouncementId <= 0 || !operation.IsValidAnnouncementType(req.Announcement.Type) {
		return NewApiResponse[ResponseEditAnnouncement](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseEditAnnouncement](req.Permission, operation.AnnouncementEdit); res != nil {
		return res
	}

	announcement, res := CallDBFunc[*operation.Announcement, ResponseEditAnnouncement](func() (*operation.Announcement, error) {
		return service.announcementOperation.GetAnnouncementById(req.AnnouncementId)
	})
	if res != nil {
		return res
	}

	updateData := map[string]interface{}{}

	if req.Content != "" && req.Content != announcement.Content {
		updateData["content"] = req.Content
	}

	if req.Type != announcement.Type {
		updateData["type"] = req.Type
	}

	if req.Important != announcement.Important {
		updateData["important"] = req.Important
	}

	if req.ForceShow != announcement.ForceShow {
		updateData["force_show"] = req.ForceShow
	}

	if res := CallDBFuncWithoutRet[ResponseEditAnnouncement](func() error {
		return service.announcementOperation.UpdateAnnouncement(announcement, updateData)
	}); res != nil {
		return res
	}

	newValue, _ := json.Marshal(updateData)
	service.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: service.auditLogOperation.NewAuditLog(
			operation.AnnouncementUpdated,
			req.Cid,
			strconv.Itoa(int(req.Announcement.ID)),
			req.Ip,
			req.UserAgent,
			&operation.ChangeDetail{
				OldValue: operation.ValueNotAvailable,
				NewValue: string(newValue),
			},
		),
	})

	data := ResponseEditAnnouncement(true)
	return NewApiResponse(SuccessEditAnnouncement, &data)
}

func (service *AnnouncementService) DeleteAnnouncement(req *RequestDeleteAnnouncement) *ApiResponse[ResponseDeleteAnnouncement] {
	if req.AnnouncementId <= 0 {
		return NewApiResponse[ResponseDeleteAnnouncement](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseDeleteAnnouncement](req.Permission, operation.AnnouncementDelete); res != nil {
		return res
	}

	announcement, res := CallDBFunc[*operation.Announcement, ResponseDeleteAnnouncement](func() (*operation.Announcement, error) {
		return service.announcementOperation.GetAnnouncementById(req.AnnouncementId)
	})
	if res != nil {
		return res
	}

	oldValue, _ := json.Marshal(announcement)

	if res := CallDBFuncWithoutRet[ResponseDeleteAnnouncement](func() error {
		return service.announcementOperation.DeleteAnnouncement(announcement)
	}); res != nil {
		return res
	}

	service.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: service.auditLogOperation.NewAuditLog(
			operation.AnnouncementDeleted,
			req.Cid,
			strconv.Itoa(int(announcement.ID)),
			req.Ip,
			req.UserAgent,
			&operation.ChangeDetail{
				OldValue: string(oldValue),
				NewValue: operation.ValueNotAvailable,
			},
		),
	})

	data := ResponseDeleteAnnouncement(true)
	return NewApiResponse(SuccessDeleteAnnouncement, &data)
}
