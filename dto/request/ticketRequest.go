package request

import (
	"github.com/google/uuid"
)

type NewTicketDTO struct {
	Name       string `json:"name" validate:"required"`
	Desc       string `json:"desc" validate:"omitempty"`
	Allocation uint   `json:"allocation" validate:"required"`
}

type PurchaseFromTicketOptionsDTO struct {
	Quantity uint      `json:"quantity" validate:"required"`
	UserID   uuid.UUID `json:"user_id" validate:"required"`
}
