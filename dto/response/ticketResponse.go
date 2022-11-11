package response

import models "samet-avci/gowit/models/ticket"

type NewTicketDTO struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Allocation uint   `json:"allocation"`
}

type GetTicketDTO struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Allocation uint   `json:"allocation"`
}

type PurchaseFromTicketOptionsDTO struct {
	Err error `json:"err"`
}

func (t *NewTicketDTO) Convert(ticket *models.Ticket) {
	t.Name = ticket.Name
	t.Desc = ticket.Desc
	t.Allocation = ticket.Allocation
	t.ID = ticket.ID
}

func (t *GetTicketDTO) Convert(ticket *models.Ticket) {
	t.Name = ticket.Name
	t.Desc = ticket.Desc
	t.Allocation = ticket.Allocation
	t.ID = ticket.ID
}
