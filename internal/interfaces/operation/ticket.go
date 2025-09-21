// Package operation
package operation

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Ticket struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Opener    int            `gorm:"index:opener;not null" json:"creator"`
	Type      int            `gorm:"not null" json:"type"`
	Title     string         `gorm:"not null" json:"title"`
	Content   string         `gorm:"not null" json:"content"`
	Reply     string         `gorm:"not null" json:"reply"`
	Closer    int            `gorm:"index:closer;not null" json:"closer"`
	CreatedAt time.Time      `json:"open_at"`
	UpdatedAt time.Time      `json:"close_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type UserTicket struct {
	ID        uint      `json:"id"`
	Type      int       `json:"type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Reply     string    `json:"reply"`
	CreatedAt time.Time `json:"open_at"`
	UpdatedAt time.Time `json:"close_at"`
}

type TicketType int

const (
	Feature     TicketType = iota // 建议
	Bug                           // bug
	Complain                      // 投诉
	Recognition                   // 表扬
	OtherType                     // 其他
)

func IsValidTicketType(s int) bool {
	return int(Feature) <= s && s <= int(Other)
}

var (
	ErrTicketNotFound      = errors.New("ticket not found")
	ErrTicketAlreadyClosed = errors.New("ticket already closed")
)

type TicketOperationInterface interface {
	NewTicket(opener int, ticketType TicketType, title string, content string) (ticket *Ticket)
	SaveTicket(ticket *Ticket) (err error)
	GetTickets(page, pageSize int) (tickets []*Ticket, total int64, err error)
	GetUserTickets(cid, page, pageSize int) (tickets []*UserTicket, total int64, err error)
	GetTicket(id uint) (ticket *Ticket, err error)
	CloseTicket(ticketId uint, closer int, content string) (err error)
	DeleteTicket(id uint) (err error)
}
