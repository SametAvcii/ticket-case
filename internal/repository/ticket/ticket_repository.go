package repository

import (
	"context"
	models "samet-avci/gowit/models/ticket"
	"sync"

	"gorm.io/gorm"
)

type Transactor interface {
	Begin(ctx context.Context) error
	Commit() error
	Rollback() error
}
type ITicketRepository interface {
	Transactor
	IsDuplicate(name string) bool
	CreateTicket(ticket models.Ticket) error
	GetTicket(id int) (models.Ticket, error)
	SellTicket(allocation int) error
	IsHaveAllocation(id int) (int, error)
}

type TicketRepository struct {
	db *gorm.DB
	mu *sync.RWMutex
}

func NewTicketRepository(db *gorm.DB, mu *sync.RWMutex) *TicketRepository {
	return &TicketRepository{db: db, mu: mu}
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

func (r *TicketRepository) CreateTicket(ticket models.Ticket) error {
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
func (r *TicketRepository) IsHaveAllocation(id int) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var allocation int
	err := r.db.Model(&models.Ticket{}).Select("allocation").Where("id = ?", id).First(&allocation).Error
	return allocation, err
}

func (r *TicketRepository) SellTicket(allocation int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	err := r.db.Model(&models.Ticket{}).Update("allocation", allocation).Error
	return err
}
