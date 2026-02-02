package order

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

func (r *Repository) FindAll() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("User").Preload("OrderItems").Preload("Payment").Find(&orders).Error
	return orders, err
}

func (r *Repository) FindByUserID(userID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("OrderItems").Preload("Payment").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (r *Repository) FindByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("User").Preload("OrderItems.Product").Preload("Payment").First(&order, "id = ?", id).Error
	return &order, err
}

func (r *Repository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *Repository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *Repository) FindProductByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, "id = ?", id).Error
	return &product, err
}

func (r *Repository) UpdateProduct(product *models.Product) error {
	return r.db.Save(product).Error
}
