// Package database
package database

import (
	"context"
	"errors"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
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

func (ticketOperation *TicketOperation) NewTicket(opener int, ticketType TicketType, title string, content string) (ticket *Ticket) {
	return &Ticket{
		Opener:  opener,
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
	err = ticketOperation.db.WithContext(ctx).First(ticket, id).Error
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
	err = ticketOperation.db.WithContext(ctx).Offset((page - 1) * pageSize).Order("created_at desc").Limit(pageSize).Find(&tickets).Error
	return
}

func (ticketOperation *TicketOperation) GetUserTickets(cid, page, pageSize int) (tickets []*Ticket, total int64, err error) {
	tickets = make([]*Ticket, 0, pageSize)
	ctx, cancel := context.WithTimeout(context.Background(), ticketOperation.queryTimeout)
	defer cancel()
	ticketOperation.db.WithContext(ctx).Model(&Ticket{}).Select("id").Where("opener = ?", cid).Count(&total)
	err = ticketOperation.db.WithContext(ctx).Offset((page-1)*pageSize).Order("created_at desc").Where("opener = ?", cid).Limit(pageSize).Find(&tickets).Error
	return
}

func (ticketOperation *TicketOperation) CloseTicket(ticketId uint, closer int, content string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), ticketOperation.queryTimeout)
	defer cancel()
	ticket, err := ticketOperation.GetTicket(ticketId)
	if err != nil {
		return err
	}
	if ticket.Closer != 0 {
		return ErrTicketAlreadyClosed
	}
	return ticketOperation.db.Clauses(clause.Locking{Strength: "UPDATE"}).WithContext(ctx).Model(&Ticket{ID: ticketId}).Updates(&Ticket{Closer: closer, Reply: content}).Error
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
