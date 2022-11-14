package ticket

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"samet-avci/gowit/dto/request"
	"samet-avci/gowit/dto/response"
	"samet-avci/gowit/mocks"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetTicket(t *testing.T) {
	t.Run("successfull get ticket", func(t *testing.T) {
		serviceMock := mocks.NewITicketService(t)

		serviceMock.On("GetTicket", context.Background(), "").Return(response.GetTicketDTO{
			ID:         1,
			Name:       "samet",
			Desc:       "successfull desc",
			Allocation: 100,
		}, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/ticket/1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler := NewTicketHandler(serviceMock)
		// Assertions
		if assert.NoError(t, handler.GetTicket(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("failed get ticket", func(t *testing.T) {
		serviceMock := mocks.NewITicketService(t)

		serviceMock.On("GetTicket", context.Background(), "").Return(response.GetTicketDTO{}, errors.New("Ticket is not found"))

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/ticket/5", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler := NewTicketHandler(serviceMock)
		// Assertions
		if assert.NoError(t, handler.GetTicket(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

}
func TestPurchaseFromTicketOption(t *testing.T) {
	t.Run("successfull", func(t *testing.T) {
		request := request.PurchaseFromTicketOptionsDTO{
			Quantity: 10,
			UserID:   uuid.MustParse("406c1d05-bbb2-4e94-b183-7d208c2692e1"),
		}

		serviceMock := mocks.NewITicketService(t)
		serviceMock.On("PurchaseFromTicketOption", context.Background(), "", request).Return(nil)

		ticketJson := `{"quantity": 10, "user_id": "406c1d05-bbb2-4e94-b183-7d208c2692e1"}`

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/ticket_options/1/purchases", strings.NewReader(ticketJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler := NewTicketHandler(serviceMock)

		if assert.NoError(t, handler.PurchaseFromTicketOption(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("not enough allocation", func(t *testing.T) {
		request := request.PurchaseFromTicketOptionsDTO{
			Quantity: 10,
			UserID:   uuid.MustParse("406c1d05-bbb2-4e94-b183-7d208c2692e1"),
		}

		serviceMock := mocks.NewITicketService(t)
		serviceMock.On("PurchaseFromTicketOption", context.Background(), "", request).Return(errors.New("error an occured while ticket selling.not enough allocations"))

		ticketJson := `{"quantity": 10, "user_id": "406c1d05-bbb2-4e94-b183-7d208c2692e1"}`

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/ticket_options/1/purchases", strings.NewReader(ticketJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler := NewTicketHandler(serviceMock)

		if assert.NoError(t, handler.PurchaseFromTicketOption(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

}

func TestCreateTicketOption(t *testing.T) {
	t.Run("successfull", func(t *testing.T) {
		request := request.NewTicketDTO{
			Name:       "successfull create",
			Desc:       "successfull desc",
			Allocation: 100,
		}

		serviceMock := mocks.NewITicketService(t)
		serviceMock.On("CreateTicketOption", context.Background(), request).Return(response.NewTicketDTO{
			ID:         1,
			Name:       "successfull create",
			Desc:       "successfull desc",
			Allocation: 100,
		}, nil)

		ticketJson := `{
			"name":"successfull create",
			"desc":"successfull desc",
			"allocation":100
		}`

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/ticket_options", strings.NewReader(ticketJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler := NewTicketHandler(serviceMock)

		if assert.NoError(t, handler.CreateTicketOption(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("service failed", func(t *testing.T) {
		request := request.NewTicketDTO{
			Name:       "successfull create",
			Desc:       "successfull desc",
			Allocation: 100,
		}

		serviceMock := mocks.NewITicketService(t)
		serviceMock.On("CreateTicketOption", context.Background(), request).Return(response.NewTicketDTO{}, errors.New("Cannot create new ticket, please try again"))

		ticketJson := `{
			"name":"successfull create",
			"desc":"successfull desc",
			"allocation":100
		}`

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/ticket_options", strings.NewReader(ticketJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler := NewTicketHandler(serviceMock)

		if assert.NoError(t, handler.CreateTicketOption(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("allocation is required", func(t *testing.T) {
		ticketJson := `{
			"name":"failed create",
			"desc":"failed desc"
		}`

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/ticket_options", strings.NewReader(ticketJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler := NewTicketHandler(nil)

		if assert.NoError(t, handler.CreateTicketOption(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("name is required", func(t *testing.T) {
		ticketJson := `{
			"desc":"failed desc",
			"allocation":100
		}`

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/ticket_options", strings.NewReader(ticketJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler := NewTicketHandler(nil)

		if assert.NoError(t, handler.CreateTicketOption(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

}
