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

func (h *TicketHandler) GetTicket(c echo.Context) error {

	id := c.Param("id")
	ctx := c.Request().Context()
	ticket, err := h.service.GetTicket(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, response.Response(http.StatusOK, ticket))
}

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
