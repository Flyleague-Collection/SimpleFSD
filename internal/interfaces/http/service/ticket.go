// Package service
package service

import "github.com/half-nothing/simple-fsd/internal/interfaces/operation"

var (
	ErrTicketNotFound      = NewApiStatus("TICKET_NOT_FOUND", "工单不存在", NotFound)
	ErrTicketAlreadyClosed = NewApiStatus("TICKET_ALREADY_CLOSED", "工单已回复", Conflict)
	SuccessGetTickets      = NewApiStatus("GET_TICKETS", "成功获取工单数据", Ok)
	SuccessGetUserTickets  = NewApiStatus("GET_USER_TICKETS", "成功获取用户工单数据", Ok)
	SuccessCreateTicket    = NewApiStatus("CREATE_TICKET", "成功创建工单", Ok)
	SuccessCloseTicket     = NewApiStatus("CLOSE_TICKET", "成功关闭工单", Ok)
	SuccessDeleteTicket    = NewApiStatus("DELETE_TICKET", "成功删除工单", Ok)
)

type TicketServiceInterface interface {
	GetTickets(req *RequestGetTickets) *ApiResponse[ResponseGetTickets]
	GetUserTickets(req *RequestGetUserTickets) *ApiResponse[ResponseGetUserTickets]
	CreateTicket(req *RequestCreateTicket) *ApiResponse[ResponseCreateTicket]
	CloseTicket(req *RequestCloseTicket) *ApiResponse[ResponseCloseTicket]
	DeleteTicket(req *RequestDeleteTicket) *ApiResponse[ResponseDeleteTicket]
}

type RequestGetTickets struct {
	JwtHeader
	Page     int `query:"page_number"`
	PageSize int `query:"page_size"`
}

type ResponseGetTickets struct {
	Items    []*operation.Ticket `json:"items"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
	Total    int64               `json:"total"`
}

type RequestGetUserTickets struct {
	JwtHeader
	Page     int `query:"page_number"`
	PageSize int `query:"page_size"`
}

type ResponseGetUserTickets struct {
	Items    []*operation.UserTicket `json:"items"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
	Total    int64                   `json:"total"`
}

type RequestCreateTicket struct {
	JwtHeader
	EchoContentHeader
	Type    int    `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ResponseCreateTicket bool

type RequestCloseTicket struct {
	JwtHeader
	EchoContentHeader
	TicketId uint   `param:"tid"`
	Reply    string `json:"reply"`
}

type ResponseCloseTicket bool

type RequestDeleteTicket struct {
	JwtHeader
	EchoContentHeader
	TicketId uint `param:"tid"`
}

type ResponseDeleteTicket bool
