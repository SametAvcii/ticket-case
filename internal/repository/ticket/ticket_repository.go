package repository

import (
	"context"
	models "samet-avci/gowit/models/ticket"

	"gorm.io/gorm"
)

type ITicketRepository interface {
	IsDuplicate(ctx context.Context, name string) bool
	CreateTicket(ctx context.Context, ticket *models.Ticket) error
	GetTicketByID(ctx context.Context, id int) (models.Ticket, error)
	UpdateAllocation(ctx context.Context, allocation, id uint) error
	SaveSoldTicket(ctx context.Context, soldTicket models.SoldTicket) error
}

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) IsDuplicate(ctx context.Context, name string) bool {

	var ticket models.Ticket
	_ = r.db.WithContext(ctx).Model(&models.Ticket{}).Debug().Where("name = ?", name).First(&ticket)
	if ticket.ID > 0 {
		return true
	}
	return false
}

func (r *TicketRepository) CreateTicket(ctx context.Context, ticket *models.Ticket) error {

	err := r.db.WithContext(ctx).Create(&ticket).Error
	return err
}

func (r *TicketRepository) GetTicketByID(ctx context.Context, id int) (models.Ticket, error) {

	var ticket models.Ticket
	err := r.db.WithContext(ctx).Model(&models.Ticket{}).Where("id = ?", id).First(&ticket).Error
	return ticket, err
}

func (r *TicketRepository) UpdateAllocation(ctx context.Context, allocation, id uint) error {

	err := r.db.WithContext(ctx).Model(&models.Ticket{}).Where("id = ?", id).Update("allocation", allocation).Error
	return err
}

func (r *TicketRepository) SaveSoldTicket(ctx context.Context, soldTicket models.SoldTicket) error {

	err := r.db.WithContext(ctx).Create(&soldTicket).Error
	return err
}
