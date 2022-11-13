package repository

import (
	"errors"
	models "samet-avci/gowit/models/ticket"
	"sync"

	"gorm.io/gorm"
)

type ITicketRepository interface {
	IsDuplicate(name string) bool
	CreateTicket(ticket *models.Ticket) error
	GetTicket(id int) (models.Ticket, error)
	SellTicket(quantity, ID int) error
	SaveSoldTicket(soldTicket models.SoldTicket) error
}

type TicketRepository struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) IsDuplicate(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	var ticket models.Ticket
	_ = r.db.Model(&models.Ticket{}).Where("name = ?", name).First(&ticket)
	if ticket.ID > 0 {
		return true
	}
	return false
}

func (r *TicketRepository) CreateTicket(ticket *models.Ticket) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	err := r.db.Create(&ticket).Error
	return err
}

func (r *TicketRepository) GetTicket(id int) (models.Ticket, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var ticket models.Ticket
	err := r.db.Model(&models.Ticket{}).Where("id = ?", id).First(&ticket).Error
	return ticket, err
}

func (r *TicketRepository) SellTicket(quantity, ID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var ticket models.Ticket
	err := r.db.Model(&models.Ticket{}).Where("id= ?", ID).Where("allocation >= ?", quantity).First(&ticket).
		Update("allocation", gorm.Expr("allocation - ?", quantity)).Error
	if err != nil {
		return errors.New("not enough allocations")
	}
	return nil
}

func (r *TicketRepository) SaveSoldTicket(soldTicket models.SoldTicket) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	err := r.db.Create(&soldTicket).Error
	return err
}
