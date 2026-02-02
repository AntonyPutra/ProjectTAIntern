package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
	CategoryID  *uuid.UUID     `gorm:"type:uuid;index" json:"category_id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:decimal(12,2);not null" json:"price"`
	Stock       int            `gorm:"type:integer;not null;default:0" json:"stock"`
	ImageURL    string         `gorm:"type:varchar(500)" json:"image_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// IsInStock checks if product has available stock
func (p *Product) IsInStock(quantity int) bool {
	return p.Stock >= quantity
}

// ReduceStock reduces product stock
func (p *Product) ReduceStock(quantity int) error {
	if !p.IsInStock(quantity) {
		return gorm.ErrRecordNotFound // or custom error
	}
	p.Stock -= quantity
	return nil
}
