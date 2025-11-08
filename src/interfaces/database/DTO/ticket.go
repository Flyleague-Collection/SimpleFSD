// Package DTO
package DTO

import "time"

type UserTicket struct {
	ID        uint      `json:"id"`
	Type      int       `json:"type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Reply     string    `json:"reply"`
	CreatedAt time.Time `json:"open_at"`
	UpdatedAt time.Time `json:"close_at"`
}
