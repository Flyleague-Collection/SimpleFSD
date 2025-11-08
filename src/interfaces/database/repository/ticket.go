// Package repository
package repository

import (
	"errors"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/DTO"
	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

type TicketType int

const (
	Feature     TicketType = iota // 建议
	Bug                           // bug
	Complain                      // 投诉
	Recognition                   // 表扬
	OtherType                     // 其他
)

var (
	ErrTicketNotFound      = errors.New("ticket not found")
	ErrTicketAlreadyClosed = errors.New("ticket already closed")
)

func IsValidTicketType(s int) bool {
	return int(Feature) <= s && s <= int(OtherType)
}

type TicketInterface interface {
	NewTicket(userId uint, ticketType TicketType, title string, content string) (ticket *entity.Ticket)
	SaveTicket(ticket *entity.Ticket) (err error)
	GetTickets(page, pageSize int) (tickets []*entity.Ticket, total int64, err error)
	GetUserTickets(uid uint, page, pageSize int) (tickets []*DTO.UserTicket, total int64, err error)
	GetTicket(id uint) (ticket *entity.Ticket, err error)
	CloseTicket(ticket *entity.Ticket, closer int, content string) (err error)
	DeleteTicket(id uint) (err error)
}
