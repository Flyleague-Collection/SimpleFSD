// Package service
// 存放 TicketServiceInterface 的实现
package service

import (
	"encoding/json"
	"fmt"

	"github.com/half-nothing/simple-fsd/internal/interfaces"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/http/service"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
)

type TicketService struct {
	logger            log.LoggerInterface
	messageQueue      queue.MessageQueueInterface
	userOperation     operation.UserOperationInterface
	ticketOperation   operation.TicketOperationInterface
	auditLogOperation operation.AuditLogOperationInterface
}

func NewTicketService(
	logger log.LoggerInterface,
	messageQueue queue.MessageQueueInterface,
	userOperation operation.UserOperationInterface,
	ticketOperation operation.TicketOperationInterface,
	auditLogOperation operation.AuditLogOperationInterface,
) *TicketService {
	return &TicketService{
		logger:            log.NewLoggerAdapter(logger, "TicketService"),
		messageQueue:      messageQueue,
		userOperation:     userOperation,
		ticketOperation:   ticketOperation,
		auditLogOperation: auditLogOperation,
	}
}

func (ticketService *TicketService) GetTickets(req *RequestGetTickets) *ApiResponse[ResponseGetTickets] {
	if req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseGetTickets](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseGetTickets](req.Permission, operation.TicketShowList); res != nil {
		return res
	}

	records, total, err := ticketService.ticketOperation.GetTickets(req.Page, req.PageSize)
	if res := CheckDatabaseError[ResponseGetTickets](err); res != nil {
		return res
	}

	return NewApiResponse(SuccessGetTickets, &ResponseGetTickets{
		Items:    records,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

func (ticketService *TicketService) GetUserTickets(req *RequestGetUserTickets) *ApiResponse[ResponseGetUserTickets] {
	if req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseGetUserTickets](ErrIllegalParam, nil)
	}

	records, total, err := ticketService.ticketOperation.GetUserTickets(req.Uid, req.Page, req.PageSize)
	if res := CheckDatabaseError[ResponseGetUserTickets](err); res != nil {
		return res
	}

	return NewApiResponse(SuccessGetUserTickets, &ResponseGetUserTickets{
		Items:    records,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

func (ticketService *TicketService) CreateTicket(req *RequestCreateTicket) *ApiResponse[ResponseCreateTicket] {
	if req.Title == "" || req.Content == "" || !operation.IsValidTicketType(req.Type) {
		return NewApiResponse[ResponseCreateTicket](ErrIllegalParam, nil)
	}

	ticketType := operation.TicketType(req.Type)

	ticket := ticketService.ticketOperation.NewTicket(req.Uid, ticketType, req.Title, req.Content)

	if res := CallDBFuncWithoutRet[ResponseCreateTicket](func() error {
		return ticketService.ticketOperation.SaveTicket(ticket)
	}); res != nil {
		return res
	}

	newValue, _ := json.Marshal(ticket)
	ticketService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: ticketService.auditLogOperation.NewAuditLog(
			operation.TicketOpen,
			req.Cid,
			fmt.Sprintf("%d(%04d)", ticket.ID, req.Cid),
			req.Ip,
			req.UserAgent,
			&operation.ChangeDetail{
				OldValue: operation.ValueNotAvailable,
				NewValue: string(newValue),
			},
		),
	})

	data := ResponseCreateTicket(true)
	return NewApiResponse(SuccessCreateTicket, &data)
}

func (ticketService *TicketService) CloseTicket(req *RequestCloseTicket) *ApiResponse[ResponseCloseTicket] {
	if req.TicketId <= 0 || req.Reply == "" {
		return NewApiResponse[ResponseCloseTicket](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseCloseTicket](req.Permission, operation.TicketReply); res != nil {
		return res
	}

	ticket, res := CallDBFunc[*operation.Ticket, ResponseCloseTicket](func() (*operation.Ticket, error) {
		return ticketService.ticketOperation.GetTicket(req.TicketId)
	})
	if res != nil {
		return res
	}

	if res := CallDBFuncWithoutRet[ResponseCloseTicket](func() error {
		return ticketService.ticketOperation.CloseTicket(ticket, req.Cid, req.Reply)
	}); res != nil {
		return res
	}

	ticketService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: ticketService.auditLogOperation.NewAuditLog(
			operation.TicketClose,
			req.Cid,
			fmt.Sprintf("%d", req.TicketId),
			req.Ip,
			req.UserAgent,
			&operation.ChangeDetail{
				OldValue: operation.ValueNotAvailable,
				NewValue: req.Reply,
			},
		),
	})

	ticketService.messageQueue.Publish(&queue.Message{
		Type: queue.SendTicketReplyEmail,
		Data: &interfaces.TicketReplyEmailData{
			User:  ticket.User,
			Title: ticket.Title,
			Reply: req.Reply,
		},
	})

	data := ResponseCloseTicket(true)
	return NewApiResponse(SuccessCloseTicket, &data)
}

func (ticketService *TicketService) DeleteTicket(req *RequestDeleteTicket) *ApiResponse[ResponseDeleteTicket] {
	if req.TicketId <= 0 {
		return NewApiResponse[ResponseDeleteTicket](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseDeleteTicket](req.Permission, operation.TicketRemove); res != nil {
		return res
	}

	if res := CallDBFuncWithoutRet[ResponseDeleteTicket](func() error {
		return ticketService.ticketOperation.DeleteTicket(req.TicketId)
	}); res != nil {
		return res
	}

	ticketService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: ticketService.auditLogOperation.NewAuditLog(
			operation.TicketDeleted,
			req.Cid,
			fmt.Sprintf("%d", req.TicketId),
			req.Ip,
			req.UserAgent,
			nil,
		),
	})

	data := ResponseDeleteTicket(true)
	return NewApiResponse(SuccessDeleteTicket, &data)
}
