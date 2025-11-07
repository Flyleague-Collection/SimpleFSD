// Package database
package database

import (
	"context"
	"errors"
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/log"
	. "github.com/half-nothing/simple-fsd/src/interfaces/operation"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketOperation struct {
	logger       log.LoggerInterface
	db           *gorm.DB
	queryTimeout time.Duration
}

func NewTicketOperation(logger log.LoggerInterface, db *gorm.DB, queryTimeout time.Duration) *TicketOperation {
	return &TicketOperation{
		logger:       logger,
		db:           db,
		queryTimeout: queryTimeout,
	}
}

func (ticketOperation *TicketOperation) NewTicket(userId uint, ticketType TicketType, title string, content string) (ticket *Ticket) {
	return &Ticket{
		UserId:  userId,
		Type:    int(ticketType),
		Title:   title,
		Content: content,
	}
}

func (ticketOperation *TicketOperation) SaveTicket(ticket *Ticket) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), ticketOperation.queryTimeout)
	defer cancel()
	if ticket.ID == 0 {
		return ticketOperation.db.WithContext(ctx).Create(ticket).Error
	}
	return ticketOperation.db.Model(&ticket).Save(ticket).Error
}

func (ticketOperation *TicketOperation) GetTicket(id uint) (ticket *Ticket, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), ticketOperation.queryTimeout)
	defer cancel()
	ticket = &Ticket{}
	err = ticketOperation.db.WithContext(ctx).Preload("User").First(ticket, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = ErrTicketNotFound
		}
	}
	return
}

func (ticketOperation *TicketOperation) GetTickets(page, pageSize int) (tickets []*Ticket, total int64, err error) {
	tickets = make([]*Ticket, 0, pageSize)
	ctx, cancel := context.WithTimeout(context.Background(), ticketOperation.queryTimeout)
	defer cancel()
	ticketOperation.db.WithContext(ctx).Model(&Ticket{}).Select("id").Count(&total)
	err = ticketOperation.db.WithContext(ctx).Preload("User").Offset((page - 1) * pageSize).Order("created_at desc").Limit(pageSize).Find(&tickets).Error
	return
}

func (ticketOperation *TicketOperation) GetUserTickets(uid uint, page, pageSize int) (tickets []*UserTicket, total int64, err error) {
	tickets = make([]*UserTicket, 0, pageSize)
	ctx, cancel := context.WithTimeout(context.Background(), ticketOperation.queryTimeout)
	defer cancel()
	ticketOperation.db.WithContext(ctx).Model(&Ticket{}).Select("id").Where("user_id = ?", uid).Count(&total)
	err = ticketOperation.db.WithContext(ctx).Model(&Ticket{}).Offset((page-1)*pageSize).Order("created_at desc").Where("user_id = ?", uid).Limit(pageSize).Find(&tickets).Error
	return
}

func (ticketOperation *TicketOperation) CloseTicket(ticket *Ticket, closer int, content string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), ticketOperation.queryTimeout)
	defer cancel()
	if ticket.Closer != 0 {
		return ErrTicketAlreadyClosed
	}
	return ticketOperation.db.Clauses(clause.Locking{Strength: "UPDATE"}).WithContext(ctx).Model(&Ticket{ID: ticket.ID}).Updates(&Ticket{Closer: closer, Reply: content}).Error
}

func (ticketOperation *TicketOperation) DeleteTicket(id uint) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), ticketOperation.queryTimeout)
	defer cancel()
	result := ticketOperation.db.WithContext(ctx).Delete(&Ticket{}, id)
	if result.RowsAffected == 0 {
		return ErrTicketNotFound
	}
	return result.Error
}
