// Package service
package service

import "github.com/half-nothing/simple-fsd/internal/interfaces/operation"

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
