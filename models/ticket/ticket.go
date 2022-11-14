package models

import (
	"github.com/google/uuid"
)

type Ticket struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Allocation uint   `json:"allocation"`
}

type SoldTicket struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	Quantity uint      `json:"quantity"`
}
