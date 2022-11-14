package ticket

import (
	"context"
	"errors"
	"samet-avci/gowit/dto/request"
	"samet-avci/gowit/dto/response"
	"samet-avci/gowit/mocks"
	models "samet-avci/gowit/models/ticket"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTicketOptionService(t *testing.T) {
	t.Run("successfull", func(t *testing.T) {
		repoMock := mocks.NewITicketRepository(t)

		request := request.NewTicketDTO{
			Name:       "success",
			Desc:       "success",
			Allocation: 100,
		}
		newTicket := models.Ticket{
			Name:       request.Name,
			Desc:       request.Desc,
			Allocation: request.Allocation,
		}

		ctx := context.Background()
		repoMock.On("IsDuplicate", ctx, request.Name).Return(false)
		repoMock.On("CreateTicket", ctx, &newTicket).Return(nil)

		service := NewTicketService(repoMock)
		newTicketDTO, err := service.CreateTicketOption(ctx, request)
		assert.Equal(t, newTicketDTO.Allocation, newTicket.Allocation)
		assert.Equal(t, err, nil)
	})

	t.Run("duplicate name", func(t *testing.T) {
		repoMock := mocks.NewITicketRepository(t)

		request := request.NewTicketDTO{
			Name:       "success",
			Desc:       "success",
			Allocation: 100,
		}

		ctx := context.Background()
		repoMock.On("IsDuplicate", ctx, request.Name).Return(true)
		service := NewTicketService(repoMock)
		_, err := service.CreateTicketOption(ctx, request)
		assert.Equal(t, err, errors.New("Ticket is already exist"))
	})

	t.Run("cannot create ticket", func(t *testing.T) {
		repoMock := mocks.NewITicketRepository(t)

		request := request.NewTicketDTO{
			Name:       "success",
			Desc:       "success",
			Allocation: 100,
		}
		newTicket := models.Ticket{
			Name:       request.Name,
			Desc:       request.Desc,
			Allocation: request.Allocation,
		}

		ctx := context.Background()
		repoMock.On("IsDuplicate", ctx, request.Name).Return(false)
		repoMock.On("CreateTicket", ctx, &newTicket).Return(errors.New("we have error"))

		service := NewTicketService(repoMock)
		_, err := service.CreateTicketOption(ctx, request)
		assert.Equal(t, err, errors.New("Cannot create new ticket, please try again"))
	})
}

func TestGetTicket(t *testing.T) {
	t.Run("successfull", func(t *testing.T) {
		repoMock := mocks.NewITicketRepository(t)

		ctx := context.Background()
		ticket := models.Ticket{
			Name:       "success",
			Desc:       "success",
			Allocation: 100,
		}
		ticket.ID = 1
		repoMock.On("GetTicketByID", ctx, 1).Return(ticket, nil)

		response := response.GetTicketDTO{
			ID:         ticket.ID,
			Name:       ticket.Name,
			Desc:       ticket.Desc,
			Allocation: ticket.Allocation,
		}
		service := NewTicketService(repoMock)
		resp, err := service.GetTicket(ctx, "1")
		assert.Equal(t, resp, response)
		assert.Equal(t, err, nil)
	})

	t.Run("failed", func(t *testing.T) {
		repoMock := mocks.NewITicketRepository(t)

		ctx := context.Background()
		repoMock.On("GetTicketByID", ctx, 1).Return(models.Ticket{}, errors.New("Ticket is not found"))

		service := NewTicketService(repoMock)
		response := response.GetTicketDTO{}
		resp, err := service.GetTicket(ctx, "1")
		assert.Equal(t, resp, response)
		assert.Equal(t, err, errors.New("Ticket is not found"))
	})
}
func TestPurchaseFromTicketOption(t *testing.T) {
	t.Run("successfull", func(t *testing.T) {
		repoMock := mocks.NewITicketRepository(t)

		request := request.PurchaseFromTicketOptionsDTO{
			Quantity: 10,
			UserID:   uuid.MustParse("406c1d05-bbb2-4e94-b183-7d208c2692e1"),
		}

		ticket := models.Ticket{
			ID:         1,
			Name:       "success",
			Desc:       "success",
			Allocation: 100,
		}
		ctx := context.Background()
		repoMock.On("GetTicketByID", ctx, 1).Return(ticket, nil)
		allocation := ticket.Allocation - request.Quantity
		repoMock.On("UpdateAllocation", ctx, allocation, ticket.ID).Return(nil)
		soldTicket := models.SoldTicket{
			UserID:   request.UserID,
			Quantity: request.Quantity,
		}
		repoMock.On("SaveSoldTicket", ctx, soldTicket).Return(nil)

		service := NewTicketService(repoMock)
		err := service.PurchaseFromTicketOption(ctx, "1", request)
		assert.Equal(t, err, nil)
	})
	t.Run("not enough allocation", func(t *testing.T) {
		repoMock := mocks.NewITicketRepository(t)

		request := request.PurchaseFromTicketOptionsDTO{
			Quantity: 101,
			UserID:   uuid.MustParse("406c1d05-bbb2-4e94-b183-7d208c2692e1"),
		}

		ticket := models.Ticket{
			ID:         1,
			Name:       "success",
			Desc:       "success",
			Allocation: 100,
		}

		ctx := context.Background()
		repoMock.On("GetTicketByID", ctx, 1).Return(ticket, nil)
		service := NewTicketService(repoMock)
		err := service.PurchaseFromTicketOption(ctx, "1", request)
		msg := "not enough ticket for this quantity " + "enough ticket is " + strconv.Itoa(int(ticket.Allocation))
		assert.Equal(t, err, errors.New(msg))
	})

	t.Run("err while get ticket", func(t *testing.T) {
		repoMock := mocks.NewITicketRepository(t)

		request := request.PurchaseFromTicketOptionsDTO{
			Quantity: 101,
			UserID:   uuid.MustParse("406c1d05-bbb2-4e94-b183-7d208c2692e1"),
		}

		ticket := models.Ticket{
			ID:         1,
			Name:       "success",
			Desc:       "success",
			Allocation: 100,
		}

		ctx := context.Background()
		repoMock.On("GetTicketByID", ctx, 1).Return(ticket, errors.New("err while get ticket"))
		service := NewTicketService(repoMock)
		err := service.PurchaseFromTicketOption(ctx, "1", request)
		assert.Equal(t, err, errors.New("error an occured while ticket selling. err while get ticket"))
	})

	t.Run("not saved sold", func(t *testing.T) {

		repoMock := mocks.NewITicketRepository(t)

		request := request.PurchaseFromTicketOptionsDTO{
			Quantity: 10,
			UserID:   uuid.MustParse("406c1d05-bbb2-4e94-b183-7d208c2692e1"),
		}

		ticket := models.Ticket{
			ID:         1,
			Name:       "success",
			Desc:       "success",
			Allocation: 100,
		}
		ctx := context.Background()
		repoMock.On("GetTicketByID", ctx, 1).Return(ticket, nil)
		allocation := ticket.Allocation - request.Quantity
		repoMock.On("UpdateAllocation", ctx, allocation, ticket.ID).Return(nil)
		soldTicket := models.SoldTicket{
			UserID:   request.UserID,
			Quantity: request.Quantity,
		}

		repoMock.On("SaveSoldTicket", ctx, soldTicket).Return(errors.New("cannot save ticket"))
		service := NewTicketService(repoMock)
		err := service.PurchaseFromTicketOption(ctx, "1", request)
		assert.Equal(t, err, errors.New("error an occured while saving selled ticket"))
	})

	t.Run("update allocation", func(t *testing.T) {
		repoMock := mocks.NewITicketRepository(t)

		request := request.PurchaseFromTicketOptionsDTO{
			Quantity: 10,
			UserID:   uuid.MustParse("406c1d05-bbb2-4e94-b183-7d208c2692e1"),
		}

		ticket := models.Ticket{
			ID:         1,
			Name:       "success",
			Desc:       "success",
			Allocation: 100,
		}
		ctx := context.Background()
		repoMock.On("GetTicketByID", ctx, 1).Return(ticket, nil)
		allocation := ticket.Allocation - request.Quantity
		repoMock.On("UpdateAllocation", ctx, allocation, ticket.ID).Return(errors.New("error an occured while ticket purchase updated."))

		service := NewTicketService(repoMock)
		err := service.PurchaseFromTicketOption(ctx, "1", request)
		assert.Equal(t, err, errors.New("error an occured while ticket purchase updated."))
	})

}
