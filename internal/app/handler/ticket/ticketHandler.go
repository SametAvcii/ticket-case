package ticket

import (
	"net/http"
	"samet-avci/gowit/dto/request"
	ticketService "samet-avci/gowit/internal/service/ticket"
	"samet-avci/gowit/internal/utils/response"
	"samet-avci/gowit/internal/utils/validate"

	"github.com/labstack/echo/v4"
)

type ITicketHandler interface {
	CreateTicketOption(c echo.Context) error
	GetTicket(c echo.Context) error
	PurchaseFromTicketOption(c echo.Context) error
}

type TicketHandler struct {
	service ticketService.ITicketService
}

func NewTicketHandler(service ticketService.ITicketService) *TicketHandler {
	return &TicketHandler{service: service}
}

// @Summary POST Create New Ticket Option
// @Description Create New Ticket Option
// @Tags Tickets
// @Success 200 {object} response.NewTicketDTO
// @Failure 404
// @Accept json
// @Param body body request.NewTicketDTO false "ticket"
// @Router /ticket_options [post]
func (h *TicketHandler) CreateTicketOption(c echo.Context) error {
	var request request.NewTicketDTO
	if validate.Validator(&c, &request) != nil {
		return nil
	}
	ctx := c.Request().Context()
	newTicket, err := h.service.CreateTicketOption(ctx, request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, response.Response(http.StatusOK, newTicket))
}

// @Summary GET Get Ticket Option
// @Description Get Ticket Option With By ID
// @Tags Tickets
// @Success 200 {object} response.GetTicketDTO
// @Failure 404
// @Accept json
// @Router /ticket_options/:id [get]
func (h *TicketHandler) GetTicket(c echo.Context) error {

	id := c.Param("id")
	ctx := c.Request().Context()
	ticket, err := h.service.GetTicket(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, response.Response(http.StatusOK, ticket))
}

// @Summary POST Purchase From Ticket Option
// @Description Purchase From Ticket Option For Sell Ticket
// @Tags Tickets
// @Success 200
// @Failure 404
// @Accept json
// @Param body body request.PurchaseFromTicketOptionsDTO false "ticket"
// @Router /ticket_options/:id/purchases [post]
func (h *TicketHandler) PurchaseFromTicketOption(c echo.Context) error {
	id := c.Param("id")
	var request request.PurchaseFromTicketOptionsDTO
	if validate.Validator(&c, &request) != nil {
		return nil
	}

	ctx := c.Request().Context()
	err := h.service.PurchaseFromTicketOption(ctx, id, request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, response.ResponseStatus(http.StatusOK))

}
