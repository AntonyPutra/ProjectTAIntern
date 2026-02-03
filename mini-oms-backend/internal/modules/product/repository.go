package product

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

func (r *Repository) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Find(&products).Error
	return products, err
}

func (r *Repository) FindByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, "id = ?", id).Error
	return &product, err
}

// FindByIDWithLock finds a product and locks the row for update (Must be called within a transaction)
func (r *Repository) FindByIDWithLock(tx *gorm.DB, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *Repository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *Repository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *Repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Product{}, "id = ?", id).Error
}
