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
	//WithTrx(*gorm.DB) *TicketService
}

type TicketService struct {
	repository repository.ITicketRepository
}

func NewTicketService(repository repository.ITicketRepository) ITicketService {
	return &TicketService{repository: repository}
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

	err := s.repository.CreateTicket(&newTicket)
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

	err := s.repository.SellTicket(int(request.Quantity), intID)
	if err != nil {
		return errors.New("error an occured while ticket selling. " + err.Error())
	}

	soldTicket := models.SoldTicket{
		UserID:   request.UserID,
		Quantity: request.Quantity,
	}
	err = s.repository.SaveSoldTicket(soldTicket)
	if err != nil {
		return errors.New("error an occured while saving selled ticket")
	}
	return nil

}
