package ticket

import (
	repository "samet-avci/gowit/internal/repository/ticket"
	ticketService "samet-avci/gowit/internal/service/ticket"

	"gorm.io/gorm"
)

func TicketInit(db *gorm.DB) ITicketHandler {
	ticketRepository := repository.NewTicketRepository(db)
	service := ticketService.NewTicketService(ticketRepository)
	handler := NewTicketHandler(service)
	return handler
}
