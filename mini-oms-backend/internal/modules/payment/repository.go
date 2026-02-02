package payment

import (
	"mini-oms-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByOrderID(orderID uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Order").First(&payment, "order_id = ?", orderID).Error
	return &payment, err
}

func (r *Repository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *Repository) FindOrderByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.First(&order, "id = ?", id).Error
	return &order, err
}
