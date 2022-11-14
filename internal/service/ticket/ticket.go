package ticket

import (
	"context"
	"errors"
	"samet-avci/gowit/dto/request"
	"samet-avci/gowit/dto/response"
	repository "samet-avci/gowit/internal/repository/ticket"
	models "samet-avci/gowit/models/ticket"
	"strconv"
	"sync"
)

type ITicketService interface {
	CreateTicketOption(ctx context.Context, request request.NewTicketDTO) (response.NewTicketDTO, error)
	GetTicket(ctx context.Context, id string) (response.GetTicketDTO, error)
	PurchaseFromTicketOption(ctx context.Context, id string, request request.PurchaseFromTicketOptionsDTO) error
	//WithTrx(*gorm.DB) *TicketService
}

type TicketService struct {
	repository repository.ITicketRepository
	m          sync.Mutex
}

func NewTicketService(repository repository.ITicketRepository) ITicketService {
	return &TicketService{repository: repository}
}

func (s *TicketService) CreateTicketOption(ctx context.Context, request request.NewTicketDTO) (response.NewTicketDTO, error) {
	s.m.Lock()
	defer s.m.Unlock()
	Isduplicated := s.repository.IsDuplicate(ctx, request.Name)
	if Isduplicated {
		return response.NewTicketDTO{}, errors.New("Ticket is already exist")
	}
	newTicket := models.Ticket{
		Name:       request.Name,
		Desc:       request.Desc,
		Allocation: request.Allocation,
	}

	err := s.repository.CreateTicket(ctx, &newTicket)
	if err != nil {
		return response.NewTicketDTO{}, errors.New("Cannot create new ticket, please try again")
	}
	var resp response.NewTicketDTO
	resp.Convert(&newTicket)

	return resp, nil
}

func (s *TicketService) GetTicket(ctx context.Context, id string) (response.GetTicketDTO, error) {
	s.m.Lock()
	defer s.m.Unlock()

	intID, _ := strconv.Atoi(id)

	ticket, err := s.repository.GetTicketByID(ctx, intID)
	if err != nil {
		return response.GetTicketDTO{}, errors.New("Ticket is not found")
	}
	var resp response.GetTicketDTO
	resp.Convert(&ticket)
	return resp, nil
}

func (s *TicketService) PurchaseFromTicketOption(ctx context.Context, id string, request request.PurchaseFromTicketOptionsDTO) error {

	s.m.Lock()
	defer s.m.Unlock()
	intID, _ := strconv.Atoi(id)

	ticket, err := s.repository.GetTicketByID(ctx, intID)
	if err == nil && ticket.Allocation < request.Quantity {
		msg := "not enough ticket for this quantity " + "enough ticket is " + strconv.Itoa(int(ticket.Allocation))
		return errors.New(msg)
	} else if err != nil {
		return errors.New("error an occured while ticket selling. " + err.Error())
	}

	allocation := ticket.Allocation - request.Quantity
	err = s.repository.UpdateAllocation(ctx, allocation, ticket.ID)
	if err != nil {
		return errors.New("error an occured while ticket purchase updated.")
	}

	soldTicket := models.SoldTicket{
		UserID:   request.UserID,
		Quantity: request.Quantity,
	}
	err = s.repository.SaveSoldTicket(ctx, soldTicket)
	if err != nil {
		return errors.New("error an occured while saving selled ticket")
	}
	return nil

}
