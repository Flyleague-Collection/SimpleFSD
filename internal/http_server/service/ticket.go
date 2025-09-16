// Package service
package service

import (
	"encoding/json"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
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
		logger:            logger,
		messageQueue:      messageQueue,
		userOperation:     userOperation,
		ticketOperation:   ticketOperation,
		auditLogOperation: auditLogOperation,
	}
}

var SuccessGetAllTickets = NewApiStatus("GET_ALL_TICKETS", "成功获取工单数据", Ok)

func (ticketService *TicketService) GetTickets(req *RequestGetTicket) *ApiResponse[ResponseGetTicket] {
	if req.Uid < 0 || req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseGetTicket](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseGetTicket](req.Permission, operation.TicketShowList); res != nil {
		return res
	}

	records, total, err := ticketService.ticketOperation.GetTickets(req.Page, req.PageSize)
	if err != nil {
		return NewApiResponse[ResponseGetTicket](ErrDatabaseFail, nil)
	}

	return NewApiResponse(SuccessGetAllTickets, &ResponseGetTicket{
		Items:    records,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

var SuccessGetUserTickets = NewApiStatus("GET_USER_TICKETS", "成功获取用户工单数据", Ok)

func (ticketService *TicketService) GetUserTicket(req *RequestGetUserTicket) *ApiResponse[ResponseGetUserTicket] {
	if req.Uid < 0 || req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseGetUserTicket](ErrIllegalParam, nil)
	}
	records, total, err := ticketService.ticketOperation.GetUserTickets(req.Cid, req.Page, req.PageSize)
	if err != nil {
		return NewApiResponse[ResponseGetUserTicket](ErrDatabaseFail, nil)
	}

	return NewApiResponse(SuccessGetUserTickets, &ResponseGetUserTicket{
		Items:    records,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

var SuccessCreateTicket = NewApiStatus("CREATE_TICKET", "成功创建工单", Ok)

func (ticketService *TicketService) CreateTicket(req *RequestCreateTicket) *ApiResponse[ResponseCreateTicket] {
	if req.Uid <= 0 || req.Content == "" || !operation.IsValidTicketType(req.Type) {
		return NewApiResponse[ResponseCreateTicket](ErrIllegalParam, nil)
	}

	ticketType := operation.ToTicketType(req.Type)

	ticket := ticketService.ticketOperation.NewTicket(req.Cid, ticketType, req.Title, req.Content)

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
				OldValue: "NOT AVAILABLE",
				NewValue: string(newValue),
			},
		),
	})

	data := ResponseCreateTicket(true)
	return NewApiResponse(SuccessCreateTicket, &data)
}

var SuccessCloseTicket = NewApiStatus("CLOSE_TICKET", "成功关闭工单", Ok)

func (ticketService *TicketService) CloseTicket(req *RequestCloseTicket) *ApiResponse[ResponseCloseTicket] {
	if req.Uid <= 0 || req.TicketId <= 0 || req.Reply == "" {
		return NewApiResponse[ResponseCloseTicket](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseCloseTicket](req.Permission, operation.TicketReply); res != nil {
		return res
	}

	if res := CallDBFuncWithoutRet[ResponseCloseTicket](func() error {
		return ticketService.ticketOperation.CloseTicket(req.TicketId, req.Cid, req.Reply)
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
				OldValue: "NOT AVAILABLE",
				NewValue: req.Reply,
			},
		),
	})

	data := ResponseCloseTicket(true)
	return NewApiResponse(SuccessCloseTicket, &data)
}

var SuccessDeleteTicket = NewApiStatus("DELETE_TICKET", "成功删除工单", Ok)

func (ticketService *TicketService) DeleteTicket(req *RequestDeleteTicket) *ApiResponse[ResponseDeleteTicket] {
	if req.Uid <= 0 || req.TicketId <= 0 {
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
