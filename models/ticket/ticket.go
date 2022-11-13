package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Allocation uint   `json:"allocation"`
}

type SoldTicket struct {
	gorm.Model
	UserID   uuid.UUID `json:"user_id"`
	Quantity uint      `json:"quantity"`
}
