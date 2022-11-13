package router

import (
	"samet-avci/gowit/internal/app/handler/ticket"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Init(router *echo.Echo, tx *gorm.DB) {
	ticketHandler := ticket.TicketInit(tx)

	router.POST("ticket_options", ticketHandler.CreateTicketOption)
	router.GET("ticket/:id", ticketHandler.GetTicket)
	router.POST("ticket_options/:id/purchases", ticketHandler.PurchaseFromTicketOption)
}
