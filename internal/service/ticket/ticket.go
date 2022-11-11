package ticket

import (
	"context"
	"errors"
	"samet-avci/gowit/dto/request"
	"samet-avci/gowit/dto/response"
	repository "samet-avci/gowit/internal/repository/ticket"
	models "samet-avci/gowit/models/ticket"
	"strconv"
)

type ITicketService interface {
	CreateTicketOption(ctx context.Context, request request.NewTicketDTO) (response.NewTicketDTO, error)
	GetTicket(ctx context.Context, id string) (response.GetTicketDTO, error)
	PurchaseFromTicketOption(ctx context.Context, id string, request request.PurchaseFromTicketOptionsDTO) error
}

type TicketService struct {
	repository repository.ITicketRepository
}

func (s *TicketService) CreateTicketOption(ctx context.Context, request request.NewTicketDTO) (response.NewTicketDTO, error) {

	Isduplicated := s.repository.IsDuplicate(request.Name)
	if Isduplicated {
		return response.NewTicketDTO{}, errors.New("Ticket is already exist")
	}
	newTicket := models.Ticket{
		Name:       request.Name,
		Desc:       request.Desc,
		Allocation: request.Allocation,
	}

	err := s.repository.CreateTicket(newTicket)
	if err != nil {
		return response.NewTicketDTO{}, errors.New("Cannot create new ticket, please try again")
	}
	var resp response.NewTicketDTO
	resp.Convert(&newTicket)

	return resp, nil
}

func (s *TicketService) GetTicket(ctx context.Context, id string) (response.GetTicketDTO, error) {
	intID, _ := strconv.Atoi(id)

	ticket, err := s.repository.GetTicket(intID)
	if err != nil {
		return response.GetTicketDTO{}, errors.New("Ticket is not found")
	}
	var resp response.GetTicketDTO
	resp.Convert(&ticket)
	return resp, nil
}

func (s *TicketService) PurchaseFromTicketOption(ctx context.Context, id string, request request.PurchaseFromTicketOptionsDTO) error {
	intID, _ := strconv.Atoi(id)

	s.repository.Begin(ctx)
	allocation, err := s.repository.IsHaveAllocation(intID)
	if err != nil {
		s.repository.Rollback()
		return errors.New("error on ocurred while getting ticket")
	}

	if int(request.Quantity) > allocation {
		s.repository.Rollback()
		return errors.New("We don't have enough ticket")
	}

	newAllocation := allocation - int(request.Quantity)
	err = s.repository.SellTicket(newAllocation)
	if err != nil {
		return errors.New("error on ocurred while getting ticket")
	}
	s.repository.Commit()
	return nil
}
